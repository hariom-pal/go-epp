package secdns

import "github.com/hariom-pal/go-epp/constants"

// Namespace is the secDNS-1.1 XML namespace.
const Namespace = constants.SecDNSNamespace

// DSData contains RFC5910 Delegation Signer data.
type DSData struct {
	KeyTag     int
	Algorithm  int
	DigestType int
	Digest     string
	KeyData    *KeyData
}

// KeyData contains RFC5910 DNSKEY public key data.
type KeyData struct {
	Flags     int
	Protocol  int
	Algorithm int
	PublicKey string
}

// Data contains DNSSEC data using either the DS Data Interface or Key Data Interface.
type Data struct {
	MaxSigLife int
	DSData     []DSData
	KeyData    []KeyData
}

// CreateRequest contains optional DNSSEC data for a domain create command.
type CreateRequest struct {
	Data
}

// UpdateRequest contains optional DNSSEC data for a domain update command.
type UpdateRequest struct {
	Urgent bool

	Add    *UpdateAdd
	Remove *UpdateRemove
	Change *UpdateChange
}

// UpdateAdd adds DNSSEC DS or key data.
type UpdateAdd struct {
	Data
}

// UpdateRemove removes DNSSEC DS or key data, or all DNSSEC data.
type UpdateRemove struct {
	All     *bool
	DSData  []DSData
	KeyData []KeyData
}

// UpdateChange changes DNSSEC metadata such as maxSigLife.
type UpdateChange struct {
	MaxSigLife int
}

// InfoData contains parsed DNSSEC data returned by domain info.
type InfoData struct {
	Data
}
