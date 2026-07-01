package main

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runContactDelete(
	client *epp.Client,
	contactID string,
) error {

	if strings.TrimSpace(contactID) == "" {
		return nil
	}

	resp, err := client.ContactDelete(types.ContactDeleteRequest{
		ContactID: contactID,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== CONTACT DELETE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	return nil
}
