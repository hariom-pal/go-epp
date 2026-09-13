package epp

import (
	"encoding/xml"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

type responseResultsXML struct {
	Response struct {
		Results []struct {
			Code int    `xml:"code,attr"`
			Msg  string `xml:"msg"`
			Lang string `xml:"lang,attr"`
		} `xml:"result"`
		TRID struct {
			ClientTRID string `xml:"clTRID"`
			ServerTRID string `xml:"svTRID"`
		} `xml:"trID"`
	} `xml:"response"`
}

func responseEnvelope(responseXML []byte) (types.Response, error) {
	var parsed responseResultsXML
	if err := xml.Unmarshal(responseXML, &parsed); err != nil {
		return types.Response{}, err
	}

	results := make([]types.Result, 0, len(parsed.Response.Results))
	for _, result := range parsed.Response.Results {
		results = append(results, types.Result{
			Code:    result.Code,
			Message: result.Msg,
			Lang:    result.Lang,
		})
	}

	resp := types.Response{
		Results:    results,
		ClientTRID: parsed.Response.TRID.ClientTRID,
		ServerTRID: parsed.Response.TRID.ServerTRID,
	}
	if len(results) > 0 {
		resp.ResultCode = results[0].Code
		resp.ResultMsg = results[0].Message
	}
	return resp, nil
}

func responseResultError(responseXML []byte) error {
	resp, err := responseEnvelope(responseXML)
	if err != nil {
		return err
	}

	for _, result := range resp.Results {
		if constants.IsSuccessResultCode(result.Code) {
			continue
		}
		return &Error{
			Kind:       ErrorKindEPPResult,
			Code:       result.Code,
			Message:    result.Message,
			ClientTRID: resp.ClientTRID,
			ServerTRID: resp.ServerTRID,
			Results:    resp.Results,
		}
	}

	return nil
}
