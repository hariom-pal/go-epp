package epp

import "encoding/xml"

// ============================================================
// DOMAIN DELETE REQUEST
// ============================================================

type domainDeleteRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainDeleteCommandXML `xml:"command"`
}

type domainDeleteCommandXML struct {
	Delete     domainDeleteXML `xml:"delete"`
	ClientTRID string          `xml:"clTRID"`
}

type domainDeleteXML struct {
	Domain domainDeleteObjectXML `xml:"domain:delete"`
}

type domainDeleteObjectXML struct {
	Name string `xml:"domain:name"`
}

// ============================================================
// DOMAIN DELETE RESPONSE
// ============================================================

type domainDeleteResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
