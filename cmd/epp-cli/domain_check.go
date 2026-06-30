package main

import (
	"fmt"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runDomainCheck(
	client *epp.Client,
	value string,
) error {

	domains := splitDomains(value)
	if len(domains) == 0 {
		return nil
	}

	resp, err := client.DomainCheck(types.DomainCheckRequest{
		Domains: domains,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN CHECK ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("----------------------------------")

	for _, result := range resp.Results {

		fmt.Printf("Domain      : %s\n", result.Domain)
		fmt.Printf("ASCII       : %s\n", result.ASCII)
		fmt.Printf("Available   : %t\n", result.Available)

		if result.Reason != "" {
			fmt.Printf("Reason      : %s\n", result.Reason)
		}

		fmt.Println("----------------------------------")
	}

	return nil
}
