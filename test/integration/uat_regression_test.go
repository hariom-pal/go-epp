package test

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	idnext "github.com/hariom-pal/go-epp/extensions/idn"
	launchext "github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/types"
)

// TestLogoutResultCodeIsSuccess covers the NIXI OT&E finding that logout was
// reported as a failure. RFC 5730 requires result code 1500 in response to
// logout, and every 1xxx code is a positive completion reply.
func TestLogoutResultCodeIsSuccess(t *testing.T) {
	for _, code := range []int{1000, 1001, 1300, 1301, 1500} {
		if !constants.IsSuccessResultCode(code) {
			t.Fatalf("result code %d should be a success code", code)
		}
	}
	for _, code := range []int{999, 2000, 2200, 2303, 2306} {
		if constants.IsSuccessResultCode(code) {
			t.Fatalf("result code %d should not be a success code", code)
		}
	}
}

// TestResultValueDiagnosticsSurface covers the NIXI OT&E finding that the
// server's <value>/<extValue> explanation of a failure was discarded, leaving
// only a generic message such as "Parameter value policy error".
func TestResultValueDiagnosticsSurface(t *testing.T) {
	response := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response>
<result code="2306"><msg>Parameter value policy error</msg>
<value><msg>address 203.0.113.11 forbidden by rule for 203.0.113.0/24</msg></value>
</result>
<result code="2005"><msg>Parameter value syntax error</msg>
<extValue><value><host:addr xmlns:host="urn:ietf:params:xml:ns:host-1.0">bad</host:addr></value><reason>invalid IP address</reason></extValue>
</result>
<trID><clTRID>CREATE-TEST</clTRID><svTRID>XYZ</svTRID></trID></response></epp>`

	cfg, requests, cleanup := startDomainCreateServer(t, response)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "reg1",
		AuthInfo:   "secret123",
	})
	if err == nil {
		t.Fatal("expected domain create to fail")
	}
	readRequest(t, requests)

	var eppErr *epp.Error
	if !errors.As(err, &eppErr) {
		t.Fatalf("expected *epp.Error, got %T", err)
	}
	if len(eppErr.Values) != 1 {
		t.Fatalf("expected 1 diagnostic value, got %d", len(eppErr.Values))
	}
	if eppErr.Values[0].Reason != "address 203.0.113.11 forbidden by rule for 203.0.113.0/24" {
		t.Fatalf("unexpected diagnostic reason: %q", eppErr.Values[0].Reason)
	}
	if !strings.Contains(err.Error(), "forbidden by rule") {
		t.Fatalf("diagnostic missing from error string: %s", err.Error())
	}

	// The second result carries an extValue, whose offending fragment and
	// reason must both survive parsing.
	if len(eppErr.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(eppErr.Results))
	}
	second := eppErr.Results[1]
	if len(second.Values) != 1 {
		t.Fatalf("expected 1 extValue on second result, got %d", len(second.Values))
	}
	if second.Values[0].Reason != "invalid IP address" {
		t.Fatalf("unexpected extValue reason: %q", second.Values[0].Reason)
	}
	if !strings.Contains(second.Values[0].Value, "host:addr") {
		t.Fatalf("extValue fragment not captured: %q", second.Values[0].Value)
	}
}

// TestTransferQueryWithDefaultedUnit covers the NIXI OT&E finding that a
// transfer query failed client-side with "period must be between 1 and 99"
// whenever the caller defaulted the period unit without setting a value.
func TestTransferQueryWithDefaultedUnit(t *testing.T) {
	_, requestXML := runDomainTransferTest(t, types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferQuery,
		Unit:       "y",
	})

	assertContains(t, requestXML, `<transfer op="query">`)
	assertNotContains(t, requestXML, `<domain:period`)
}

// TestDomainCreateIDNExtension covers the NIXI OT&E finding that IDN
// registration was impossible: the registry rejects an IDN create without an
// idn-1.0 table tag, and the SDK had no way to send one.
func TestDomainCreateIDNExtension(t *testing.T) {
	const ascii = "xn--i1b0a6ag6lba6bvd.xn--h2brj9c"

	cfg, requests, cleanup := startDomainCreateServer(t, domainCreateResponse(ascii))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if _, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "जांचसीएससी.भारत",
		Period:     1,
		Unit:       "y",
		Registrant: "reg1",
		AuthInfo:   "secret123",
		IDN: &idnext.CreateRequest{
			Data: idnext.Data{Table: "hi", UName: "जांचसीएससी.भारत"},
		},
	}); err != nil {
		t.Fatalf("IDN domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertContains(t, requestXML, `xmlns:idn="urn:ietf:params:xml:ns:idn-1.0"`)
	assertContains(t, requestXML, `<idn:table>hi</idn:table>`)
	assertContains(t, requestXML, `<idn:uname>जांचसीएससी.भारत</idn:uname>`)
	assertContains(t, requestXML, `<domain:name>`+ascii+`</domain:name>`)
}

// TestAmbiguousTransformCarriesClientTRID covers the finding that an ambiguous
// transform reported no transaction identifier. The SDK directs callers to
// reconcile a lost transform against the registry's transaction record, which
// is impossible without the identifier the command was sent with.
func TestAmbiguousTransformCarriesClientTRID(t *testing.T) {
	cfg, requests, cleanup := startSilentServer(t)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()
	cfg.Timeout.Read = 1

	_, err = client.DomainCreate(types.DomainCreateRequest{
		Domain: "example.in", Period: 1, Unit: "y",
		Registrant: "reg1", AuthInfo: "secret123",
	})
	if err == nil {
		t.Fatal("expected the transform to fail with no response")
	}

	var ambiguous *epp.AmbiguousTransformError
	if !errors.As(err, &ambiguous) {
		t.Fatalf("expected *epp.AmbiguousTransformError, got %T: %v", err, err)
	}
	if ambiguous.ClientTRID == "" {
		t.Fatal("an ambiguous transform must report the clTRID it was sent with")
	}

	sent := string(<-requests)
	if !strings.Contains(sent, "<clTRID>"+ambiguous.ClientTRID+"</clTRID>") {
		t.Fatalf("reported clTRID %q does not match the command that was sent", ambiguous.ClientTRID)
	}
	if !strings.Contains(err.Error(), ambiguous.ClientTRID) {
		t.Fatalf("the error string omits the clTRID needed for reconciliation: %s", err.Error())
	}
}

// TestCancellationWatcherIsJoined covers the cross-command hazard in the
// per-command cancellation watcher. The watcher aborts a blocked read by
// pushing the connection's read deadline into the past; because it selects on
// both ctx.Done() and command completion, and Go picks randomly between ready
// cases, a watcher left running after its command returned could clobber the
// deadline of the next command on the same session. Joining it before the
// command returns is what makes that impossible, and the observable invariant
// is that no goroutine survives the call.
func TestCancellationWatcherIsJoined(t *testing.T) {
	cfg, requests, cleanup := startDomainCreateServer(t, domainCreateResponse("example.in"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	// Warm up so connection-related goroutines already exist.
	_, _ = client.DomainCheck(types.DomainCheckRequest{Domains: []string{"example.in"}})
	runtime.Gosched()
	baseline := runtime.NumGoroutine()

	for i := 0; i < 50; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		_, _ = client.DomainCheck(types.DomainCheckRequest{Domains: []string{"example.in"}})
		_, _ = client.DomainCheckContext(ctx, types.DomainCheckRequest{Domains: []string{"example.in"}})
		cancel()

		// No settling: the watcher must already be finished when the command
		// returns, so the count cannot have grown.
		if current := runtime.NumGoroutine(); current > baseline {
			t.Fatalf("iteration %d: %d goroutines outlived the command (baseline %d)", i, current, baseline)
		}
	}
	drain(requests)
}

// drain empties the recorded request channel so the test server can exit.
func drain(requests <-chan []byte) {
	for {
		select {
		case <-requests:
		default:
			return
		}
	}
}

// TestIDNInfoDataElementSpellings covers the NIXI OT&E finding that domain
// info silently dropped the IDN table and U-label: the SDK parsed only
// <idn:infData> while the registry sends <idn:data>. Both spellings are in use,
// so both must be accepted.
func TestIDNInfoDataElementSpellings(t *testing.T) {
	for _, element := range []string{"data", "infData"} {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response>
<result code="1000"><msg>Command completed successfully</msg></result>
<resData><domain:infData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
<domain:name>xn--11bbsu6b1ak.xn--h2brj9c</domain:name>
<domain:roid>DO_TEST-NIXI</domain:roid>
<domain:registrant>reg1</domain:registrant>
</domain:infData></resData>
<extension><idn:` + element + ` xmlns:idn="urn:ietf:params:xml:ns:idn-1.0">
<idn:table>hin_deva</idn:table><idn:uname>परखजडमक.भारत</idn:uname>
</idn:` + element + `></extension>
<trID><clTRID>INFO-TEST</clTRID><svTRID>SV-1</svTRID></trID></response></epp>`

		cfg, requests, cleanup := startDomainCreateServer(t, response)

		client, err := epp.Connect(cfg)
		if err != nil {
			cleanup()
			t.Fatalf("connect failed: %v", err)
		}

		info, err := client.DomainInfo(types.DomainInfoRequest{Domain: "xn--11bbsu6b1ak.xn--h2brj9c"})
		readRequest(t, requests)
		client.Close()
		cleanup()

		if err != nil {
			t.Fatalf("<idn:%s>: domain info failed: %v", element, err)
		}
		if info.Result.IDN.Table != "hin_deva" {
			t.Fatalf("<idn:%s>: expected table hin_deva, got %q", element, info.Result.IDN.Table)
		}
		if info.Result.IDN.UName != "परखजडमक.भारत" {
			t.Fatalf("<idn:%s>: expected the U-label, got %q", element, info.Result.IDN.UName)
		}
	}
}

// TestMessageQueueOnOrdinaryResponse covers RFC 5730 section 2.9.2.3: a server
// may report queued messages with <msgQ> on any response, not only a poll. The
// SDK read it only on poll responses, so the signal was lost everywhere else.
func TestMessageQueueOnOrdinaryResponse(t *testing.T) {
	response := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response>
<result code="1000"><msg>Command completed successfully</msg></result>
<msgQ count="2" id="4455"/>
<resData><domain:infData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
<domain:name>example.in</domain:name><domain:roid>DO_TEST-NIXI</domain:roid>
</domain:infData></resData>
<trID><clTRID>INFO-TEST</clTRID><svTRID>SV-1</svTRID></trID></response></epp>`

	cfg, requests, cleanup := startDomainCreateServer(t, response)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	info, err := client.DomainInfo(types.DomainInfoRequest{Domain: "example.in"})
	if err != nil {
		t.Fatalf("domain info failed: %v", err)
	}
	readRequest(t, requests)

	if info.MessageQueue.Count != 2 {
		t.Fatalf("expected a queue count of 2, got %d", info.MessageQueue.Count)
	}
	if info.MessageQueue.ID != "4455" {
		t.Fatalf("expected message id 4455, got %q", info.MessageQueue.ID)
	}
}

// TestClientSideValidationIsNotAnEPPResult covers the finding that every
// locally rejected request was reported with Kind "epp_result" and an EPP
// result code, making it indistinguishable from a registry rejection. A
// registrar keying retry, alerting or reconciliation on "the registry answered"
// would misattribute a request that never left the process.
func TestClientSideValidationIsNotAnEPPResult(t *testing.T) {
	client := &epp.Client{}

	cases := []struct {
		name string
		call func() error
	}{
		{"domain create without a domain", func() error {
			_, err := client.DomainCreate(types.DomainCreateRequest{})
			return err
		}},
		{"domain transfer without an operation", func() error {
			_, err := client.DomainTransfer(types.DomainTransferRequest{DomainName: "example.in"})
			return err
		}},
		{"poll ack without a message ID", func() error {
			_, err := client.Poll(types.PollRequest{Operation: constants.PollAcknowledge})
			return err
		}},
		{"contact create without an ID", func() error {
			_, err := client.ContactCreate(types.ContactCreateRequest{})
			return err
		}},
		{"host create without a name", func() error {
			_, err := client.HostCreate(types.HostCreateRequest{})
			return err
		}},
		{"domain check with no names", func() error {
			_, err := client.DomainCheck(types.DomainCheckRequest{})
			return err
		}},
	}

	for _, testCase := range cases {
		err := testCase.call()
		if err == nil {
			t.Fatalf("%s: expected rejection", testCase.name)
		}

		var eppErr *epp.Error
		if !errors.As(err, &eppErr) {
			t.Fatalf("%s: expected *epp.Error, got %T", testCase.name, err)
		}
		if eppErr.Kind != epp.ErrorKindValidation {
			t.Fatalf("%s: expected kind %q, got %q", testCase.name, epp.ErrorKindValidation, eppErr.Kind)
		}
		if !eppErr.IsValidationError() {
			t.Fatalf("%s: IsValidationError returned false", testCase.name)
		}
		if eppErr.ServerTRID != "" || eppErr.ClientTRID != "" {
			t.Fatalf("%s: a local rejection must carry no TRIDs", testCase.name)
		}
		if len(eppErr.Results) != 0 {
			t.Fatalf("%s: a local rejection must carry no server results", testCase.name)
		}
		if strings.HasPrefix(err.Error(), string(epp.ErrorKindEPPResult)) {
			t.Fatalf("%s: error string claims to be a registry result: %s", testCase.name, err.Error())
		}
	}
}

// TestServerResultKeepsEPPResultKind confirms a genuine registry rejection is
// still classified as an EPP result, so the two remain distinguishable.
func TestServerResultKeepsEPPResultKind(t *testing.T) {
	response := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response>
<result code="2303"><msg>Object does not exist</msg></result>
<trID><clTRID>CREATE-TEST</clTRID><svTRID>SV-1</svTRID></trID></response></epp>`

	cfg, requests, cleanup := startDomainCreateServer(t, response)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	_, err = client.DomainCreate(types.DomainCreateRequest{
		Domain: "example.in", Period: 1, Unit: "y",
		Registrant: "reg1", AuthInfo: "secret123",
	})
	readRequest(t, requests)

	var eppErr *epp.Error
	if !errors.As(err, &eppErr) {
		t.Fatalf("expected *epp.Error, got %T", err)
	}
	if eppErr.Kind != epp.ErrorKindEPPResult {
		t.Fatalf("expected kind %q, got %q", epp.ErrorKindEPPResult, eppErr.Kind)
	}
	if eppErr.IsValidationError() {
		t.Fatal("a registry rejection must not be classified as client-side validation")
	}
	if eppErr.ServerTRID != "SV-1" {
		t.Fatalf("expected the server TRID to survive, got %q", eppErr.ServerTRID)
	}
}

// TestLaunchRawXMLIsNotWrapped covers the NIXI OT&E finding that the launch
// extension emitted a bogus <ChoiceXML> element in the EPP namespace instead
// of the caller's signed mark, because encoding/xml honours the ",innerxml"
// flag on marshal only for string and []byte fields. Any other type is
// silently emitted as an ordinary element named after the Go field, so signed
// marks were both mis-wrapped and dropped, making sunrise registration
// impossible at any registry. See RFC 8334 section 3.1.
func TestLaunchRawXMLIsNotWrapped(t *testing.T) {
	const signedMark = `<smd:encodedSignedMark xmlns:smd="urn:ietf:params:xml:ns:signedMark-1.0">BASE64DATA</smd:encodedSignedMark>`

	cfg, requests, cleanup := startDomainCreateServer(t, domainCreateResponse("example.in"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if _, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "reg1",
		AuthInfo:   "secret123",
		Launch: &launchext.CreateRequest{
			Phase:                launchext.Phase{Value: launchext.PhaseSunrise},
			EncodedSignedMarkXML: []string{signedMark},
		},
	}); err != nil {
		t.Fatalf("launch domain create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertNotContains(t, requestXML, "ChoiceXML")
	assertContains(t, requestXML, signedMark)
	assertContains(t, requestXML, `<launch:phase>sunrise</launch:phase>`)

	// RFC 8334 orders the launch:create children as phase, then the mark
	// choice, then any notices.
	phase := strings.Index(requestXML, "<launch:phase>")
	mark := strings.Index(requestXML, "encodedSignedMark")
	if phase < 0 || mark < 0 || phase > mark {
		t.Fatalf("launch children are out of schema order: phase=%d mark=%d", phase, mark)
	}
}

// TestLaunchCodeMarkRawXML covers the same defect on the codeMark path.
func TestLaunchCodeMarkRawXML(t *testing.T) {
	const markXML = `<mark:mark xmlns:mark="urn:ietf:params:xml:ns:mark-1.0">MARKBODY</mark:mark>`

	cfg, requests, cleanup := startDomainCreateServer(t, domainCreateResponse("example.in"))
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	if _, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "reg1",
		AuthInfo:   "secret123",
		Launch: &launchext.CreateRequest{
			Phase: launchext.Phase{Value: launchext.PhaseSunrise},
			CodeMarks: []launchext.CodeMark{{
				Code:    &launchext.Code{Value: "49FD46E6C4B45C55", ValidatorID: "tmch"},
				MarkXML: markXML,
			}},
		},
	}); err != nil {
		t.Fatalf("launch codeMark create failed: %v", err)
	}

	requestXML := readRequest(t, requests)

	assertNotContains(t, requestXML, "MarkXML")
	assertContains(t, requestXML, markXML)
	assertContains(t, requestXML, `<launch:code validatorID="tmch">49FD46E6C4B45C55</launch:code>`)
}

// TestDomainCreateIDNRequiresTable keeps the extension from being emitted
// without the one field a registry cannot derive from the A-label.
func TestDomainCreateIDNRequiresTable(t *testing.T) {
	if idnext.ValidCreate(&idnext.CreateRequest{Data: idnext.Data{UName: "जांच.भारत"}}) {
		t.Fatal("IDN create without a table should be rejected")
	}
	if !idnext.ValidCreate(nil) {
		t.Fatal("absent IDN extension should be valid")
	}
	if idnext.NewCreate(nil) != nil {
		t.Fatal("absent IDN extension should produce no XML")
	}
}
