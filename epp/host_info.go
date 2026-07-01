package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func (c *Client) HostInfo(
	req types.HostInfoRequest,
) (*types.HostInfoResponse, error) {

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

	request := hostInfoRequestXML{
		XMLNS:     constants.EPPNamespace,
		HostXMLNS: constants.HostNamespace,
		Command: hostInfoCommandXML{
			ClientTRID: c.nextTRID("INFO"),
			Info: hostInfoXML{
				Host: hostInfoObjectXML{
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

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	var response hostInfoResponseXML

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

	unicode, err := idn.ToUnicode(info.Name)
	if err != nil {
		unicode = info.Name
	}

	resp := &types.HostInfoResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Host: types.HostInfo{
			HostName:  unicode,
			ASCIIName: strings.TrimSpace(info.Name),
			ROID:      strings.TrimSpace(info.ROID),

			Statuses:  make([]string, 0, len(info.Statuses)),
			Addresses: make([]types.HostAddress, 0, len(info.Addresses)),

			ClientID:  strings.TrimSpace(info.ClientID),
			CreatedBy: strings.TrimSpace(info.CreatedBy),
			UpdatedBy: strings.TrimSpace(info.UpdatedBy),
		},
	}

	if createdDate := parseEPPDateTime(info.CreatedDate); createdDate != nil {
		resp.Host.CreatedDate = *createdDate
	}

	if updatedDate := parseEPPDateTime(info.UpdatedDate); updatedDate != nil {
		resp.Host.UpdatedDate = *updatedDate
	}

	if transferDate := parseEPPDateTime(info.TransferDate); transferDate != nil {
		resp.Host.TransferDate = *transferDate
	}

	for _, status := range info.Statuses {
		statusValue := strings.TrimSpace(status.Value)
		if statusValue == "" {
			continue
		}

		resp.Host.Statuses = append(resp.Host.Statuses, statusValue)
	}

	for _, address := range info.Addresses {
		ip := strings.TrimSpace(address.Address)
		if ip == "" {
			continue
		}

		resp.Host.Addresses = append(resp.Host.Addresses, types.HostAddress{
			IPVersion: strings.TrimSpace(address.IPVersion),
			Address:   ip,
		})
	}

	return resp, nil
}
