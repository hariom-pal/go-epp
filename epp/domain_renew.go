package epp

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

const domainRenewDateLayout = "2006-01-02"

func (c *Client) DomainRenew(
	req types.DomainRenewRequest,
) (*types.DomainRenewResponse, error) {

	requestXML, err := buildDomainRenewRequestXML(
		req,
		c.nextTRID("RENEW"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseDomainRenewResponseXML(responseXML)
}

func buildDomainRenewRequestXML(
	req types.DomainRenewRequest,
	clientTRID string,
) ([]byte, error) {

	domain := strings.TrimSpace(req.DomainName)
	if domain == "" {
		domain = strings.TrimSpace(req.Domain)
	}
	domain = strings.TrimSuffix(domain, ".")

	if domain == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "domain name is required",
		}
	}

	ascii, err := idn.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	if req.CurrentExpiryDate.IsZero() {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "current expiry date is required",
		}
	}

	period, err := domainRenewPeriod(req)
	if err != nil {
		return nil, err
	}

	request := domainRenewRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainRenewCommandXML{
			ClientTRID: clientTRID,
			Renew: domainRenewXML{
				Domain: domainRenewObjectXML{
					Name:              ascii,
					CurrentExpiryDate: req.CurrentExpiryDate.Format(domainRenewDateLayout),
					Period: domainCreatePeriodXML{
						Unit:  period.Unit,
						Value: period.Value,
					},
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

func parseDomainRenewResponseXML(
	responseXML []byte,
) (*types.DomainRenewResponse, error) {

	var response domainRenewResponseXML

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

	renewData := response.Response.ResData.RenewData

	unicode, err := idn.ToUnicode(renewData.Name)
	if err != nil {
		unicode = renewData.Name
	}

	resp := &types.DomainRenewResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Result: types.DomainRenewResult{
			Domain:     unicode,
			DomainName: unicode,
		},
	}

	if newExpiryDate := parseDomainRenewDateTime(renewData.NewExpiryDate); newExpiryDate != nil {
		resp.Result.NewExpiryDate = *newExpiryDate
	}

	return resp, nil
}

func domainRenewPeriod(
	req types.DomainRenewRequest,
) (types.Period, error) {

	period := req.PeriodInfo
	if period.Value == 0 && period.Unit == "" {
		period = types.Period{
			Value: req.Period,
			Unit:  req.Unit,
		}
	}

	if period.Value < 1 {
		return types.Period{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "period must be greater than 0",
		}
	}

	period.Unit = strings.ToLower(strings.TrimSpace(period.Unit))
	if period.Unit == "" {
		period.Unit = domainCreatePeriodUnitYears
	}

	if period.Unit != domainCreatePeriodUnitYears &&
		period.Unit != domainCreatePeriodUnitMonths {

		return types.Period{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "period unit must be y or m",
		}
	}

	return period, nil
}

func parseDomainRenewDateTime(
	value string,
) *time.Time {

	if parsed := parseEPPDateTime(value); parsed != nil {
		return parsed
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parsed, err := time.Parse(domainRenewDateLayout, value)
	if err != nil {
		return nil
	}

	return &parsed
}
