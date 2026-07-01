package epp

import (
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

	requestXML, err := buildPollRequestXML(
		req,
		c.nextTRID("POLL"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
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
			return "", "", &Error{
				Code:    constants.ResultParameterError,
				Message: "message ID is required for poll ack",
			}
		}

		return constants.PollAcknowledge, messageID, nil
	default:
		return "", "", &Error{
			Code:    constants.ResultParameterError,
			Message: "poll operation must be request or ack",
		}
	}
}

func parsePollResponseXML(
	responseXML []byte,
) (*types.PollResponse, error) {

	var response pollResponseXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, err
	}

	if !constants.IsSuccessResultCode(response.Response.Result.Code) {
		return nil, &Error{
			Code:       response.Response.Result.Code,
			Message:    response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		}
	}

	queue := response.Response.MessageQueue
	message := strings.TrimSpace(queue.Message)

	resp := &types.PollResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
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
