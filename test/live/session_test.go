//go:build live

package live

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

// TestGreeting covers RFC 5730 section 2.4: the server sends a greeting
// immediately on connection, before any command.
func TestGreeting(t *testing.T) {
	client := connect(t)

	raw := client.Greeting()
	if len(raw) == 0 {
		t.Fatal("empty greeting")
	}
	if !strings.Contains(string(raw), "<greeting>") {
		t.Fatalf("greeting is not an EPP greeting: %s", truncate(string(raw)))
	}

	greeting, err := client.GreetingInfo()
	if err != nil {
		t.Fatalf("GreetingInfo failed: %v", err)
	}
	if greeting.ServerID == "" {
		t.Fatal("greeting has no svID")
	}
	if len(greeting.Versions) == 0 || greeting.Versions[0] != "1.0" {
		t.Fatalf("unexpected greeting versions: %v", greeting.Versions)
	}
	if len(greeting.SupportedObjects) == 0 {
		t.Fatal("greeting advertises no object URIs")
	}
	t.Logf("svID=%q objects=%v extensions=%v", greeting.ServerID, greeting.SupportedObjects, greeting.SupportedExtensions)
}

// TestHello covers RFC 5730 section 2.3: <hello> elicits a fresh greeting.
func TestHello(t *testing.T) {
	client := connect(t)

	greeting, err := client.Hello()
	if err != nil {
		t.Fatalf("hello failed: %v", err)
	}
	if greeting.ServerID == "" {
		t.Fatal("hello greeting has no svID")
	}
	if greeting.ServerDate == nil {
		t.Fatal("hello greeting has no svDate")
	}
	if time.Since(*greeting.ServerDate) > 24*time.Hour {
		t.Fatalf("hello svDate is implausible: %s", greeting.ServerDate)
	}
}

// TestLoginLogout covers RFC 5730 sections 2.9.1.1 and 2.9.1.2. Logout must
// succeed with result code 1500.
func TestLoginLogout(t *testing.T) {
	client := connect(t)

	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if err := client.Logout(); err != nil {
		t.Fatalf("logout failed (RFC 5730 requires 1500 to be a success): %v", err)
	}
}

// TestHelloAfterLogin confirms <hello> is valid at any point in a session.
func TestHelloAfterLogin(t *testing.T) {
	client := session(t)

	if _, err := client.Hello(); err != nil {
		t.Fatalf("hello after login failed: %v", err)
	}
}

// TestLoginWrongPassword covers RFC 5730 result code 2200.
func TestLoginWrongPassword(t *testing.T) {
	cfg := loadConfig(t)
	// RFC 5730 pwType is 6-16 characters; a longer value would be rejected as
	// a schema error (2001) instead of exercising authentication.
	cfg.Authentication.Password = "Wr0ngPass99"

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	err = client.Login()
	eppErr := asEPPError(t, err, "login with wrong password")
	if eppErr.Code != 2200 && eppErr.Code != constants.ResultAuthenticationError {
		t.Fatalf("expected 2200/2501 authentication error, got %d (%s)", eppErr.Code, eppErr.Message)
	}
	if strings.Contains(err.Error(), cfg.Authentication.Password) {
		t.Fatal("error string leaked the password")
	}
}

// TestLoginWrongClientID covers an unknown client identifier.
func TestLoginWrongClientID(t *testing.T) {
	cfg := loadConfig(t)
	// RFC 5730 clIDType is 3-16 characters.
	cfg.Authentication.Username = "nosuchreg" + tag()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	err = client.Login()
	eppErr := asEPPError(t, err, "login with wrong client ID")
	if eppErr.Code < 2000 {
		t.Fatalf("expected an error result, got %d (%s)", eppErr.Code, eppErr.Message)
	}
	t.Logf("wrong client ID rejected with %d: %s", eppErr.Code, eppErr.Message)
}

// TestDoubleLogin covers RFC 5730: a second login on an established session
// must be refused by the server.
func TestDoubleLogin(t *testing.T) {
	client := session(t)

	err := client.Login()
	if err == nil {
		t.Fatal("expected second login on the same session to be refused")
	}
	t.Logf("second login rejected: %v", err)
}

// TestCommandBeforeLogin covers RFC 5730 result code 2002, since only hello
// and login are valid before authentication.
func TestCommandBeforeLogin(t *testing.T) {
	client := connect(t)

	_, err := client.DomainCheck(domainCheckRequest(domainName("pre")))
	eppErr := asEPPError(t, err, "command before login")
	if eppErr.Code < 2000 {
		t.Fatalf("expected an error result, got %d", eppErr.Code)
	}
	t.Logf("command before login rejected with %d: %s", eppErr.Code, eppErr.Message)
}

// TestOperationAfterLogout confirms the session is unusable once ended.
func TestOperationAfterLogout(t *testing.T) {
	client := connect(t)
	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if err := client.Logout(); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	_, err := client.DomainCheck(domainCheckRequest(domainName("post")))
	if err == nil {
		t.Fatal("expected a command after logout to fail")
	}
	t.Logf("command after logout failed as expected: %v", err)
}

// TestOperationAfterClose covers the SDK's own closed-session guard.
func TestOperationAfterClose(t *testing.T) {
	client := connect(t)
	if err := client.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	_, err := client.DomainCheck(domainCheckRequest(domainName("closed")))
	if err == nil {
		t.Fatal("expected a command on a closed session to fail")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindSession {
		t.Fatalf("expected a session SDKError, got %T: %v", err, err)
	}
}

// TestDoubleClose confirms Close is idempotent.
func TestDoubleClose(t *testing.T) {
	client := connect(t)

	if err := client.Close(); err != nil {
		t.Fatalf("first close failed: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("second close should be a no-op, got: %v", err)
	}
}

// TestReconnect covers session recovery: a new session must be usable and
// must leave the original client untouched.
func TestReconnect(t *testing.T) {
	client := session(t)

	next, err := client.Reconnect()
	if err != nil {
		t.Fatalf("reconnect failed: %v", err)
	}
	defer func() { _ = next.Close() }()

	if next == client {
		t.Fatal("reconnect returned the same client")
	}
	if len(next.Greeting()) == 0 {
		t.Fatal("reconnected session has no greeting")
	}
	if err := next.Login(); err != nil {
		t.Fatalf("login on reconnected session failed: %v", err)
	}
	defer func() { _ = next.Logout() }()

	if _, err := next.DomainCheck(domainCheckRequest(domainName("recon"))); err != nil {
		t.Fatalf("command on reconnected session failed: %v", err)
	}
	// The original session must still work.
	if _, err := client.DomainCheck(domainCheckRequest(domainName("orig"))); err != nil {
		t.Fatalf("original session broke after reconnect: %v", err)
	}
}

// TestContextCancellation confirms an in-flight command honours cancellation
// and reports it as a cancellation rather than a protocol failure.
func TestContextCancellation(t *testing.T) {
	client := session(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.DomainCheckContext(ctx, domainCheckRequest(domainName("cancel")))
	if err == nil {
		t.Fatal("expected a cancelled context to fail the command")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindCancellation {
		t.Fatalf("expected a cancellation SDKError, got %T: %v", err, err)
	}
}

// TestContextDeadline confirms an impossible deadline is reported as a timeout.
func TestContextDeadline(t *testing.T) {
	client := session(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	_, err := client.DomainCheckContext(ctx, domainCheckRequest(domainName("deadline")))
	if err == nil {
		t.Fatal("expected an expired deadline to fail the command")
	}
	var sdkErr *epp.SDKError
	if !asSDKError(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTimeout {
		t.Fatalf("expected a timeout SDKError, got %T: %v", err, err)
	}
}

func truncate(value string) string {
	if len(value) <= 200 {
		return value
	}
	return value[:200] + "..."
}

// TestClientTRIDIsUniquePerCommand covers RFC 5730 section 2.5: the client
// transaction identifier must be unique per command so a caller can correlate
// a response, and reconcile an ambiguous transform against the registry's
// transaction record, without ambiguity.
func TestClientTRIDIsUniquePerCommand(t *testing.T) {
	client := session(t)

	seen := map[string]bool{}
	for i := 0; i < 25; i++ {
		resp, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("tr")}})
		if err != nil {
			t.Fatalf("domain check failed: %v", err)
		}
		if resp.ClientTRID == "" {
			t.Fatal("the registry echoed no clTRID, so the command sent none")
		}
		if seen[resp.ClientTRID] {
			t.Fatalf("clTRID %q was reused, so responses cannot be correlated", resp.ClientTRID)
		}
		seen[resp.ClientTRID] = true

		if resp.ServerTRID == "" {
			t.Fatal("the registry returned no svTRID")
		}
	}
}

// TestClientTRIDUniqueAcrossSessions confirms identifiers do not collide
// between two sessions of the same registrar running concurrently.
func TestClientTRIDUniqueAcrossSessions(t *testing.T) {
	first := session(t)

	second, err := first.Reconnect()
	if err != nil {
		t.Fatalf("reconnect failed: %v", err)
	}
	defer func() { _ = second.Close() }()
	if err := second.Login(); err != nil {
		t.Fatalf("login on second session failed: %v", err)
	}
	defer func() { _ = second.Logout() }()

	seen := map[string]bool{}
	for i := 0; i < 10; i++ {
		for _, client := range []*epp.Client{first, second} {
			resp, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("ts")}})
			if err != nil {
				t.Fatalf("domain check failed: %v", err)
			}
			if seen[resp.ClientTRID] {
				t.Fatalf("clTRID %q collided across sessions", resp.ClientTRID)
			}
			seen[resp.ClientTRID] = true
		}
	}
}

// TestLoginEmptyCredentials covers client-side credential validation: an empty
// client ID or password is unambiguously a configuration mistake, and the SDK
// must say so without a round trip that reports a schema error instead.
func TestLoginEmptyCredentials(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		username string
		password string
	}{
		{"empty password", "csc_egov_in_a", ""},
		{"empty client ID", "", "whatever99"},
		{"whitespace password", "csc_egov_in_a", "   "},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := loadConfig(t)
			cfg.Authentication.Username = testCase.username
			cfg.Authentication.Password = testCase.password

			client, err := epp.Connect(cfg)
			if err != nil {
				t.Fatalf("connect failed: %v", err)
			}
			defer func() { _ = client.Close() }()

			err = client.Login()
			if err == nil {
				t.Fatal("expected an empty credential to be rejected")
			}
			eppErr := asEPPError(t, err, "login with an empty credential")
			if !eppErr.IsValidationError() {
				t.Fatalf("empty credentials reached the registry and were rejected with %d: %s",
					eppErr.Code, eppErr.Message)
			}
		})
	}
}
