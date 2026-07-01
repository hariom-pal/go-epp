package types

import "time"

// Response contains result and transaction identifiers common to EPP responses.
type Response struct {
	ResultCode int
	ResultMsg  string

	ClientTRID string
	ServerTRID string
}

// PostalInfo contains reusable RFC5733 contact postal information.
type PostalInfo struct {
	Type string

	Name         string
	Organization string

	Street        []string
	City          string
	StateProvince string
	PostalCode    string
	CountryCode   string
}

// Phone contains a telephone number and optional extension.
type Phone struct {
	Number    string
	Extension string
}

// HostAddress contains an IP address and its EPP IP version value.
type HostAddress struct {
	IPVersion string
	Address   string
}

// Period contains a reusable EPP period value and unit.
type Period struct {
	Value int
	Unit  string
}

// TransferData contains object-agnostic transfer response data.
type TransferData struct {
	ObjectName string

	TransferStatus string

	RequestedBy   string
	RequestedDate time.Time

	ActionBy   string
	ActionDate time.Time

	ExpiryDate time.Time
}

//
// ============================================================
// POLL
// ============================================================
//

// PollRequest requests the next queued EPP service message or acknowledges one.
type PollRequest struct {
	Operation string
	MessageID string
}

// PollResponse contains the response for an RFC5730 poll command.
type PollResponse struct {
	Response

	ResultMessage string
	MessageQueue  MessageQueue
	Message       string
}

// MessageQueue contains RFC5730 message queue metadata and message text.
type MessageQueue struct {
	Count   int
	ID      string
	Date    *time.Time
	Message string
}
