package epp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// HostDelete deletes a host object.
func (c *Client) HostDelete(
	req types.HostDeleteRequest,
) (*types.HostDeleteResponse, error) {
	return c.HostDeleteContext(context.Background(), req)
}

// HostDeleteContext deletes a host object.
func (c *Client) HostDeleteContext(
	ctx context.Context,
	req types.HostDeleteRequest,
) (*types.HostDeleteResponse, error) {
	host := strings.TrimSpace(req.HostName)
	host = strings.TrimSuffix(host, ".")

	if host == "" {
		return nil, newValidationError(constants.ResultParameterError, "host name is required")
	}

	ascii, err := idn.ToASCII(host)
	if err != nil {
		return nil, err
	}

	request := hostDeleteRequestXML{
		XMLNS:     constants.EPPNamespace,
		HostXMLNS: constants.HostNamespace,
		Command: hostDeleteCommandXML{
			ClientTRID: c.nextTRID("DELETE"),
			Delete: hostDeleteXML{
				Host: hostDeleteObjectXML{
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
		return nil, err
	}

	requestXML = append([]byte(xml.Header), requestXML...)

	responseXML, err := c.executeCommandContext(ctx, requestXML, "host.delete", true)
	if err != nil {
		return nil, err
	}

	var response hostDeleteResponseXML

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

	return &types.HostDeleteResponse{
		Response: commonResponse,
	}, nil
}
