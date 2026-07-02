package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/fee"
	"github.com/hariom-pal/go-epp/types"
)

func TestFeeDomainCheckExtensionXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, feeDomainCheckResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{"example.com"},
		Fee: &fee.CheckRequest{
			Domains: []fee.CheckDomain{
				{
					Name:     "example.com",
					Currency: "USD",
					Command: fee.Command{
						Name:  fee.CommandCreate,
						Phase: "sunrise",
					},
					Period: &fee.Period{
						Value: 1,
						Unit:  "y",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain check failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<extension>`)
	assertContains(t, requestXML, `<fee:check xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">`)
	assertContains(t, requestXML, `<fee:domain>`)
	assertContains(t, requestXML, `<fee:name>example.com</fee:name>`)
	assertContains(t, requestXML, `<fee:currency>USD</fee:currency>`)
	assertContains(t, requestXML, `<fee:command phase="sunrise">create</fee:command>`)
	assertContains(t, requestXML, `<fee:period unit="y">1</fee:period>`)

	if len(resp.Fee.Results) != 1 {
		t.Fatalf("expected one fee result, got %d", len(resp.Fee.Results))
	}

	result := resp.Fee.Results[0]
	if result.Name != "example.com" ||
		result.Currency != "USD" ||
		result.Command.Name != fee.CommandCreate ||
		result.Command.Phase != "sunrise" ||
		result.Class != "premium-tier1" ||
		len(result.Fees) != 2 ||
		len(result.Credits) != 1 {

		t.Fatalf("unexpected fee check result: %#v", result)
	}

	if result.Period == nil ||
		result.Period.Value != 1 ||
		result.Period.Unit != "y" {

		t.Fatalf("unexpected fee period: %#v", result.Period)
	}

	assertFeeAmount(t, result.Fees[0], "5.00", "Application Fee", "0", "", fee.AppliedImmediate)
	assertFeeAmount(t, result.Fees[1], "7.00", "Registration Fee", "1", "P5D", "")
	assertFeeAmount(t, result.Credits[0], "-1.00", "Promo Credit", "", "", "")
}

func TestFeeDomainCreateExtensionXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, feeDomainCreateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.com",
		Period:     1,
		Unit:       "y",
		Registrant: "REG123",
		AuthInfo:   "secret",
		Fee: &fee.TransformRequest{
			Currency: "USD",
			Fees: []fee.Fee{
				{Amount: "5.00"},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<fee:create xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">`)
	assertContains(t, requestXML, `<fee:currency>USD</fee:currency>`)
	assertContains(t, requestXML, `<fee:fee>5.00</fee:fee>`)

	if resp.Result.Fee.Currency != "USD" ||
		resp.Result.Fee.Balance != "-5.00" ||
		resp.Result.Fee.CreditLimit != "1000.00" ||
		len(resp.Result.Fee.Fees) != 1 {

		t.Fatalf("unexpected create fee data: %#v", resp.Result.Fee)
	}

	assertFeeAmount(t, resp.Result.Fee.Fees[0], "5.00", "", "", "P5D", "")
}

func TestFeeDomainRenewExtensionXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, feeDomainRenewResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName:        "example.com",
		CurrentExpiryDate: time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
		Period:            1,
		Unit:              "y",
		Fee: &fee.TransformRequest{
			Currency: "USD",
			Fees: []fee.Fee{
				{Amount: "8.00"},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain renew failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<fee:renew xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">`)
	assertContains(t, requestXML, `<fee:fee>8.00</fee:fee>`)

	if resp.Result.Fee.Currency != "USD" ||
		resp.Result.Fee.Balance != "1000.00" ||
		len(resp.Result.Fee.Fees) != 1 {

		t.Fatalf("unexpected renew fee data: %#v", resp.Result.Fee)
	}
}

func TestFeeDomainTransferExtensionXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, feeDomainTransferResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferRequest,
		AuthInfo:   "secret123",
		Fee: &fee.TransformRequest{
			Currency: "USD",
			Fees: []fee.Fee{
				{Amount: "5.00"},
			},
		},
	})
	if err != nil {
		t.Fatalf("domain transfer failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<fee:transfer xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">`)
	assertContains(t, requestXML, `<fee:currency>USD</fee:currency>`)
	assertContains(t, requestXML, `<fee:fee>5.00</fee:fee>`)

	if resp.TransferData.Fee.Currency != "USD" ||
		resp.TransferData.Fee.Balance != "995.00" ||
		resp.TransferData.Fee.CreditLimit != "1000.00" ||
		len(resp.TransferData.Fee.Fees) != 1 {

		t.Fatalf("unexpected transfer fee data: %#v", resp.TransferData.Fee)
	}

	assertFeeAmount(t, resp.TransferData.Fee.Fees[0], "5.00", "", "", "P5D", "")
}

func assertFeeAmount(
	t *testing.T,
	got fee.Amount,
	amount string,
	description string,
	refundable string,
	gracePeriod string,
	applied string,
) {
	t.Helper()

	if got.Amount != amount ||
		got.Description != description ||
		got.Refundable != refundable ||
		got.GracePeriod != gracePeriod ||
		got.Applied != applied {

		t.Fatalf("unexpected amount: %#v", got)
	}
}

func feeDomainCheckResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:cd>
                    <domain:name avail="1">example.com</domain:name>
                </domain:cd>
            </domain:chkData>
        </resData>
        <extension>
            <fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">
                <fee:cd>
                    <fee:name>example.com</fee:name>
                    <fee:currency>USD</fee:currency>
                    <fee:command phase="sunrise">create</fee:command>
                    <fee:period unit="y">1</fee:period>
                    <fee:fee description="Application Fee" refundable="0" applied="immediate">5.00</fee:fee>
                    <fee:fee description="Registration Fee" refundable="1" grace-period="P5D">7.00</fee:fee>
                    <fee:credit description="Promo Credit">-1.00</fee:credit>
                    <fee:class>premium-tier1</fee:class>
                </fee:cd>
            </fee:chkData>
        </extension>
        <trID>
            <clTRID>DOMAIN-CHECK-FEE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-CHECK-FEE</svTRID>
        </trID>
    </response>
</epp>`
}

func feeDomainCreateResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:creData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>example.com</domain:name>
                <domain:crDate>2026-06-30T09:30:00Z</domain:crDate>
                <domain:exDate>2027-06-30T09:30:00Z</domain:exDate>
            </domain:creData>
        </resData>
        <extension>
            <fee:creData xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">
                <fee:currency>USD</fee:currency>
                <fee:fee grace-period="P5D">5.00</fee:fee>
                <fee:balance>-5.00</fee:balance>
                <fee:creditLimit>1000.00</fee:creditLimit>
            </fee:creData>
        </extension>
        <trID>
            <clTRID>CREATE-FEE-TEST</clTRID>
            <svTRID>SERVER-FEE-TEST</svTRID>
        </trID>
    </response>
</epp>`
}

func feeDomainRenewResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:renData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>example.com</domain:name>
                <domain:exDate>2028-06-30T09:30:00Z</domain:exDate>
            </domain:renData>
        </resData>
        <extension>
            <fee:renData xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">
                <fee:currency>USD</fee:currency>
                <fee:fee grace-period="P5D">8.00</fee:fee>
                <fee:balance>1000.00</fee:balance>
            </fee:renData>
        </extension>
        <trID>
            <clTRID>DOMAIN-RENEW-FEE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-RENEW-FEE</svTRID>
        </trID>
    </response>
</epp>`
}

func feeDomainTransferResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:trnData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>example.in</domain:name>
                <domain:trStatus>pending</domain:trStatus>
                <domain:reID>ClientX</domain:reID>
                <domain:reDate>2026-07-01T10:30:00Z</domain:reDate>
                <domain:acID>ClientY</domain:acID>
                <domain:acDate>2026-07-06T10:30:00Z</domain:acDate>
                <domain:exDate>2027-07-01T10:30:00Z</domain:exDate>
            </domain:trnData>
        </resData>
        <extension>
            <fee:trnData xmlns:fee="urn:ietf:params:xml:ns:fee-0.7">
                <fee:currency>USD</fee:currency>
                <fee:fee grace-period="P5D">5.00</fee:fee>
                <fee:balance>995.00</fee:balance>
                <fee:creditLimit>1000.00</fee:creditLimit>
            </fee:trnData>
        </extension>
        <trID>
            <clTRID>DOMAIN-TRANSFER-FEE-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-TRANSFER-FEE</svTRID>
        </trID>
    </response>
</epp>`
}
