package test

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestHostCheckXMLGenerationParsingAndMultipleNames(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, hostCheckResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostCheck(types.HostCheckRequest{
		Hosts: []string{
			"ns1.example.in",
			"ns2.example.in",
		},
	})
	if err != nil {
		t.Fatalf("host check failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:host="urn:ietf:params:xml:ns:host-1.0"`)
	assertContains(t, requestXML, `<host:check>`)
	assertContains(t, requestXML, `<host:name>ns1.example.in</host:name>`)
	assertContains(t, requestXML, `<host:name>ns2.example.in</host:name>`)
	assertContains(t, requestXML, `<clTRID>CHECK-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "HOST-CHECK-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-HOST-CHECK" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(resp.Results))
	}

	assertHostCheckResult(t, resp.Results[0], "ns1.example.in", "ns1.example.in", true, "")
	assertHostCheckResult(t, resp.Results[1], "ns2.example.in", "ns2.example.in", false, "In use")
}

func TestHostCheckUnicodeConversion(t *testing.T) {
	unicodeHost := "ns1.भारत.भारत"
	asciiHost, err := idn.ToASCII(unicodeHost)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, hostCheckUnicodeResponse(asciiHost))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostCheck(types.HostCheckRequest{
		Hosts: []string{
			unicodeHost,
		},
	})
	if err != nil {
		t.Fatalf("host check failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<host:name>"+asciiHost+"</host:name>")

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}

	assertHostCheckResult(t, resp.Results[0], unicodeHost, asciiHost, true, "")
}

func TestHostCheckEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostCheck(types.HostCheckRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertHostCheckResult(
	t *testing.T,
	result types.HostCheckResult,
	hostName string,
	asciiName string,
	available bool,
	reason string,
) {
	t.Helper()

	if result.HostName != hostName {
		t.Fatalf("expected host name %s, got %s", hostName, result.HostName)
	}

	if result.ASCIIName != asciiName {
		t.Fatalf("expected ASCII name %s, got %s", asciiName, result.ASCIIName)
	}

	if result.Available != available {
		t.Fatalf("expected availability %t for %s, got %t", available, hostName, result.Available)
	}

	if result.Reason != reason {
		t.Fatalf("expected reason %q for %s, got %q", reason, hostName, result.Reason)
	}
}

func hostCheckResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <host:chkData xmlns:host="urn:ietf:params:xml:ns:host-1.0">
                <host:cd>
                    <host:name avail="1">ns1.example.in</host:name>
                </host:cd>
                <host:cd>
                    <host:name avail="0">ns2.example.in</host:name>
                    <host:reason>In use</host:reason>
                </host:cd>
            </host:chkData>
        </resData>
        <trID>
            <clTRID>HOST-CHECK-TEST</clTRID>
            <svTRID>SERVER-HOST-CHECK</svTRID>
        </trID>
    </response>
</epp>`
}

func hostCheckUnicodeResponse(host string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <host:chkData xmlns:host="urn:ietf:params:xml:ns:host-1.0">
                <host:cd>
                    <host:name avail="1">` + host + `</host:name>
                </host:cd>
            </host:chkData>
        </resData>
        <trID>
            <clTRID>HOST-CHECK-UNICODE-TEST</clTRID>
            <svTRID>SERVER-HOST-CHECK-UNICODE</svTRID>
        </trID>
    </response>
</epp>`
}
