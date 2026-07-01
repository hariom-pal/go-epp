package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestHostUpdateXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, hostUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostUpdate(types.HostUpdateRequest{
		HostName: "ns1.example.in",
		AddAddresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.10"},
			{IPVersion: "v6", Address: "2001:db8::10"},
		},
		RemoveAddresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.1"},
			{IPVersion: "v6", Address: "2001:db8::1"},
		},
		AddStatuses: []string{
			constants.HostStatusClientUpdateProhibited,
		},
		RemoveStatuses: []string{
			constants.HostStatusClientDeleteProhibited,
		},
		NewHostName: "ns2.example.in",
	})
	if err != nil {
		t.Fatalf("host update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:host="urn:ietf:params:xml:ns:host-1.0"`)
	assertContains(t, requestXML, `<host:update>`)
	assertContains(t, requestXML, `<host:name>ns1.example.in</host:name>`)
	assertContains(t, requestXML, `<host:add>`)
	assertContains(t, requestXML, `<host:addr ip="v4">192.0.2.10</host:addr>`)
	assertContains(t, requestXML, `<host:addr ip="v6">2001:db8::10</host:addr>`)
	assertContains(t, requestXML, `<host:status s="clientUpdateProhibited"></host:status>`)
	assertContains(t, requestXML, `<host:rem>`)
	assertContains(t, requestXML, `<host:addr ip="v4">192.0.2.1</host:addr>`)
	assertContains(t, requestXML, `<host:addr ip="v6">2001:db8::1</host:addr>`)
	assertContains(t, requestXML, `<host:status s="clientDeleteProhibited"></host:status>`)
	assertContains(t, requestXML, `<host:chg>`)
	assertContains(t, requestXML, `<host:name>ns2.example.in</host:name>`)
	assertContains(t, requestXML, `<clTRID>UPDATE-`)

	assertHostUpdateResponse(t, resp)
}

func TestHostUpdateAddressAddRemoveOnly(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, hostUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.HostUpdate(types.HostUpdateRequest{
		HostName: "ns1.example.in",
		AddAddresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.10"},
		},
		RemoveAddresses: []types.HostAddress{
			{IPVersion: "v6", Address: "2001:db8::1"},
		},
	})
	if err != nil {
		t.Fatalf("host update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<host:add>`)
	assertContains(t, requestXML, `<host:addr ip="v4">192.0.2.10</host:addr>`)
	assertContains(t, requestXML, `<host:rem>`)
	assertContains(t, requestXML, `<host:addr ip="v6">2001:db8::1</host:addr>`)
	assertNotContains(t, requestXML, `<host:chg>`)
}

func TestHostUpdateUnicodeConversion(t *testing.T) {
	currentHost := "ns1.भारत.भारत"
	newHost := "ns2.भारत.भारत"

	currentASCII, err := idn.ToASCII(currentHost)
	if err != nil {
		t.Fatalf("failed to prepare current punycode: %v", err)
	}

	newASCII, err := idn.ToASCII(newHost)
	if err != nil {
		t.Fatalf("failed to prepare new punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, hostUpdateResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.HostUpdate(types.HostUpdateRequest{
		HostName:    currentHost,
		NewHostName: newHost,
	})
	if err != nil {
		t.Fatalf("host update failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<host:name>"+currentASCII+"</host:name>")
	assertContains(t, requestXML, "<host:name>"+newASCII+"</host:name>")
}

func TestHostUpdateEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostUpdate(types.HostUpdateRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestHostUpdateEmptyOperationValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostUpdate(types.HostUpdateRequest{
		HostName: "ns1.example.in",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestHostUpdateInvalidIP(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostUpdate(types.HostUpdateRequest{
		HostName: "ns1.example.in",
		AddAddresses: []types.HostAddress{
			{IPVersion: "v4", Address: "not-an-ip"},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestHostUpdateInvalidIPVersion(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostUpdate(types.HostUpdateRequest{
		HostName: "ns1.example.in",
		AddAddresses: []types.HostAddress{
			{IPVersion: "v5", Address: "192.0.2.10"},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestHostUpdateEmptyStatusValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostUpdate(types.HostUpdateRequest{
		HostName:    "ns1.example.in",
		AddStatuses: []string{" "},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertHostUpdateResponse(
	t *testing.T,
	resp *types.HostUpdateResponse,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "HOST-UPDATE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-HOST-UPDATE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}
}

func hostUpdateResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>HOST-UPDATE-TEST</clTRID>
            <svTRID>SERVER-HOST-UPDATE</svTRID>
        </trID>
    </response>
</epp>`
}
