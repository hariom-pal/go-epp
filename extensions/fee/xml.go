package fee

import "encoding/xml"

// CheckExtensionXML wraps a fee-0.7 check command extension.
type CheckExtensionXML struct {
	Check CheckXML `xml:"fee:check"`
}

// CheckXML contains fee-0.7 check request domains.
type CheckXML struct {
	XMLNS   string           `xml:"xmlns:fee,attr"`
	Domains []CheckDomainXML `xml:"fee:domain"`
}

// CheckDomainXML contains one fee-0.7 domain check request.
type CheckDomainXML struct {
	Name     string     `xml:"fee:name"`
	Currency string     `xml:"fee:currency,omitempty"`
	Command  CommandXML `xml:"fee:command"`
	Period   *PeriodXML `xml:"fee:period,omitempty"`
}

// TransformExtensionXML wraps a fee-0.7 transform command extension.
type TransformExtensionXML struct {
	Create   *TransformXML `xml:"fee:create,omitempty"`
	Renew    *TransformXML `xml:"fee:renew,omitempty"`
	Transfer *TransformXML `xml:"fee:transfer,omitempty"`
}

// TransformXML contains fee-0.7 transform request data.
type TransformXML struct {
	XMLNS    string      `xml:"xmlns:fee,attr"`
	Currency string      `xml:"fee:currency,omitempty"`
	Fees     []AmountXML `xml:"fee:fee"`
	Credits  []AmountXML `xml:"fee:credit,omitempty"`
}

// CommandXML contains a fee-0.7 command value and launch metadata.
type CommandXML struct {
	Phase    string `xml:"phase,attr,omitempty"`
	Subphase string `xml:"subphase,attr,omitempty"`
	Value    string `xml:",chardata"`
}

// PeriodXML contains a fee-0.7 period.
type PeriodXML struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

// AmountXML contains a fee-0.7 fee or credit amount.
type AmountXML struct {
	Description string `xml:"description,attr,omitempty"`
	Refundable  string `xml:"refundable,attr,omitempty"`
	GracePeriod string `xml:"grace-period,attr,omitempty"`
	Applied     string `xml:"applied,attr,omitempty"`
	Value       string `xml:",chardata"`
}

// CheckDataXML contains fee-0.7 check response data.
type CheckDataXML struct {
	XMLName xml.Name         `xml:"urn:ietf:params:xml:ns:fee-0.7 chkData"`
	Results []CheckResultXML `xml:"cd"`
}

// CheckResultXML contains fee-0.7 response data for one checked domain.
type CheckResultXML struct {
	Name     string      `xml:"name"`
	Currency string      `xml:"currency"`
	Command  CommandXML  `xml:"command"`
	Period   *PeriodXML  `xml:"period"`
	Fees     []AmountXML `xml:"fee"`
	Credits  []AmountXML `xml:"credit"`
	Class    string      `xml:"class"`
	Reason   string      `xml:"reason"`
}

// InfoDataXML contains fee-0.7 info response data.
type InfoDataXML struct {
	XMLName  xml.Name    `xml:"urn:ietf:params:xml:ns:fee-0.7 infData"`
	Currency string      `xml:"currency"`
	Command  CommandXML  `xml:"command"`
	Period   *PeriodXML  `xml:"period"`
	Fees     []AmountXML `xml:"fee"`
	Credits  []AmountXML `xml:"credit"`
	Class    string      `xml:"class"`
}

// TransformDataXML contains fee-0.7 transform response data.
type TransformDataXML struct {
	XMLName     xml.Name    `xml:""`
	Currency    string      `xml:"currency"`
	Period      *PeriodXML  `xml:"period"`
	Fees        []AmountXML `xml:"fee"`
	Credits     []AmountXML `xml:"credit"`
	Balance     string      `xml:"balance"`
	CreditLimit string      `xml:"creditLimit"`
}
