package epp

import "encoding/xml"

// ============================================================
// HOST DELETE REQUEST
// ============================================================

type hostDeleteRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS     string `xml:"xmlns,attr"`
	HostXMLNS string `xml:"xmlns:host,attr"`

	Command hostDeleteCommandXML `xml:"command"`
}

type hostDeleteCommandXML struct {
	Delete     hostDeleteXML `xml:"delete"`
	ClientTRID string        `xml:"clTRID"`
}

type hostDeleteXML struct {
	Host hostDeleteObjectXML `xml:"host:delete"`
}

type hostDeleteObjectXML struct {
	Name string `xml:"host:name"`
}

// ============================================================
// HOST DELETE RESPONSE
// ============================================================

type hostDeleteResponseXML struct {
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
