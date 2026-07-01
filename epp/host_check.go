package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func (c *Client) HostCheck(
	req types.HostCheckRequest,
) (*types.HostCheckResponse, error) {

	if len(req.Hosts) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one host name is required",
		}
	}

	names := make([]string, 0, len(req.Hosts))

	for _, host := range req.Hosts {
		host = strings.TrimSpace(host)
		host = strings.TrimSuffix(host, ".")

		if host == "" {
			continue
		}

		ascii, err := idn.ToASCII(host)
		if err != nil {
			return nil, err
		}

		names = append(names, ascii)
	}

	if len(names) == 0 {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "no valid host name supplied",
		}
	}

	request := hostCheckRequestXML{
		XMLNS:     constants.EPPNamespace,
		HostXMLNS: constants.HostNamespace,
		Command: hostCheckCommandXML{
			ClientTRID: c.nextTRID("CHECK"),
			Check: hostCheckXML{
				Host: hostCheckNamesXML{
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

	var response hostCheckResponseXML

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

	resp := &types.HostCheckResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Results: make([]types.HostCheckResult, 0, len(response.Response.ResData.CheckData.CD)),
	}

	for _, cd := range response.Response.ResData.CheckData.CD {
		unicode, err := idn.ToUnicode(cd.Name.Value)
		if err != nil {
			unicode = cd.Name.Value
		}

		resp.Results = append(resp.Results, types.HostCheckResult{
			HostName:  unicode,
			ASCIIName: cd.Name.Value,
			Available: cd.Name.Available == 1,
			Reason:    cd.Reason,
		})
	}

	return resp, nil
}
