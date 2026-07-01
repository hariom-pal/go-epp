package epp

import "encoding/xml"

// ============================================================
// CONTACT CREATE REQUEST
// ============================================================

type contactCreateRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS        string `xml:"xmlns,attr"`
	ContactXMLNS string `xml:"xmlns:contact,attr"`

	Command contactCreateCommandXML `xml:"command"`
}

type contactCreateCommandXML struct {
	Create     contactCreateXML `xml:"create"`
	ClientTRID string           `xml:"clTRID"`
}

type contactCreateXML struct {
	Contact contactCreateObjectXML `xml:"contact:create"`
}

type contactCreateObjectXML struct {
	ID string `xml:"contact:id"`

	PostalInfo []contactCreatePostalInfoXML `xml:"contact:postalInfo"`

	Voice *contactCreatePhoneXML `xml:"contact:voice,omitempty"`
	Fax   *contactCreatePhoneXML `xml:"contact:fax,omitempty"`
	Email string                 `xml:"contact:email"`

	AuthInfo contactCreateAuthInfoXML `xml:"contact:authInfo"`

	Disclosure *contactCreateDisclosureXML `xml:"contact:disclose,omitempty"`
}

type contactCreatePostalInfoXML struct {
	Type string `xml:"type,attr"`

	Name string `xml:"contact:name"`
	Org  string `xml:"contact:org,omitempty"`

	Addr contactCreateAddressXML `xml:"contact:addr"`
}

type contactCreateAddressXML struct {
	Street []string `xml:"contact:street,omitempty"`
	City   string   `xml:"contact:city"`
	SP     string   `xml:"contact:sp,omitempty"`
	PC     string   `xml:"contact:pc,omitempty"`
	CC     string   `xml:"contact:cc"`
}

type contactCreatePhoneXML struct {
	Extension string `xml:"x,attr,omitempty"`
	Number    string `xml:",chardata"`
}

type contactCreateAuthInfoXML struct {
	Password string `xml:"contact:pw"`
}

type contactCreateDisclosureXML struct {
	Flag int `xml:"flag,attr"`

	Names         []contactCreateDisclosurePostalXML `xml:"contact:name,omitempty"`
	Organizations []contactCreateDisclosurePostalXML `xml:"contact:org,omitempty"`
	Addresses     []contactCreateDisclosurePostalXML `xml:"contact:addr,omitempty"`

	Voice *struct{} `xml:"contact:voice,omitempty"`
	Fax   *struct{} `xml:"contact:fax,omitempty"`
	Email *struct{} `xml:"contact:email,omitempty"`
}

type contactCreateDisclosurePostalXML struct {
	Type string `xml:"type,attr"`
}

// ============================================================
// CONTACT CREATE RESPONSE
// ============================================================

type contactCreateResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			CreateData struct {
				ID          string `xml:"id"`
				CreatedDate string `xml:"crDate"`
			} `xml:"creData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
