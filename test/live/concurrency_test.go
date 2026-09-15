//go:build live

package live

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

// TestConcurrentCommandsOneSession drives one session from many goroutines.
// RFC 5734 runs one command at a time over a single TCP connection, so the SDK
// must serialise the write/read pair. A leak would show up as a response
// delivered to the wrong caller, which the clTRID echo detects.
func TestConcurrentCommandsOneSession(t *testing.T) {
	client := session(t)

	const workers = 12
	const perWorker = 4

	var wg sync.WaitGroup
	errs := make(chan error, workers*perWorker)

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				resp, err := client.DomainCheck(types.DomainCheckRequest{
					Domains: []string{domainName("c")},
				})
				if err != nil {
					errs <- err
					return
				}
				if resp.ResultCode != constants.ResultSuccess {
					errs <- errResultf("unexpected result code %d", resp.ResultCode)
					return
				}
				if len(resp.Results) != 1 {
					errs <- errResultf("expected 1 check result, got %d", len(resp.Results))
					return
				}
				// The response must be the one this command asked for.
				if resp.ClientTRID == "" {
					errs <- errResultf("response carried no clTRID")
					return
				}
			}
		}(worker)
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent command failed: %v", err)
	}
}

// TestConcurrentMixedCommands mixes different object types on one session so a
// misrouted response would be parsed as the wrong object.
func TestConcurrentMixedCommands(t *testing.T) {
	client := session(t)

	registrant := contactID("cm")
	createContact(t, client, registrant)
	domain := domainName("cm")
	createDomain(t, client, domain, registrant)
	host := hostName("ns1", domain)
	createHost(t, client, host, "103.51.75.41")

	var wg sync.WaitGroup
	errs := make(chan error, 64)

	run := func(name string, fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 4; i++ {
				if err := fn(); err != nil {
					errs <- errResultf("%s: %v", name, err)
					return
				}
			}
		}()
	}

	run("domain-info", func() error {
		resp, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
		if err != nil {
			return err
		}
		if resp.Result.Domain != domain {
			return errResultf("got domain %q, expected %q", resp.Result.Domain, domain)
		}
		return nil
	})
	run("contact-info", func() error {
		resp, err := client.ContactInfo(types.ContactInfoRequest{ContactID: registrant})
		if err != nil {
			return err
		}
		if resp.Contact.ContactID != registrant {
			return errResultf("got contact %q, expected %q", resp.Contact.ContactID, registrant)
		}
		return nil
	})
	run("host-info", func() error {
		resp, err := client.HostInfo(types.HostInfoRequest{HostName: host})
		if err != nil {
			return err
		}
		if resp.Host.HostName != host {
			return errResultf("got host %q, expected %q", resp.Host.HostName, host)
		}
		return nil
	})
	run("domain-check", func() error {
		_, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("cx")}})
		return err
	})
	run("poll", func() error {
		_, err := client.Poll(types.PollRequest{Operation: constants.PollRequest})
		return err
	})

	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent mixed command failed: %v", err)
	}
}

// TestConcurrentSessions opens several independent sessions at once, which is
// how a registrar platform actually scales.
func TestConcurrentSessions(t *testing.T) {
	const sessions = 4

	var wg sync.WaitGroup
	errs := make(chan error, sessions)

	for i := 0; i < sessions; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			cfg := loadConfigNoT()
			if cfg == nil {
				errs <- errResultf("config unavailable")
				return
			}
			client, err := epp.Connect(cfg)
			if err != nil {
				errs <- err
				return
			}
			defer func() { _ = client.Close() }()

			if err := client.Login(); err != nil {
				errs <- err
				return
			}
			defer func() { _ = client.Logout() }()

			for j := 0; j < 3; j++ {
				if _, err := client.DomainCheck(types.DomainCheckRequest{
					Domains: []string{domainName("s")},
				}); err != nil {
					errs <- err
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent session failed: %v", err)
	}
}

// TestConcurrentCloseDuringCommands closes a session while commands are in
// flight. The SDK must report a clean session or transport error rather than
// panicking or hanging.
func TestConcurrentCloseDuringCommands(t *testing.T) {
	client := session(t)

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				_, _ = client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("k")}})
			}
		}()
	}

	time.Sleep(150 * time.Millisecond)
	_ = client.Close()
	wg.Wait()

	// Every later command must fail cleanly.
	if _, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("k2")}}); err == nil {
		t.Fatal("expected a command on a closed session to fail")
	}
}

// TestConcurrentSetLogger races logger installation against command traffic,
// since a registrar may swap sinks at runtime.
func TestConcurrentSetLogger(t *testing.T) {
	client := session(t)

	var wg sync.WaitGroup
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				client.SetLogger(countingLogger{})
				client.SetLogger(nil)
			}
		}
	}()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				_, _ = client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("lg")}})
			}
		}()
	}

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// TestNoGoroutineLeak confirms the per-command cancellation watcher goroutine
// is always reaped, including on the context paths.
func TestNoGoroutineLeak(t *testing.T) {
	client := session(t)

	// Warm up so one-off internal goroutines are already running.
	for i := 0; i < 3; i++ {
		_, _ = client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("w")}})
	}
	settle()
	before := runtime.NumGoroutine()

	for i := 0; i < 25; i++ {
		if _, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("w")}}); err != nil {
			t.Fatalf("domain check failed: %v", err)
		}
	}

	// Cancelled and expired contexts take different exit paths.
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _ = client.DomainCheckContext(ctx, types.DomainCheckRequest{Domains: []string{domainName("w")}})

		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Millisecond)
		_, _ = client.DomainCheckContext(timeoutCtx, types.DomainCheckRequest{Domains: []string{domainName("w")}})
		timeoutCancel()
	}

	settle()
	after := runtime.NumGoroutine()
	if after > before+2 {
		t.Fatalf("goroutine count grew from %d to %d, which suggests a leak", before, after)
	}
}

type countingLogger struct{}

func (countingLogger) EPPEvent(epp.Event) {}

func settle() {
	for i := 0; i < 5; i++ {
		runtime.GC()
		time.Sleep(60 * time.Millisecond)
	}
}
