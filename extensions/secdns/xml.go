package secdns

import "encoding/xml"

// CreateXML contains RFC5910 secDNS create extension XML.
type CreateXML struct {
	XMLNS      string       `xml:"xmlns:secDNS,attr"`
	MaxSigLife int          `xml:"secDNS:maxSigLife,omitempty"`
	DSData     []DSDataXML  `xml:"secDNS:dsData,omitempty"`
	KeyData    []KeyDataXML `xml:"secDNS:keyData,omitempty"`
}

// UpdateXML contains RFC5910 secDNS update extension XML.
type UpdateXML struct {
	XMLNS  string           `xml:"xmlns:secDNS,attr"`
	Urgent string           `xml:"urgent,attr,omitempty"`
	Remove *UpdateRemoveXML `xml:"secDNS:rem,omitempty"`
	Add    *UpdateAddXML    `xml:"secDNS:add,omitempty"`
	Change *UpdateChangeXML `xml:"secDNS:chg,omitempty"`
}

// UpdateAddXML contains secDNS add XML.
type UpdateAddXML struct {
	MaxSigLife int          `xml:"secDNS:maxSigLife,omitempty"`
	DSData     []DSDataXML  `xml:"secDNS:dsData,omitempty"`
	KeyData    []KeyDataXML `xml:"secDNS:keyData,omitempty"`
}

// UpdateRemoveXML contains secDNS remove XML.
type UpdateRemoveXML struct {
	All     *bool        `xml:"secDNS:all,omitempty"`
	DSData  []DSDataXML  `xml:"secDNS:dsData,omitempty"`
	KeyData []KeyDataXML `xml:"secDNS:keyData,omitempty"`
}

// UpdateChangeXML contains secDNS change XML.
type UpdateChangeXML struct {
	MaxSigLife int `xml:"secDNS:maxSigLife,omitempty"`
}

// InfoDataXML contains secDNS info response XML.
type InfoDataXML struct {
	XMLName    xml.Name         `xml:"urn:ietf:params:xml:ns:secDNS-1.1 infData"`
	MaxSigLife int              `xml:"maxSigLife"`
	DSData     []InfoDSDataXML  `xml:"dsData"`
	KeyData    []InfoKeyDataXML `xml:"keyData"`
}

// DSDataXML contains secDNS DS data XML.
type DSDataXML struct {
	KeyTag     int         `xml:"secDNS:keyTag"`
	Algorithm  int         `xml:"secDNS:alg"`
	DigestType int         `xml:"secDNS:digestType"`
	Digest     string      `xml:"secDNS:digest"`
	KeyData    *KeyDataXML `xml:"secDNS:keyData,omitempty"`
}

// KeyDataXML contains secDNS key data XML.
type KeyDataXML struct {
	Flags     int    `xml:"secDNS:flags"`
	Protocol  int    `xml:"secDNS:protocol"`
	Algorithm int    `xml:"secDNS:alg"`
	PublicKey string `xml:"secDNS:pubKey"`
}

// InfoDSDataXML contains secDNS DS data from an info response.
type InfoDSDataXML struct {
	KeyTag     int             `xml:"keyTag"`
	Algorithm  int             `xml:"alg"`
	DigestType int             `xml:"digestType"`
	Digest     string          `xml:"digest"`
	KeyData    *InfoKeyDataXML `xml:"keyData,omitempty"`
}

// InfoKeyDataXML contains secDNS key data from an info response.
type InfoKeyDataXML struct {
	Flags     int    `xml:"flags"`
	Protocol  int    `xml:"protocol"`
	Algorithm int    `xml:"alg"`
	PublicKey string `xml:"pubKey"`
}
