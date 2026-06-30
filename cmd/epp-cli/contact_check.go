package main

import (
	"fmt"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runContactCheck(
	client *epp.Client,
	ids []string,
) error {

	if len(ids) == 0 {
		return nil
	}

	resp, err := client.ContactCheck(types.ContactCheckRequest{
		IDs: ids,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== CONTACT CHECK ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("-----------------------------------")

	for _, result := range resp.Results {
		fmt.Printf("Contact ID  : %s\n", result.ContactID)
		fmt.Printf("Available   : %t\n", result.Available)

		if result.Reason != "" {
			fmt.Printf("Reason      : %s\n", result.Reason)
		}

		fmt.Println("-----------------------------------")
	}

	return nil
}
