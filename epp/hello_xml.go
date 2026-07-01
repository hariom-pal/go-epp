package epp

import "encoding/xml"

const helloRequestXML = `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <hello/>
</epp>`

type greetingEnvelopeXML struct {
	XMLName  xml.Name    `xml:"epp"`
	Greeting greetingXML `xml:"greeting"`
}

type greetingXML struct {
	ServerID   string             `xml:"svID"`
	ServerDate string             `xml:"svDate"`
	Service    greetingServiceXML `xml:"svcMenu"`
	DCP        greetingDCPXML     `xml:"dcp"`
}

type greetingServiceXML struct {
	Versions   []string `xml:"version"`
	Languages  []string `xml:"lang"`
	Objects    []string `xml:"objURI"`
	Extensions struct {
		URIs []string `xml:"extURI"`
	} `xml:"svcExtension"`
}

type greetingDCPXML struct {
	Access     greetingAnyElementListXML `xml:"access"`
	Statements []greetingDCPStatementXML `xml:"statement"`
}

type greetingDCPStatementXML struct {
	Purpose   greetingAnyElementListXML `xml:"purpose"`
	Recipient greetingAnyElementListXML `xml:"recipient"`
	Retention greetingAnyElementListXML `xml:"retention"`
	Expiry    struct {
		Absolute string `xml:"absolute"`
		Relative string `xml:"relative"`
	} `xml:"expiry"`
}

type greetingAnyElementXML struct {
	XMLName xml.Name
}

type greetingAnyElementListXML struct {
	Elements []greetingAnyElementXML `xml:",any"`
}
