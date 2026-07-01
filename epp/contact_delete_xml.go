package epp

import "encoding/xml"

// ============================================================
// CONTACT DELETE REQUEST
// ============================================================

type contactDeleteRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS        string `xml:"xmlns,attr"`
	ContactXMLNS string `xml:"xmlns:contact,attr"`

	Command contactDeleteCommandXML `xml:"command"`
}

type contactDeleteCommandXML struct {
	Delete     contactDeleteXML `xml:"delete"`
	ClientTRID string           `xml:"clTRID"`
}

type contactDeleteXML struct {
	Contact contactDeleteObjectXML `xml:"contact:delete"`
}

type contactDeleteObjectXML struct {
	ID string `xml:"contact:id"`
}

// ============================================================
// CONTACT DELETE RESPONSE
// ============================================================

type contactDeleteResponseXML struct {
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
