package epp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	launchext "github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// DomainDelete deletes a domain object.
func (c *Client) DomainDelete(
	req types.DomainDeleteRequest,
) (*types.DomainDeleteResponse, error) {
	return c.DomainDeleteContext(context.Background(), req)
}

// DomainDeleteContext deletes a domain object.
func (c *Client) DomainDeleteContext(
	ctx context.Context,
	req types.DomainDeleteRequest,
) (*types.DomainDeleteResponse, error) {
	domain, requestXML, err := buildDomainDeleteRequestXML(
		req,
		c.nextTRID("DELETE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.executeCommandContext(ctx, requestXML, "domain.delete", true)
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
		return "", nil, newValidationError(constants.ResultParameterError, "domain name is required")
	}

	ascii, err := idn.ToASCII(domain)
	if err != nil {
		return "", nil, err
	}

	if !launchext.ValidDelete(req.Launch) {
		return "", nil, newValidationError(constants.ResultParameterError, "invalid launch delete extension")
	}

	request := domainDeleteRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainDeleteCommandXML{
			ClientTRID: clientTRID,
			Extension:  domainDeleteExtension(req),
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

func domainDeleteExtension(
	req types.DomainDeleteRequest,
) *domainDeleteExtensionXML {

	launchDelete := launchext.NewDelete(req.Launch)
	if launchDelete == nil {
		return nil
	}

	return &domainDeleteExtensionXML{
		LaunchDelete: launchDelete,
	}
}

func parseDomainDeleteResponseXML(
	responseXML []byte,
	domain string,
) (*types.DomainDeleteResponse, error) {

	var response domainDeleteResponseXML

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

	return &types.DomainDeleteResponse{
		Response: commonResponse,
		Result: types.DomainDeleteResult{
			Domain:     domain,
			DomainName: domain,
		},
	}, nil
}
