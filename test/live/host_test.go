//go:build live

package live

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

// createHost registers a host object and deletes it during cleanup.
func createHost(t *testing.T, client *epp.Client, host, ipv4 string) {
	t.Helper()

	if _, err := client.HostCreate(types.HostCreateRequest{
		HostName:  host,
		Addresses: []types.HostAddress{{Address: ipv4, IPVersion: "v4"}},
	}); err != nil {
		t.Fatalf("host create %s failed: %v", host, err)
	}
	t.Cleanup(func() {
		_, _ = client.HostDelete(types.HostDeleteRequest{HostName: host})
	})
}

// TestHostLifecycle covers RFC 5732 check, create, info, update and delete.
func TestHostLifecycle(t *testing.T) {
	client := session(t)

	registrant := contactID("hr")
	createContact(t, client, registrant)
	domain := domainName("hl")
	createDomain(t, client, domain, registrant)

	host := hostName("ns1", domain)

	check, err := client.HostCheck(types.HostCheckRequest{Hosts: []string{host}})
	if err != nil {
		t.Fatalf("host check failed: %v", err)
	}
	if !check.Results[0].Available {
		t.Fatalf("host %s unexpectedly already exists", host)
	}

	createHost(t, client, host, "103.51.75.21")

	check, err = client.HostCheck(types.HostCheckRequest{Hosts: []string{host}})
	if err != nil {
		t.Fatalf("host re-check failed: %v", err)
	}
	if check.Results[0].Available {
		t.Fatalf("host %s still available after create", host)
	}

	info, err := client.HostInfo(types.HostInfoRequest{HostName: host})
	if err != nil {
		t.Fatalf("host info failed: %v", err)
	}
	if info.Host.HostName != host {
		t.Fatalf("unexpected host name %q", info.Host.HostName)
	}
	if info.Host.ROID == "" {
		t.Fatal("host info returned no ROID")
	}
	if len(info.Host.Addresses) != 1 || info.Host.Addresses[0].Address != "103.51.75.21" {
		t.Fatalf("unexpected addresses: %v", info.Host.Addresses)
	}

	// Add an IPv6 address (RFC 5732 allows both families on one host).
	if _, err := client.HostUpdate(types.HostUpdateRequest{
		HostName:     host,
		AddAddresses: []types.HostAddress{{Address: "2401:4900:1c00::21", IPVersion: "v6"}},
	}); err != nil {
		t.Fatalf("host update add address failed: %v", err)
	}

	info, err = client.HostInfo(types.HostInfoRequest{HostName: host})
	if err != nil {
		t.Fatalf("host info after update failed: %v", err)
	}
	if len(info.Host.Addresses) != 2 {
		t.Fatalf("expected 2 addresses after update, got %v", info.Host.Addresses)
	}

	// Remove the IPv6 address again.
	if _, err := client.HostUpdate(types.HostUpdateRequest{
		HostName:        host,
		RemoveAddresses: []types.HostAddress{{Address: "2401:4900:1c00::21", IPVersion: "v6"}},
	}); err != nil {
		t.Fatalf("host update remove address failed: %v", err)
	}

	// Status update.
	if _, err := client.HostUpdate(types.HostUpdateRequest{
		HostName:    host,
		AddStatuses: []string{"clientUpdateProhibited"},
	}); err != nil {
		t.Fatalf("host update add status failed: %v", err)
	}
	if _, err := client.HostUpdate(types.HostUpdateRequest{
		HostName:       host,
		RemoveStatuses: []string{"clientUpdateProhibited"},
	}); err != nil {
		t.Fatalf("host update remove status failed: %v", err)
	}

	if _, err := client.HostDelete(types.HostDeleteRequest{HostName: host}); err != nil {
		t.Fatalf("host delete failed: %v", err)
	}

	_, err = client.HostInfo(types.HostInfoRequest{HostName: host})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "host info after delete")
}

// TestHostCreateDuplicate covers RFC 5730 result code 2302.
func TestHostCreateDuplicate(t *testing.T) {
	client := session(t)

	registrant := contactID("hd")
	createContact(t, client, registrant)
	domain := domainName("hd")
	createDomain(t, client, domain, registrant)

	host := hostName("ns1", domain)
	createHost(t, client, host, "103.51.75.22")

	_, err := client.HostCreate(types.HostCreateRequest{
		HostName:  host,
		Addresses: []types.HostAddress{{Address: "103.51.75.22", IPVersion: "v4"}},
	})
	requireResultCode(t, err, constants.ResultObjectExists, "duplicate host create")
}

// TestHostInfoUnknown covers RFC 5730 result code 2303.
func TestHostInfoUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.HostInfo(types.HostInfoRequest{HostName: "ns1.no-such-domain-" + tag() + ".in"})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "info on unknown host")
}

// TestHostDeleteUnknown covers deleting a host that does not exist.
func TestHostDeleteUnknown(t *testing.T) {
	client := shared(t)

	_, err := client.HostDelete(types.HostDeleteRequest{HostName: "ns9.no-such-domain-" + tag() + ".in"})
	requireResultCode(t, err, constants.ResultObjectDoesNotExist, "delete of unknown host")
}

// TestHostCreateReservedAddress covers registry IP policy: RFC 5737
// documentation ranges must be refused, and the SDK must surface the reason.
func TestHostCreateReservedAddress(t *testing.T) {
	client := session(t)

	registrant := contactID("hp")
	createContact(t, client, registrant)
	domain := domainName("hp")
	createDomain(t, client, domain, registrant)

	_, err := client.HostCreate(types.HostCreateRequest{
		HostName:  hostName("ns1", domain),
		Addresses: []types.HostAddress{{Address: "203.0.113.11", IPVersion: "v4"}},
	})
	if err == nil {
		t.Fatal("expected a reserved documentation address to be refused")
	}
	eppErr := asEPPError(t, err, "host create with reserved address")
	if len(eppErr.Values) == 0 {
		t.Fatal("registry policy rejection carried no diagnostic detail")
	}
	t.Logf("reserved address refused with %d: %s values=%v", eppErr.Code, eppErr.Message, eppErr.Values)
}

// TestHostCreateInvalidAddress covers SDK-side IP validation.
func TestHostCreateInvalidAddress(t *testing.T) {
	client := shared(t)

	for _, address := range []types.HostAddress{
		{Address: "not-an-ip", IPVersion: "v4"},
		{Address: "999.999.999.999", IPVersion: "v4"},
		{Address: "2401:4900:1c00::1", IPVersion: "v4"},
		{Address: "103.51.75.1", IPVersion: "v6"},
	} {
		_, err := client.HostCreate(types.HostCreateRequest{
			HostName:  "ns1.example-" + tag() + ".in",
			Addresses: []types.HostAddress{address},
		})
		if err == nil {
			t.Fatalf("expected address %q (%s) to be rejected", address.Address, address.IPVersion)
		}
	}
}

// TestHostCreateMissingName covers SDK-side name validation.
func TestHostCreateMissingName(t *testing.T) {
	client := shared(t)

	if _, err := client.HostCreate(types.HostCreateRequest{}); err == nil {
		t.Fatal("expected a host create with no name to be rejected")
	}
}

// TestHostDeleteLinked records how the registry treats a host that is in use
// as a nameserver. RFC 5732 section 3.2.2 makes rejecting such a delete a
// SHOULD rather than a MUST, and NIXI chooses to cascade: the delete succeeds
// and the host is silently detached from the domain. A registrar must not rely
// on the registry to protect delegation, so this test pins the behaviour and
// verifies the SDK reports the `linked` status that signals it.
func TestHostDeleteLinked(t *testing.T) {
	client := session(t)

	registrant := contactID("hx")
	createContact(t, client, registrant)
	domain := domainName("hx")
	createDomain(t, client, domain, registrant)

	host := hostName("ns1", domain)
	createHost(t, client, host, "103.51.75.23")

	if _, err := client.DomainUpdate(types.DomainUpdateRequest{
		Domain:         domain,
		AddNameServers: []string{host},
	}); err != nil {
		t.Fatalf("attach nameserver failed: %v", err)
	}

	info, err := client.HostInfo(types.HostInfoRequest{HostName: host})
	if err != nil {
		t.Fatalf("host info failed: %v", err)
	}
	if !containsString(info.Host.Statuses, "linked") {
		t.Fatalf("expected the in-use host to report linked, got %v", info.Host.Statuses)
	}

	_, err = client.HostDelete(types.HostDeleteRequest{HostName: host})
	if err != nil {
		// A registry that enforces the SHOULD is equally correct.
		eppErr := asEPPError(t, err, "delete of linked host")
		t.Logf("registry refused the linked host delete with %d: %s", eppErr.Code, eppErr.Message)
		if _, err := client.DomainUpdate(types.DomainUpdateRequest{
			Domain:            domain,
			RemoveNameServers: []string{host},
		}); err != nil {
			t.Fatalf("detach nameserver failed: %v", err)
		}
		return
	}

	t.Log("registry cascaded the delete; confirming the domain was detached")
	domainInfo, err := client.DomainInfo(types.DomainInfoRequest{Domain: domain, Hosts: "all"})
	if err != nil {
		t.Fatalf("domain info after host delete failed: %v", err)
	}
	if containsString(domainInfo.Result.NameServers, host) {
		t.Fatalf("host was deleted but the domain still delegates to it: %v", domainInfo.Result.NameServers)
	}
}

// TestHostCheckMultiple covers a multi-name check in one command.
func TestHostCheckMultiple(t *testing.T) {
	client := shared(t)

	hosts := []string{
		"ns1.multi-" + tag() + ".in",
		"ns2.multi-" + tag() + ".in",
	}
	check, err := client.HostCheck(types.HostCheckRequest{Hosts: hosts})
	if err != nil {
		t.Fatalf("multi host check failed: %v", err)
	}
	if len(check.Results) != len(hosts) {
		t.Fatalf("expected %d results, got %d", len(hosts), len(check.Results))
	}
}
