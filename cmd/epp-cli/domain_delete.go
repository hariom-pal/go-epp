package main

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runDomainDelete(
	client *epp.Client,
	domainName string,
) error {

	if strings.TrimSpace(domainName) == "" {
		return nil
	}

	resp, err := client.DomainDelete(types.DomainDeleteRequest{
		DomainName: domainName,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN DELETE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	return nil
}
