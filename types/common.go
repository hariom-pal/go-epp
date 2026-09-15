package types

import "time"

// Response contains result and transaction identifiers common to EPP responses.
type Response struct {
	ResultCode int
	ResultMsg  string
	Results    []Result

	// MessageQueue reports the server message queue. RFC 5730 section 2.9.2.3
	// allows <msgQ> on any response, not just a poll, so a registrar can learn
	// that messages are waiting from an ordinary command. Count is zero when
	// the server sent no <msgQ>.
	MessageQueue MessageQueue

	ClientTRID string
	ServerTRID string
}

// Result contains one RFC5730 result element.
type Result struct {
	Code    int
	Message string
	Lang    string

	// Values carries the server diagnostics from <value> and <extValue>
	// elements. RFC 5730 section 2.6 uses these to identify exactly which
	// part of a command a server objected to, so they are the most useful
	// detail available when a command fails.
	Values []ResultValue
}

// ResultValue contains one RFC5730 <value> or <extValue> diagnostic.
type ResultValue struct {
	// Value is the offending command fragment reported by the server.
	Value string

	// Reason explains the failure. It is populated from <extValue><reason>,
	// or from the <msg> some servers nest inside a plain <value>.
	Reason string
}

//
// ============================================================
// GREETING
// ============================================================
//

// Greeting contains RFC5730 server greeting data.
type Greeting struct {
	ServerID string

	ServerDate *time.Time

	Versions            []string
	Languages           []string
	SupportedObjects    []string
	SupportedExtensions []string

	DCP GreetingDCP
}

// GreetingDCP contains RFC5730 data collection policy metadata.
type GreetingDCP struct {
	Access     string
	Statements []GreetingDCPStatement
}

// GreetingDCPStatement contains one RFC5730 data collection policy statement.
type GreetingDCPStatement struct {
	Purposes   []string
	Recipients []string
	Retentions []string

	ExpiryAbsolute *time.Time
	ExpiryRelative string
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
