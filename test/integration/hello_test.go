package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/epp"
)

func TestHelloXMLGenerationAndGreetingParsing(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, helloGreetingResponse())
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	greeting, err := client.Hello()
	if err != nil {
		t.Fatalf("hello failed: %v", err)
	}

	requestXML := readRequest(t, requests)
	expectedXML := `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <hello/>
</epp>`
	if requestXML != expectedXML {
		t.Fatalf("unexpected hello XML:\n%s", requestXML)
	}

	if greeting.ServerID != "hello EPP server" {
		t.Fatalf("unexpected server ID: %s", greeting.ServerID)
	}

	expectedDate := time.Date(2026, 7, 1, 12, 30, 0, 0, time.UTC)
	if greeting.ServerDate == nil ||
		!greeting.ServerDate.Equal(expectedDate) {

		t.Fatalf("unexpected server date: %#v", greeting.ServerDate)
	}

	assertStringSlice(t, greeting.Versions, []string{"1.0"})
	assertStringSlice(t, greeting.Languages, []string{"en", "hi"})
	assertStringSlice(t, greeting.SupportedObjects, []string{
		"urn:ietf:params:xml:ns:domain-1.0",
		"urn:ietf:params:xml:ns:contact-1.0",
		"urn:ietf:params:xml:ns:host-1.0",
	})
	assertStringSlice(t, greeting.SupportedExtensions, []string{
		"urn:ietf:params:xml:ns:secDNS-1.1",
		"urn:ietf:params:xml:ns:rgp-1.0",
	})

	if greeting.DCP.Access != "all" {
		t.Fatalf("unexpected DCP access: %s", greeting.DCP.Access)
	}

	if len(greeting.DCP.Statements) != 1 {
		t.Fatalf("expected one DCP statement, got %d", len(greeting.DCP.Statements))
	}

	statement := greeting.DCP.Statements[0]
	assertStringSlice(t, statement.Purposes, []string{"admin", "prov"})
	assertStringSlice(t, statement.Recipients, []string{"ours", "public"})
	assertStringSlice(t, statement.Retentions, []string{"stated"})
}

func TestGreetingInfoParsesInitialGreeting(t *testing.T) {
	cfg, _, cleanup := startDomainCreateServer(t, simpleEPPResponse("UNUSED", "UNUSED"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	greeting, err := client.GreetingInfo()
	if err != nil {
		t.Fatalf("parse initial greeting failed: %v", err)
	}

	if greeting.ServerID != "test EPP server" {
		t.Fatalf("unexpected server ID: %s", greeting.ServerID)
	}

	assertStringSlice(t, greeting.Versions, []string{"1.0"})
	assertStringSlice(t, greeting.Languages, []string{"en"})
	assertStringSlice(t, greeting.SupportedObjects, []string{
		"urn:ietf:params:xml:ns:domain-1.0",
	})
}

func assertStringSlice(
	t *testing.T,
	got []string,
	want []string,
) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func helloGreetingResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <greeting>
        <svID>hello EPP server</svID>
        <svDate>2026-07-01T12:30:00Z</svDate>
        <svcMenu>
            <version>1.0</version>
            <lang>en</lang>
            <lang>hi</lang>
            <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:host-1.0</objURI>
            <svcExtension>
                <extURI>urn:ietf:params:xml:ns:secDNS-1.1</extURI>
                <extURI>urn:ietf:params:xml:ns:rgp-1.0</extURI>
            </svcExtension>
        </svcMenu>
        <dcp>
            <access><all/></access>
            <statement>
                <purpose>
                    <admin/>
                    <prov/>
                </purpose>
                <recipient>
                    <ours/>
                    <public/>
                </recipient>
                <retention>
                    <stated/>
                </retention>
            </statement>
        </dcp>
    </greeting>
</epp>`
}
