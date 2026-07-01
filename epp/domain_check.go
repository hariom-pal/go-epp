package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// DomainCheck checks availability for one or more domains.
func (c *Client) DomainCheck(
	req types.DomainCheckRequest,
) (*types.DomainCheckResponse, error) {

	if len(req.Domains) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one domain is required",
		}
	}

	names := make([]string, 0, len(req.Domains))

	for _, domain := range req.Domains {

		domain = strings.TrimSpace(domain)
		domain = strings.TrimSuffix(domain, ".")

		if domain == "" {
			continue
		}

		ascii, err := idn.ToASCII(domain)
		if err != nil {
			return nil, err
		}

		names = append(names, ascii)
	}

	if len(names) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "no valid domain supplied",
		}
	}

	request := domainCheckRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainCheckCommandXML{
			ClientTRID: c.nextTRID("CHECK"),
			Check: domainCheckXML{
				Domain: domainCheckNamesXML{
					Names: names,
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

	var response domainCheckResponseXML

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

	resp := &types.DomainCheckResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Results: make([]types.DomainCheckResult, 0, len(response.Response.ResData.CheckData.CD)),
	}

	for _, cd := range response.Response.ResData.CheckData.CD {

		unicode, err := idn.ToUnicode(cd.Name.Value)
		if err != nil {
			unicode = cd.Name.Value
		}

		resp.Results = append(resp.Results, types.DomainCheckResult{
			Domain:    unicode,
			ASCII:     cd.Name.Value,
			Available: cd.Name.Available == 1,
			Reason:    cd.Reason,
		})
	}

	return resp, nil
}
