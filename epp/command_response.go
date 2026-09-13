package epp

import (
	"encoding/xml"
)

func parseCommandResponse(responseXML []byte) error {
	var response commandResponseXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return err
	}

	return responseResultError(responseXML)
}

type commandResponseXML struct {
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
