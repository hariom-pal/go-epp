package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

func (c *Client) ContactDelete(
	req types.ContactDeleteRequest,
) (*types.ContactDeleteResponse, error) {

	contactID := strings.TrimSpace(req.ContactID)
	if contactID == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "contact ID is required",
		}
	}

	request := contactDeleteRequestXML{
		XMLNS:        constants.EPPNamespace,
		ContactXMLNS: constants.ContactNamespace,
		Command: contactDeleteCommandXML{
			ClientTRID: c.nextTRID("DELETE"),
			Delete: contactDeleteXML{
				Contact: contactDeleteObjectXML{
					ID: contactID,
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

	var response contactDeleteResponseXML

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

	return &types.ContactDeleteResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
	}, nil
}
