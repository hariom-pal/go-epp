//go:build live

package live

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// TestHostRename covers the RFC 5732 chg element, which renames a host object.
func TestHostRename(t *testing.T) {
	client := session(t)

	registrant := contactID("hn")
	createContact(t, client, registrant)
	domain := domainName("hn")
	createDomain(t, client, domain, registrant)

	original := hostName("ns1", domain)
	renamed := hostName("ns9", domain)
	createHost(t, client, original, "103.51.75.51")

	if _, err := client.HostUpdate(types.HostUpdateRequest{
		HostName:    original,
		NewHostName: renamed,
	}); err != nil {
		t.Fatalf("host rename failed: %v", err)
	}

	if _, err := client.HostInfo(types.HostInfoRequest{HostName: original}); err == nil {
		t.Fatal("the original host name still resolves after a rename")
	}

	info, err := client.HostInfo(types.HostInfoRequest{HostName: renamed})
	if err != nil {
		t.Fatalf("info on the renamed host failed: %v", err)
	}
	if info.Host.HostName != renamed {
		t.Fatalf("expected %q, got %q", renamed, info.Host.HostName)
	}
	t.Cleanup(func() {
		_, _ = client.HostDelete(types.HostDeleteRequest{HostName: renamed})
	})
}

// TestDomainChangeRegistrant covers the RFC 5731 chg element for a domain.
func TestDomainChangeRegistrant(t *testing.T) {
	client := session(t)

	first := contactID("g1")
	second := contactID("g2")
	createContact(t, client, first)
	createContact(t, client, second)

	domain := domainName("cg")
	createDomain(t, client, domain, first)

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:        domain,
		NewRegistrant: second,
	}); err != nil {
		t.Fatalf("registrant change failed: %v", err)
	}

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if info.Result.Registrant != second {
		t.Fatalf("registrant did not change: %q", info.Result.Registrant)
	}
}

// TestDomainChangeAuthInfo covers rotating a domain's authInfo, which is how a
// registrar revokes a transfer code.
func TestDomainChangeAuthInfo(t *testing.T) {
	client := session(t)

	registrant := contactID("ga")
	createContact(t, client, registrant)
	domain := domainName("ga")
	createDomain(t, client, domain, registrant)

	const rotated = "uAtR0t#2026x"
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:   domain,
		AuthInfo: rotated,
	}); err != nil {
		t.Fatalf("authInfo rotation failed: %v", err)
	}

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if info.Result.AuthInfo != rotated {
		t.Fatalf("authInfo did not rotate, got %q", info.Result.AuthInfo)
	}
	if info.Result.AuthInfo == domainAuthInfo {
		t.Fatal("the old authInfo is still in force")
	}
}

// TestDomainContactSwap covers adding and removing typed contacts.
func TestDomainContactSwap(t *testing.T) {
	client := session(t)

	first := contactID("k1")
	second := contactID("k2")
	createContact(t, client, first)
	createContact(t, client, second)

	domain := domainName("ks")
	createDomain(t, client, domain, first)

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:             domain,
		AddTechContacts:    []string{second},
		RemoveTechContacts: []string{first},
	}); err != nil {
		t.Fatalf("tech contact swap failed: %v", err)
	}

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	for _, contact := range info.Result.Contacts {
		if contact.Type == "tech" && contact.ID != second {
			t.Fatalf("tech contact is %q, expected %q", contact.ID, second)
		}
	}
}

// TestDomainInfoHostsVariants covers the RFC 5731 hosts attribute, which
// selects which delegation data the registry returns.
func TestDomainInfoHostsVariants(t *testing.T) {
	client := session(t)

	registrant := contactID("hv")
	createContact(t, client, registrant)
	domain := domainName("hv")
	createDomain(t, client, domain, registrant)

	host := hostName("ns1", domain)
	createHost(t, client, host, "103.51.75.52")
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:         domain,
		AddNameServers: []string{host},
	}); err != nil {
		t.Fatalf("attach nameserver failed: %v", err)
	}

	for _, hosts := range []string{"all", "del", "sub", "none"} {
		info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain, Hosts: hosts})
		if err != nil {
			t.Fatalf("domain info hosts=%q failed: %v", hosts, err)
		}
		t.Logf("hosts=%-5s nameservers=%v", hosts, info.Result.NameServers)
	}

	// Detach so cleanup can remove the host.
	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:            domain,
		RemoveNameServers: []string{host},
	}); err != nil {
		t.Fatalf("detach nameserver failed: %v", err)
	}
}

// TestContactLocalizedPostalInfo covers the RFC 5733 loc postalInfo form
// alongside the int form.
func TestContactLocalizedPostalInfo(t *testing.T) {
	client := session(t)

	id := contactID("lp")
	request := newContact(id)
	request.LocalizedPostalInfo = &types.PostalInfo{
		Type:          "loc",
		Name:          "CSC UAT Contact",
		Organization:  "CSC e-Governance Services India Limited",
		Street:        []string{"3rd Floor Electronics Niketan"},
		City:          "New Delhi",
		StateProvince: "Delhi",
		PostalCode:    "110003",
		CountryCode:   "IN",
	}

	if _, err := client.ContactCreate(request); err != nil {
		t.Fatalf("contact create with localized postalInfo failed: %v", err)
	}
	t.Cleanup(func() {
		_, _ = client.ContactDelete(types.ContactDeleteRequest{ContactID: id})
	})

	info, err := client.ContactInfo(types.ContactInfoRequest{ContactID: id})
	if err != nil {
		t.Fatalf("contact info failed: %v", err)
	}
	if info.Contact.InternationalPostalInfo == nil && info.Contact.LocalizedPostalInfo == nil {
		t.Fatal("neither postalInfo form came back")
	}
}

// TestContactDisclosure covers the RFC 5733 disclose element, which carries a
// contact's data-disclosure preferences.
func TestContactDisclosure(t *testing.T) {
	client := session(t)

	id := contactID("dc")
	request := newContact(id)
	request.Disclosure = &types.ContactDisclosure{
		Flag:  false,
		Voice: true,
		Email: true,
	}

	if _, err := client.ContactCreate(request); err != nil {
		eppErr := asEPPError(t, err, "contact create with disclose")
		t.Logf("registry refused the disclose element with %d: %s", eppErr.Code, eppErr.Message)
		return
	}
	t.Cleanup(func() {
		_, _ = client.ContactDelete(types.ContactDeleteRequest{ContactID: id})
	})

	if _, err := client.ContactInfo(types.ContactInfoRequest{ContactID: id}); err != nil {
		t.Fatalf("contact info failed: %v", err)
	}
}

// TestContactStatusUpdate covers adding and removing a client status on a
// contact.
func TestContactStatusUpdate(t *testing.T) {
	client := session(t)

	id := contactID("cs")
	createContact(t, client, id)

	if _, err := client.ContactUpdate(types.ContactUpdateRequest{
		ContactID:   id,
		AddStatuses: []string{"clientUpdateProhibited"},
	}); err != nil {
		t.Fatalf("contact add status failed: %v", err)
	}

	info, err := client.ContactInfo(types.ContactInfoRequest{ContactID: id})
	if err != nil {
		t.Fatalf("contact info failed: %v", err)
	}
	if !containsString(info.Contact.Statuses, "clientUpdateProhibited") {
		t.Fatalf("status not applied: %v", info.Contact.Statuses)
	}

	if _, err := client.ContactUpdate(types.ContactUpdateRequest{
		ContactID:      id,
		RemoveStatuses: []string{"clientUpdateProhibited"},
	}); err != nil {
		t.Fatalf("contact remove status failed: %v", err)
	}
}

// TestDomainCreateWithHostAttr covers the RFC 5731 hostAttr delegation form,
// where nameservers are supplied inline instead of as host objects.
func TestDomainCreateWithHostAttr(t *testing.T) {
	client := session(t)

	registrant := contactID("ha")
	createContact(t, client, registrant)
	domain := domainName("ha")

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain: domain, Period: 1, Unit: "y",
		Registrant: registrant, AuthInfo: domainAuthInfo,
		AdminContacts: []string{registrant}, TechContacts: []string{registrant},
		BillingContacts: []string{registrant},
		NameServerInfo: []types.DomainNameServer{{
			HostName:  hostName("ns1", domain),
			Addresses: []types.DomainHostAddress{{IP: "103.51.75.53", Version: "v4"}},
		}},
	})
	if err != nil {
		eppErr := asEPPError(t, err, "domain create with hostAttr")
		t.Logf("registry does not accept the hostAttr form: %d %s values=%v",
			eppErr.Code, eppErr.Message, eppErr.Values)
		return
	}
	t.Cleanup(func() {
		_, _ = client.DomainDelete(types.DomainDeleteRequest{Domain: domain})
	})

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain, Hosts: "all"})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if len(info.Result.NameServers) == 0 {
		t.Fatal("hostAttr delegation did not register any nameserver")
	}
}

// TestContactTransferOperations covers every contact transfer operation the
// SDK supports, confirming the registry parses each rather than rejecting the
// SDK's XML as malformed.
func TestContactTransferOperations(t *testing.T) {
	client := session(t)

	id := contactID("ct")
	createContact(t, client, id)

	for _, op := range []string{
		constants.TransferQuery,
		constants.TransferRequest,
		constants.TransferApprove,
		constants.TransferReject,
		constants.TransferCancel,
	} {
		request := types.ContactTransferRequest{ContactID: id, Operation: op}
		if op == constants.TransferRequest {
			request.AuthInfo = contactAuthInfo
		}
		_, err := client.ContactTransfer(request)
		if err == nil {
			t.Logf("contact transfer %s was accepted", op)
			continue
		}
		eppErr := asEPPError(t, err, "contact transfer "+op)
		if eppErr.Code == constants.ResultSyntaxError || eppErr.Code == 2000 {
			t.Fatalf("registry rejected the SDK's contact transfer %s XML as malformed: %d %s",
				op, eppErr.Code, eppErr.Message)
		}
		t.Logf("contact transfer %-8s refused with %d: %s", op, eppErr.Code, eppErr.Message)
	}
}

// TestDomainDeleteRestoreWindow records how a delete behaves for a domain that
// is past the add grace period, which is the only way RFC 3915 redemption can
// be reached. OT&E purges add-grace deletes immediately, so this documents the
// boundary rather than asserting a redemption that the environment cannot
// produce.
func TestDomainDeleteRestoreWindow(t *testing.T) {
	client := session(t)

	registrant := contactID("dw")
	createContact(t, client, registrant)
	domain := domainName("dw")
	createDomain(t, client, domain, registrant)

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	if !containsString(info.Result.RGPStatuses, "addPeriod") {
		t.Fatalf("expected addPeriod, got %v", info.Result.RGPStatuses)
	}

	if _, err := client.DomainDelete(types.DomainDeleteRequest{Domain: domain}); err != nil {
		t.Fatalf("domain delete failed: %v", err)
	}

	_, err = client.DomainInfo(types.DomainInfoRequest{Domain: domain})
	if err == nil {
		t.Log("domain persisted after delete, so a grace period applies and restore is reachable")
		return
	}
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "info after add-grace delete")
	t.Log("add grace period delete purges immediately; redemption needs a domain older than that window")
}
