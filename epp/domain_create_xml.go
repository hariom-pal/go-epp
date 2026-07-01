package epp

import (
	"encoding/xml"

	feeext "github.com/hariom-pal/go-epp/extensions/fee"
	secdnsext "github.com/hariom-pal/go-epp/extensions/secdns"
)

// ============================================================
// DOMAIN CREATE REQUEST
// ============================================================

type domainCreateRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainCreateCommandXML `xml:"command"`
}

type domainCreateCommandXML struct {
	Create     domainCreateXML           `xml:"create"`
	Extension  *domainCreateExtensionXML `xml:"extension,omitempty"`
	ClientTRID string                    `xml:"clTRID"`
}

type domainCreateExtensionXML struct {
	FeeCreate    *feeext.TransformXML `xml:"fee:create,omitempty"`
	SecDNSCreate *secdnsext.CreateXML `xml:"secDNS:create,omitempty"`
}

type domainCreateXML struct {
	Domain domainCreateObjectXML `xml:"domain:create"`
}

type domainCreateObjectXML struct {
	Name string `xml:"domain:name"`

	Period *domainCreatePeriodXML `xml:"domain:period,omitempty"`

	NameServers *domainCreateNameServersXML `xml:"domain:ns,omitempty"`

	Registrant string `xml:"domain:registrant"`

	Contacts []domainCreateContactXML `xml:"domain:contact,omitempty"`

	AuthInfo domainCreateAuthInfoXML `xml:"domain:authInfo"`
}

type domainCreatePeriodXML struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

type domainCreateNameServersXML struct {
	HostObjects []string                  `xml:"domain:hostObj,omitempty"`
	HostAttrs   []domainCreateHostAttrXML `xml:"domain:hostAttr,omitempty"`
}

type domainCreateHostAttrXML struct {
	HostName  string                    `xml:"domain:hostName"`
	HostAddrs []domainCreateHostAddrXML `xml:"domain:hostAddr,omitempty"`
}

type domainCreateHostAddrXML struct {
	Version string `xml:"ip,attr,omitempty"`
	Value   string `xml:",chardata"`
}

type domainCreateContactXML struct {
	Type string `xml:"type,attr"`
	ID   string `xml:",chardata"`
}

type domainCreateAuthInfoXML struct {
	Password string `xml:"domain:pw"`
}

// ============================================================
// DOMAIN CREATE RESPONSE
// ============================================================

type domainCreateResponseXML struct {
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
				ExpiryDate  string `xml:"exDate"`
			} `xml:"creData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`

		Extension struct {
			FeeCreateData feeext.TransformDataXML `xml:"urn:ietf:params:xml:ns:fee-0.7 creData"`
		} `xml:"extension"`
	} `xml:"response"`
}
