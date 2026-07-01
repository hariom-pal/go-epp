package epp

import "encoding/xml"

// ============================================================
// POLL REQUEST
// ============================================================

type pollRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS string `xml:"xmlns,attr"`

	Command pollCommandXML `xml:"command"`
}

type pollCommandXML struct {
	Poll       pollXML `xml:"poll"`
	ClientTRID string  `xml:"clTRID"`
}

type pollXML struct {
	Operation string `xml:"op,attr"`
	MessageID string `xml:"msgID,attr,omitempty"`
}

// ============================================================
// POLL RESPONSE
// ============================================================

type pollResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		MessageQueue struct {
			Count   int    `xml:"count,attr"`
			ID      string `xml:"id,attr"`
			Date    string `xml:"qDate"`
			Message string `xml:"msg"`
		} `xml:"msgQ"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
