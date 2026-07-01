package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runDomainTransfer(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.DomainTransferName) == "" {
		return nil
	}

	resp, err := client.DomainTransfer(types.DomainTransferRequest{
		DomainName: options.DomainTransferName,
		Operation:  options.DomainTransferOperation,
		AuthInfo:   options.CreateAuthInfo,
		Period:     options.CreatePeriod,
		Unit:       options.CreateUnit,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN TRANSFER ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("-------------------------------------")

	fmt.Printf("Domain      : %s\n", resp.TransferData.DomainName)
	fmt.Printf("Transfer Status: %s\n", resp.TransferData.TransferStatus)
	fmt.Printf("Requested By: %s\n", resp.TransferData.RequestedBy)

	if !resp.TransferData.RequestedDate.IsZero() {
		fmt.Printf("Requested Date: %s\n", resp.TransferData.RequestedDate.Format(time.RFC3339))
	}

	fmt.Printf("Action By   : %s\n", resp.TransferData.ActionBy)

	if !resp.TransferData.ActionDate.IsZero() {
		fmt.Printf("Action Date : %s\n", resp.TransferData.ActionDate.Format(time.RFC3339))
	}

	if !resp.TransferData.ExpiryDate.IsZero() {
		fmt.Printf("Expiry Date : %s\n", resp.TransferData.ExpiryDate.Format(time.RFC3339))
	}

	return nil
}
