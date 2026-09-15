//go:build live

package live

import (
	"strings"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/extensions/fee"
	idnext "github.com/hariom-pal/go-epp/extensions/idn"
	"github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/extensions/rgp"
	"github.com/hariom-pal/go-epp/extensions/secdns"
	"github.com/hariom-pal/go-epp/types"
)

// testDS is a syntactically valid RFC 5910 DS record (algorithm 8, SHA-256).
func testDS() secdns.DSData {
	return secdns.DSData{
		KeyTag:     12345,
		Algorithm:  8,
		DigestType: 2,
		Digest:     "49FD46E6C4B45C55D4AC69CBD3CD34AC1AFE51DE0EE1E5A2B3C4D5E6F708192A",
	}
}

// TestSecDNSLifecycle covers RFC 5910: adding, reading back and removing DS
// data through the secDNS-1.1 extension.
func TestSecDNSLifecycle(t *testing.T) {
	client := session(t)

	registrant := contactID("sd")
	createContact(t, client, registrant)
	domain := domainName("sd")
	createDomain(t, client, domain, registrant)

	ds := testDS()

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		SecDNS: &secdns.UpdateRequest{
			Add: &secdns.UpdateAdd{Data: secdns.Data{DSData: []secdns.DSData{ds}}},
		},
	}); err != nil {
		t.Fatalf("secDNS add failed: %v", err)
	}

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if len(info.Result.SecDNS.DSData) != 1 {
		t.Fatalf("expected 1 DS record, got %d", len(info.Result.SecDNS.DSData))
	}
	got := info.Result.SecDNS.DSData[0]
	if got.KeyTag != ds.KeyTag || got.Algorithm != ds.Algorithm || got.DigestType != ds.DigestType {
		t.Fatalf("DS record did not round-trip: %+v", got)
	}
	if !strings.EqualFold(got.Digest, ds.Digest) {
		t.Fatalf("DS digest did not round-trip: %q", got.Digest)
	}

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		SecDNS: &secdns.UpdateRequest{
			Remove: &secdns.UpdateRemove{DSData: []secdns.DSData{ds}},
		},
	}); err != nil {
		t.Fatalf("secDNS remove failed: %v", err)
	}

	info, err = client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info after remove failed: %v", err)
	}
	if len(info.Result.SecDNS.DSData) != 0 {
		t.Fatalf("expected no DS records after remove, got %d", len(info.Result.SecDNS.DSData))
	}
}

// TestSecDNSOnCreate covers RFC 5910 DS data supplied at domain create time.
func TestSecDNSOnCreate(t *testing.T) {
	client := session(t)

	registrant := contactID("sc")
	createContact(t, client, registrant)
	domain := domainName("sc")

	if _, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domain, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
		SecDNS:          &secdns.CreateRequest{Data: secdns.Data{DSData: []secdns.DSData{testDS()}}},
	}); err != nil {
		t.Fatalf("domain create with secDNS failed: %v", err)
	}
	t.Cleanup(func() {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: domain})
	})

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if len(info.Result.SecDNS.DSData) != 1 {
		t.Fatalf("expected DS data from create, got %d records", len(info.Result.SecDNS.DSData))
	}
}

// TestSecDNSRemoveAll covers the RFC 5910 <secDNS:rem><secDNS:all> form.
func TestSecDNSRemoveAll(t *testing.T) {
	client := session(t)

	registrant := contactID("sa")
	createContact(t, client, registrant)
	domain := domainName("sa")
	createDomain(t, client, domain, registrant)

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		SecDNS: &secdns.UpdateRequest{
			Add: &secdns.UpdateAdd{Data: secdns.Data{DSData: []secdns.DSData{testDS()}}},
		},
	}); err != nil {
		t.Fatalf("secDNS add failed: %v", err)
	}

	all := true
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		SecDNS: &secdns.UpdateRequest{Remove: &secdns.UpdateRemove{All: &all}},
	}); err != nil {
		t.Fatalf("secDNS remove-all failed: %v", err)
	}

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if len(info.Result.SecDNS.DSData) != 0 {
		t.Fatalf("remove-all left %d DS records", len(info.Result.SecDNS.DSData))
	}
}

// TestSecDNSInvalidDigest covers RFC 5910 digest validation: a SHA-256 digest
// must be 64 hex characters, so a short one must be refused.
func TestSecDNSInvalidDigest(t *testing.T) {
	client := session(t)

	registrant := contactID("si")
	createContact(t, client, registrant)
	domain := domainName("si")
	createDomain(t, client, domain, registrant)

	_, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		SecDNS: &secdns.UpdateRequest{
			Add: &secdns.UpdateAdd{Data: secdns.Data{DSData: []secdns.DSData{{
				KeyTag: 12345, Algorithm: 8, DigestType: 2, Digest: "DEADBEEF",
			}}}},
		},
	})
	if err == nil {
		t.Fatal("expected an invalid DS digest to be refused")
	}
	eppErr := asEPPError(t, err, "secDNS invalid digest")
	t.Logf("invalid digest refused with %d: %s values=%v", eppErr.Code, eppErr.Message, eppErr.Values)
}

// TestFeeCheck covers RFC 8748-style fee data using the fee-0.7 profile that
// this registry advertises.
func TestFeeCheck(t *testing.T) {
	client := shared(t)

	domain := domainName("fe")
	resp, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{domain},
		Fee: &fee.CheckRequest{Domains: []fee.CheckDomain{{
			Name:     domain,
			Currency: "INR",
			Command:  fee.Command{Name: fee.CommandCreate},
			Period:   &fee.Period{Value: 1, Unit: "y"},
		}}},
	})
	if err != nil {
		t.Fatalf("fee check failed: %v", err)
	}
	if len(resp.Fee.Results) != 1 {
		t.Fatalf("expected 1 fee result, got %d", len(resp.Fee.Results))
	}
	result := resp.Fee.Results[0]
	if result.Currency == "" {
		t.Fatal("fee result carries no currency")
	}
	if len(result.Fees) == 0 {
		t.Fatal("fee result carries no fee amount")
	}
	t.Logf("fee: %s %s command=%s class=%s", result.Fees[0].Amount, result.Currency, result.Command.Name, result.Class)
}

// TestFeeCheckCommands covers fee data for each transform command.
func TestFeeCheckCommands(t *testing.T) {
	client := shared(t)

	for _, command := range []string{fee.CommandCreate, fee.CommandRenew, fee.CommandTransfer, fee.CommandRestore} {
		domain := domainName("f" + command[:2])
		resp, err := client.DomainCheck(types.DomainCheckRequest{
			Domains: []string{domain},
			Fee: &fee.CheckRequest{Domains: []fee.CheckDomain{{
				Name:     domain,
				Currency: "INR",
				Command:  fee.Command{Name: command},
				Period:   &fee.Period{Value: 1, Unit: "y"},
			}}},
		})
		if err != nil {
			t.Fatalf("fee check for %s failed: %v", command, err)
		}
		if len(resp.Fee.Results) == 0 {
			t.Fatalf("fee check for %s returned no results", command)
		}
		amounts := resp.Fee.Results[0].Fees
		if len(amounts) == 0 {
			t.Logf("fee for %-8s: none returned (reason=%q)", command, resp.Fee.Results[0].Reason)
			continue
		}
		t.Logf("fee for %-8s: %s %s", command, amounts[0].Amount, resp.Fee.Results[0].Currency)
	}
}

// TestFeeOnCreate covers fee data returned by a transform command.
func TestFeeOnCreate(t *testing.T) {
	client := session(t)

	registrant := contactID("fc")
	createContact(t, client, registrant)
	domain := domainName("fc")

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domain, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
		Fee:             &fee.TransformRequest{Currency: "INR", Fees: []fee.Fee{{Amount: "500"}}},
	})
	if err != nil {
		t.Fatalf("domain create with fee extension failed: %v", err)
	}
	t.Cleanup(func() {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: domain})
	})

	if len(resp.Result.Fee.Fees) == 0 {
		t.Fatal("create response carried no fee data")
	}
	t.Logf("create fee: %s %s", resp.Result.Fee.Fees[0].Amount, resp.Result.Fee.Currency)
}

// TestIDNCreate covers the idn-1.0 extension: registries require the language
// table tag for an IDN registration.
func TestIDNCreate(t *testing.T) {
	client := session(t)

	registrant := contactID("in")
	createContact(t, client, registrant)

	unicode := "जांच" + devanagariTag() + ".भारत"

	check, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{unicode}})
	if err != nil {
		t.Fatalf("IDN check failed: %v", err)
	}
	if !check.Results[0].Available {
		t.Skipf("IDN test name unavailable: %s", check.Results[0].Reason)
	}

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: unicode, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
		IDN:             &idnext.CreateRequest{Data: idnext.Data{Table: "hi", UName: unicode}},
	})
	if err != nil {
		t.Fatalf("IDN domain create failed: %v", err)
	}
	t.Cleanup(func() {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: unicode})
	})
	if resp.Result.Domain != unicode {
		t.Fatalf("create returned %q, expected the U-label %q", resp.Result.Domain, unicode)
	}

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: unicode})
	if err != nil {
		t.Fatalf("IDN domain info failed: %v", err)
	}
	if info.Result.ASCII == "" || !strings.HasPrefix(info.Result.ASCII, "xn--") {
		t.Fatalf("expected an A-label in info, got %q", info.Result.ASCII)
	}
	if info.Result.IDN.Table == "" {
		t.Fatal("domain info returned no IDN table for an IDN registration")
	}
	if info.Result.IDN.UName != unicode {
		t.Fatalf("expected the registry U-label %q, got %q", unicode, info.Result.IDN.UName)
	}
	t.Logf("IDN registered: %s (%s) table=%q uname=%q",
		info.Result.Domain, info.Result.ASCII, info.Result.IDN.Table, info.Result.IDN.UName)
}

// TestIDNCreateWithoutTable covers the registry requirement that an IDN create
// carry a table tag, and confirms the SDK surfaces the reason.
func TestIDNCreateWithoutTable(t *testing.T) {
	client := session(t)

	registrant := contactID("iw")
	createContact(t, client, registrant)

	unicode := "परख" + devanagariTag() + ".भारत"
	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: unicode, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
	})
	if err == nil {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: unicode})
		t.Skip("registry accepted an IDN create without a table tag")
	}
	eppErr := asEPPError(t, err, "IDN create without table")
	if len(eppErr.Values) == 0 {
		t.Fatal("IDN rejection carried no diagnostic detail")
	}
	t.Logf("IDN without table refused with %d: %v", eppErr.Code, eppErr.Values)
}

// TestRGPStatusAfterDelete covers RFC 3915 grace period reporting. Whether a
// deleted domain enters redemption is registry policy; the SDK must report
// whatever grace period the registry assigns.
func TestRGPStatusAfterDelete(t *testing.T) {
	client := session(t)

	registrant := contactID("rg")
	createContact(t, client, registrant)
	domain := domainName("rg")
	createDomain(t, client, domain, registrant)

	// A freshly created domain is in the add grace period.
	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if !containsString(info.Result.RGPStatuses, rgp.StatusAddPeriod) {
		t.Fatalf("expected addPeriod on a new domain, got %v", info.Result.RGPStatuses)
	}
	if len(info.Result.RGP.Statuses) == 0 {
		t.Fatal("structured RGP data is empty while RGPStatuses is populated")
	}

	if _, err := client.DomainDelete(types.DomainDeleteRequest{Domain: domain}); err != nil {
		t.Fatalf("domain delete failed: %v", err)
	}

	_, err = client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err == nil {
		t.Log("domain survived delete; it is in a grace period")
		return
	}
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "info after add-grace-period delete")
	t.Log("add grace period delete purged the domain immediately, so redemption is unreachable")
}

// TestRGPRestoreEncoding covers RFC 3915 restore request and report encoding.
// The registry must parse both and reject them on business grounds (the domain
// is not in redemption) rather than as a syntax error, which is what proves the
// SDK emits valid rgp-1.0 XML.
func TestRGPRestoreEncoding(t *testing.T) {
	client := session(t)

	registrant := contactID("rr")
	createContact(t, client, registrant)
	domain := domainName("rr")
	createDomain(t, client, domain, registrant)

	_, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		RGP:    &rgp.UpdateRequest{Restore: &rgp.Restore{Operation: rgp.OperationRequest}},
	})
	assertNotSyntaxError(t, err, "rgp restore request")

	_, err = client.DomainUpdate(types.DomainUpdateRequest{
		Domain: domain,
		RGP: &rgp.UpdateRequest{Restore: &rgp.Restore{
			Operation: rgp.OperationReport,
			Report: &rgp.RestoreReport{
				PreData:       "Pre-delete registration data",
				PostData:      "Post-restore registration data",
				DeleteTime:    time.Now().UTC().Add(-time.Hour).Format("2006-01-02T15:04:05Z"),
				RestoreTime:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
				RestoreReason: rgp.Text{Lang: "en", Value: "Registrant error"},
				Statements: []rgp.Text{
					{Lang: "en", Value: "This registrar has not restored the domain in order to assume the rights to use or sell the name."},
					{Lang: "en", Value: "The information in this report is true to the best of this registrar's knowledge."},
				},
			},
		}},
	})
	assertNotSyntaxError(t, err, "rgp restore report")
}

// TestLaunchEncoding covers the launch-1.0 extension. No launch phase is open
// in OT&E, so the test asserts the registry parses the extension rather than
// rejecting it as malformed XML.
func TestLaunchEncoding(t *testing.T) {
	client := session(t)

	registrant := contactID("lc")
	createContact(t, client, registrant)

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domainName("lc"), Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
		Launch: &launch.CreateRequest{
			Phase: launch.Phase{Value: launch.PhaseClaims},
			Notices: []launch.Notice{{
				ID:           "370d0b7c9223372036854775807",
				ValidatorID:  "tmch",
				NotAfter:     time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z"),
				AcceptedDate: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			}},
		},
	})
	if err == nil {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: domainName("lc")})
		t.Fatal("a claims-phase create unexpectedly succeeded outside a launch phase")
	}
	assertNotSyntaxError(t, err, "launch claims create")

	_, err = client.DomainInfo(types.DomainInfoRequest{
		Domain: domainName("lc"),
		Launch: &launch.InfoRequest{Phase: launch.Phase{Value: launch.PhaseClaims}},
	})
	assertNotSyntaxError(t, err, "launch info")
}

// assertNotSyntaxError fails if the registry rejected a command as malformed,
// which would mean the SDK produced invalid XML for the extension under test.
func assertNotSyntaxError(t *testing.T, err error, context string) {
	t.Helper()

	if err == nil {
		t.Logf("%s: accepted by the registry", context)
		return
	}
	eppErr := asEPPError(t, err, context)
	switch eppErr.Code {
	case constants.ResultSyntaxError, 2000:
		t.Fatalf("%s: registry rejected the SDK's XML as malformed (%d: %s) values=%v",
			context, eppErr.Code, eppErr.Message, eppErr.Values)
	}
	t.Logf("%s: parsed by the registry, refused on business grounds with %d: %s", context, eppErr.Code, eppErr.Message)
}

// TestSecDNSIncompleteRecords covers RFC 5910 section 4: a dsData record must
// carry keyTag, alg, digestType and digest, and a keyData record must carry
// flags, protocol, alg and pubKey. An incomplete record is schema-invalid, so
// the SDK must refuse it locally rather than spend a round trip being told so.
func TestSecDNSIncompleteRecords(t *testing.T) {
	client := session(t)

	registrant := contactID("sx")
	createContact(t, client, registrant)
	domain := domainName("sx")
	createDomain(t, client, domain, registrant)

	cases := []struct {
		name string
		ds   secdns.DSData
	}{
		{"empty digest", secdns.DSData{KeyTag: 12345, Algorithm: 8, DigestType: 2}},
		{"zero digest type", secdns.DSData{KeyTag: 12345, Algorithm: 8, Digest: testDS().Digest}},
		{"zero algorithm", secdns.DSData{KeyTag: 12345, DigestType: 2, Digest: testDS().Digest}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := client.DomainUpdate(types.DomainUpdateRequest{
				Domain: domain,
				SecDNS: &secdns.UpdateRequest{
					Add: &secdns.UpdateAdd{Data: secdns.Data{DSData: []secdns.DSData{testCase.ds}}},
				},
			})
			if err == nil {
				t.Fatalf("an incomplete DS record was accepted by the registry")
			}
			eppErr := asEPPError(t, err, "secDNS incomplete record")
			if !eppErr.IsValidationError() {
				t.Fatalf("the SDK sent a schema-invalid DS record and let the registry reject it with %d: %s",
					eppErr.Code, eppErr.Message)
			}
		})
	}
}
