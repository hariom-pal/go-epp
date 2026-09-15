package epp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ContactTransfer performs a contact transfer query, request, approve, cancel, or reject command.
func (c *Client) ContactTransfer(
	req types.ContactTransferRequest,
) (*types.ContactTransferResponse, error) {
	return c.ContactTransferContext(context.Background(), req)
}

// ContactTransferContext performs a contact transfer query, request, approve, cancel, or reject command.
func (c *Client) ContactTransferContext(
	ctx context.Context,
	req types.ContactTransferRequest,
) (*types.ContactTransferResponse, error) {
	requestXML, err := buildContactTransferRequestXML(req, c.nextTRID("TRANSFER"))
	if err != nil {
		return nil, err
	}

	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	transform := operation != constants.TransferQuery
	responseXML, err := c.executeCommandContext(ctx, requestXML, "contact.transfer."+operation, transform)
	if err != nil {
		return nil, err
	}

	return parseContactTransferResponseXML(responseXML)
}

func buildContactTransferRequestXML(
	req types.ContactTransferRequest,
	clientTRID string,
) ([]byte, error) {
	contactID := strings.TrimSpace(req.ContactID)
	if contactID == "" {
		return nil, newValidationError(constants.ResultParameterError, "contact ID is required")
	}

	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	if operation == "" {
		return nil, newValidationError(constants.ResultParameterError, "transfer operation is required")
	}
	if !constants.IsTransferOperation(operation) {
		return nil, newValidationError(constants.ResultParameterError, "invalid transfer operation")
	}

	authInfo := strings.TrimSpace(req.AuthInfo)
	if operation == constants.TransferRequest && authInfo == "" {
		return nil, newValidationError(constants.ResultParameterError, "authInfo is required for transfer request")
	}

	var authInfoXML *contactCreateAuthInfoXML
	if authInfo != "" {
		authInfoXML = &contactCreateAuthInfoXML{Password: authInfo}
	}

	request := contactTransferRequestXML{
		XMLNS:        constants.EPPNamespace,
		ContactXMLNS: constants.ContactNamespace,
		Command: contactTransferCommandXML{
			ClientTRID: clientTRID,
			Transfer: contactTransferXML{
				Operation: operation,
				Contact: contactTransferObjectXML{
					ID:       contactID,
					AuthInfo: authInfoXML,
				},
			},
		},
	}

	requestXML, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), requestXML...), nil
}

func parseContactTransferResponseXML(
	responseXML []byte,
) (*types.ContactTransferResponse, error) {
	var response contactTransferResponseXML
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

	data := response.Response.ResData.TransferData
	transferData := types.ContactTransferData{
		TransferData: types.TransferData{
			ObjectName:     strings.TrimSpace(data.ID),
			TransferStatus: strings.TrimSpace(data.TransferStatus),
			RequestedBy:    strings.TrimSpace(data.RequestedBy),
			ActionBy:       strings.TrimSpace(data.ActionBy),
		},
		ContactID: strings.TrimSpace(data.ID),
	}
	if requestedDate := parseEPPDateTime(data.RequestedDate); requestedDate != nil {
		transferData.RequestedDate = *requestedDate
	}
	if actionDate := parseEPPDateTime(data.ActionDate); actionDate != nil {
		transferData.ActionDate = *actionDate
	}

	return &types.ContactTransferResponse{
		Response:     commonResponse,
		TransferData: transferData,
		Result: types.ContactTransferResult{
			ContactTransferData: transferData,
			ContactID:           transferData.ContactID,
			Status:              transferData.TransferStatus,
		},
	}, nil
}
