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
