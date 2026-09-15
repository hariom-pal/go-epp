package idn

// DataXML contains idn-1.0 extension XML for a command.
type DataXML struct {
	XMLNS string `xml:"xmlns:idn,attr"`
	Table string `xml:"idn:table"`
	UName string `xml:"idn:uname,omitempty"`
}
