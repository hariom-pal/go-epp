package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainUpdateXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainUpdate(types.DomainUpdateRequest{
		DomainName:        "example.in",
		AddNameServers:    []string{"ns1.example.in", "ns2.example.in"},
		RemoveNameServers: []string{"old-ns.example.in"},
		AddContacts: []types.DomainContact{
			{Type: "admin", ID: "CNT001"},
			{Type: "custom", ID: "CNT002"},
		},
		RemoveContacts: []types.DomainContact{
			{Type: "tech", ID: "CNT003"},
		},
		AddStatuses:    []string{"clientTransferProhibited"},
		RemoveStatuses: []string{"clientHold"},
		Registrant:     "REG123",
		AuthInfo:       "newSecret",
	})
	if err != nil {
		t.Fatalf("domain update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"`)
	assertContains(t, requestXML, `<domain:update>`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
	assertContains(t, requestXML, `<domain:add>`)
	assertContains(t, requestXML, `<domain:hostObj>ns1.example.in</domain:hostObj>`)
	assertContains(t, requestXML, `<domain:hostObj>ns2.example.in</domain:hostObj>`)
	assertContains(t, requestXML, `<domain:contact type="admin">CNT001</domain:contact>`)
	assertContains(t, requestXML, `<domain:contact type="custom">CNT002</domain:contact>`)
	assertContains(t, requestXML, `<domain:status s="clientTransferProhibited"></domain:status>`)
	assertContains(t, requestXML, `<domain:rem>`)
	assertContains(t, requestXML, `<domain:hostObj>old-ns.example.in</domain:hostObj>`)
	assertContains(t, requestXML, `<domain:contact type="tech">CNT003</domain:contact>`)
	assertContains(t, requestXML, `<domain:status s="clientHold"></domain:status>`)
	assertContains(t, requestXML, `<domain:chg>`)
	assertContains(t, requestXML, `<domain:registrant>REG123</domain:registrant>`)
	assertContains(t, requestXML, `<domain:authInfo>`)
	assertContains(t, requestXML, `<domain:pw>newSecret</domain:pw>`)
	assertContains(t, requestXML, `<clTRID>UPDATE-`)

	assertDomainUpdateResponse(t, resp, "example.in")
}

func TestDomainUpdateUnicodeConversion(t *testing.T) {
	unicodeDomain := "भारत.भारत"
	asciiDomain, err := idn.ToASCII(unicodeDomain)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, domainUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainUpdate(types.DomainUpdateRequest{
		DomainName: unicodeDomain,
		AuthInfo:   "newSecret",
	})
	if err != nil {
		t.Fatalf("domain update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<domain:name>"+asciiDomain+"</domain:name>")
	assertDomainUpdateResponse(t, resp, unicodeDomain)
}

func TestDomainUpdateEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainUpdate(types.DomainUpdateRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainUpdateNoOperationValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainUpdate(types.DomainUpdateRequest{
		DomainName: "example.in",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertDomainUpdateResponse(
	t *testing.T,
	resp *types.DomainUpdateResponse,
	domain string,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "DOMAIN-UPDATE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-DOMAIN-UPDATE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.Result.Domain != domain {
		t.Fatalf("unexpected result domain: %s", resp.Result.Domain)
	}
}

func domainUpdateResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>DOMAIN-UPDATE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-UPDATE</svTRID>
        </trID>
    </response>
</epp>`
}
