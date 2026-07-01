package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainCheckXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainCheckResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{
			"available.in",
			"taken.in",
		},
	})
	if err != nil {
		t.Fatalf("domain check failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"`)
	assertContains(t, requestXML, `<domain:check>`)
	assertContains(t, requestXML, `<domain:name>available.in</domain:name>`)
	assertContains(t, requestXML, `<domain:name>taken.in</domain:name>`)
	assertContains(t, requestXML, `<clTRID>CHECK-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "DOMAIN-CHECK-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-DOMAIN-CHECK" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(resp.Results))
	}

	if resp.Results[0].Domain != "available.in" ||
		resp.Results[0].ASCII != "available.in" ||
		!resp.Results[0].Available {

		t.Fatalf("unexpected available result: %+v", resp.Results[0])
	}

	if resp.Results[1].Domain != "taken.in" ||
		resp.Results[1].ASCII != "taken.in" ||
		resp.Results[1].Available ||
		resp.Results[1].Reason != "Already registered" {

		t.Fatalf("unexpected unavailable result: %+v", resp.Results[1])
	}
}

func TestDomainCheckUnicodeConversion(t *testing.T) {
	unicodeDomain := "भारत.भारत"
	asciiDomain, err := idn.ToASCII(unicodeDomain)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, domainCheckUnicodeResponse(asciiDomain))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{unicodeDomain},
	})
	if err != nil {
		t.Fatalf("domain check failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<domain:name>"+asciiDomain+"</domain:name>")

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}

	if resp.Results[0].Domain != unicodeDomain ||
		resp.Results[0].ASCII != asciiDomain {

		t.Fatalf("unexpected unicode result: %+v", resp.Results[0])
	}
}

func TestDomainCheckEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainCheck(types.DomainCheckRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func domainCheckResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:cd>
                    <domain:name avail="1">available.in</domain:name>
                </domain:cd>
                <domain:cd>
                    <domain:name avail="0">taken.in</domain:name>
                    <domain:reason>Already registered</domain:reason>
                </domain:cd>
            </domain:chkData>
        </resData>
        <trID>
            <clTRID>DOMAIN-CHECK-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-CHECK</svTRID>
        </trID>
    </response>
</epp>`
}

func domainCheckUnicodeResponse(domain string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:cd>
                    <domain:name avail="1">` + domain + `</domain:name>
                </domain:cd>
            </domain:chkData>
        </resData>
        <trID>
            <clTRID>DOMAIN-CHECK-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-CHECK</svTRID>
        </trID>
    </response>
</epp>`
}
