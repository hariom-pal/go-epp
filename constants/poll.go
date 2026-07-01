package constants

const (
	// PollRequest requests the next message from the server message queue.
	PollRequest = "request"

	// PollAcknowledge acknowledges and dequeues a server message.
	PollAcknowledge = "ack"
)

// IsPollOperation reports whether operation is a supported poll operation.
func IsPollOperation(operation string) bool {
	switch operation {
	case PollRequest,
		PollAcknowledge:
		return true
	default:
		return false
	}
}
