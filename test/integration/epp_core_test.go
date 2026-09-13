package test

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"crypto/tls"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
	"github.com/hariom-pal/go-epp/types"
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

func TestReadFrameRejectsOversizedLengthBeforeAllocation(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	done := make(chan error, 1)
	go func() {
		header := make([]byte, 4)
		binary.BigEndian.PutUint32(header, 1024)
		_, err := serverConn.Write(header)
		done <- err
	}()

	_, err := epp.ReadFrameWithMax(clientConn, 64)
	if !errors.Is(err, epp.ErrFrameTooLarge) {
		t.Fatalf("expected ErrFrameTooLarge, got %T: %v", err, err)
	}

	if err := <-done; err != nil {
		t.Fatalf("write oversized frame header failed: %v", err)
	}
}

func TestReadFrameRejectsTruncatedFrame(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go func() {
		defer serverConn.Close()
		header := make([]byte, 4)
		binary.BigEndian.PutUint32(header, 12)
		_, _ = serverConn.Write(append(header, []byte("short")...))
	}()

	_, err := epp.ReadFrameWithMax(clientConn, epp.DefaultMaxFrameSize)
	var sdkErr *epp.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindFraming {
		t.Fatalf("expected framing error, got %T: %v", err, err)
	}
}

func TestConnectRejectsInsecureSkipVerifyWithoutExplicitOptIn(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 700},
		TLS: config.TLSConfig{
			CertFile:           "client.crt",
			KeyFile:            "client.key",
			InsecureSkipVerify: true,
		},
	}

	_, err := epp.Connect(cfg)
	var sdkErr *epp.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindConfiguration {
		t.Fatalf("expected configuration error, got %T: %v", err, err)
	}
}

func TestNewClientDoesNotFallbackToPlainTCP(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer listener.Close()

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: listener.Addr().(*net.TCPAddr).Port,
		},
	}

	_, err = epp.NewClient(cfg)
	var sdkErr *epp.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindConfiguration {
		t.Fatalf("expected configuration error before TCP dial, got %T: %v", err, err)
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

func TestExecuteOnClosedSessionReturnsSessionError(t *testing.T) {
	cfg, _, cleanup := startSequentialEPPServer(t, nil)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	_, err = client.Execute([]byte("<epp/>"))
	var sdkErr *epp.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindSession {
		t.Fatalf("expected session error, got %T: %v", err, err)
	}
}

func TestExecuteContextCancellationInterruptsRead(t *testing.T) {
	cfg, requests, cleanup := startBlockingEPPServer(t)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err = client.ExecuteContext(ctx, []byte("<epp>cancel</epp>"))
	var sdkErr *epp.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindTimeout {
		t.Fatalf("expected timeout error, got %T: %v", err, err)
	}
	if got := readRequest(t, requests); got != "<epp>cancel</epp>" {
		t.Fatalf("unexpected request: %s", got)
	}
}

func TestTransformReadFailureReturnsAmbiguousOutcome(t *testing.T) {
	cfg, requests, cleanup := startDisconnectAfterRequestServer(t)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainCreate(types.DomainCreateRequest{
		Domain:      "example.in",
		Period:      1,
		Unit:        "y",
		Registrant:  "REG123",
		AuthInfo:    "secret",
		NameServers: []string{"ns1.example.in"},
	})
	var ambiguous *epp.AmbiguousTransformError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("expected ambiguous transform error, got %T: %v", err, err)
	}
	if got := readRequest(t, requests); !strings.Contains(got, "<domain:create") {
		t.Fatalf("expected domain create request, got: %s", got)
	}
}

func TestExecuteSerializesConcurrentCommands(t *testing.T) {
	cfg, requests, cleanup := startSequentialEPPServer(t, []string{
		simpleEPPResponse("ONE", "SERVER-ONE"),
		simpleEPPResponse("TWO", "SERVER-TWO"),
	})
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := client.Execute([]byte(fmt.Sprintf("<epp>%d</epp>", i)))
			errs <- err
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent execute failed: %v", err)
		}
	}

	_ = readRequest(t, requests)
	_ = readRequest(t, requests)
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
	if strings.Contains(loginXML, `<svcExtension>`) {
		t.Fatalf("default login should not advertise extensions blindly: %s", loginXML)
	}
	if strings.Contains(loginXML, constants.ContactNamespace) ||
		strings.Contains(loginXML, constants.HostNamespace) {

		t.Fatalf("default login should not advertise object URIs missing from greeting: %s", loginXML)
	}
	assertContains(t, loginXML, `<clTRID>LOGIN-`)

	logoutXML := readRequest(t, requests)
	assertContains(t, logoutXML, `<logout/>`)
	assertContains(t, logoutXML, `<clTRID>LOGOUT-`)
}

func TestLoginAdvertisesSelectedGreetingSupportedExtensions(t *testing.T) {
	cfg, requests, cleanup := startSequentialEPPServerWithGreeting(t, greetingWithExtensions(
		constants.SecDNSNamespace,
		constants.FeeNamespace,
	), []string{simpleEPPResponse("LOGIN-TEST", "SERVER-LOGIN")})
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{Username: "ote-user", Password: "ote-password"}
	cfg.Login.ExtensionURIs = []string{constants.SecDNSNamespace}

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	loginXML := readRequest(t, requests)
	assertContains(t, loginXML, `<svcExtension>`)
	assertContains(t, loginXML, `<extURI>`+constants.SecDNSNamespace+`</extURI>`)
	if strings.Contains(loginXML, constants.FeeNamespace) {
		t.Fatalf("unexpected unrequested fee extension in login XML: %s", loginXML)
	}
}

func TestLoginRejectsRequestedUnsupportedExtensionWhenStrict(t *testing.T) {
	cfg, _, cleanup := startSequentialEPPServer(t, nil)
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{Username: "ote-user", Password: "ote-password"}
	cfg.Login.ExtensionURIs = []string{constants.FeeNamespace}
	cfg.Login.RequireSupportedExtensions = true

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	err = client.Login()
	var sdkErr *epp.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Kind != epp.ErrorKindConfiguration {
		t.Fatalf("expected configuration error, got %T: %v", err, err)
	}
}

func TestLoggerReceivesRedactedEvents(t *testing.T) {
	cfg, _, cleanup := startSequentialEPPServer(t, []string{
		simpleEPPResponse("LOGIN-TEST", "SERVER-LOGIN"),
	})
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{Username: "ote-user", Password: "secret-password"}

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	logger := &eventRecorder{}
	client.SetLogger(logger)

	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	for _, event := range logger.events {
		if strings.Contains(fmt.Sprint(event), "secret-password") {
			t.Fatalf("event leaked password: %+v", event)
		}
	}
	if !slices.ContainsFunc(logger.events, func(event epp.Event) bool { return event.Type == epp.EventLogin }) {
		t.Fatalf("expected login event, got %+v", logger.events)
	}
}

func TestLoginEscapesCredentials(t *testing.T) {
	cfg, requests, cleanup := startSequentialEPPServer(t, []string{
		simpleEPPResponse("LOGIN-TEST", "SERVER-LOGIN"),
	})
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{
		Username: `ote-user&<"`,
		Password: `p&ss<word>`,
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if err := client.Login(); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	loginXML := readRequest(t, requests)
	assertContains(t, loginXML, `<clID>ote-user&amp;&lt;&#34;</clID>`)
	assertContains(t, loginXML, `<pw>p&amp;ss&lt;word&gt;</pw>`)
}

func TestLoginReturnsEPPError(t *testing.T) {
	cfg, _, cleanup := startSequentialEPPServer(t, []string{
		eppErrorResponse(constants.ResultAuthenticationError, "Authentication failed", "LOGIN-TEST", "SERVER-LOGIN"),
	})
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{
		Username: "ote-user",
		Password: "bad-password",
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	err = client.Login()
	assertEPPErrorCode(t, err, constants.ResultAuthenticationError)
}

func TestCommandResponseChecksAllResults(t *testing.T) {
	cfg, _, cleanup := startSequentialEPPServer(t, []string{
		multiResultResponse(),
	})
	defer cleanup()

	cfg.Authentication = config.AuthenticationConfig{
		Username: "ote-user",
		Password: "bad-password",
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	err = client.Login()
	var eppErr *epp.Error
	if !errors.As(err, &eppErr) {
		t.Fatalf("expected EPP error, got %T: %v", err, err)
	}
	if eppErr.Code != constants.ResultParameterError {
		t.Fatalf("expected second result error code, got %d", eppErr.Code)
	}
	if len(eppErr.Results) != 2 {
		t.Fatalf("expected both results preserved, got %+v", eppErr.Results)
	}
}

func TestEPPErrorHelpers(t *testing.T) {
	err := &epp.Error{
		Code:       constants.ResultObjectExists,
		Message:    "object exists",
		ClientTRID: "CLIENT-TRID",
		ServerTRID: "SERVER-TRID",
	}

	if !strings.Contains(err.Error(), "EPP result [2302]") {
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
	return startSequentialEPPServerWithGreeting(t, domainCreateGreeting(), responses)
}

func startSequentialEPPServerWithGreeting(
	t *testing.T,
	greeting string,
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

		if err := epp.WriteFrame(conn, []byte(greeting)); err != nil {
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
			AllowInsecure:      true,
		},
	}

	cleanup := func() {
		_ = listener.Close()
		<-done
	}

	return cfg, requests, cleanup
}

func startBlockingEPPServer(t *testing.T) (*config.Config, <-chan []byte, func()) {
	t.Helper()

	cfg, requests, cleanup := startSequentialEPPServerWithHandler(t, func(conn net.Conn, requests chan<- []byte) {
		requestXML, err := epp.ReadFrame(conn)
		if err == nil {
			requests <- requestXML
		}
		<-time.After(200 * time.Millisecond)
	})
	return cfg, requests, cleanup
}

func startDisconnectAfterRequestServer(t *testing.T) (*config.Config, <-chan []byte, func()) {
	t.Helper()

	return startSequentialEPPServerWithHandler(t, func(conn net.Conn, requests chan<- []byte) {
		requestXML, err := epp.ReadFrame(conn)
		if err == nil {
			requests <- requestXML
		}
	})
}

func startSequentialEPPServerWithHandler(
	t *testing.T,
	handler func(net.Conn, chan<- []byte),
) (*config.Config, chan []byte, func()) {
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
		handler(conn, requests)
	}()

	addr := listener.Addr().(*net.TCPAddr)
	cfg := &config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: addr.Port},
		TLS: config.TLSConfig{
			CertFile:           certFile,
			KeyFile:            keyFile,
			InsecureSkipVerify: true,
			AllowInsecure:      true,
		},
	}
	cleanup := func() {
		_ = listener.Close()
		<-done
	}
	return cfg, requests, cleanup
}

func greetingWithExtensions(extensions ...string) string {
	var builder strings.Builder
	builder.WriteString(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>test EPP server</svID><svDate>2026-06-30T09:30:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><svcExtension>`)
	for _, extension := range extensions {
		builder.WriteString(`<extURI>`)
		builder.WriteString(extension)
		builder.WriteString(`</extURI>`)
	}
	builder.WriteString(`</svcExtension></svcMenu></greeting></epp>`)
	return builder.String()
}

type eventRecorder struct {
	events []epp.Event
}

func (r *eventRecorder) EPPEvent(event epp.Event) {
	r.events = append(r.events, event)
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

func eppErrorResponse(
	code int,
	message string,
	clientTRID string,
	serverTRID string,
) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="` + fmt.Sprint(code) + `">
            <msg>` + message + `</msg>
        </result>
        <trID>
            <clTRID>` + clientTRID + `</clTRID>
            <svTRID>` + serverTRID + `</svTRID>
        </trID>
    </response>
</epp>`
}

func multiResultResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <result code="2005">
            <msg lang="en">Parameter value policy error</msg>
        </result>
        <trID>
            <clTRID>LOGIN-TEST</clTRID>
            <svTRID>SERVER-LOGIN</svTRID>
        </trID>
    </response>
</epp>`
}
