package epp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

const pollOperationRequestXML = "req"

// Poll requests the next queued RFC5730 service message or acknowledges one.
func (c *Client) Poll(
	req types.PollRequest,
) (*types.PollResponse, error) {
	return c.PollContext(context.Background(), req)
}

// PollContext requests the next queued RFC5730 service message or acknowledges one.
func (c *Client) PollContext(
	ctx context.Context,
	req types.PollRequest,
) (*types.PollResponse, error) {
	requestXML, err := buildPollRequestXML(
		req,
		c.nextTRID("POLL"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.ExecuteContext(ctx, requestXML)
	if err != nil {
		return nil, err
	}

	return parsePollResponseXML(responseXML)
}

func buildPollRequestXML(
	req types.PollRequest,
	clientTRID string,
) ([]byte, error) {

	operation, messageID, err := pollRequestValues(req)
	if err != nil {
		return nil, err
	}

	request := pollRequestXML{
		XMLNS: constants.EPPNamespace,
		Command: pollCommandXML{
			ClientTRID: clientTRID,
			Poll: pollXML{
				Operation: operation,
				MessageID: messageID,
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

func pollRequestValues(
	req types.PollRequest,
) (string, string, error) {

	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	switch operation {
	case constants.PollRequest, pollOperationRequestXML:
		return pollOperationRequestXML, "", nil
	case constants.PollAcknowledge:
		messageID := strings.TrimSpace(req.MessageID)
		if messageID == "" {
			return "", "", newValidationError(constants.ResultParameterError, "message ID is required for poll ack")
		}

		return constants.PollAcknowledge, messageID, nil
	default:
		return "", "", newValidationError(constants.ResultParameterError, "poll operation must be request or ack")
	}
}

func parsePollResponseXML(
	responseXML []byte,
) (*types.PollResponse, error) {

	var response pollResponseXML

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

	queue := response.Response.MessageQueue
	message := strings.TrimSpace(queue.Message)

	resp := &types.PollResponse{
		Response:      commonResponse,
		ResultMessage: response.Response.Result.Msg,
		MessageQueue: types.MessageQueue{
			Count:   queue.Count,
			ID:      strings.TrimSpace(queue.ID),
			Message: message,
		},
		Message: message,
	}

	if date := parseEPPDateTime(queue.Date); date != nil {
		resp.MessageQueue.Date = date
	}

	return resp, nil
}
