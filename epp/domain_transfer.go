package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// DomainTransfer performs a domain transfer query, request, approve, cancel, or reject command.
func (c *Client) DomainTransfer(
	req types.DomainTransferRequest,
) (*types.DomainTransferResponse, error) {

	requestXML, err := buildDomainTransferRequestXML(
		req,
		c.nextTRID("TRANSFER"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseDomainTransferResponseXML(responseXML)
}

func buildDomainTransferRequestXML(
	req types.DomainTransferRequest,
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

	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	if operation == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "transfer operation is required",
		}
	}

	if !constants.IsTransferOperation(operation) {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "invalid transfer operation",
		}
	}

	authInfo := strings.TrimSpace(req.AuthInfo)
	if operation == constants.TransferRequest && authInfo == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "authInfo is required for transfer request",
		}
	}

	period, err := domainTransferPeriod(req)
	if err != nil {
		return nil, err
	}

	var periodXML *domainCreatePeriodXML
	if period.Value > 0 {
		periodXML = &domainCreatePeriodXML{
			Unit:  period.Unit,
			Value: period.Value,
		}
	}

	var authInfoXML *domainCreateAuthInfoXML
	if authInfo != "" {
		authInfoXML = &domainCreateAuthInfoXML{
			Password: authInfo,
		}
	}

	request := domainTransferRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainTransferCommandXML{
			ClientTRID: clientTRID,
			Transfer: domainTransferXML{
				Operation: operation,
				Domain: domainTransferObjectXML{
					Name:     ascii,
					Period:   periodXML,
					AuthInfo: authInfoXML,
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

func parseDomainTransferResponseXML(
	responseXML []byte,
) (*types.DomainTransferResponse, error) {

	var response domainTransferResponseXML

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

	transferData := domainTransferDataFromXML(
		response.Response.ResData.TransferData,
	)

	return &types.DomainTransferResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		TransferData: transferData,
		Result: types.DomainTransferResult{
			DomainTransferData: transferData,
			Domain:             transferData.DomainName,
			Status:             transferData.TransferStatus,
		},
	}, nil
}

func domainTransferPeriod(
	req types.DomainTransferRequest,
) (types.Period, error) {

	period := req.PeriodInfo
	if period.Value == 0 && period.Unit == "" {
		period = types.Period{
			Value: req.Period,
			Unit:  req.Unit,
		}
	}

	return domainPeriod(period.Value, period.Unit, false)
}

func domainTransferDataFromXML(
	data struct {
		Name           string `xml:"name"`
		TransferStatus string `xml:"trStatus"`
		RequestedBy    string `xml:"reID"`
		RequestedDate  string `xml:"reDate"`
		ActionBy       string `xml:"acID"`
		ActionDate     string `xml:"acDate"`
		ExpiryDate     string `xml:"exDate"`
	},
) types.DomainTransferData {

	domainName := strings.TrimSpace(data.Name)
	unicode, err := idn.ToUnicode(domainName)
	if err != nil {
		unicode = domainName
	}

	result := types.DomainTransferData{
		TransferData: types.TransferData{
			ObjectName:     unicode,
			TransferStatus: strings.TrimSpace(data.TransferStatus),
			RequestedBy:    strings.TrimSpace(data.RequestedBy),
			ActionBy:       strings.TrimSpace(data.ActionBy),
		},
		DomainName: unicode,
	}

	if requestedDate := parseEPPDateTime(data.RequestedDate); requestedDate != nil {
		result.RequestedDate = *requestedDate
	}

	if actionDate := parseEPPDateTime(data.ActionDate); actionDate != nil {
		result.ActionDate = *actionDate
	}

	if expiryDate := parseDomainRenewDateTime(data.ExpiryDate); expiryDate != nil {
		result.ExpiryDate = *expiryDate
	}

	return result
}
