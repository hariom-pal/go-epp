//go:build live

package live

import (
	"testing"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// TestPollRequestEmptyQueue covers RFC 5730 section 2.9.2.3: an empty message
// queue is reported with result code 1300, which is a success.
func TestPollRequestEmptyQueue(t *testing.T) {
	client := shared(t)

	resp, err := client.Poll(types.PollRequest{Operation: constants.PollRequest})
	if err != nil {
		t.Fatalf("poll request failed: %v", err)
	}
	switch resp.ResultCode {
	case constants.ResultNoMessages:
		if resp.MessageQueue.Count != 0 {
			t.Fatalf("1300 reported but queue count is %d", resp.MessageQueue.Count)
		}
		if resp.MessageQueue.ID != "" {
			t.Fatalf("1300 reported but a message ID was returned: %q", resp.MessageQueue.ID)
		}
	case constants.ResultAckToDequeue:
		if resp.MessageQueue.ID == "" {
			t.Fatal("1301 reported but no message ID was returned")
		}
		t.Logf("queue has %d messages, head id=%s", resp.MessageQueue.Count, resp.MessageQueue.ID)
	default:
		t.Fatalf("unexpected poll result code %d", resp.ResultCode)
	}
}

// TestPollUnsetOperation confirms an unset operation is refused locally rather
// than defaulting. RFC 5730 makes the op attribute mandatory, and defaulting it
// could silently acknowledge and destroy a queued message.
func TestPollUnsetOperation(t *testing.T) {
	client := shared(t)

	_, err := client.Poll(types.PollRequest{})
	if err == nil {
		t.Fatal("expected an unset poll operation to be rejected")
	}
	eppErr := asEPPError(t, err, "poll with unset operation")
	if !eppErr.IsValidationError() {
		t.Fatalf("a locally rejected request must not be reported as a registry result: kind=%q", eppErr.Kind)
	}
	if eppErr.ServerTRID != "" {
		t.Fatalf("a locally rejected request must carry no svTRID, got %q", eppErr.ServerTRID)
	}
}

// TestPollAckWithoutMessageID covers SDK-side validation: RFC 5730 requires
// msgID on an acknowledgement, so the SDK must refuse without a round trip.
func TestPollAckWithoutMessageID(t *testing.T) {
	client := shared(t)

	_, err := client.Poll(types.PollRequest{Operation: constants.PollAcknowledge})
	if err == nil {
		t.Fatal("expected a poll ack without a message ID to be rejected")
	}
	t.Logf("poll ack without msgID rejected: %v", err)
}

// TestPollAckUnknownMessageID covers acknowledging a message that is not in
// the queue.
func TestPollAckUnknownMessageID(t *testing.T) {
	client := shared(t)

	_, err := client.Poll(types.PollRequest{
		Operation: constants.PollAcknowledge,
		MessageID: "999999999",
	})
	if err == nil {
		t.Fatal("expected an ack for an unknown message ID to be refused")
	}
	eppErr := asEPPError(t, err, "poll ack with unknown message ID")
	if eppErr.Code < 2000 {
		t.Fatalf("expected an error result, got %d", eppErr.Code)
	}
	t.Logf("unknown msgID refused with %d: %s", eppErr.Code, eppErr.Message)
}

// TestPollAckMalformedMessageID covers a non-numeric message identifier.
func TestPollAckMalformedMessageID(t *testing.T) {
	client := shared(t)

	_, err := client.Poll(types.PollRequest{
		Operation: constants.PollAcknowledge,
		MessageID: "not-a-message-id",
	})
	if err == nil {
		t.Fatal("expected a malformed message ID to be refused")
	}
	eppErr := asEPPError(t, err, "poll ack with malformed message ID")
	t.Logf("malformed msgID refused with %d: %s", eppErr.Code, eppErr.Message)
}

// TestPollInvalidOperation covers SDK-side operation validation.
func TestPollInvalidOperation(t *testing.T) {
	client := shared(t)

	_, err := client.Poll(types.PollRequest{Operation: "dequeue-everything"})
	if err == nil {
		t.Fatal("expected an invalid poll operation to be rejected")
	}
}

// TestPollDrainAndAck exercises the full RFC 5730 poll loop when the queue has
// messages: read the head, acknowledge it, and confirm the count falls.
func TestPollDrainAndAck(t *testing.T) {
	client := shared(t)

	resp, err := client.Poll(types.PollRequest{Operation: constants.PollRequest})
	if err != nil {
		t.Fatalf("poll request failed: %v", err)
	}
	if resp.ResultCode == constants.ResultNoMessages {
		t.Skip("poll queue is empty; acknowledgement needs a counterparty action such as an inbound transfer")
	}

	before := resp.MessageQueue.Count
	messageID := resp.MessageQueue.ID

	ack, err := client.Poll(types.PollRequest{
		Operation: constants.PollAcknowledge,
		MessageID: messageID,
	})
	if err != nil {
		t.Fatalf("poll ack for %s failed: %v", messageID, err)
	}
	if ack.ResultCode != constants.ResultSuccess {
		t.Fatalf("unexpected ack result code %d", ack.ResultCode)
	}
	if ack.MessageQueue.Count >= before {
		t.Fatalf("queue did not shrink after ack: %d -> %d", before, ack.MessageQueue.Count)
	}
	t.Logf("acknowledged %s, queue %d -> %d", messageID, before, ack.MessageQueue.Count)
}
