package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactCheck checks availability for one or more contact IDs.
func (c *Client) ContactCheck(
	req types.ContactCheckRequest,
) (*types.ContactCheckResponse, error) {

	if len(req.IDs) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one contact ID is required",
		}
	}

	ids := make([]string, 0, len(req.IDs))

	for _, id := range req.IDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "no valid contact ID supplied",
		}
	}

	request := contactCheckRequestXML{
		XMLNS:        constants.EPPNamespace,
		ContactXMLNS: constants.ContactNamespace,
		Command: contactCheckCommandXML{
			ClientTRID: c.nextTRID("CHECK"),
			Check: contactCheckXML{
				Contact: contactCheckIDsXML{
					IDs: ids,
				},
			},
		},
	}

	requestXML, err := xml.MarshalIndent(
		request,
		"",
		"    ",
	)
	if err != nil {
		return nil, err
	}

	requestXML = append([]byte(xml.Header), requestXML...)

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	var response contactCheckResponseXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, err
	}

	if response.Response.Result.Code != constants.ResultSuccess &&
		response.Response.Result.Code != constants.ResultSuccessPending {

		return nil, &Error{
			Code:       response.Response.Result.Code,
			Message:    response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		}
	}

	resp := &types.ContactCheckResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Results: make([]types.ContactCheckResult, 0, len(response.Response.ResData.CheckData.CD)),
	}

	for _, cd := range response.Response.ResData.CheckData.CD {
		resp.Results = append(resp.Results, types.ContactCheckResult{
			ContactID: cd.ID.Value,
			ID:        cd.ID.Value,
			Available: cd.ID.Available == 1,
			Reason:    cd.Reason,
		})
	}

	return resp, nil
}
