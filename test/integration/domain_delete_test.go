package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainDeleteXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainDeleteResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainDelete(types.DomainDeleteRequest{
		DomainName: "example.in",
	})
	if err != nil {
		t.Fatalf("domain delete failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"`)
	assertContains(t, requestXML, `<domain:delete>`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
	assertContains(t, requestXML, `<clTRID>DELETE-`)

	assertDomainDeleteResponse(t, resp, "example.in")
}

func TestDomainDeleteUnicodeConversion(t *testing.T) {
	unicodeDomain := "भारत.भारत"
	asciiDomain, err := idn.ToASCII(unicodeDomain)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, domainDeleteResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainDelete(types.DomainDeleteRequest{
		DomainName: unicodeDomain,
	})
	if err != nil {
		t.Fatalf("domain delete failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<domain:name>"+asciiDomain+"</domain:name>")
	assertDomainDeleteResponse(t, resp, unicodeDomain)
}

func TestDomainDeleteEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainDelete(types.DomainDeleteRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertDomainDeleteResponse(
	t *testing.T,
	resp *types.DomainDeleteResponse,
	domain string,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "DOMAIN-DELETE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-DOMAIN-DELETE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.Result.DomainName != domain ||
		resp.Result.Domain != domain {

		t.Fatalf("unexpected result domain: %#v", resp.Result)
	}
}

func domainDeleteResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>DOMAIN-DELETE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-DELETE</svTRID>
        </trID>
    </response>
</epp>`
}
