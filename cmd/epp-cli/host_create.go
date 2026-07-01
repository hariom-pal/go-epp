package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runHostCreate(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.HostCreateName) == "" {
		return nil
	}

	addresses := make([]types.HostAddress, 0,
		len(options.HostCreateIPv4)+len(options.HostCreateIPv6),
	)

	for _, address := range options.HostCreateIPv4 {
		addresses = append(addresses, types.HostAddress{
			IPVersion: "v4",
			Address:   address,
		})
	}

	for _, address := range options.HostCreateIPv6 {
		addresses = append(addresses, types.HostAddress{
			IPVersion: "v6",
			Address:   address,
		})
	}

	resp, err := client.HostCreate(types.HostCreateRequest{
		HostName:  options.HostCreateName,
		Addresses: addresses,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== HOST CREATE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("---------------------------------")

	fmt.Printf("Host Name   : %s\n", resp.Result.HostName)

	if !resp.Result.CreatedDate.IsZero() {
		fmt.Printf("Created Date: %s\n", resp.Result.CreatedDate.Format(time.RFC3339))
	}

	return nil
}
