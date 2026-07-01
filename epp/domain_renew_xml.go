package epp

import "encoding/xml"

// ============================================================
// DOMAIN RENEW REQUEST
// ============================================================

type domainRenewRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainRenewCommandXML `xml:"command"`
}

type domainRenewCommandXML struct {
	Renew      domainRenewXML `xml:"renew"`
	ClientTRID string         `xml:"clTRID"`
}

type domainRenewXML struct {
	Domain domainRenewObjectXML `xml:"domain:renew"`
}

type domainRenewObjectXML struct {
	Name              string                `xml:"domain:name"`
	CurrentExpiryDate string                `xml:"domain:curExpDate"`
	Period            domainCreatePeriodXML `xml:"domain:period"`
}

// ============================================================
// DOMAIN RENEW RESPONSE
// ============================================================

type domainRenewResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			RenewData struct {
				Name          string `xml:"name"`
				NewExpiryDate string `xml:"exDate"`
			} `xml:"renData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
