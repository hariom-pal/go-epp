package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainInfoXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainInfoResponse("example.in"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainInfo(types.DomainInfoRequest{
		Domain: "example.in",
		Hosts:  constants.HostsAll,
	})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"`)
	assertContains(t, requestXML, `<domain:info>`)
	assertContains(t, requestXML, `<domain:name hosts="all">example.in</domain:name>`)
	assertContains(t, requestXML, `<clTRID>INFO-`)

	assertDomainInfoResponse(t, resp, "example.in", "example.in")
}

func TestDomainInfoUnicodeConversion(t *testing.T) {
	unicodeDomain := "भारत.भारत"
	asciiDomain, err := idn.ToASCII(unicodeDomain)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, domainInfoResponse(asciiDomain))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainInfo(types.DomainInfoRequest{
		Domain: unicodeDomain,
	})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<domain:name hosts="all">`+asciiDomain+`</domain:name>`)

	if resp.Result.Domain != unicodeDomain {
		t.Fatalf("unexpected unicode domain: %s", resp.Result.Domain)
	}

	if resp.Result.ASCII != asciiDomain {
		t.Fatalf("unexpected ASCII domain: %s", resp.Result.ASCII)
	}
}

func TestDomainInfoEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainInfo(types.DomainInfoRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainInfoInvalidHostsValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainInfo(types.DomainInfoRequest{
		Domain: "example.in",
		Hosts:  "invalid",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertDomainInfoResponse(
	t *testing.T,
	resp *types.DomainInfoResponse,
	domain string,
	ascii string,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "DOMAIN-INFO-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-DOMAIN-INFO" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	result := resp.Result
	if result.Domain != domain || result.ASCII != ascii {
		t.Fatalf("unexpected domain identity: %+v", result)
	}

	if result.ROID != "D123456-IN" ||
		result.Registrant != "REG001" ||
		result.Registrar != "Registrar-OTE" {

		t.Fatalf("unexpected domain ownership fields: %+v", result)
	}

	if len(result.Statuses) != 2 ||
		result.Statuses[0] != constants.StatusOK ||
		result.Statuses[1] != constants.StatusClientTransferProhibited {

		t.Fatalf("unexpected statuses: %+v", result.Statuses)
	}

	if len(result.Contacts) != 3 ||
		result.Contacts[0] != (types.DomainContact{Type: "admin", ID: "CNT001"}) ||
		result.Contacts[1] != (types.DomainContact{Type: "tech", ID: "CNT002"}) ||
		result.Contacts[2] != (types.DomainContact{Type: "billing", ID: "CNT003"}) {

		t.Fatalf("unexpected contacts: %+v", result.Contacts)
	}

	if len(result.NameServers) != 2 ||
		result.NameServers[0] != "ns1.example.in" ||
		result.NameServers[1] != "ns2.example.in" {

		t.Fatalf("unexpected name servers: %+v", result.NameServers)
	}

	if len(result.NameServerInfo) != 2 ||
		result.NameServerInfo[1].HostName != "ns2.example.in" ||
		len(result.NameServerInfo[1].Addresses) != 2 {

		t.Fatalf("unexpected name server info: %+v", result.NameServerInfo)
	}

	if result.AuthInfo != "secret" ||
		result.AuthInfoROID != "D123456-IN" {

		t.Fatalf("unexpected auth info fields: %+v", result)
	}

	assertDomainInfoDate(t, result.CreatedDate, time.Date(2026, 6, 30, 9, 30, 0, 0, time.UTC))
	assertDomainInfoDate(t, result.UpdatedDate, time.Date(2026, 7, 1, 10, 30, 0, 0, time.UTC))
	assertDomainInfoDate(t, result.ExpiryDate, time.Date(2027, 6, 30, 9, 30, 0, 0, time.UTC))

	if len(result.RGPStatuses) != 1 ||
		result.RGPStatuses[0] != "addPeriod" {

		t.Fatalf("unexpected RGP statuses: %+v", result.RGPStatuses)
	}
}

func assertDomainInfoDate(
	t *testing.T,
	actual *time.Time,
	expected time.Time,
) {
	t.Helper()

	if actual == nil {
		t.Fatalf("expected date %s, got nil", expected.Format(time.RFC3339))
	}

	if !actual.Equal(expected) {
		t.Fatalf("expected date %s, got %s", expected.Format(time.RFC3339), actual.Format(time.RFC3339))
	}
}

func domainInfoResponse(domain string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:infData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>` + domain + `</domain:name>
                <domain:roid>D123456-IN</domain:roid>
                <domain:registrant>REG001</domain:registrant>
                <domain:status s="ok">Active</domain:status>
                <domain:status s="clientTransferProhibited">Transfer locked</domain:status>
                <domain:contact type="admin">CNT001</domain:contact>
                <domain:contact type="tech">CNT002</domain:contact>
                <domain:contact type="billing">CNT003</domain:contact>
                <domain:ns>
                    <domain:hostObj>ns1.example.in</domain:hostObj>
                    <domain:hostAttr>
                        <domain:hostName>ns2.example.in</domain:hostName>
                        <domain:hostAddr ip="v4">192.0.2.1</domain:hostAddr>
                        <domain:hostAddr ip="v6">2001:db8::1</domain:hostAddr>
                    </domain:hostAttr>
                </domain:ns>
                <domain:clID>Registrar-OTE</domain:clID>
                <domain:crID>Registrar-OTE</domain:crID>
                <domain:crDate>2026-06-30T09:30:00Z</domain:crDate>
                <domain:upID>Registrar-OTE</domain:upID>
                <domain:upDate>2026-07-01T10:30:00Z</domain:upDate>
                <domain:exDate>2027-06-30T09:30:00Z</domain:exDate>
                <domain:authInfo>
                    <domain:pw roid="D123456-IN">secret</domain:pw>
                </domain:authInfo>
            </domain:infData>
        </resData>
        <extension>
            <rgp:infData xmlns:rgp="urn:ietf:params:xml:ns:rgp-1.0">
                <rgp:rgpStatus s="addPeriod"/>
            </rgp:infData>
        </extension>
        <trID>
            <clTRID>DOMAIN-INFO-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-INFO</svTRID>
        </trID>
    </response>
</epp>`
}
