package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestContactCheckXMLGenerationParsingAndMultipleIDs(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactCheckResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactCheck(types.ContactCheckRequest{
		IDs: []string{
			"CNT001",
			"CNT002",
			"CNT003",
		},
	})
	if err != nil {
		t.Fatalf("contact check failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"`)
	assertContains(t, requestXML, `<contact:check>`)
	assertContains(t, requestXML, `<contact:id>CNT001</contact:id>`)
	assertContains(t, requestXML, `<contact:id>CNT002</contact:id>`)
	assertContains(t, requestXML, `<contact:id>CNT003</contact:id>`)
	assertContains(t, requestXML, `<clTRID>CHECK-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "CONTACT-CHECK-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-CONTACT-CHECK" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if len(resp.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(resp.Results))
	}

	assertContactCheckResult(t, resp.Results[0], "CNT001", true, "")
	assertContactCheckResult(t, resp.Results[1], "CNT002", false, "In use")
	assertContactCheckResult(t, resp.Results[2], "CNT003", true, "")
}

func TestContactCheckEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.ContactCheck(types.ContactCheckRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertContactCheckResult(
	t *testing.T,
	result types.ContactCheckResult,
	contactID string,
	available bool,
	reason string,
) {
	t.Helper()

	if result.ContactID != contactID {
		t.Fatalf("expected contact ID %s, got %s", contactID, result.ContactID)
	}

	if result.ID != contactID {
		t.Fatalf("expected compatibility ID %s, got %s", contactID, result.ID)
	}

	if result.Available != available {
		t.Fatalf("expected availability %t for %s, got %t", available, contactID, result.Available)
	}

	if result.Reason != reason {
		t.Fatalf("expected reason %q for %s, got %q", reason, contactID, result.Reason)
	}
}

func contactCheckResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <contact:chkData xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                <contact:cd>
                    <contact:id avail="1">CNT001</contact:id>
                </contact:cd>
                <contact:cd>
                    <contact:id avail="0">CNT002</contact:id>
                    <contact:reason>In use</contact:reason>
                </contact:cd>
                <contact:cd>
                    <contact:id avail="1">CNT003</contact:id>
                </contact:cd>
            </contact:chkData>
        </resData>
        <trID>
            <clTRID>CONTACT-CHECK-TEST</clTRID>
            <svTRID>SERVER-CONTACT-CHECK</svTRID>
        </trID>
    </response>
</epp>`
}
