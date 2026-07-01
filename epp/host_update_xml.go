package epp

import "encoding/xml"

// ============================================================
// HOST UPDATE REQUEST
// ============================================================

type hostUpdateRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS     string `xml:"xmlns,attr"`
	HostXMLNS string `xml:"xmlns:host,attr"`

	Command hostUpdateCommandXML `xml:"command"`
}

type hostUpdateCommandXML struct {
	Update     hostUpdateXML `xml:"update"`
	ClientTRID string        `xml:"clTRID"`
}

type hostUpdateXML struct {
	Host hostUpdateObjectXML `xml:"host:update"`
}

type hostUpdateObjectXML struct {
	Name string `xml:"host:name"`

	Add    *hostUpdateListXML   `xml:"host:add,omitempty"`
	Remove *hostUpdateListXML   `xml:"host:rem,omitempty"`
	Change *hostUpdateChangeXML `xml:"host:chg,omitempty"`
}

type hostUpdateListXML struct {
	Addresses []hostUpdateAddressXML `xml:"host:addr,omitempty"`
	Statuses  []hostUpdateStatusXML  `xml:"host:status,omitempty"`
}

type hostUpdateAddressXML struct {
	IPVersion string `xml:"ip,attr"`
	Address   string `xml:",chardata"`
}

type hostUpdateStatusXML struct {
	Status string `xml:"s,attr"`
}

type hostUpdateChangeXML struct {
	Name string `xml:"host:name"`
}

// ============================================================
// HOST UPDATE RESPONSE
// ============================================================

type hostUpdateResponseXML struct {
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
