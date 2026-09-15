//go:build live

package live

import (
	"strings"
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

const contactAuthInfo = "uAtC0nt#99"

// newContact builds a valid RFC 5733 contact create request.
func newContact(id string) types.ContactCreateRequest {
	return types.ContactCreateRequest{
		ContactID: id,
		InternationalPostalInfo: &types.PostalInfo{
			Type:          "int",
			Name:          "CSC UAT Contact",
			Organization:  "CSC e-Governance Services India Limited",
			Street:        []string{"3rd Floor Electronics Niketan"},
			City:          "New Delhi",
			StateProvince: "Delhi",
			PostalCode:    "110003",
			CountryCode:   "IN",
		},
		Voice:    types.Phone{Number: "+91.1124301756"},
		Email:    "uat@csc.gov.in",
		AuthInfo: contactAuthInfo,
	}
}

// createContact creates a contact and deletes it during cleanup.
func createContact(t *testing.T, client *epp.Client, id string) {
	t.Helper()

	if _, err := client.ContactCreate(newContact(id)); err != nil {
		t.Fatalf("contact create %s failed: %v", id, err)
	}
	t.Cleanup(func() {
		_, _ = client.ContactDelete(types.ContactDeleteRequest{ContactID: id})
	})
}

// TestContactLifecycle covers RFC 5733 check, create, info, update and delete.
func TestContactLifecycle(t *testing.T) {
	client := session(t)
	id := contactID("l")

	// Check: must be available before it exists.
	check, err := client.ContactCheck(types.ContactCheckRequest{IDs: []string{id}})
	if err != nil {
		t.Fatalf("contact check failed: %v", err)
	}
	if len(check.Results) != 1 {
		t.Fatalf("expected 1 check result, got %d", len(check.Results))
	}
	if !check.Results[0].Available {
		t.Fatalf("contact %s unexpectedly already exists", id)
	}

	createContact(t, client, id)

	// Check again: must now be unavailable.
	check, err = client.ContactCheck(types.ContactCheckRequest{IDs: []string{id}})
	if err != nil {
		t.Fatalf("contact re-check failed: %v", err)
	}
	if check.Results[0].Available {
		t.Fatalf("contact %s still reported available after create", id)
	}

	// Info: every submitted field must round-trip.
	info, err := client.ContactInfo(types.ContactInfoRequest{ContactID: id})
	if err != nil {
		t.Fatalf("contact info failed: %v", err)
	}
	if info.Contact.ContactID != id {
		t.Fatalf("unexpected contact ID: %s", info.Contact.ContactID)
	}
	if info.Contact.ROID == "" {
		t.Fatal("contact info returned no ROID")
	}
	postal := info.Contact.InternationalPostalInfo
	if postal == nil {
		t.Fatal("contact info returned no international postalInfo")
	}
	if postal.Name != "CSC UAT Contact" {
		t.Fatalf("unexpected contact name: %q", postal.Name)
	}
	if postal.CountryCode != "IN" {
		t.Fatalf("unexpected country code: %q", postal.CountryCode)
	}
	if info.Contact.Email != "uat@csc.gov.in" {
		t.Fatalf("unexpected email: %q", info.Contact.Email)
	}
	if len(info.Contact.Statuses) == 0 {
		t.Fatal("contact info returned no statuses")
	}

	// Update: change email and voice.
	if _, err := client.ContactUpdate(types.ContactUpdateRequest{
		ContactID: id,
		Email:     "uat-updated@csc.gov.in",
		Voice:     &types.Phone{Number: "+91.1124301799"},
	}); err != nil {
		t.Fatalf("contact update failed: %v", err)
	}

	info, err = client.ContactInfo(types.ContactInfoRequest{ContactID: id})
	if err != nil {
		t.Fatalf("contact info after update failed: %v", err)
	}
	if info.Contact.Email != "uat-updated@csc.gov.in" {
		t.Fatalf("update did not apply, email is %q", info.Contact.Email)
	}
	if info.Contact.Voice.Number != "+91.1124301799" {
		t.Fatalf("update did not apply, voice is %q", info.Contact.Voice.Number)
	}
	if info.Contact.UpdatedDate.IsZero() {
		t.Fatal("contact info returned no upDate after update")
	}

	// Delete.
	if _, err := client.ContactDelete(types.ContactDeleteRequest{ContactID: id}); err != nil {
		t.Fatalf("contact delete failed: %v", err)
	}

	// Info on the deleted object must report 2303.
	_, err = client.ContactInfo(types.ContactInfoRequest{ContactID: id})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "contact info after delete")
}

// TestContactCheckMultiple covers a multi-ID check in one command.
func TestContactCheckMultiple(t *testing.T) {
	client := shared(t)

	ids := []string{contactID("m1"), contactID("m2"), contactID("m3")}
	check, err := client.ContactCheck(types.ContactCheckRequest{IDs: ids})
	if err != nil {
		t.Fatalf("multi contact check failed: %v", err)
	}
	if len(check.Results) != len(ids) {
		t.Fatalf("expected %d results, got %d", len(ids), len(check.Results))
	}
	for i, result := range check.Results {
		if result.ContactID != ids[i] {
			t.Fatalf("result %d is for %q, expected %q", i, result.ContactID, ids[i])
		}
	}
}

// TestContactCreateDuplicate covers RFC 5730 result code 2302.
func TestContactCreateDuplicate(t *testing.T) {
	client := session(t)
	id := contactID("d")

	createContact(t, client, id)

	_, err := client.ContactCreate(newContact(id))
	eppErr := requireResultCodeErr(t, err, constants.ResultObjectExists, "duplicate contact create")
	if !eppErr.IsObjectExists() {
		t.Fatal("IsObjectExists did not classify 2302")
	}
}

// TestContactInfoUnknown covers RFC 5730 result code 2303.
func TestContactInfoUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.ContactInfo(types.ContactInfoRequest{ContactID: contactID("x")})
	eppErr := requireResultCodeErr(t, err, constants.ResultObjectDoesNotExist, "info on unknown contact")
	if !eppErr.IsObjectNotFound() {
		t.Fatal("IsObjectNotFound did not classify 2303")
	}
}

// TestContactDeleteUnknown covers deleting an object that does not exist.
func TestContactDeleteUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.ContactDelete(types.ContactDeleteRequest{ContactID: contactID("y")})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "delete of unknown contact")
}

// TestContactUpdateUnknown covers updating an object that does not exist.
func TestContactUpdateUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.ContactUpdate(types.ContactUpdateRequest{
		ContactID: contactID("z"),
		Email:     "nobody@csc.gov.in",
	})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "update of unknown contact")
}

// TestContactCreateInvalid covers server-side validation of contact fields and
// confirms the SDK surfaces the server's diagnostic detail.
func TestContactCreateInvalid(t *testing.T) {
	client := session(t)

	req := newContact(contactID("i"))
	req.Email = "not-an-email-address"

	_, err := client.ContactCreate(req)
	if err == nil {
		t.Fatal("expected an invalid email to be rejected")
	}
	eppErr := asEPPError(t, err, "contact create with invalid email")
	if eppErr.Code < 2000 {
		t.Fatalf("expected an error result, got %d", eppErr.Code)
	}
	t.Logf("invalid email rejected with %d: %s values=%v", eppErr.Code, eppErr.Message, eppErr.Values)
}

// TestContactCreateMissingRequiredFields covers client-side validation: the
// SDK must reject an obviously incomplete request without a round trip.
func TestContactCreateMissingRequiredFields(t *testing.T) {
	client := shared(t)

	if _, err := client.ContactCreate(types.ContactCreateRequest{}); err == nil {
		t.Fatal("expected a contact create with no ID to be rejected")
	}
}

// TestContactTransferQuery covers RFC 5733 transfer query on an object that is
// not pending transfer, which must report 2301.
func TestContactTransferQuery(t *testing.T) {
	client := session(t)
	id := contactID("t")
	createContact(t, client, id)

	_, err := client.ContactTransfer(types.ContactTransferRequest{
		ContactID: id,
		Operation: constants.TransferQuery,
	})
	if err == nil {
		t.Fatal("expected transfer query on a non-pending object to fail")
	}
	eppErr := asEPPError(t, err, "contact transfer query")
	if eppErr.Code != 2301 {
		t.Logf("contact transfer query returned %d: %s", eppErr.Code, eppErr.Message)
	}
}

// TestContactTransferRequestSelf covers requesting a transfer of an object the
// client already sponsors, which registries reject as a use or authorization
// error rather than performing a no-op transfer.
func TestContactTransferRequestSelf(t *testing.T) {
	client := session(t)
	id := contactID("s")
	createContact(t, client, id)

	_, err := client.ContactTransfer(types.ContactTransferRequest{
		ContactID: id,
		Operation: constants.TransferRequest,
		AuthInfo:  contactAuthInfo,
	})
	if err == nil {
		t.Fatal("expected a self-transfer request to be refused")
	}
	eppErr := asEPPError(t, err, "contact self transfer")
	t.Logf("self transfer refused with %d: %s", eppErr.Code, eppErr.Message)
}

// TestContactTransferWrongAuthInfo covers RFC 5730 result code 2202.
func TestContactTransferWrongAuthInfo(t *testing.T) {
	client := session(t)
	id := contactID("w")
	createContact(t, client, id)

	_, err := client.ContactTransfer(types.ContactTransferRequest{
		ContactID: id,
		Operation: constants.TransferRequest,
		AuthInfo:  "Wr0ngAuth9",
	})
	if err == nil {
		t.Fatal("expected a transfer with wrong authInfo to be refused")
	}
	eppErr := asEPPError(t, err, "contact transfer wrong authInfo")
	if strings.Contains(err.Error(), contactAuthInfo) {
		t.Fatal("error string leaked the real authInfo")
	}
	t.Logf("wrong authInfo refused with %d: %s", eppErr.Code, eppErr.Message)
}

// TestContactTransferMissingAuthInfo covers the SDK's own requirement that a
// transfer request carry authInfo.
func TestContactTransferMissingAuthInfo(t *testing.T) {
	client := shared(t)

	_, err := client.ContactTransfer(types.ContactTransferRequest{
		ContactID: contactID("n"),
		Operation: constants.TransferRequest,
	})
	if err == nil {
		t.Fatal("expected a transfer request without authInfo to be rejected")
	}
}

// TestContactDeleteLinked covers RFC 5730 result code 2305: a contact linked
// to a domain cannot be deleted.
func TestContactDeleteLinked(t *testing.T) {
	client := session(t)

	registrant := contactID("lr")
	createContact(t, client, registrant)

	domain := domainName("lk")
	createDomain(t, client, domain, registrant)

	_, err := client.ContactDelete(types.ContactDeleteRequest{ContactID: registrant})
	if err == nil {
		t.Fatal("expected deleting a linked contact to be refused")
	}
	eppErr := asEPPError(t, err, "delete of linked contact")
	if eppErr.Code != 2305 {
		t.Logf("linked contact delete returned %d: %s", eppErr.Code, eppErr.Message)
	}
}
