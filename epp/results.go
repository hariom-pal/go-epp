package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

type responseResultsXML struct {
	Response struct {
		Results []struct {
			Code  int    `xml:"code,attr"`
			Msg   string `xml:"msg"`
			Lang  string `xml:"lang,attr"`
			Value []struct {
				Inner string `xml:",innerxml"`
				Msg   string `xml:"msg"`
			} `xml:"value"`
			ExtValue []struct {
				Value struct {
					Inner string `xml:",innerxml"`
				} `xml:"value"`
				Reason string `xml:"reason"`
			} `xml:"extValue"`
		} `xml:"result"`
		MessageQueue struct {
			Count int    `xml:"count,attr"`
			ID    string `xml:"id,attr"`
			Date  string `xml:"qDate"`
			Msg   string `xml:"msg"`
		} `xml:"msgQ"`

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
		values := make([]types.ResultValue, 0, len(result.Value)+len(result.ExtValue))
		for _, value := range result.Value {
			// Some registries nest an explanatory <msg> inside <value>; when
			// they do, that text is the reason and should not be duplicated
			// into Value as raw XML.
			if strings.TrimSpace(value.Msg) != "" {
				values = append(values, types.ResultValue{Reason: strings.TrimSpace(value.Msg)})
				continue
			}
			values = append(values, types.ResultValue{Value: strings.TrimSpace(value.Inner)})
		}
		for _, extValue := range result.ExtValue {
			values = append(values, types.ResultValue{
				Value:  strings.TrimSpace(extValue.Value.Inner),
				Reason: strings.TrimSpace(extValue.Reason),
			})
		}
		results = append(results, types.Result{
			Code:    result.Code,
			Message: result.Msg,
			Lang:    result.Lang,
			Values:  values,
		})
	}

	queue := parsed.Response.MessageQueue
	resp := types.Response{
		Results: results,
		MessageQueue: types.MessageQueue{
			Count:   queue.Count,
			ID:      strings.TrimSpace(queue.ID),
			Date:    parseEPPDateTime(queue.Date),
			Message: strings.TrimSpace(queue.Msg),
		},
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
			Values:     result.Values,
			ClientTRID: resp.ClientTRID,
			ServerTRID: resp.ServerTRID,
			Results:    resp.Results,
		}
	}

	return nil
}
