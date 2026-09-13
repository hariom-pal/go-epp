package epp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactDelete deletes a contact object.
func (c *Client) ContactDelete(
	req types.ContactDeleteRequest,
) (*types.ContactDeleteResponse, error) {
	return c.ContactDeleteContext(context.Background(), req)
}

// ContactDeleteContext deletes a contact object.
func (c *Client) ContactDeleteContext(
	ctx context.Context,
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

	responseXML, err := c.executeCommandContext(ctx, requestXML, "contact.delete", true)
	if err != nil {
		return nil, err
	}

	var response contactDeleteResponseXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, err
	}

	commonResponse, err := responseEnvelope(responseXML)
	if err != nil {
		return nil, err
	}
	if err := responseResultError(responseXML); err != nil {
		return nil, err
	}

	return &types.ContactDeleteResponse{
		Response: commonResponse,
	}, nil
}
