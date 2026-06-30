package epp

import "encoding/xml"

// ============================================================
// DOMAIN CHECK REQUEST
// ============================================================

type domainCheckRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainCheckCommandXML `xml:"command"`
}

type domainCheckCommandXML struct {
	Check   domainCheckXML `xml:"check"`
	ClientTRID string       `xml:"clTRID"`
}

type domainCheckXML struct {
	Domain domainCheckNamesXML `xml:"domain:check"`
}

type domainCheckNamesXML struct {
	Names []string `xml:"domain:name"`
}

// ============================================================
// DOMAIN CHECK RESPONSE
// ============================================================

type domainCheckResponseXML struct {
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
