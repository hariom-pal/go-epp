package rgp

import "encoding/xml"

// UpdateXML contains RFC3915 RGP update extension XML.
type UpdateXML struct {
	XMLNS   string      `xml:"xmlns:rgp,attr"`
	Restore *RestoreXML `xml:"rgp:restore"`
}

// RestoreXML contains an RGP restore operation.
type RestoreXML struct {
	Operation string     `xml:"op,attr"`
	Report    *ReportXML `xml:"rgp:report,omitempty"`
}

// ReportXML contains an RGP restore report.
type ReportXML struct {
	PreData       MixedXML  `xml:"rgp:preData"`
	PostData      MixedXML  `xml:"rgp:postData"`
	DeleteTime    string    `xml:"rgp:delTime"`
	RestoreTime   string    `xml:"rgp:resTime"`
	RestoreReason TextXML   `xml:"rgp:resReason"`
	Statements    []TextXML `xml:"rgp:statement"`
	Other         *MixedXML `xml:"rgp:other,omitempty"`
}

// MixedXML contains free text or XML markup allowed by RFC3915.
type MixedXML struct {
	Text string `xml:",chardata"`
	XML  string `xml:",innerxml"`
}

// TextXML contains report text and an optional language attribute.
type TextXML struct {
	Lang string `xml:"lang,attr,omitempty"`
	Text string `xml:",chardata"`
	XML  string `xml:",innerxml"`
}

// InfoDataXML contains RGP info response XML.
type InfoDataXML struct {
	XMLName  xml.Name    `xml:"urn:ietf:params:xml:ns:rgp-1.0 infData"`
	Statuses []StatusXML `xml:"rgpStatus"`
}

// UpdateDataXML contains RGP update response XML.
type UpdateDataXML struct {
	XMLName  xml.Name    `xml:"urn:ietf:params:xml:ns:rgp-1.0 upData"`
	Statuses []StatusXML `xml:"rgpStatus"`
}

// StatusXML contains an RGP response status.
type StatusXML struct {
	Status string `xml:"s,attr"`
	Lang   string `xml:"lang,attr,omitempty"`
	Text   string `xml:",chardata"`
}
