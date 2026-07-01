package epp

import "encoding/xml"

// ============================================================
// HOST CHECK REQUEST
// ============================================================

type hostCheckRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS     string `xml:"xmlns,attr"`
	HostXMLNS string `xml:"xmlns:host,attr"`

	Command hostCheckCommandXML `xml:"command"`
}

type hostCheckCommandXML struct {
	Check      hostCheckXML `xml:"check"`
	ClientTRID string       `xml:"clTRID"`
}

type hostCheckXML struct {
	Host hostCheckNamesXML `xml:"host:check"`
}

type hostCheckNamesXML struct {
	Names []string `xml:"host:name"`
}

// ============================================================
// HOST CHECK RESPONSE
// ============================================================

type hostCheckResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			CheckData struct {
				CD []struct {
					Name struct {
						Available int    `xml:"avail,attr"`
						Value     string `xml:",chardata"`
					} `xml:"name"`

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
