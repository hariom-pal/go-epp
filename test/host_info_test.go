package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestHostInfoXMLGenerationAndParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, hostInfoResponse("ns1.example.in"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostInfo(types.HostInfoRequest{
		HostName: "ns1.example.in",
	})
	if err != nil {
		t.Fatalf("host info failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:host="urn:ietf:params:xml:ns:host-1.0"`)
	assertContains(t, requestXML, `<host:info>`)
	assertContains(t, requestXML, `<host:name>ns1.example.in</host:name>`)
	assertContains(t, requestXML, `<clTRID>INFO-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "HOST-INFO-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-HOST-INFO" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	host := resp.Host
	if host.HostName != "ns1.example.in" {
		t.Fatalf("unexpected host name: %s", host.HostName)
	}

	if host.ASCIIName != "ns1.example.in" {
		t.Fatalf("unexpected ASCII name: %s", host.ASCIIName)
	}

	if host.ROID != "NS1_EXAMPLE_IN-ROID" {
		t.Fatalf("unexpected ROID: %s", host.ROID)
	}

	if len(host.Statuses) != 2 ||
		host.Statuses[0] != "linked" ||
		host.Statuses[1] != "serverUpdateProhibited" {

		t.Fatalf("unexpected statuses: %#v", host.Statuses)
	}

	if len(host.Addresses) != 2 {
		t.Fatalf("expected 2 addresses, got %#v", host.Addresses)
	}

	assertHostAddress(t, host.Addresses[0], "v4", "192.0.2.2")
	assertHostAddress(t, host.Addresses[1], "v6", "2001:db8::2")

	if host.ClientID != "ClientY" ||
		host.CreatedBy != "ClientX" ||
		host.UpdatedBy != "ClientZ" {

		t.Fatalf("unexpected client fields: clID=%s crID=%s upID=%s",
			host.ClientID,
			host.CreatedBy,
			host.UpdatedBy,
		)
	}

	if !host.CreatedDate.Equal(time.Date(2026, 7, 1, 9, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected created date: %s", host.CreatedDate.Format(time.RFC3339))
	}

	if !host.UpdatedDate.Equal(time.Date(2026, 7, 1, 10, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected updated date: %s", host.UpdatedDate.Format(time.RFC3339))
	}

	if !host.TransferDate.Equal(time.Date(2026, 7, 1, 11, 30, 0, 0, time.UTC)) {
		t.Fatalf("unexpected transfer date: %s", host.TransferDate.Format(time.RFC3339))
	}
}

func TestHostInfoUnicodeConversion(t *testing.T) {
	unicodeHost := "ns1.भारत.भारत"
	asciiHost, err := idn.ToASCII(unicodeHost)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, hostInfoMinimalResponse(asciiHost))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.HostInfo(types.HostInfoRequest{
		HostName: unicodeHost,
	})
	if err != nil {
		t.Fatalf("host info failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<host:name>"+asciiHost+"</host:name>")

	if resp.Host.HostName != unicodeHost {
		t.Fatalf("unexpected unicode host name: %s", resp.Host.HostName)
	}

	if resp.Host.ASCIIName != asciiHost {
		t.Fatalf("unexpected ASCII host name: %s", resp.Host.ASCIIName)
	}
}

func TestHostInfoEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.HostInfo(types.HostInfoRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func assertHostAddress(
	t *testing.T,
	address types.HostAddress,
	ipVersion string,
	ip string,
) {
	t.Helper()

	if address.IPVersion != ipVersion {
		t.Fatalf("expected IP version %s, got %s", ipVersion, address.IPVersion)
	}

	if address.Address != ip {
		t.Fatalf("expected IP address %s, got %s", ip, address.Address)
	}
}

func hostInfoResponse(host string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <host:infData xmlns:host="urn:ietf:params:xml:ns:host-1.0">
                <host:name>` + host + `</host:name>
                <host:roid>NS1_EXAMPLE_IN-ROID</host:roid>
                <host:status s="linked"/>
                <host:status s="serverUpdateProhibited"/>
                <host:addr ip="v4">192.0.2.2</host:addr>
                <host:addr ip="v6">2001:db8::2</host:addr>
                <host:clID>ClientY</host:clID>
                <host:crID>ClientX</host:crID>
                <host:crDate>2026-07-01T09:30:00Z</host:crDate>
                <host:upID>ClientZ</host:upID>
                <host:upDate>2026-07-01T10:30:00Z</host:upDate>
                <host:trDate>2026-07-01T11:30:00Z</host:trDate>
            </host:infData>
        </resData>
        <trID>
            <clTRID>HOST-INFO-TEST</clTRID>
            <svTRID>SERVER-HOST-INFO</svTRID>
        </trID>
    </response>
</epp>`
}

func hostInfoMinimalResponse(host string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <host:infData xmlns:host="urn:ietf:params:xml:ns:host-1.0">
                <host:name>` + host + `</host:name>
                <host:roid>HOST-ROID</host:roid>
                <host:clID>ClientY</host:clID>
                <host:crID>ClientX</host:crID>
                <host:crDate>2026-07-01T09:30:00Z</host:crDate>
            </host:infData>
        </resData>
        <trID>
            <clTRID>HOST-INFO-UNICODE-TEST</clTRID>
            <svTRID>SERVER-HOST-INFO-UNICODE</svTRID>
        </trID>
    </response>
</epp>`
}
