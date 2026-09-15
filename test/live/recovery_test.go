//go:build live

package live

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

// TestRecoverAfterUnexpectedDisconnect covers the production case where the
// registry or an intermediary drops the connection. The failing command must
// report a transport error, and the caller must be able to recover by opening
// a new session and continuing.
func TestRecoverAfterUnexpectedDisconnect(t *testing.T) {
	client := session(t)

	// Confirm the session works first.
	if _, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("r1")}}); err != nil {
		t.Fatalf("pre-disconnect command failed: %v", err)
	}

	// Drop the connection underneath the client.
	if err := client.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	_, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("r2")}})
	if err == nil {
		t.Fatal("expected a command on a dropped connection to fail")
	}

	// Recover.
	recovered, err := client.Reconnect()
	if err != nil {
		t.Fatalf("reconnect after disconnect failed: %v", err)
	}
	defer func() { _ = recovered.Close() }()
	if err := recovered.Login(); err != nil {
		t.Fatalf("login after reconnect failed: %v", err)
	}
	defer func() { _ = recovered.Logout() }()

	if _, err := recovered.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("r3")}}); err != nil {
		t.Fatalf("command after recovery failed: %v", err)
	}
}

// TestReconnectPreservesLogger confirms an event sink survives recovery, so a
// registrar does not lose telemetry after a reconnect.
func TestReconnectPreservesLogger(t *testing.T) {
	client := session(t)

	sink := &recordingLogger{}
	client.SetLogger(sink)

	recovered, err := client.Reconnect()
	if err != nil {
		t.Fatalf("reconnect failed: %v", err)
	}
	defer func() { _ = recovered.Close() }()

	if err := recovered.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}
	defer func() { _ = recovered.Logout() }()

	if _, err := recovered.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("rl")}}); err != nil {
		t.Fatalf("domain check failed: %v", err)
	}
	if sink.count() == 0 {
		t.Fatal("the recovered session emitted no events, so the logger was lost")
	}
}

// TestTransformIdempotencyReconciliation covers the reconciliation contract
// the SDK documents for a lost transform: after an uncertain outcome, the
// caller determines the truth from registry state, and a blind retry of a
// create surfaces as 2302 rather than silently duplicating.
func TestTransformIdempotencyReconciliation(t *testing.T) {
	client := session(t)

	registrant := contactID("ri")
	createContact(t, client, registrant)
	domain := domainName("ri")

	created := createDomain(t, client, domain, registrant)
	if created.Result.Domain != domain {
		t.Fatalf("create returned %q", created.Result.Domain)
	}

	// Reconcile: info is authoritative about whether the transform landed.
	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("reconciliation info failed: %v", err)
	}
	if info.Result.Domain != domain {
		t.Fatalf("reconciliation returned %q", info.Result.Domain)
	}

	// A blind retry must be refused, not silently duplicated.
	_, err = client.DomainCreate(types.DomainCreateRequest{
		Domain: domain, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
	})
	requireResultCode(t, err, constants.ResultObjectExists, "blind retry of a create")

	// The object must be unchanged by the refused retry.
	after, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("post-retry info failed: %v", err)
	}
	if !after.Result.ExpiryDate.Equal(*info.Result.ExpiryDate) {
		t.Fatalf("a refused retry changed the expiry: %s -> %s", info.Result.ExpiryDate, after.Result.ExpiryDate)
	}
}

// TestRenewIsNotIdempotent pins a property registrars must design around: a
// renew that is retried after an uncertain outcome extends the term twice
// unless the caller reconciles first. curExpDate is the guard RFC 5731
// provides, and it must be respected.
func TestRenewIsNotIdempotent(t *testing.T) {
	client := session(t)

	registrant := contactID("rn")
	createContact(t, client, registrant)
	domain := domainName("rn")
	created := createDomain(t, client, domain, registrant)

	expiry := created.Result.ExpiryDate

	first, err := client.DomainRenew(types.DomainRenewRequest{
		Domain: domain, CurrentExpiryDate: expiry, Period: 1, Unit: "y",
	})
	if err != nil {
		t.Fatalf("first renew failed: %v", err)
	}

	// Replaying the same command, with the now-stale curExpDate, must be
	// refused. That is what makes a retry safe.
	_, err = client.DomainRenew(types.DomainRenewRequest{
		Domain: domain, CurrentExpiryDate: expiry, Period: 1, Unit: "y",
	})
	if err == nil {
		t.Fatal("a replayed renew with a stale curExpDate must be refused")
	}

	after, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("info after renew failed: %v", err)
	}
	if !after.Result.ExpiryDate.Equal(first.Result.NewExpiryDate) {
		t.Fatalf("expiry drifted after the refused replay: %s vs %s",
			after.Result.ExpiryDate, first.Result.NewExpiryDate)
	}
}

// TestSessionSurvivesRegistryErrors confirms a business-level rejection does
// not poison the session: the next command must still work.
func TestSessionSurvivesRegistryErrors(t *testing.T) {
	client := session(t)

	for i := 0; i < 5; i++ {
		_, err := client.DomainInfo(types.DomainInfoRequest{Domain: domainName("se")})
		requireResultCode(t, err, constants.ResultObjectDoesNotExist, "info on unknown domain")

		if _, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{domainName("sc")}}); err != nil {
			t.Fatalf("iteration %d: session broken after an error response: %v", i, err)
		}
	}
}

type recordingLogger struct {
	events []epp.Event
}

func (l *recordingLogger) EPPEvent(event epp.Event) {
	l.events = append(l.events, event)
}

func (l *recordingLogger) count() int {
	return len(l.events)
}
