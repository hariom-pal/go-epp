package epp

import (
	"context"
	"encoding/xml"
	"net"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

const (
	hostAddressIPv4 = "v4"
	hostAddressIPv6 = "v6"
)

// HostCreate creates a host object.
func (c *Client) HostCreate(
	req types.HostCreateRequest,
) (*types.HostCreateResponse, error) {
	return c.HostCreateContext(context.Background(), req)
}

// HostCreateContext creates a host object.
func (c *Client) HostCreateContext(
	ctx context.Context,
	req types.HostCreateRequest,
) (*types.HostCreateResponse, error) {
	requestXML, err := buildHostCreateRequestXML(
		req,
		c.nextTRID("CREATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.executeCommandContext(ctx, requestXML, "host.create", true)
	if err != nil {
		return nil, err
	}

	return parseHostCreateResponseXML(responseXML)
}

func buildHostCreateRequestXML(
	req types.HostCreateRequest,
	clientTRID string,
) ([]byte, error) {

	host := strings.TrimSpace(req.HostName)
	host = strings.TrimSuffix(host, ".")

	if host == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "host name is required",
		}
	}

	ascii, err := idn.ToASCII(host)
	if err != nil {
		return nil, err
	}

	addresses, err := hostCreateAddresses(req.Addresses)
	if err != nil {
		return nil, err
	}

	request := hostCreateRequestXML{
		XMLNS:     constants.EPPNamespace,
		HostXMLNS: constants.HostNamespace,
		Command: hostCreateCommandXML{
			ClientTRID: clientTRID,
			Create: hostCreateXML{
				Host: hostCreateObjectXML{
					Name:      ascii,
					Addresses: addresses,
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

func parseHostCreateResponseXML(
	responseXML []byte,
) (*types.HostCreateResponse, error) {

	var response hostCreateResponseXML

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

	createData := response.Response.ResData.CreateData

	unicode, err := idn.ToUnicode(createData.Name)
	if err != nil {
		unicode = createData.Name
	}

	resp := &types.HostCreateResponse{
		Response: commonResponse,
		Result: types.HostCreateResult{
			HostName: unicode,
		},
	}

	if createdDate := parseEPPDateTime(createData.CreatedDate); createdDate != nil {
		resp.Result.CreatedDate = *createdDate
	}

	return resp, nil
}

func hostCreateAddresses(
	addresses []types.HostAddress,
) ([]hostCreateAddressXML, error) {

	result := make([]hostCreateAddressXML, 0, len(addresses))

	for _, address := range addresses {
		parsed, err := hostCreateAddress(address)
		if err != nil {
			return nil, err
		}

		result = append(result, parsed)
	}

	return result, nil
}

func hostCreateAddress(
	address types.HostAddress,
) (hostCreateAddressXML, error) {

	ipVersion := strings.ToLower(strings.TrimSpace(address.IPVersion))
	if ipVersion != hostAddressIPv4 &&
		ipVersion != hostAddressIPv6 {

		return hostCreateAddressXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "host address IP version must be v4 or v6",
		}
	}

	ipAddress := strings.TrimSpace(address.Address)
	parsedIP := net.ParseIP(ipAddress)
	if parsedIP == nil {
		return hostCreateAddressXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "host address must be a valid IP address",
		}
	}

	if ipVersion == hostAddressIPv4 && parsedIP.To4() == nil {
		return hostCreateAddressXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "host address must be a valid IPv4 address",
		}
	}

	if ipVersion == hostAddressIPv6 && parsedIP.To4() != nil {
		return hostCreateAddressXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "host address must be a valid IPv6 address",
		}
	}

	return hostCreateAddressXML{
		IPVersion: ipVersion,
		Address:   ipAddress,
	}, nil
}
