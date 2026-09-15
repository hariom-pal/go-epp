package test

import (
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
	"github.com/hariom-pal/go-epp/types"
)

// TestMalformedServerResponses covers hostile or broken server output. A
// registrar SDK must surface a clear error for each and must never panic,
// hang, or silently report success.
func TestMalformedServerResponses(t *testing.T) {
	cases := []struct {
		name     string
		response string
	}{
		{"not XML at all", `this is not xml`},
		{"truncated XML", `<?xml version="1.0"?><epp><response><result code="1000">`},
		{"empty payload", ``},
		{"wrong root element", `<html><body>502 Bad Gateway</body></html>`},
		{"no result element", `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><trID><svTRID>X</svTRID></trID></response></epp>`},
		{"non-numeric result code", `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="abc"><msg>?</msg></result></response></epp>`},
		{"greeting instead of response", `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>x</svID></greeting></epp>`},
		{"mismatched tags", `<epp><response><result code="1000"></response></result></epp>`},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, requests, cleanup := startDomainCreateServer(t, testCase.response)
			defer cleanup()

			client, err := epp.Connect(cfg)
			if err != nil {
				t.Fatalf("connect failed: %v", err)
			}
			defer client.Close()

			resp, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{"example.in"}})
			readRequest(t, requests)

			if err == nil {
				// Parsing may succeed on structurally valid but meaningless
				// XML; it must never claim a successful EPP result.
				if resp != nil && resp.ResultCode == 1000 {
					t.Fatalf("malformed response reported a successful EPP result: %q", testCase.response)
				}
				return
			}
			if strings.Contains(err.Error(), "panic") {
				t.Fatalf("malformed response caused a panic: %v", err)
			}
		})
	}
}

// TestResponseWithMismatchedClientTRID covers a server echoing a transaction
// identifier that does not belong to the command just sent. The SDK exposes
// the echoed value so a caller can detect the mismatch.
func TestResponseWithMismatchedClientTRID(t *testing.T) {
	response := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response>
<result code="1000"><msg>Command completed successfully</msg></result>
<resData><domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
<domain:cd><domain:name avail="1">example.in</domain:name></domain:cd>
</domain:chkData></resData>
<trID><clTRID>SOMEONE-ELSES-TRID</clTRID><svTRID>SV-1</svTRID></trID></response></epp>`

	cfg, requests, cleanup := startDomainCreateServer(t, response)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCheck(types.DomainCheckRequest{Domains: []string{"example.in"}})
	if err != nil {
		t.Fatalf("domain check failed: %v", err)
	}
	sent := readRequest(t, requests)

	if resp.ClientTRID != "SOMEONE-ELSES-TRID" {
		t.Fatalf("the echoed clTRID must be reported verbatim, got %q", resp.ClientTRID)
	}
	if strings.Contains(sent, "SOMEONE-ELSES-TRID") {
		t.Fatal("the sent command should not contain the server's echoed value")
	}
}

// TestMaxFrameSizeIsEnforced covers the configured wire-level cap, which stops
// a hostile or broken peer from forcing a large allocation.
func TestMaxFrameSizeIsEnforced(t *testing.T) {
	big := strings.Repeat("A", 200000)
	response := `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="1000"><msg>` +
		big + `</msg></result></response></epp>`

	cfg, requests, cleanup := startDomainCreateServer(t, response)
	defer cleanup()
	cfg.Transport.MaxFrameSize = 4096

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainCheck(types.DomainCheckRequest{Domains: []string{"example.in"}})
	readRequest(t, requests)

	if err == nil {
		t.Fatal("expected an oversized frame to be rejected")
	}
	if !errors.Is(err, epp.ErrFrameTooLarge) {
		t.Fatalf("expected ErrFrameTooLarge, got %v", err)
	}
}

// TestGreetingOverMaxFrameSizeIsRejected covers the same cap on the greeting,
// which is read before any command is sent.
func TestGreetingOverMaxFrameSizeIsRejected(t *testing.T) {
	cfg, _, cleanup := startDomainCreateServer(t, domainCreateResponse("example.in"))
	defer cleanup()
	cfg.Transport.MaxFrameSize = 16

	_, err := epp.Connect(cfg)
	if err == nil {
		t.Fatal("expected an oversized greeting to be rejected")
	}
	if !errors.Is(err, epp.ErrFrameTooLarge) {
		t.Fatalf("expected ErrFrameTooLarge, got %v", err)
	}
}

// startSilentServer accepts a connection, sends a greeting, records the first
// command and then never replies, which is how a transform outcome becomes
// ambiguous in production.
func startSilentServer(t *testing.T) (*config.Config, <-chan []byte, func()) {
	t.Helper()

	certFile, keyFile, cert := writeTestCertificate(t)

	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}

	requests := make(chan []byte, 1)
	done := make(chan struct{})
	release := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		if err := epp.WriteFrame(conn, []byte(domainCreateGreeting())); err != nil {
			return
		}
		requestXML, err := epp.ReadFrame(conn)
		if err != nil {
			return
		}
		requests <- requestXML
		<-release
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("expected TCP listener address, got %T", listener.Addr())
	}

	cfg := &config.Config{
		Server:  config.ServerConfig{Host: "127.0.0.1", Port: addr.Port},
		Timeout: config.TimeoutConfig{Connect: 5, Read: 1, Write: 5},
		TLS: config.TLSConfig{
			CertFile:           certFile,
			KeyFile:            keyFile,
			InsecureSkipVerify: true,
			AllowInsecure:      true,
		},
	}

	cleanup := func() {
		close(release)
		_ = listener.Close()
		<-done
	}

	return cfg, requests, cleanup
}

// asSDK reports whether err is an *epp.SDKError and stores it in target.
func asSDK(err error, target **epp.SDKError) bool {
	return errors.As(err, target)
}
