package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestHostDeleteXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, hostDeleteResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostDelete(types.HostDeleteRequest{
		HostName: "ns1.example.in",
	})
	if err != nil {
		t.Fatalf("host delete failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:host="urn:ietf:params:xml:ns:host-1.0"`)
	assertContains(t, requestXML, `<host:delete>`)
	assertContains(t, requestXML, `<host:name>ns1.example.in</host:name>`)
	assertContains(t, requestXML, `<clTRID>DELETE-`)

	assertHostDeleteResponse(t, resp)
}

func TestHostDeleteUnicodeConversion(t *testing.T) {
	unicodeHost := "ns1.भारत.भारत"
	asciiHost, err := idn.ToASCII(unicodeHost)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, hostDeleteResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostDelete(types.HostDeleteRequest{
		HostName: unicodeHost,
	})
	if err != nil {
		t.Fatalf("host delete failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<host:name>"+asciiHost+"</host:name>")
	assertHostDeleteResponse(t, resp)
}

func TestHostDeleteEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostDelete(types.HostDeleteRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertHostDeleteResponse(
	t *testing.T,
	resp *types.HostDeleteResponse,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "HOST-DELETE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-HOST-DELETE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}
}

func hostDeleteResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>HOST-DELETE-TEST</clTRID>
            <svTRID>SERVER-HOST-DELETE</svTRID>
        </trID>
    </response>
</epp>`
}
