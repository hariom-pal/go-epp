package epp

import (
	"encoding/xml"

	feeext "github.com/hariom-pal/go-epp/extensions/fee"
	secdnsext "github.com/hariom-pal/go-epp/extensions/secdns"
)

// ============================================================
// DOMAIN INFO REQUEST
// ============================================================

type domainInfoRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainInfoCommandXML `xml:"command"`
}

type domainInfoCommandXML struct {
	Info       domainInfoXML `xml:"info"`
	ClientTRID string        `xml:"clTRID"`
}

type domainInfoXML struct {
	Domain domainInfoObjectXML `xml:"domain:info"`
}

type domainInfoObjectXML struct {
	Name domainInfoNameXML `xml:"domain:name"`
}

type domainInfoNameXML struct {
	Hosts string `xml:"hosts,attr,omitempty"`
	Value string `xml:",chardata"`
}

// ============================================================
// DOMAIN INFO RESPONSE
// ============================================================

type domainInfoResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			InfoData struct {
				Name       string `xml:"name"`
				ROID       string `xml:"roid"`
				Registrant string `xml:"registrant"`

				Statuses []struct {
					Value string `xml:"s,attr"`
					Lang  string `xml:"lang,attr"`
					Text  string `xml:",chardata"`
				} `xml:"status"`

				Contacts []struct {
					Type string `xml:"type,attr"`
					ID   string `xml:",chardata"`
				} `xml:"contact"`

				NameServers struct {
					HostObjects []string `xml:"hostObj"`
					HostAttrs   []struct {
						HostName string `xml:"hostName"`
						HostAddr []struct {
							Version string `xml:"ip,attr"`
							IP      string `xml:",chardata"`
						} `xml:"hostAddr"`
					} `xml:"hostAttr"`
				} `xml:"ns"`

				Registrar    string `xml:"clID"`
				CreatedBy    string `xml:"crID"`
				CreatedDate  string `xml:"crDate"`
				UpdatedBy    string `xml:"upID"`
				UpdatedDate  string `xml:"upDate"`
				ExpiryDate   string `xml:"exDate"`
				TransferDate string `xml:"trDate"`
				AuthInfo     struct {
					Password struct {
						ROID  string `xml:"roid,attr"`
						Value string `xml:",chardata"`
					} `xml:"pw"`
				} `xml:"authInfo"`
			} `xml:"infData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`

		Extension struct {
			RGPInfoData struct {
				Statuses []struct {
					Value string `xml:"s,attr"`
				} `xml:"rgpStatus"`
			} `xml:"urn:ietf:params:xml:ns:rgp-1.0 infData"`

			SecDNSInfoData secdnsext.InfoDataXML `xml:"urn:ietf:params:xml:ns:secDNS-1.1 infData"`

			FeeInfoData feeext.InfoDataXML `xml:"urn:ietf:params:xml:ns:fee-0.7 infData"`

			LaunchInfoData struct {
				Phase  string `xml:"phase"`
				Status struct {
					Value string `xml:"s,attr"`
					Lang  string `xml:"lang,attr"`
					Text  string `xml:",chardata"`
				} `xml:"status"`
				ApplicationID string `xml:"applicationID"`
			} `xml:"urn:ietf:params:xml:ns:launch-1.0 infData"`

			IDNInfoData struct {
				Table string `xml:"table"`
			} `xml:"urn:ietf:params:xml:ns:idn-1.0 infData"`
		} `xml:"extension"`
	} `xml:"response"`
}
