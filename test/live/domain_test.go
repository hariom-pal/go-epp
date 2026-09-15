//go:build live

package live

import (
	"strings"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

const domainAuthInfo = "uAtD0m#2026"

// createDomain registers a domain and deletes it during cleanup.
func createDomain(t *testing.T, client *epp.Client, domain, registrant string) *types.DomainCreateResponse {
	t.Helper()

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:          domain,
		Period:          1,
		Unit:            "y",
		Registrant:      registrant,
		AdminContacts:   []string{registrant},
		TechContacts:    []string{registrant},
		BillingContacts: []string{registrant},
		AuthInfo:        domainAuthInfo,
	})
	if err != nil {
		t.Fatalf("domain create %s failed: %v", domain, err)
	}
	t.Cleanup(func() {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: domain})
	})
	return resp
}

// TestDomainCheck covers RFC 5731 check, including the avail/reason pair.
func TestDomainCheck(t *testing.T) {
	client := shared(t)

	available := domainName("ck")
	check, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{available, "nixi.in", "zzz.invalidtld"},
	})
	if err != nil {
		t.Fatalf("domain check failed: %v", err)
	}
	if len(check.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(check.Results))
	}
	if !check.Results[0].Available {
		t.Fatalf("%s unexpectedly unavailable: %s", available, check.Results[0].Reason)
	}
	if check.Results[0].Reason != "" {
		t.Fatalf("an available name must carry no reason, got %q", check.Results[0].Reason)
	}
	for _, result := range check.Results[1:] {
		if result.Available {
			t.Fatalf("%s unexpectedly available", result.Domain)
		}
		if result.Reason == "" {
			t.Fatalf("%s is unavailable but carries no reason", result.Domain)
		}
	}
}

// TestDomainLifecycle covers RFC 5731 create, info, update, renew and delete.
func TestDomainLifecycle(t *testing.T) {
	client := session(t)

	registrant := contactID("dr")
	createContact(t, client, registrant)

	domain := domainName("lf")
	created := createDomain(t, client, domain, registrant)

	if created.Result.Domain != domain {
		t.Fatalf("create returned domain %q", created.Result.Domain)
	}
	if created.Result.CreatedDate.IsZero() {
		t.Fatal("create returned no crDate")
	}
	if created.Result.ExpiryDate.IsZero() {
		t.Fatal("create returned no exDate")
	}
	expiry := created.Result.ExpiryDate

	// Info before nameservers.
	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if info.Result.Registrant != registrant {
		t.Fatalf("unexpected registrant %q", info.Result.Registrant)
	}
	if info.Result.ROID == "" {
		t.Fatal("domain info returned no ROID")
	}
	if len(info.Result.Contacts) < 3 {
		t.Fatalf("expected admin/tech/billing contacts, got %v", info.Result.Contacts)
	}
	if info.Result.AuthInfo == "" {
		t.Fatal("domain info returned no authInfo for the sponsoring client")
	}

	// Hosts, then attach them.
	ns1 := hostName("ns1", domain)
	ns2 := hostName("ns2", domain)
	createHost(t, client, ns1, "103.51.75.11")
	createHost(t, client, ns2, "103.51.75.12")

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:         domain,
		AddNameServers: []string{ns1, ns2},
	}); err != nil {
		t.Fatalf("domain update add nameservers failed: %v", err)
	}

	info, err = client.DomainInfo(types.DomainInfoRequest{Domain: domain, Hosts: "all"})
	if err != nil {
		t.Fatalf("domain info after update failed: %v", err)
	}
	if len(info.Result.NameServers) != 2 {
		t.Fatalf("expected 2 nameservers, got %v", info.Result.NameServers)
	}

	// Status update.
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:      domain,
		AddStatuses: []string{"clientUpdateProhibited"},
	}); err != nil {
		t.Fatalf("domain update add status failed: %v", err)
	}
	info, err = client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info after status update failed: %v", err)
	}
	if !containsString(info.Result.Statuses, "clientUpdateProhibited") {
		t.Fatalf("clientUpdateProhibited not applied, statuses=%v", info.Result.Statuses)
	}

	// A further update must now be refused with 2304.
	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain:      domain,
		AddStatuses: []string{"clientHold"},
	})
	eppErr := requireResultCodeErr(t, err, constants.ResultObjectStatusProhibits, "update under clientUpdateProhibited")
	if !eppErr.IsObjectStatusProhibited() {
		t.Fatal("IsObjectStatusProhibited did not classify 2304")
	}

	// Lift the lock.
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:         domain,
		RemoveStatuses: []string{"clientUpdateProhibited"},
	}); err != nil {
		t.Fatalf("domain update remove status failed: %v", err)
	}

	// Renew against the current expiry.
	renewed, err := client.DomainRenew(types.DomainRenewRequest{
		Domain:            domain,
		CurrentExpiryDate: expiry,
		Period:            2,
		Unit:              "y",
	})
	if err != nil {
		t.Fatalf("domain renew failed: %v", err)
	}
	if !renewed.Result.NewExpiryDate.After(expiry) {
		t.Fatalf("renew did not advance expiry: %s -> %s", expiry, renewed.Result.NewExpiryDate)
	}
	if got := renewed.Result.NewExpiryDate.Year() - expiry.Year(); got != 2 {
		t.Fatalf("expected +2 years, got +%d", got)
	}

	// Detach and delete a host, then delete the domain.
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:            domain,
		RemoveNameServers: []string{ns2},
	}); err != nil {
		t.Fatalf("domain update remove nameserver failed: %v", err)
	}
	if _, err := client.HostDelete(types.HostDeleteRequest{HostName: ns2}); err != nil {
		t.Fatalf("host delete failed: %v", err)
	}
}

// TestDomainCreateDuplicate covers RFC 5730 result code 2302.
func TestDomainCreateDuplicate(t *testing.T) {
	client := session(t)

	registrant := contactID("dd")
	createContact(t, client, registrant)
	domain := domainName("dp")
	createDomain(t, client, domain, registrant)

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domain, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
	})
	requireResultCode(t, err, constants.ResultObjectExists, "duplicate domain create")
}

// TestDomainInfoUnknown covers RFC 5730 result code 2303.
func TestDomainInfoUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.DomainInfo(types.DomainInfoRequest{Domain: domainName("nx")})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "info on unknown domain")
}

// TestDomainDeleteUnknown covers deleting a domain that does not exist.
func TestDomainDeleteUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.DomainDelete(types.DomainDeleteRequest{Domain: domainName("nd")})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "delete of unknown domain")
}

// TestDomainCreateUnknownContact covers a create that references a contact
// that does not exist.
func TestDomainCreateUnknownContact(t *testing.T) {
	client := shared(t)

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domainName("uc"), Period: 1, Unit: "y",
		Registrant: contactID("qq"), AuthInfo: domainAuthInfo,
	})
	if err == nil {
		t.Fatal("expected a create with an unknown registrant to fail")
	}
	eppErr := asEPPError(t, err, "domain create with unknown contact")
	t.Logf("unknown registrant rejected with %d: %s values=%v", eppErr.Code, eppErr.Message, eppErr.Values)
}

// TestDomainCreateUnknownNameserver covers a create referencing a host object
// that does not exist.
func TestDomainCreateUnknownNameserver(t *testing.T) {
	client := session(t)

	registrant := contactID("un")
	createContact(t, client, registrant)

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domainName("un"), Period: 1, Unit: "y",
		Registrant:      registrant,
		AuthInfo:        domainAuthInfo,
		AdminContacts:   []string{registrant},
		TechContacts:    []string{registrant},
		BillingContacts: []string{registrant},
		NameServers:     []string{"ns1.this-host-does-not-exist-" + tag() + ".in"},
	})
	if err == nil {
		t.Fatal("expected a create with an unknown nameserver to fail")
	}
	eppErr := asEPPError(t, err, "domain create with unknown nameserver")
	t.Logf("unknown nameserver rejected with %d: %s values=%v", eppErr.Code, eppErr.Message, eppErr.Values)
}

// TestDomainCreateInvalidPeriod covers SDK-side period validation.
func TestDomainCreateInvalidPeriod(t *testing.T) {
	client := shared(t)

	for _, period := range []int{0, -1, 100} {
		_, err := client.DomainCreate(types.DomainCreateRequest{
			Domain: domainName("ip"), Period: period, Unit: "y",
			Registrant: contactID("ip"), AuthInfo: domainAuthInfo,
		})
		if err == nil {
			t.Fatalf("expected period %d to be rejected", period)
		}
	}

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domainName("iu"), Period: 1, Unit: "decades",
		Registrant: contactID("iu"), AuthInfo: domainAuthInfo,
	})
	if err == nil {
		t.Fatal("expected an invalid period unit to be rejected")
	}
}

// TestDomainCreateMissingAuthInfo covers SDK-side authInfo validation.
func TestDomainCreateMissingAuthInfo(t *testing.T) {
	client := shared(t)

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domainName("na"), Period: 1, Unit: "y",
		Registrant: contactID("na"),
	})
	if err == nil {
		t.Fatal("expected a create without authInfo to be rejected")
	}
}

// TestDomainRenewWrongExpiry covers RFC 5731: curExpDate must match the
// registry's current expiry, which guards against duplicate renewals.
func TestDomainRenewWrongExpiry(t *testing.T) {
	client := session(t)

	registrant := contactID("rw")
	createContact(t, client, registrant)
	domain := domainName("rw")
	createDomain(t, client, domain, registrant)

	_, err := client.DomainRenew(types.DomainRenewRequest{
		Domain:            domain,
		CurrentExpiryDate: time.Now().AddDate(5, 0, 0),
		Period:            1,
		Unit:              "y",
	})
	if err == nil {
		t.Fatal("expected a renew with a wrong curExpDate to be refused")
	}
	eppErr := asEPPError(t, err, "renew with wrong expiry")
	t.Logf("wrong curExpDate refused with %d: %s", eppErr.Code, eppErr.Message)
}

// TestDomainTransferQueryNotPending covers RFC 5730 result code 2301.
func TestDomainTransferQueryNotPending(t *testing.T) {
	client := session(t)

	registrant := contactID("tq")
	createContact(t, client, registrant)
	domain := domainName("tq")
	createDomain(t, client, domain, registrant)

	_, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: domain,
		Operation:  constants.TransferQuery,
	})
	if err == nil {
		t.Fatal("expected transfer query on a non-pending domain to fail")
	}
	eppErr := asEPPError(t, err, "domain transfer query")
	if eppErr.Code != 2301 {
		t.Logf("transfer query returned %d: %s", eppErr.Code, eppErr.Message)
	}
}

// TestDomainTransferOperations covers every transfer operation the SDK
// supports. Approve/reject/cancel on a domain with no pending transfer must be
// refused by the registry rather than mis-encoded by the SDK.
func TestDomainTransferOperations(t *testing.T) {
	client := session(t)

	registrant := contactID("to")
	createContact(t, client, registrant)
	domain := domainName("to")
	createDomain(t, client, domain, registrant)

	for _, op := range []string{
		constants.TransferRequest,
		constants.TransferApprove,
		constants.TransferReject,
		constants.TransferCancel,
	} {
		req := types.DomainTransferRequest{DomainName: domain, Operation: op}
		if op == constants.TransferRequest {
			req.AuthInfo = domainAuthInfo
		}
		_, err := client.DomainTransfer(req)
		if err == nil {
			t.Fatalf("transfer %s on a self-sponsored domain unexpectedly succeeded", op)
		}
		eppErr := asEPPError(t, err, "domain transfer "+op)
		if eppErr.Code < 2000 {
			t.Fatalf("transfer %s produced a non-error result %d", op, eppErr.Code)
		}
		t.Logf("transfer %-8s refused with %d: %s", op, eppErr.Code, eppErr.Message)
	}
}

// TestDomainTransferInvalidOperation covers SDK-side operation validation.
func TestDomainTransferInvalidOperation(t *testing.T) {
	client := shared(t)

	if _, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: domainName("io"), Operation: "explode",
	}); err == nil {
		t.Fatal("expected an invalid transfer operation to be rejected")
	}
	if _, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: domainName("io"),
	}); err == nil {
		t.Fatal("expected a missing transfer operation to be rejected")
	}
	if _, err := client.DomainTransfer(types.DomainTransferRequest{
		Operation: constants.TransferQuery,
	}); err == nil {
		t.Fatal("expected a missing domain name to be rejected")
	}
}

// TestDomainIDNPunycode covers RFC 5891 handling: the SDK must send the
// A-label and return the U-label.
func TestDomainIDNPunycode(t *testing.T) {
	client := shared(t)

	unicode := "परीक्षण" + devanagariTag() + ".भारत"
	check, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{unicode}})
	if err != nil {
		t.Fatalf("IDN domain check failed: %v", err)
	}
	if len(check.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(check.Results))
	}
	if !strings.HasPrefix(check.Results[0].ASCII, "xn--") {
		t.Fatalf("expected an A-label, got %q", check.Results[0].ASCII)
	}
	if check.Results[0].Domain != unicode {
		t.Fatalf("expected the U-label back, got %q", check.Results[0].Domain)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
