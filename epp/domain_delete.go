package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// DomainDelete deletes a domain object.
func (c *Client) DomainDelete(
	req types.DomainDeleteRequest,
) (*types.DomainDeleteResponse, error) {

	domain, requestXML, err := buildDomainDeleteRequestXML(
		req,
		c.nextTRID("DELETE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseDomainDeleteResponseXML(responseXML, domain)
}

func buildDomainDeleteRequestXML(
	req types.DomainDeleteRequest,
	clientTRID string,
) (string, []byte, error) {

	domain := strings.TrimSpace(req.DomainName)
	if domain == "" {
		domain = strings.TrimSpace(req.Domain)
	}
	domain = strings.TrimSuffix(domain, ".")

	if domain == "" {
		return "", nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "domain name is required",
		}
	}

	ascii, err := idn.ToASCII(domain)
	if err != nil {
		return "", nil, err
	}

	request := domainDeleteRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainDeleteCommandXML{
			ClientTRID: clientTRID,
			Delete: domainDeleteXML{
				Domain: domainDeleteObjectXML{
					Name: ascii,
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
		return "", nil, err
	}

	requestXML = append([]byte(xml.Header), requestXML...)

	return domain, requestXML, nil
}

func parseDomainDeleteResponseXML(
	responseXML []byte,
	domain string,
) (*types.DomainDeleteResponse, error) {

	var response domainDeleteResponseXML

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

	return &types.DomainDeleteResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Result: types.DomainDeleteResult{
			Domain:     domain,
			DomainName: domain,
		},
	}, nil
}
