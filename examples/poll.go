package examples

import (
	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func pollRequestExample(client *epp.Client) (*types.PollResponse, error) {
	return client.Poll(types.PollRequest{
		Operation: constants.PollRequest,
	})
}

func pollAcknowledgeExample(client *epp.Client, messageID string) (*types.PollResponse, error) {
	return client.Poll(types.PollRequest{
		Operation: constants.PollAcknowledge,
		MessageID: messageID,
	})
}
