package epp

import "encoding/xml"

// ============================================================
// HOST CREATE REQUEST
// ============================================================

type hostCreateRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS     string `xml:"xmlns,attr"`
	HostXMLNS string `xml:"xmlns:host,attr"`

	Command hostCreateCommandXML `xml:"command"`
}

type hostCreateCommandXML struct {
	Create     hostCreateXML `xml:"create"`
	ClientTRID string        `xml:"clTRID"`
}

type hostCreateXML struct {
	Host hostCreateObjectXML `xml:"host:create"`
}

type hostCreateObjectXML struct {
	Name      string                 `xml:"host:name"`
	Addresses []hostCreateAddressXML `xml:"host:addr,omitempty"`
}

type hostCreateAddressXML struct {
	IPVersion string `xml:"ip,attr"`
	Address   string `xml:",chardata"`
}

// ============================================================
// HOST CREATE RESPONSE
// ============================================================

type hostCreateResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			CreateData struct {
				Name        string `xml:"name"`
				CreatedDate string `xml:"crDate"`
			} `xml:"creData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
