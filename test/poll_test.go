package test

import (
	"testing"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func TestPollRequestXMLAndParsing(t *testing.T) {
	resp, requestXML := runPollTest(t, types.PollRequest{
		Operation: constants.PollRequest,
		MessageID: "ignored-message-id",
	}, pollMessageResponse())

	assertContains(t, requestXML, `<poll op="req"></poll>`)
	assertNotContains(t, requestXML, `msgID=`)
	assertContains(t, requestXML, `<clTRID>POLL-`)

	assertPollMessageResponse(t, resp)
}

func TestPollAckXMLAndParsing(t *testing.T) {
	resp, requestXML := runPollTest(t, types.PollRequest{
		Operation: constants.PollAcknowledge,
		MessageID: "12345",
	}, pollAckResponse())

	assertContains(t, requestXML, `<poll op="ack" msgID="12345"></poll>`)
	assertContains(t, requestXML, `<clTRID>POLL-`)

	if resp.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.MessageQueue.Count != 0 {
		t.Fatalf("unexpected queue count: %d", resp.MessageQueue.Count)
	}
}

func TestPollEmptyQueueParsing(t *testing.T) {
	resp, requestXML := runPollTest(t, types.PollRequest{
		Operation: constants.PollRequest,
	}, pollEmptyQueueResponse())

	assertContains(t, requestXML, `<poll op="req"></poll>`)

	if resp.ResultCode != constants.ResultNoMessages {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.MessageQueue.Count != 0 {
		t.Fatalf("unexpected queue count: %d", resp.MessageQueue.Count)
	}

	if resp.MessageQueue.ID != "" ||
		resp.MessageQueue.Date != nil ||
		resp.Message != "" {

		t.Fatalf("expected empty queue data, got %#v", resp)
	}
}

func TestPollAckRequiresMessageID(t *testing.T) {
	client := &epp.Client{}

	_, err := client.Poll(types.PollRequest{
		Operation: constants.PollAcknowledge,
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func TestPollInvalidOperation(t *testing.T) {
	client := &epp.Client{}

	_, err := client.Poll(types.PollRequest{
		Operation: "invalid",
	})
	assertEPPErrorCode(t, err, constants.ResultParameterError)
}

func runPollTest(
	t *testing.T,
	req types.PollRequest,
	responseXML string,
) (*types.PollResponse, string) {
	t.Helper()

	cfg, requests, cleanup := startDomainCreateServer(t, responseXML)
	defer cleanup()

	client, err := epp.Connect(cfg)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer client.Close()

	resp, err := client.Poll(req)
	if err != nil {
		t.Fatalf("poll failed: %v", err)
	}

	return resp, readRequest(t, requests)
}

func assertPollMessageResponse(
	t *testing.T,
	resp *types.PollResponse,
) {
	t.Helper()

	if resp.ResultCode != constants.ResultAckToDequeue {
		t.Fatalf("unexpected result code: %d", resp.ResultCode)
	}

	if resp.ResultMsg != "Command completed successfully; ack to dequeue" {
		t.Fatalf("unexpected result message: %s", resp.ResultMsg)
	}

	if resp.ResultMessage != resp.ResultMsg {
		t.Fatalf("unexpected result message alias: %s", resp.ResultMessage)
	}

	if resp.ClientTRID != "POLL-TEST" {
		t.Fatalf("unexpected client TRID: %s", resp.ClientTRID)
	}

	if resp.ServerTRID != "SERVER-POLL" {
		t.Fatalf("unexpected server TRID: %s", resp.ServerTRID)
	}

	if resp.MessageQueue.Count != 2 {
		t.Fatalf("unexpected queue count: %d", resp.MessageQueue.Count)
	}

	if resp.MessageQueue.ID != "12345" {
		t.Fatalf("unexpected message ID: %s", resp.MessageQueue.ID)
	}

	expectedDate := time.Date(2026, 7, 1, 10, 30, 0, 0, time.UTC)
	if resp.MessageQueue.Date == nil ||
		!resp.MessageQueue.Date.Equal(expectedDate) {

		t.Fatalf("unexpected queue date: %#v", resp.MessageQueue.Date)
	}

	if resp.Message != "Transfer requested for example.in" ||
		resp.MessageQueue.Message != resp.Message {

		t.Fatalf("unexpected poll message: %#v", resp.MessageQueue)
	}
}

func pollMessageResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1301">
            <msg>Command completed successfully; ack to dequeue</msg>
        </result>
        <msgQ count="2" id="12345">
            <qDate>2026-07-01T10:30:00Z</qDate>
            <msg>Transfer requested for example.in</msg>
        </msgQ>
        <trID>
            <clTRID>POLL-TEST</clTRID>
            <svTRID>SERVER-POLL</svTRID>
        </trID>
    </response>
</epp>`
}

func pollAckResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1000">
            <msg>Command completed successfully</msg>
        </result>
        <msgQ count="0" id="12345"/>
        <trID>
            <clTRID>POLL-ACK-TEST</clTRID>
            <svTRID>SERVER-POLL-ACK</svTRID>
        </trID>
    </response>
</epp>`
}

func pollEmptyQueueResponse() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <response>
        <result code="1300">
            <msg>Command completed successfully; no messages</msg>
        </result>
        <trID>
            <clTRID>POLL-EMPTY-TEST</clTRID>
            <svTRID>SERVER-POLL-EMPTY</svTRID>
        </trID>
    </response>
</epp>`
}
