package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestContactDeleteXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, contactDeleteResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.ContactDelete(types.ContactDeleteRequest{
		ContactID: "CNT001",
	})
	if err != nil {
		t.Fatalf("contact delete failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"`)
	assertContains(t, requestXML, `<contact:delete>`)
	assertContains(t, requestXML, `<contact:id>CNT001</contact:id>`)
	assertContains(t, requestXML, `<clTRID>DELETE-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "CONTACT-DELETE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-CONTACT-DELETE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}
}

func TestContactDeleteEmptyContactIDValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.ContactDelete(types.ContactDeleteRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func contactDeleteResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>CONTACT-DELETE-TEST</clTRID>
            <svTRID>SERVER-CONTACT-DELETE</svTRID>
        </trID>
    </response>
</epp>`
}
