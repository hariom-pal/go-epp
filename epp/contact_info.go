package epp

import (
	"context"
	"encoding/xml"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactInfo retrieves RFC5733 information for a contact.
func (c *Client) ContactInfo(
	req types.ContactInfoRequest,
) (*types.ContactInfoResponse, error) {
	return c.ContactInfoContext(context.Background(), req)
}

// ContactInfoContext retrieves RFC5733 information for a contact.
func (c *Client) ContactInfoContext(
	ctx context.Context,
	req types.ContactInfoRequest,
) (*types.ContactInfoResponse, error) {
	contactID := strings.TrimSpace(req.ContactID)
	if contactID == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "contact ID is required",
		}
	}

	request := contactInfoRequestXML{
		XMLNS:        constants.EPPNamespace,
		ContactXMLNS: constants.ContactNamespace,
		Command: contactInfoCommandXML{
			ClientTRID: c.nextTRID("INFO"),
			Info: contactInfoXML{
				Contact: contactInfoObjectXML{
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

	responseXML, err := c.executeCommandContext(ctx, requestXML, "contact.info", false)
	if err != nil {
		return nil, err
	}

	var response contactInfoResponseXML

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

	info := response.Response.ResData.InfoData

	resp := &types.ContactInfoResponse{
		Response: commonResponse,
		Contact: types.ContactInfo{
			ContactID: info.ID,
			ROID:      info.ROID,
			Statuses:  make([]string, 0, len(info.Statuses)),
			Voice: types.Phone{
				Number:    strings.TrimSpace(info.Voice.Number),
				Extension: info.Voice.Extension,
			},
			Fax: types.Phone{
				Number:    strings.TrimSpace(info.Fax.Number),
				Extension: info.Fax.Extension,
			},
			Email:        strings.TrimSpace(info.Email),
			ClientID:     info.ClientID,
			CreatedBy:    info.CreatedBy,
			UpdatedBy:    info.UpdatedBy,
			CreatedDate:  parseContactInfoDateTime(info.CreatedDate),
			UpdatedDate:  parseContactInfoDateTime(info.UpdatedDate),
			TransferDate: parseContactInfoDateTime(info.TransferDate),
			AuthInfo:     strings.TrimSpace(info.AuthInfo.Password),
		},
	}

	for _, status := range info.Statuses {
		if status.Value == "" {
			continue
		}
		resp.Contact.Statuses = append(resp.Contact.Statuses, status.Value)
	}

	for _, postalInfo := range info.PostalInfo {
		parsed := types.PostalInfo{
			Type:          postalInfo.Type,
			Name:          strings.TrimSpace(postalInfo.Name),
			Organization:  strings.TrimSpace(postalInfo.Org),
			Street:        trimContactInfoStrings(postalInfo.Addr.Street),
			City:          strings.TrimSpace(postalInfo.Addr.City),
			StateProvince: strings.TrimSpace(postalInfo.Addr.SP),
			PostalCode:    strings.TrimSpace(postalInfo.Addr.PC),
			CountryCode:   strings.TrimSpace(postalInfo.Addr.CC),
		}

		switch parsed.Type {
		case "int":
			resp.Contact.InternationalPostalInfo = &parsed
		case "loc":
			resp.Contact.LocalizedPostalInfo = &parsed
		}
	}

	return resp, nil
}

func parseContactInfoDateTime(value string) time.Time {
	parsed := parseEPPDateTime(value)
	if parsed == nil {
		return time.Time{}
	}

	return *parsed
}

func trimContactInfoStrings(values []string) []string {
	result := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		result = append(result, value)
	}

	return result
}
