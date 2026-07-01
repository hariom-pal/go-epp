package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainRenewXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainRenewResponse("example.in"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName:        "example.in",
		CurrentExpiryDate: time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
		PeriodInfo: types.Period{
			Value: 1,
			Unit:  "m",
		},
	})
	if err != nil {
		t.Fatalf("domain renew failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"`)
	assertContains(t, requestXML, `<domain:renew>`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
	assertContains(t, requestXML, `<domain:curExpDate>2027-06-30</domain:curExpDate>`)
	assertContains(t, requestXML, `<domain:period unit="m">1</domain:period>`)
	assertContains(t, requestXML, `<clTRID>RENEW-`)

	assertDomainRenewResponse(t, resp, "example.in")
}

func TestDomainRenewUnicodeConversion(t *testing.T) {
	unicodeDomain := "भारत.भारत"
	asciiDomain, err := idn.ToASCII(unicodeDomain)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, domainRenewResponse(asciiDomain))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName:        unicodeDomain,
		CurrentExpiryDate: time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
		Period:            1,
		Unit:              "y",
	})
	if err != nil {
		t.Fatalf("domain renew failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<domain:name>"+asciiDomain+"</domain:name>")

	if resp.Result.DomainName != unicodeDomain {
		t.Fatalf("unexpected unicode domain: %s", resp.Result.DomainName)
	}
}

func TestDomainRenewEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainRenew(types.DomainRenewRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainRenewMissingExpiryDate(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName: "example.in",
		Period:     1,
		Unit:       "y",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainRenewInvalidPeriod(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName:        "example.in",
		CurrentExpiryDate: time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
		Period:            0,
		Unit:              "y",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainRenewRejectsPeriodAboveMaximum(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName:        "example.in",
		CurrentExpiryDate: time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
		Period:            100,
		Unit:              "y",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertDomainRenewResponse(
	t *testing.T,
	resp *types.DomainRenewResponse,
	domain string,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "DOMAIN-RENEW-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-DOMAIN-RENEW" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.Result.DomainName != domain {
		t.Fatalf("unexpected domain name: %s", resp.Result.DomainName)
	}

	expectedExpiry := time.Date(2028, 6, 30, 9, 30, 0, 0, time.UTC)
	if !resp.Result.NewExpiryDate.Equal(expectedExpiry) {
		t.Fatalf("unexpected new expiry date: %s", resp.Result.NewExpiryDate.Format(time.RFC3339))
	}
}

func domainRenewResponse(domain string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:renData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>` + domain + `</domain:name>
                <domain:exDate>2028-06-30T09:30:00Z</domain:exDate>
            </domain:renData>
        </resData>
        <trID>
            <clTRID>DOMAIN-RENEW-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-RENEW</svTRID>
        </trID>
    </response>
</epp>`
}
