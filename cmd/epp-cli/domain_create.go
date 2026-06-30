package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runDomainCreate(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.CreateDomain) == "" {
		return nil
	}

	resp, err := client.DomainCreate(types.DomainCreateRequest{
		Domain:          options.CreateDomain,
		Period:          options.CreatePeriod,
		Unit:            options.CreateUnit,
		Registrant:      options.CreateRegistrant,
		AdminContacts:   options.CreateAdminContacts,
		TechContacts:    options.CreateTechContacts,
		BillingContacts: options.CreateBillingContacts,
		NameServers:     options.CreateNameServers,
		AuthInfo:        options.CreateAuthInfo,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN CREATE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("-----------------------------------")

	fmt.Printf("Domain      : %s\n", resp.Result.Domain)

	if !resp.Result.CreatedDate.IsZero() {
		fmt.Printf("Created Date: %s\n", resp.Result.CreatedDate.Format(time.RFC3339))
	}

	if !resp.Result.ExpiryDate.IsZero() {
		fmt.Printf("Expiry Date : %s\n", resp.Result.ExpiryDate.Format(time.RFC3339))
	}

	return nil
}
