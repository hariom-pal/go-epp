package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

const domainRenewCLIDateLayout = "2006-01-02"

func runDomainRenew(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.DomainRenewName) == "" {
		return nil
	}

	currentExpiryDate, err := time.Parse(
		domainRenewCLIDateLayout,
		strings.TrimSpace(options.DomainRenewCurrentExpiryDate),
	)
	if err != nil {
		return fmt.Errorf("cur-expiry must be YYYY-MM-DD: %w", err)
	}

	resp, err := client.DomainRenew(types.DomainRenewRequest{
		DomainName:        options.DomainRenewName,
		CurrentExpiryDate: currentExpiryDate,
		Period:            options.CreatePeriod,
		Unit:              options.CreateUnit,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN RENEW ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("----------------------------------")

	fmt.Printf("Domain Name : %s\n", resp.Result.DomainName)

	if !resp.Result.NewExpiryDate.IsZero() {
		fmt.Printf("New Expiry Date: %s\n", resp.Result.NewExpiryDate.Format(time.RFC3339))
	}

	return nil
}
