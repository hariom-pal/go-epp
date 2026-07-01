package epp

import "encoding/xml"

// ============================================================
// CONTACT INFO REQUEST
// ============================================================

type contactInfoRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS        string `xml:"xmlns,attr"`
	ContactXMLNS string `xml:"xmlns:contact,attr"`

	Command contactInfoCommandXML `xml:"command"`
}

type contactInfoCommandXML struct {
	Info       contactInfoXML `xml:"info"`
	ClientTRID string         `xml:"clTRID"`
}

type contactInfoXML struct {
	Contact contactInfoObjectXML `xml:"contact:info"`
}

type contactInfoObjectXML struct {
	ID string `xml:"contact:id"`
}

// ============================================================
// CONTACT INFO RESPONSE
// ============================================================

type contactInfoResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			InfoData struct {
				ID   string `xml:"id"`
				ROID string `xml:"roid"`

				Statuses []struct {
					Value string `xml:"s,attr"`
				} `xml:"status"`

				PostalInfo []struct {
					Type string `xml:"type,attr"`
					Name string `xml:"name"`
					Org  string `xml:"org"`
					Addr struct {
						Street []string `xml:"street"`
						City   string   `xml:"city"`
						SP     string   `xml:"sp"`
						PC     string   `xml:"pc"`
						CC     string   `xml:"cc"`
					} `xml:"addr"`
				} `xml:"postalInfo"`

				Voice struct {
					Extension string `xml:"x,attr"`
					Number    string `xml:",chardata"`
				} `xml:"voice"`

				Fax struct {
					Extension string `xml:"x,attr"`
					Number    string `xml:",chardata"`
				} `xml:"fax"`

				Email string `xml:"email"`

				ClientID     string `xml:"clID"`
				CreatedBy    string `xml:"crID"`
				CreatedDate  string `xml:"crDate"`
				UpdatedBy    string `xml:"upID"`
				UpdatedDate  string `xml:"upDate"`
				TransferDate string `xml:"trDate"`

				AuthInfo struct {
					Password string `xml:"pw"`
				} `xml:"authInfo"`
			} `xml:"infData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
