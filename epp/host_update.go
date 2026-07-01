package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// HostUpdate updates a host object.
func (c *Client) HostUpdate(
	req types.HostUpdateRequest,
) (*types.HostUpdateResponse, error) {

	requestXML, err := buildHostUpdateRequestXML(
		req,
		c.nextTRID("UPDATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseHostUpdateResponseXML(responseXML)
}

func buildHostUpdateRequestXML(
	req types.HostUpdateRequest,
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

	add, err := hostUpdateList(req.AddAddresses, req.AddStatuses)
	if err != nil {
		return nil, err
	}

	remove, err := hostUpdateList(req.RemoveAddresses, req.RemoveStatuses)
	if err != nil {
		return nil, err
	}

	change, err := hostUpdateChange(req.NewHostName)
	if err != nil {
		return nil, err
	}

	if add == nil &&
		remove == nil &&
		change == nil {

		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one host update operation is required",
		}
	}

	request := hostUpdateRequestXML{
		XMLNS:     constants.EPPNamespace,
		HostXMLNS: constants.HostNamespace,
		Command: hostUpdateCommandXML{
			ClientTRID: clientTRID,
			Update: hostUpdateXML{
				Host: hostUpdateObjectXML{
					Name:   ascii,
					Add:    add,
					Remove: remove,
					Change: change,
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

func parseHostUpdateResponseXML(
	responseXML []byte,
) (*types.HostUpdateResponse, error) {

	var response hostUpdateResponseXML

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

	return &types.HostUpdateResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
	}, nil
}

func hostUpdateList(
	addresses []types.HostAddress,
	statuses []string,
) (*hostUpdateListXML, error) {

	if len(addresses) == 0 &&
		len(statuses) == 0 {

		return nil, nil
	}

	result := &hostUpdateListXML{
		Addresses: make([]hostUpdateAddressXML, 0, len(addresses)),
		Statuses:  make([]hostUpdateStatusXML, 0, len(statuses)),
	}

	for _, address := range addresses {
		parsed, err := hostUpdateAddress(address)
		if err != nil {
			return nil, err
		}

		result.Addresses = append(result.Addresses, parsed)
	}

	for _, status := range statuses {
		status = strings.TrimSpace(status)
		if status == "" {
			return nil, &Error{
				Code:    constants.ResultParameterError,
				Message: "host status is required",
			}
		}

		result.Statuses = append(result.Statuses, hostUpdateStatusXML{
			Status: status,
		})
	}

	return result, nil
}

func hostUpdateAddress(
	address types.HostAddress,
) (hostUpdateAddressXML, error) {

	parsed, err := hostCreateAddress(address)
	if err != nil {
		return hostUpdateAddressXML{}, err
	}

	return hostUpdateAddressXML(parsed), nil
}

func hostUpdateChange(
	hostName string,
) (*hostUpdateChangeXML, error) {

	host := strings.TrimSpace(hostName)
	host = strings.TrimSuffix(host, ".")

	if host == "" {
		return nil, nil
	}

	ascii, err := idn.ToASCII(host)
	if err != nil {
		return nil, err
	}

	return &hostUpdateChangeXML{
		Name: ascii,
	}, nil
}
