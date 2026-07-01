package launch

import "encoding/xml"

// CreateXML contains RFC8334 launch create extension XML.
type CreateXML struct {
	XMLNS     string      `xml:"xmlns:launch,attr"`
	XMLNSSMD  string      `xml:"xmlns:smd,attr,omitempty"`
	XMLNSMark string      `xml:"xmlns:mark,attr,omitempty"`
	Type      string      `xml:"type,attr,omitempty"`
	Phase     PhaseXML    `xml:"launch:phase"`
	ChoiceXML RawXML      `xml:",innerxml"`
	Notices   []NoticeXML `xml:"launch:notice,omitempty"`
}

// InfoXML contains RFC8334 launch info extension XML.
type InfoXML struct {
	XMLNS         string   `xml:"xmlns:launch,attr"`
	IncludeMark   string   `xml:"includeMark,attr,omitempty"`
	Phase         PhaseXML `xml:"launch:phase"`
	ApplicationID string   `xml:"launch:applicationID,omitempty"`
}

// UpdateXML contains RFC8334 launch update extension XML.
type UpdateXML struct {
	XMLNS         string   `xml:"xmlns:launch,attr"`
	Phase         PhaseXML `xml:"launch:phase"`
	ApplicationID string   `xml:"launch:applicationID"`
}

// DeleteXML contains RFC8334 launch delete extension XML.
type DeleteXML struct {
	XMLNS         string   `xml:"xmlns:launch,attr"`
	Phase         PhaseXML `xml:"launch:phase"`
	ApplicationID string   `xml:"launch:applicationID"`
}

// PhaseXML contains a launch phase.
type PhaseXML struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:",chardata"`
}

// CodeMarkXML contains a launch codeMark element.
type CodeMarkXML struct {
	XMLName xml.Name `xml:"launch:codeMark"`
	Code    *CodeXML `xml:"launch:code,omitempty"`
	MarkXML RawXML   `xml:",innerxml"`
}

// CodeXML contains a launch code.
type CodeXML struct {
	ValidatorID string `xml:"validatorID,attr,omitempty"`
	Value       string `xml:",chardata"`
}

// NoticeXML contains a launch claims notice acknowledgement.
type NoticeXML struct {
	NoticeID     NoticeIDXML `xml:"launch:noticeID"`
	NotAfter     string      `xml:"launch:notAfter"`
	AcceptedDate string      `xml:"launch:acceptedDate"`
}

// NoticeIDXML contains a claims notice identifier.
type NoticeIDXML struct {
	ValidatorID string `xml:"validatorID,attr,omitempty"`
	Value       string `xml:",chardata"`
}

// IDDataXML contains launch create response XML.
type IDDataXML struct {
	XMLName       xml.Name `xml:"urn:ietf:params:xml:ns:launch-1.0 creData"`
	Phase         PhaseXML `xml:"phase"`
	ApplicationID string   `xml:"applicationID"`
}

// InfoDataXML contains launch info response XML.
type InfoDataXML struct {
	XMLName       xml.Name  `xml:"urn:ietf:params:xml:ns:launch-1.0 infData"`
	Phase         PhaseXML  `xml:"phase"`
	ApplicationID string    `xml:"applicationID"`
	Status        StatusXML `xml:"status"`
	RawXML        string    `xml:",innerxml"`
}

// StatusXML contains launch status response XML.
type StatusXML struct {
	Status string `xml:"s,attr"`
	Name   string `xml:"name,attr,omitempty"`
	Lang   string `xml:"lang,attr,omitempty"`
	Text   string `xml:",chardata"`
}

// RawXML carries already-formed XML for schema-extension content.
type RawXML struct {
	Value string `xml:",innerxml"`
}
