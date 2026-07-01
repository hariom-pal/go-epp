package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestHostCreateXMLGenerationAndParsing(t *testing.T) {
	resp, requestXML := runHostCreateTest(t, types.HostCreateRequest{
		HostName: "ns1.example.in",
		Addresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.1"},
			{IPVersion: "v6", Address: "2001:db8::1"},
		},
	}, "ns1.example.in")

	assertContains(t, requestXML, `xmlns:host="urn:ietf:params:xml:ns:host-1.0"`)
	assertContains(t, requestXML, `<host:create>`)
	assertContains(t, requestXML, `<host:name>ns1.example.in</host:name>`)
	assertContains(t, requestXML, `<host:addr ip="v4">192.0.2.1</host:addr>`)
	assertContains(t, requestXML, `<host:addr ip="v6">2001:db8::1</host:addr>`)
	assertContains(t, requestXML, `<clTRID>CREATE-`)

	assertHostCreateResponse(t, resp, "ns1.example.in")
}

func TestHostCreateIPv4Only(t *testing.T) {
	_, requestXML := runHostCreateTest(t, types.HostCreateRequest{
		HostName: "ns1.example.in",
		Addresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.1"},
		},
	}, "ns1.example.in")

	assertContains(t, requestXML, `<host:addr ip="v4">192.0.2.1</host:addr>`)
	assertNotContains(t, requestXML, `<host:addr ip="v6">`)
}

func TestHostCreateIPv6Only(t *testing.T) {
	_, requestXML := runHostCreateTest(t, types.HostCreateRequest{
		HostName: "ns1.example.in",
		Addresses: []types.HostAddress{
			{IPVersion: "v6", Address: "2001:db8::1"},
		},
	}, "ns1.example.in")

	assertContains(t, requestXML, `<host:addr ip="v6">2001:db8::1</host:addr>`)
	assertNotContains(t, requestXML, `<host:addr ip="v4">`)
}

func TestHostCreateUnicodeConversion(t *testing.T) {
	unicodeHost := "ns1.भारत.भारत"
	asciiHost, err := idn.ToASCII(unicodeHost)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	resp, requestXML := runHostCreateTest(t, types.HostCreateRequest{
		HostName: unicodeHost,
		Addresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.1"},
		},
	}, asciiHost)

	assertContains(t, requestXML, "<host:name>"+asciiHost+"</host:name>")

	if resp.Result.HostName != unicodeHost {
		t.Fatalf("unexpected unicode host name: %s", resp.Result.HostName)
	}
}

func TestHostCreateEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostCreate(types.HostCreateRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestHostCreateInvalidIP(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostCreate(types.HostCreateRequest{
		HostName: "ns1.example.in",
		Addresses: []types.HostAddress{
			{IPVersion: "v4", Address: "not-an-ip"},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestHostCreateInvalidIPVersion(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostCreate(types.HostCreateRequest{
		HostName: "ns1.example.in",
		Addresses: []types.HostAddress{
			{IPVersion: "v5", Address: "192.0.2.1"},
		},
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func runHostCreateTest(
	t *testing.T,
	req types.HostCreateRequest,
	responseHost string,
) (*types.HostCreateResponse, string) {
	t.Helper()

	cfg, requests, cleanup := startDomainCreateServer(t, hostCreateResponse(responseHost))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostCreate(req)
	if err != nil {
		t.Fatalf("host create failed: %v", err)
	}

	return resp, readRequest(t, requests)
}

func assertHostCreateResponse(
	t *testing.T,
	resp *types.HostCreateResponse,
	host string,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "HOST-CREATE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-HOST-CREATE" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.Result.HostName != host {
		t.Fatalf("unexpected host name: %s", resp.Result.HostName)
	}

	expectedCreated := time.Date(2026, 7, 1, 12, 30, 0, 0, time.UTC)
	if !resp.Result.CreatedDate.Equal(expectedCreated) {
		t.Fatalf("unexpected created date: %s", resp.Result.CreatedDate.Format(time.RFC3339))
	}
}

func hostCreateResponse(host string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <host:creData xmlns:host="urn:ietf:params:xml:ns:host-1.0">
                <host:name>` + host + `</host:name>
                <host:crDate>2026-07-01T12:30:00Z</host:crDate>
            </host:creData>
        </resData>
        <trID>
            <clTRID>HOST-CREATE-TEST</clTRID>
            <svTRID>SERVER-HOST-CREATE</svTRID>
        </trID>
    </response>
</epp>`
}
