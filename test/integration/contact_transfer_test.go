package test

import (
	"strings"
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestContactTransferRequestXMLAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactTransferResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactTransfer(types.ContactTransferRequest{
		ContactID: "CNT123",
		Operation: constants.TransferRequest,
		AuthInfo:  "secret",
	})
	if err != nil {
		t.Fatalf("contact transfer failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<transfer op="request">`)
	assertContains(t, requestXML, `<contact:transfer>`)
	assertContains(t, requestXML, `<contact:id>CNT123</contact:id>`)
	assertContains(t, requestXML, `<contact:pw>secret</contact:pw>`)

	if resp.TransferData.ContactID != "CNT123" ||
		resp.TransferData.TransferStatus != "pending" ||
		resp.Result.Status != "pending" {
		t.Fatalf("unexpected contact transfer response: %+v", resp)
	}
}

func TestContactTransferQueryXML(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactTransferResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.ContactTransfer(types.ContactTransferRequest{
		ContactID: "CNT123",
		Operation: constants.TransferQuery,
	})
	if err != nil {
		t.Fatalf("contact transfer query failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	assertContains(t, requestXML, `<transfer op="query">`)
	if strings.Contains(requestXML, `<contact:authInfo>`) {
		t.Fatalf("query transfer should not include authInfo: %s", requestXML)
	}
}

func TestContactTransferRequestRequiresAuthInfo(t *testing.T) {
	var client epp.Client
	_, err := client.ContactTransfer(types.ContactTransferRequest{
		ContactID: "CNT123",
		Operation: constants.TransferRequest,
	})
	if err == nil {
		t.Fatal("expected authInfo validation error")
	}
}

func contactTransferResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1001">
            <msg>Command completed successfully; action pending</msg>
        </result>
        <resData>
            <contact:trnData xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:id>CNT123</contact:id>
                <contact:trStatus>pending</contact:trStatus>
                <contact:reID>Registrar-A</contact:reID>
                <contact:reDate>2026-06-30T09:30:00Z</contact:reDate>
                <contact:acID>Registrar-B</contact:acID>
                <contact:acDate>2026-07-05T09:30:00Z</contact:acDate>
            </contact:trnData>
        </resData>
        <trID>
            <clTRID>TRANSFER-TEST</clTRID>
            <svTRID>SERVER-TEST</svTRID>
        </trID>
    </response>
</epp>`
}
