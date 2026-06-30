package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func TestDomainCreateXMLGenerationAndParsing(t *testing.T) {
	responseXML := domainCreateResponse("example.in")

	cfg, requests, cleanup := startDomainCreateServer(t, responseXML)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "ABC123",
		Contacts: []types.DomainContact{
			{Type: "custom", ID: "CNT000"},
		},
		AdminContacts:   []string{"CNT001", "CNT004"},
		TechContacts:    []string{"CNT002"},
		BillingContacts: []string{"CNT003"},
		NameServers:     []string{"ns1.example.in"},
		NameServerInfo: []types.DomainNameServer{
			{
				HostName: "ns2.example.in",
				Addresses: []types.DomainHostAddress{
					{Version: "v4", IP: "1.2.3.4"},
				},
			},
		},
		AuthInfo: "mySecret",
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `<domain:create>`)
	assertContains(t, requestXML, `<domain:name>example.in</domain:name>`)
	assertContains(t, requestXML, `<domain:period unit="y">1</domain:period>`)
	assertContains(t, requestXML, `<domain:registrant>ABC123</domain:registrant>`)
	assertContains(t, requestXML, `<domain:contact type="custom">CNT000</domain:contact>`)
	assertContains(t, requestXML, `<domain:contact type="admin">CNT001</domain:contact>`)
	assertContains(t, requestXML, `<domain:contact type="admin">CNT004</domain:contact>`)
	assertContains(t, requestXML, `<domain:contact type="tech">CNT002</domain:contact>`)
	assertContains(t, requestXML, `<domain:contact type="billing">CNT003</domain:contact>`)
	assertContains(t, requestXML, `<domain:hostObj>ns1.example.in</domain:hostObj>`)
	assertContains(t, requestXML, `<domain:hostAttr>`)
	assertContains(t, requestXML, `<domain:hostName>ns2.example.in</domain:hostName>`)
	assertContains(t, requestXML, `<domain:hostAddr ip="v4">1.2.3.4</domain:hostAddr>`)
	assertContains(t, requestXML, `<domain:authInfo>`)
	assertContains(t, requestXML, `<domain:pw>mySecret</domain:pw>`)
	assertContains(t, requestXML, `<clTRID>CREATE-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ClientTRID != "CREATE-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-TEST" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.Result.Domain != "example.in" {
		t.Fatalf("unexpected domain: %s", resp.Result.Domain)
	}

	expectedCreated := time.Date(2026, 6, 30, 9, 30, 0, 0, time.UTC)
	if !resp.Result.CreatedDate.Equal(expectedCreated) {
		t.Fatalf("unexpected created date: %s", resp.Result.CreatedDate.Format(time.RFC3339))
	}

	expectedExpiry := time.Date(2027, 6, 30, 9, 30, 0, 0, time.UTC)
	if !resp.Result.ExpiryDate.Equal(expectedExpiry) {
		t.Fatalf("unexpected expiry date: %s", resp.Result.ExpiryDate.Format(time.RFC3339))
	}
}

func TestDomainCreateUnicodeConversion(t *testing.T) {
	unicodeDomain := "भारत.भारत"
	asciiDomain, err := idn.ToASCII(unicodeDomain)
	if err != nil {
		t.Fatalf("failed to prepare expected punycode: %v", err)
	}

	cfg, requests, cleanup := startDomainCreateServer(t, domainCreateResponse(asciiDomain))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     unicodeDomain,
		Period:     1,
		Unit:       "y",
		Registrant: "ABC123",
		AuthInfo:   "mySecret",
	})
	if err != nil {
		t.Fatalf("domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, "<domain:name>"+asciiDomain+"</domain:name>")

	if resp.Result.Domain != unicodeDomain {
		t.Fatalf("unexpected unicode domain: %s", resp.Result.Domain)
	}
}

func TestDomainCreateEmptyRequestValidation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainCreate(types.DomainCreateRequest{})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainCreateInvalidPeriod(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     0,
		Unit:       "y",
		Registrant: "ABC123",
		AuthInfo:   "mySecret",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestDomainCreateInvalidContacts(t *testing.T) {
	client := &epp.Client{}

	_, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "ABC123",
		Contacts: []types.DomainContact{
			{ID: "CNT001"},
		},
		AuthInfo: "mySecret",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func startDomainCreateServer(
	t *testing.T,
	responseXML string,
) (*config.Config, <-chan []byte, func()) {
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

		_ = epp.WriteFrame(conn, []byte(responseXML))
	}()

	addr := listener.Addr().(*net.TCPAddr)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: addr.Port,
		},
		TLS: config.TLSConfig{
			CertFile:           certFile,
			KeyFile:            keyFile,
			InsecureSkipVerify: true,
		},
	}

	cleanup := func() {
		_ = listener.Close()
		<-done
	}

	return cfg, requests, cleanup
}

func writeTestCertificate(t *testing.T) (string, string, tls.Certificate) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key failed: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "127.0.0.1",
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}

	der, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		t.Fatalf("create certificate failed: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	dir := t.TempDir()
	certFile := dir + "/cert.pem"
	keyFile := dir + "/key.pem"

	if err := os.WriteFile(certFile, certPEM, 0600); err != nil {
		t.Fatalf("write cert failed: %v", err)
	}

	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		t.Fatalf("write key failed: %v", err)
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("load key pair failed: %v", err)
	}

	return certFile, keyFile, cert
}

func readRequest(t *testing.T, requests <-chan []byte) string {
	t.Helper()

	select {
	case request := <-requests:
		return string(request)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for create request")
	}

	return ""
}

func assertContains(t *testing.T, value string, expected string) {
	t.Helper()

	if !strings.Contains(value, expected) {
		t.Fatalf("expected XML to contain %q\nXML:\n%s", expected, value)
	}
}

func assertEPPErrorCode(t *testing.T, err error, code int) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error code %d, got nil", code)
	}

	var eppErr *epp.Error
	if !errors.As(err, &eppErr) {
		t.Fatalf("expected EPP error, got %T: %v", err, err)
	}

	if eppErr.Code != code {
		t.Fatalf("expected error code %d, got %d", code, eppErr.Code)
	}
}

func domainCreateGreeting() string {
	return `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>test EPP server</svID><svDate>2026-06-30T09:30:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI></svcMenu></greeting></epp>`
}

func domainCreateResponse(domain string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <resData>
            <domain:creData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
                <domain:name>` + domain + `</domain:name>
                <domain:crDate>2026-06-30T09:30:00Z</domain:crDate>
                <domain:exDate>2027-06-30T09:30:00Z</domain:exDate>
            </domain:creData>
        </resData>
        <trID>
            <clTRID>CREATE-TEST</clTRID>
            <svTRID>SERVER-TEST</svTRID>
        </trID>
    </response>
</epp>`
}
