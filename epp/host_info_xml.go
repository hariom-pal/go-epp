package epp

import "encoding/xml"

// ============================================================
// HOST INFO REQUEST
// ============================================================

type hostInfoRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS     string `xml:"xmlns,attr"`
	HostXMLNS string `xml:"xmlns:host,attr"`

	Command hostInfoCommandXML `xml:"command"`
}

type hostInfoCommandXML struct {
	Info       hostInfoXML `xml:"info"`
	ClientTRID string      `xml:"clTRID"`
}

type hostInfoXML struct {
	Host hostInfoObjectXML `xml:"host:info"`
}

type hostInfoObjectXML struct {
	Name string `xml:"host:name"`
}

// ============================================================
// HOST INFO RESPONSE
// ============================================================

type hostInfoResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			InfoData struct {
				Name string `xml:"name"`
				ROID string `xml:"roid"`

				Statuses []struct {
					Value string `xml:"s,attr"`
				} `xml:"status"`

				Addresses []struct {
					IPVersion string `xml:"ip,attr"`
					Address   string `xml:",chardata"`
				} `xml:"addr"`

				ClientID     string `xml:"clID"`
				CreatedBy    string `xml:"crID"`
				CreatedDate  string `xml:"crDate"`
				UpdatedBy    string `xml:"upID"`
				UpdatedDate  string `xml:"upDate"`
				TransferDate string `xml:"trDate"`
			} `xml:"infData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
