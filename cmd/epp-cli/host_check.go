package main

import (
	"fmt"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runHostCheck(
	client *epp.Client,
	names []string,
) error {

	if len(names) == 0 {
		return nil
	}

	resp, err := client.HostCheck(types.HostCheckRequest{
		Hosts: names,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== HOST CHECK ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("--------------------------------")

	for _, result := range resp.Results {
		fmt.Printf("Host Name   : %s\n", result.HostName)
		fmt.Printf("ASCII Name  : %s\n", result.ASCIIName)
		fmt.Printf("Available   : %t\n", result.Available)

		if result.Reason != "" {
			fmt.Printf("Reason      : %s\n", result.Reason)
		}

		fmt.Println("--------------------------------")
	}

	return nil
}
