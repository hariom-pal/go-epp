package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runHostInfo(
	client *epp.Client,
	hostName string,
) error {

	if strings.TrimSpace(hostName) == "" {
		return nil
	}

	resp, err := client.HostInfo(types.HostInfoRequest{
		HostName: hostName,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== HOST INFO ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("--------------------------------")

	fmt.Printf("Host Name   : %s\n", resp.Host.HostName)
	fmt.Printf("ASCII Name  : %s\n", resp.Host.ASCIIName)
	fmt.Printf("ROID        : %s\n", resp.Host.ROID)

	for _, status := range resp.Host.Statuses {
		fmt.Printf("Status      : %s\n", status)
	}

	for _, address := range resp.Host.Addresses {
		fmt.Printf("IP Version  : %s\n", address.IPVersion)
		fmt.Printf("Address     : %s\n", address.Address)
	}

	fmt.Printf("Created By  : %s\n", resp.Host.CreatedBy)
	fmt.Printf("Updated By  : %s\n", resp.Host.UpdatedBy)
	fmt.Printf("Client ID   : %s\n", resp.Host.ClientID)

	if !resp.Host.CreatedDate.IsZero() {
		fmt.Printf("Created Date: %s\n", resp.Host.CreatedDate.Format(time.RFC3339))
	}

	if !resp.Host.UpdatedDate.IsZero() {
		fmt.Printf("Updated Date: %s\n", resp.Host.UpdatedDate.Format(time.RFC3339))
	}

	if !resp.Host.TransferDate.IsZero() {
		fmt.Printf("Transfer Date: %s\n", resp.Host.TransferDate.Format(time.RFC3339))
	}

	return nil
}
