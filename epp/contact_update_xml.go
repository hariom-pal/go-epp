package epp

import "encoding/xml"

// ============================================================
// CONTACT UPDATE REQUEST
// ============================================================

type contactUpdateRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS        string `xml:"xmlns,attr"`
	ContactXMLNS string `xml:"xmlns:contact,attr"`

	Command contactUpdateCommandXML `xml:"command"`
}

type contactUpdateCommandXML struct {
	Update     contactUpdateXML `xml:"update"`
	ClientTRID string           `xml:"clTRID"`
}

type contactUpdateXML struct {
	Contact contactUpdateObjectXML `xml:"contact:update"`
}

type contactUpdateObjectXML struct {
	ID string `xml:"contact:id"`

	Add    *contactUpdateStatusListXML `xml:"contact:add,omitempty"`
	Remove *contactUpdateStatusListXML `xml:"contact:rem,omitempty"`
	Change *contactUpdateChangeXML     `xml:"contact:chg,omitempty"`
}

type contactUpdateStatusListXML struct {
	Statuses []contactUpdateStatusXML `xml:"contact:status"`
}

type contactUpdateStatusXML struct {
	Status string `xml:"s,attr"`
}

type contactUpdateChangeXML struct {
	PostalInfo []contactCreatePostalInfoXML `xml:"contact:postalInfo,omitempty"`

	Voice *contactCreatePhoneXML `xml:"contact:voice,omitempty"`
	Fax   *contactCreatePhoneXML `xml:"contact:fax,omitempty"`
	Email string                 `xml:"contact:email,omitempty"`

	AuthInfo *contactCreateAuthInfoXML `xml:"contact:authInfo,omitempty"`

	Disclosure *contactCreateDisclosureXML `xml:"contact:disclose,omitempty"`
}

// ============================================================
// CONTACT UPDATE RESPONSE
// ============================================================

type contactUpdateResponseXML struct {
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
