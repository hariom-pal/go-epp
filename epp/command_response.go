package epp

import (
	"encoding/xml"

	"github.com/hariom-pal/go-epp/constants"
)

func parseCommandResponse(responseXML []byte) error {
	var response commandResponseXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return err
	}

	if !constants.IsSuccessResultCode(response.Response.Result.Code) {
		return &Error{
			Code:       response.Response.Result.Code,
			Message:    response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		}
	}

	return nil
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
