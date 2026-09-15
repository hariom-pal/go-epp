package epp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactUpdate updates a contact object.
func (c *Client) ContactUpdate(
	req types.ContactUpdateRequest,
) (*types.ContactUpdateResponse, error) {
	return c.ContactUpdateContext(context.Background(), req)
}

// ContactUpdateContext updates a contact object.
func (c *Client) ContactUpdateContext(
	ctx context.Context,
	req types.ContactUpdateRequest,
) (*types.ContactUpdateResponse, error) {
	requestXML, err := buildContactUpdateRequestXML(
		req,
		c.nextTRID("UPDATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.executeCommandContext(ctx, requestXML, "contact.update", true)
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
		return nil, newValidationError(constants.ResultParameterError, "contact ID is required")
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

		return nil, newValidationError(constants.ResultParameterError, "at least one contact update operation is required")
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

	commonResponse, err := responseEnvelope(responseXML)
	if err != nil {
		return nil, err
	}
	if err := responseResultError(responseXML); err != nil {
		return nil, err
	}

	return &types.ContactUpdateResponse{
		Response: commonResponse,
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
			return nil, newValidationError(constants.ResultParameterError, "contact status is required")
		}
		if !constants.IsContactStatus(value) {
			return nil, newValidationError(constants.ResultParameterError, "invalid contact status")
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
