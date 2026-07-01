package epp

import (
	"encoding/xml"

	secdnsext "github.com/hariom-pal/go-epp/extensions/secdns"
)

// ============================================================
// DOMAIN UPDATE REQUEST
// ============================================================

type domainUpdateRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainUpdateCommandXML `xml:"command"`
}

type domainUpdateCommandXML struct {
	Update     domainUpdateXML           `xml:"update"`
	Extension  *domainUpdateExtensionXML `xml:"extension,omitempty"`
	ClientTRID string                    `xml:"clTRID"`
}

type domainUpdateExtensionXML struct {
	SecDNSUpdate *secdnsext.UpdateXML `xml:"secDNS:update,omitempty"`
}

type domainUpdateXML struct {
	Domain domainUpdateObjectXML `xml:"domain:update"`
}

type domainUpdateObjectXML struct {
	Name string `xml:"domain:name"`

	Add    *domainUpdateAddRemoveXML `xml:"domain:add,omitempty"`
	Remove *domainUpdateAddRemoveXML `xml:"domain:rem,omitempty"`
	Change *domainUpdateChangeXML    `xml:"domain:chg,omitempty"`
}

type domainUpdateAddRemoveXML struct {
	NameServers *domainCreateNameServersXML `xml:"domain:ns,omitempty"`
	Contacts    []domainCreateContactXML    `xml:"domain:contact,omitempty"`
	Statuses    []domainUpdateStatusXML     `xml:"domain:status,omitempty"`
}

type domainUpdateStatusXML struct {
	Status string `xml:"s,attr"`
	Lang   string `xml:"lang,attr,omitempty"`
	Text   string `xml:",chardata"`
}

type domainUpdateChangeXML struct {
	Registrant string                   `xml:"domain:registrant,omitempty"`
	AuthInfo   *domainCreateAuthInfoXML `xml:"domain:authInfo,omitempty"`
}

// ============================================================
// DOMAIN UPDATE RESPONSE
// ============================================================

type domainUpdateResponseXML struct {
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
