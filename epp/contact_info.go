package epp

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

func (c *Client) ContactInfo(
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

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	var response contactInfoResponseXML

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

	info := response.Response.ResData.InfoData

	resp := &types.ContactInfoResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
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
