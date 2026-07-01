package epp

import (
	"encoding/xml"

	feeext "github.com/hariom-pal/go-epp/extensions/fee"
)

// ============================================================
// DOMAIN TRANSFER REQUEST
// ============================================================

type domainTransferRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS       string `xml:"xmlns,attr"`
	DomainXMLNS string `xml:"xmlns:domain,attr"`

	Command domainTransferCommandXML `xml:"command"`
}

type domainTransferCommandXML struct {
	Transfer   domainTransferXML             `xml:"transfer"`
	Extension  *feeext.TransformExtensionXML `xml:"extension,omitempty"`
	ClientTRID string                        `xml:"clTRID"`
}

type domainTransferXML struct {
	Operation string                  `xml:"op,attr"`
	Domain    domainTransferObjectXML `xml:"domain:transfer"`
}

type domainTransferObjectXML struct {
	Name     string                   `xml:"domain:name"`
	Period   *domainCreatePeriodXML   `xml:"domain:period,omitempty"`
	AuthInfo *domainCreateAuthInfoXML `xml:"domain:authInfo,omitempty"`
}

// ============================================================
// DOMAIN TRANSFER RESPONSE
// ============================================================

type domainTransferResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			TransferData struct {
				Name           string `xml:"name"`
				TransferStatus string `xml:"trStatus"`
				RequestedBy    string `xml:"reID"`
				RequestedDate  string `xml:"reDate"`
				ActionBy       string `xml:"acID"`
				ActionDate     string `xml:"acDate"`
				ExpiryDate     string `xml:"exDate"`
			} `xml:"trnData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`

		Extension struct {
			FeeTransferData feeext.TransformDataXML `xml:"urn:ietf:params:xml:ns:fee-0.7 trnData"`
		} `xml:"extension"`
	} `xml:"response"`
}
