package epp

import "encoding/xml"

// ============================================================
// CONTACT TRANSFER REQUEST
// ============================================================

type contactTransferRequestXML struct {
	XMLName xml.Name `xml:"epp"`

	XMLNS        string `xml:"xmlns,attr"`
	ContactXMLNS string `xml:"xmlns:contact,attr"`

	Command contactTransferCommandXML `xml:"command"`
}

type contactTransferCommandXML struct {
	Transfer   contactTransferXML `xml:"transfer"`
	ClientTRID string             `xml:"clTRID"`
}

type contactTransferXML struct {
	Operation string                   `xml:"op,attr"`
	Contact   contactTransferObjectXML `xml:"contact:transfer"`
}

type contactTransferObjectXML struct {
	ID       string                    `xml:"contact:id"`
	AuthInfo *contactCreateAuthInfoXML `xml:"contact:authInfo,omitempty"`
}

// ============================================================
// CONTACT TRANSFER RESPONSE
// ============================================================

type contactTransferResponseXML struct {
	XMLName xml.Name `xml:"epp"`

	Response struct {
		Result struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
		} `xml:"result"`

		ResData struct {
			TransferData struct {
				ID             string `xml:"id"`
				TransferStatus string `xml:"trStatus"`
				RequestedBy    string `xml:"reID"`
				RequestedDate  string `xml:"reDate"`
				ActionBy       string `xml:"acID"`
				ActionDate     string `xml:"acDate"`
			} `xml:"trnData"`
		} `xml:"resData"`

		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}
