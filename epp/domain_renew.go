package epp

import (
	"context"
	"encoding/xml"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	feeext "github.com/hariom-pal/go-epp/extensions/fee"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

const domainRenewDateLayout = "2006-01-02"

// DomainRenew renews a domain registration.
func (c *Client) DomainRenew(
	req types.DomainRenewRequest,
) (*types.DomainRenewResponse, error) {
	return c.DomainRenewContext(context.Background(), req)
}

// DomainRenewContext renews a domain registration.
func (c *Client) DomainRenewContext(
	ctx context.Context,
	req types.DomainRenewRequest,
) (*types.DomainRenewResponse, error) {
	requestXML, err := buildDomainRenewRequestXML(
		req,
		c.nextTRID("RENEW"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.executeCommandContext(ctx, requestXML, "domain.renew", true)
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
		return nil, newValidationError(constants.ResultParameterError, "domain name is required")
	}

	ascii, err := idn.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	if req.CurrentExpiryDate.IsZero() {
		return nil, newValidationError(constants.ResultParameterError, "current expiry date is required")
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
			Extension:  feeext.NewTransformExtension(feeext.CommandRenew, req.Fee),
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

	commonResponse, err := responseEnvelope(responseXML)
	if err != nil {
		return nil, err
	}
	if err := responseResultError(responseXML); err != nil {
		return nil, err
	}

	renewData := response.Response.ResData.RenewData

	unicode, err := idn.ToUnicode(renewData.Name)
	if err != nil {
		unicode = renewData.Name
	}

	resp := &types.DomainRenewResponse{
		Response: commonResponse,
		Result: types.DomainRenewResult{
			Domain:     unicode,
			DomainName: unicode,
			Fee:        feeext.TransformDataFromXML(response.Response.Extension.FeeRenewData),
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

	return domainPeriod(period.Value, period.Unit, true)
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
