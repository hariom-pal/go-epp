package test

import (
	"encoding/binary"
	"net"
	"strings"
	"testing"

	"crypto/tls"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
)

func TestFrameReadWrite(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	payload := []byte("<epp>frame-test</epp>")
	done := make(chan error, 1)

	go func() {
		done <- epp.WriteFrame(serverConn, payload)
	}()

	got, err := epp.ReadFrame(clientConn)
	if err != nil {
		t.Fatalf("read frame failed: %v", err)
	}

	if err := <-done; err != nil {
		t.Fatalf("write frame failed: %v", err)
	}

	if string(got) != string(payload) {
		t.Fatalf("unexpected frame payload: %s", string(got))
	}
}

func TestReadFrameRejectsInvalidLength(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	done := make(chan error, 1)
	go func() {
		header := make([]byte, 4)
		binary.BigEndian.PutUint32(header, 3)
		_, err := serverConn.Write(header)
		done <- err
	}()

	_, err := epp.ReadFrame(clientConn)
	if err == nil {
		t.Fatal("expected invalid frame length error")
	}

	if err := <-done; err != nil {
		t.Fatalf("write invalid frame failed: %v", err)
	}
}

func TestConnectGreetingAndExecute(t *testing.T) {
	responseXML := simpleEPPResponse("EXEC-TEST", "SERVER-EXEC")
	cfg, requests, cleanup := startSequentialEPPServer(t, []string{responseXML})
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if string(client.Greeting()) != domainCreateGreeting() {
		t.Fatalf("unexpected greeting: %s", string(client.Greeting()))
	}

	response, err := client.Execute([]byte("<epp>execute-test</epp>"))
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}

	if string(response) != responseXML {
		t.Fatalf("unexpected execute response: %s", string(response))
	}

	requestXML := readRequest(t, requests)
	if requestXML != "<epp>execute-test</epp>" {
		t.Fatalf("unexpected execute request: %s", requestXML)
	}
}

func TestLoginAndLogout(t *testing.T) {
	cfg, requests, cleanup := startSequentialEPPServer(t, []string{
		simpleEPPResponse("LOGIN-TEST", "SERVER-LOGIN"),
		simpleEPPResponse("LOGOUT-TEST", "SERVER-LOGOUT"),
	})
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{
		Username: "ote-user",
		Password: "ote-password",
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if err := client.Logout(); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	loginXML := readRequest(t, requests)
	assertContains(t, loginXML, `<login>`)
	assertContains(t, loginXML, `<clID>ote-user</clID>`)
	assertContains(t, loginXML, `<pw>ote-password</pw>`)
	assertContains(t, loginXML, `<clTRID>LOGIN-`)

	logoutXML := readRequest(t, requests)
	assertContains(t, logoutXML, `<logout/>`)
	assertContains(t, logoutXML, `<clTRID>LOGOUT-`)
}

func TestEPPErrorHelpers(t *testing.T) {
	err := &epp.Error{
		Code:       constants.ResultObjectExists,
		Message:    "object exists",
		ClientTRID: "CLIENT-TRID",
		ServerTRID: "SERVER-TRID",
	}

	if !strings.Contains(err.Error(), "EPP Error [2302]") {
		t.Fatalf("unexpected error string: %s", err.Error())
	}

	if err.IsSuccess() {
		t.Fatal("object exists should not be success")
	}

	if !err.IsObjectExists() {
		t.Fatal("expected object exists helper to match")
	}

	if !(&epp.Error{Code: constants.ResultSuccess}).IsSuccess() {
		t.Fatal("expected success helper to match result 1000")
	}

	if !(&epp.Error{Code: constants.ResultObjectDoesNotExist}).IsObjectNotFound() {
		t.Fatal("expected object not found helper to match")
	}

	if !(&epp.Error{Code: constants.ResultAuthenticationError}).IsAuthenticationError() {
		t.Fatal("expected authentication helper to match")
	}

	if !(&epp.Error{Code: constants.ResultAuthorizationError}).IsAuthorizationError() {
		t.Fatal("expected authorization helper to match")
	}

	if !(&epp.Error{Code: constants.ResultObjectStatusProhibits}).IsObjectStatusProhibited() {
		t.Fatal("expected status prohibited helper to match")
	}
}

func startSequentialEPPServer(
	t *testing.T,
	responses []string,
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

	requests := make(chan []byte, len(responses))
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

		for _, response := range responses {
			requestXML, err := epp.ReadFrame(conn)
			if err != nil {
				return
			}

			requests <- requestXML

			if err := epp.WriteFrame(conn, []byte(response)); err != nil {
				return
			}
		}
	}()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("expected TCP listener address, got %T", listener.Addr())
	}

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

func simpleEPPResponse(
	clientTRID string,
	serverTRID string,
) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <trID>
            <clTRID>` + clientTRID + `</clTRID>
            <svTRID>` + serverTRID + `</svTRID>
        </trID>
    </response>
</epp>`
}
