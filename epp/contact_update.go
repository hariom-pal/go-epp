package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactUpdate updates a contact object.
func (c *Client) ContactUpdate(
	req types.ContactUpdateRequest,
) (*types.ContactUpdateResponse, error) {

	requestXML, err := buildContactUpdateRequestXML(
		req,
		c.nextTRID("UPDATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseContactUpdateResponseXML(responseXML)
}

func buildContactUpdateRequestXML(
	req types.ContactUpdateRequest,
	clientTRID string,
) ([]byte, error) {

	contactID := strings.TrimSpace(req.ContactID)
	if contactID == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "contact ID is required",
		}
	}

	add, err := contactUpdateStatuses(req.AddStatuses)
	if err != nil {
		return nil, err
	}

	remove, err := contactUpdateStatuses(req.RemoveStatuses)
	if err != nil {
		return nil, err
	}

	change, err := contactUpdateChange(req)
	if err != nil {
		return nil, err
	}

	if add == nil &&
		remove == nil &&
		change == nil {

		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one contact update operation is required",
		}
	}

	request := contactUpdateRequestXML{
		XMLNS:        constants.EPPNamespace,
		ContactXMLNS: constants.ContactNamespace,
		Command: contactUpdateCommandXML{
			ClientTRID: clientTRID,
			Update: contactUpdateXML{
				Contact: contactUpdateObjectXML{
					ID:     contactID,
					Add:    add,
					Remove: remove,
					Change: change,
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

	return requestXML, nil
}

func parseContactUpdateResponseXML(
	responseXML []byte,
) (*types.ContactUpdateResponse, error) {

	var response contactUpdateResponseXML

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

	return &types.ContactUpdateResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
	}, nil
}

func contactUpdateStatuses(
	values []string,
) (*contactUpdateStatusListXML, error) {

	if len(values) == 0 {
		return nil, nil
	}

	result := &contactUpdateStatusListXML{
		Statuses: make([]contactUpdateStatusXML, 0, len(values)),
	}

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			return nil, &Error{
				Code:    constants.ResultParameterError,
				Message: "contact status is required",
			}
		}
		if !constants.IsContactStatus(value) {
			return nil, &Error{
				Code:    constants.ResultParameterError,
				Message: "invalid contact status",
			}
		}
		result.Statuses = append(result.Statuses, contactUpdateStatusXML{
			Status: value,
		})
	}

	return result, nil
}

func contactUpdateChange(
	req types.ContactUpdateRequest,
) (*contactUpdateChangeXML, error) {

	change := &contactUpdateChangeXML{}

	if req.InternationalPostalInfo != nil {
		postalInfo, err := contactCreatePostalInfoXMLFromTypes("int", req.InternationalPostalInfo)
		if err != nil {
			return nil, err
		}
		change.PostalInfo = append(change.PostalInfo, postalInfo)
	}

	if req.LocalizedPostalInfo != nil {
		postalInfo, err := contactCreatePostalInfoXMLFromTypes("loc", req.LocalizedPostalInfo)
		if err != nil {
			return nil, err
		}
		change.PostalInfo = append(change.PostalInfo, postalInfo)
	}

	if req.Voice != nil {
		voice, err := contactCreatePhone(*req.Voice, true, "voice is required")
		if err != nil {
			return nil, err
		}
		change.Voice = voice
	}

	if req.Fax != nil {
		fax, err := contactCreatePhone(*req.Fax, true, "fax is required")
		if err != nil {
			return nil, err
		}
		change.Fax = fax
	}

	if email := strings.TrimSpace(req.Email); email != "" {
		change.Email = email
	}

	if authInfo := strings.TrimSpace(req.AuthInfo); authInfo != "" {
		change.AuthInfo = &contactCreateAuthInfoXML{
			Password: authInfo,
		}
	}

	change.Disclosure = contactCreateDisclosure(req.Disclosure)

	if len(change.PostalInfo) == 0 &&
		change.Voice == nil &&
		change.Fax == nil &&
		change.Email == "" &&
		change.AuthInfo == nil &&
		change.Disclosure == nil {

		return nil, nil
	}

	return change, nil
}
