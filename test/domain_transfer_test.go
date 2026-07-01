package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainTransferRequestXML(t *testing.T) {
	resp, requestXML := runDomainTransferTest(t, types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferRequest,
		AuthInfo:   "secret123",
		PeriodInfo: types.Period{
			Value: 1,
			Unit:  "y",
		},
	})

	assertContains(t, requestXML, `xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"`)
	assertContains(t, requestXML, `<transfer op="request">`)
	assertContains(t, requestXML, `<domain:transfer>`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
	assertContains(t, requestXML, `<domain:period unit="y">1</domain:period>`)
	assertContains(t, requestXML, `<domain:authInfo>`)
	assertContains(t, requestXML, `<domain:pw>secret123</domain:pw>`)
	assertContains(t, requestXML, `<clTRID>TRANSFER-`)

	assertDomainTransferResponse(t, resp)
}

func TestDomainTransferQueryXML(t *testing.T) {
	_, requestXML := runDomainTransferTest(t, types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferQuery,
	})

	assertContains(t, requestXML, `<transfer op="query">`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
	assertNotContains(t, requestXML, `<domain:authInfo>`)
	assertNotContains(t, requestXML, `<domain:period`)
}

func TestDomainTransferApproveXML(t *testing.T) {
	_, requestXML := runDomainTransferTest(t, types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferApprove,
	})

	assertContains(t, requestXML, `<transfer op="approve">`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
}

func TestDomainTransferCancelXML(t *testing.T) {
	_, requestXML := runDomainTransferTest(t, types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferCancel,
	})

	assertContains(t, requestXML, `<transfer op="cancel">`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
}

func TestDomainTransferRejectXML(t *testing.T) {
	_, requestXML := runDomainTransferTest(t, types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferReject,
	})

	assertContains(t, requestXML, `<transfer op="reject">`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
}

func TestDomainTransferInvalidOperation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  "invalid",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainTransferEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainTransfer(types.DomainTransferRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainTransferRequestRequiresAuthInfo(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferRequest,
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func runDomainTransferTest(
	t *testing.T,
	req types.DomainTransferRequest,
) (*types.DomainTransferResponse, string) {
	t.Helper()

	cfg, requests, cleanup := startDomainCreateServer(t, domainTransferResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainTransfer(req)
	if err != nil {
		t.Fatalf("domain transfer failed: %v", err)
	}

	return resp, readRequest(t, requests)
}

func assertDomainTransferResponse(
	t *testing.T,
	resp *types.DomainTransferResponse,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "DOMAIN-TRANSFER-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-DOMAIN-TRANSFER" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.TransferData.DomainName != "example.in" {
		t.Fatalf("unexpected domain name: %s", resp.TransferData.DomainName)
	}

	if resp.TransferData.ObjectName != "example.in" {
		t.Fatalf("unexpected object name: %s", resp.TransferData.ObjectName)
	}

	if resp.TransferData.TransferStatus != "pending" {
		t.Fatalf("unexpected transfer status: %s", resp.TransferData.TransferStatus)
	}

	if resp.TransferData.RequestedBy != "ClientX" ||
		resp.TransferData.ActionBy != "ClientY" {

		t.Fatalf("unexpected transfer clients: reID=%s acID=%s",
			resp.TransferData.RequestedBy,
			resp.TransferData.ActionBy,
		)
	}

	if !resp.TransferData.RequestedDate.Equal(time.Date(2026, 7, 1, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected requested date: %s", resp.TransferData.RequestedDate.Format(time.RFC3339))
	}

	if !resp.TransferData.ActionDate.Equal(time.Date(2026, 7, 6, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected action date: %s", resp.TransferData.ActionDate.Format(time.RFC3339))
	}

	if !resp.TransferData.ExpiryDate.Equal(time.Date(2027, 7, 1, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected expiry date: %s", resp.TransferData.ExpiryDate.Format(time.RFC3339))
	}

	if resp.Result.Domain != "example.in" ||
		resp.Result.Status != "pending" {

		t.Fatalf("unexpected compatibility result: %#v", resp.Result)
	}
}

func domainTransferResponse() string {
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
        <trID>
            <clTRID>DOMAIN-TRANSFER-TEST</clTRID>
            <svTRID>SERVER-DOMAIN-TRANSFER</svTRID>
        </trID>
    </response>
</epp>`
}
