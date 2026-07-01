package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runPoll(
	client *epp.Client,
	options cliOptions,
) error {

	if !options.Poll && strings.TrimSpace(options.PollAckID) == "" {
		return nil
	}

	req := types.PollRequest{
		Operation: constants.PollRequest,
	}

	if messageID := strings.TrimSpace(options.PollAckID); messageID != "" {
		req.Operation = constants.PollAcknowledge
		req.MessageID = messageID
	}

	resp, err := client.Poll(req)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== POLL ==========")

	fmt.Printf("Result Code    : %d\n", resp.ResultCode)
	fmt.Printf("Result Message : %s\n", resp.ResultMessage)
	fmt.Printf("ClientTRID     : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID     : %s\n", resp.ServerTRID)
	fmt.Printf("Queue Count    : %d\n", resp.MessageQueue.Count)
	fmt.Printf("Message ID     : %s\n", resp.MessageQueue.ID)

	if resp.MessageQueue.Date != nil {
		fmt.Printf("Queue Date     : %s\n", resp.MessageQueue.Date.Format(time.RFC3339))
	}

	fmt.Printf("Message        : %s\n", resp.Message)

	return nil
}
