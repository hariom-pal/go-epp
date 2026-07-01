package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactCreate creates a contact object.
func (c *Client) ContactCreate(
	req types.ContactCreateRequest,
) (*types.ContactCreateResponse, error) {

	requestXML, err := buildContactCreateRequestXML(
		req,
		c.nextTRID("CREATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseContactCreateResponseXML(responseXML)
}

func buildContactCreateRequestXML(
	req types.ContactCreateRequest,
	clientTRID string,
) ([]byte, error) {

	contactID := strings.TrimSpace(req.ContactID)
	if contactID == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "contact ID is required",
		}
	}

	postalInfo, err := contactCreatePostalInfo(req)
	if err != nil {
		return nil, err
	}

	voice, err := contactCreatePhone(req.Voice, true, "voice is required")
	if err != nil {
		return nil, err
	}

	fax, err := contactCreatePhone(req.Fax, false, "")
	if err != nil {
		return nil, err
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "email is required",
		}
	}

	authInfo := strings.TrimSpace(req.AuthInfo)
	if authInfo == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "authInfo is required",
		}
	}

	request := contactCreateRequestXML{
		XMLNS:        constants.EPPNamespace,
		ContactXMLNS: constants.ContactNamespace,
		Command: contactCreateCommandXML{
			ClientTRID: clientTRID,
			Create: contactCreateXML{
				Contact: contactCreateObjectXML{
					ID:         contactID,
					PostalInfo: postalInfo,
					Voice:      voice,
					Fax:        fax,
					Email:      email,
					AuthInfo:   contactCreateAuthInfoXML{Password: authInfo},
					Disclosure: contactCreateDisclosure(req.Disclosure),
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

func parseContactCreateResponseXML(
	responseXML []byte,
) (*types.ContactCreateResponse, error) {

	var response contactCreateResponseXML

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

	createData := response.Response.ResData.CreateData

	resp := &types.ContactCreateResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Result: types.ContactCreateResult{
			ContactID: strings.TrimSpace(createData.ID),
		},
	}

	if createdDate := parseEPPDateTime(createData.CreatedDate); createdDate != nil {
		resp.Result.CreatedDate = *createdDate
	}

	return resp, nil
}

func contactCreatePostalInfo(
	req types.ContactCreateRequest,
) ([]contactCreatePostalInfoXML, error) {

	postalInfo := make([]contactCreatePostalInfoXML, 0, 2)

	if req.InternationalPostalInfo != nil {
		next, err := contactCreatePostalInfoXMLFromTypes("int", req.InternationalPostalInfo)
		if err != nil {
			return nil, err
		}
		postalInfo = append(postalInfo, next)
	}

	if req.LocalizedPostalInfo != nil {
		next, err := contactCreatePostalInfoXMLFromTypes("loc", req.LocalizedPostalInfo)
		if err != nil {
			return nil, err
		}
		postalInfo = append(postalInfo, next)
	}

	if len(postalInfo) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one postalInfo is required",
		}
	}

	return postalInfo, nil
}

func contactCreatePostalInfoXMLFromTypes(
	postalType string,
	info *types.PostalInfo,
) (contactCreatePostalInfoXML, error) {

	name := strings.TrimSpace(info.Name)
	if name == "" {
		return contactCreatePostalInfoXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "postalInfo name is required",
		}
	}

	city := strings.TrimSpace(info.City)
	if city == "" {
		return contactCreatePostalInfoXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "postalInfo city is required",
		}
	}

	countryCode := strings.TrimSpace(info.CountryCode)
	if countryCode == "" {
		return contactCreatePostalInfoXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "postalInfo country code is required",
		}
	}

	return contactCreatePostalInfoXML{
		Type: postalType,
		Name: name,
		Org:  strings.TrimSpace(info.Organization),
		Addr: contactCreateAddressXML{
			Street: trimContactInfoStrings(info.Street),
			City:   city,
			SP:     strings.TrimSpace(info.StateProvince),
			PC:     strings.TrimSpace(info.PostalCode),
			CC:     countryCode,
		},
	}, nil
}

func contactCreatePhone(
	phone types.Phone,
	required bool,
	requiredMessage string,
) (*contactCreatePhoneXML, error) {

	number := strings.TrimSpace(phone.Number)
	if number == "" {
		if !required {
			return nil, nil
		}

		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: requiredMessage,
		}
	}

	return &contactCreatePhoneXML{
		Extension: strings.TrimSpace(phone.Extension),
		Number:    number,
	}, nil
}

func contactCreateDisclosure(
	disclosure *types.ContactDisclosure,
) *contactCreateDisclosureXML {

	if disclosure == nil {
		return nil
	}

	flag := 0
	if disclosure.Flag {
		flag = 1
	}

	result := &contactCreateDisclosureXML{
		Flag:          flag,
		Names:         contactCreateDisclosurePostal(disclosure.Name),
		Organizations: contactCreateDisclosurePostal(disclosure.Organization),
		Addresses:     contactCreateDisclosurePostal(disclosure.Address),
	}

	if disclosure.Voice {
		result.Voice = &struct{}{}
	}

	if disclosure.Fax {
		result.Fax = &struct{}{}
	}

	if disclosure.Email {
		result.Email = &struct{}{}
	}

	return result
}

func contactCreateDisclosurePostal(
	postal types.ContactDisclosurePostal,
) []contactCreateDisclosurePostalXML {

	result := make([]contactCreateDisclosurePostalXML, 0, 2)

	if postal.International {
		result = append(result, contactCreateDisclosurePostalXML{Type: "int"})
	}

	if postal.Localized {
		result = append(result, contactCreateDisclosurePostalXML{Type: "loc"})
	}

	return result
}
