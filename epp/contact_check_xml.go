package epp

import "encoding/xml"

// ============================================================
// CONTACT CHECK REQUEST
// ============================================================

type contactCheckRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS        string `xml:"xmlns,attr"`
	ContactXMLNS string `xml:"xmlns:contact,attr"`

	Command contactCheckCommandXML `xml:"command"`
}

type contactCheckCommandXML struct {
	Check      contactCheckXML `xml:"check"`
	ClientTRID string          `xml:"clTRID"`
}

type contactCheckXML struct {
	Contact contactCheckIDsXML `xml:"contact:check"`
}

type contactCheckIDsXML struct {
	IDs []string `xml:"contact:id"`
}

// ============================================================
// CONTACT CHECK RESPONSE
// ============================================================

type contactCheckResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			CheckData struct {
				CD []struct {
					ID struct {
						Available int    `xml:"avail,attr"`
						Value     string `xml:",chardata"`
					} `xml:"id"`

					Reason string `xml:"reason"`
				} `xml:"cd"`
			} `xml:"chkData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
