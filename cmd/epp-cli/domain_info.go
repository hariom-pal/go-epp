package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runDomainInfo(
	client *epp.Client,
	domain string,
	hosts string,
) error {

	if strings.TrimSpace(domain) == "" {
		return nil
	}

	resp, err := client.DomainInfo(types.DomainInfoRequest{
		Domain: domain,
		Hosts:  hosts,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN INFO ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("---------------------------------")

	fmt.Printf("Domain      : %s\n", resp.Result.Domain)
	fmt.Printf("ASCII       : %s\n", resp.Result.ASCII)
	fmt.Printf("ROID        : %s\n", resp.Result.ROID)
	fmt.Printf("Registrant  : %s\n", resp.Result.Registrant)
	fmt.Printf("Registrar   : %s\n", resp.Result.Registrar)
	fmt.Printf("Created By  : %s\n", resp.Result.CreatedBy)
	fmt.Printf("Updated By  : %s\n", resp.Result.UpdatedBy)

	if resp.Result.AuthInfo != "" {
		fmt.Println("AuthInfo    : <redacted>")
	}

	if resp.Result.AuthInfoROID != "" {
		fmt.Printf("AuthInfoROID: %s\n", resp.Result.AuthInfoROID)
	}

	for _, contact := range resp.Result.Contacts {
		fmt.Printf("Contact     : %s=%s\n", contact.Type, contact.ID)
	}

	for _, nameServer := range resp.Result.NameServerInfo {
		fmt.Printf("Name Server : %s\n", nameServer.HostName)
		for _, address := range nameServer.Addresses {
			fmt.Printf("Host Address: %s=%s\n", address.Version, address.IP)
		}
	}

	for _, status := range resp.Result.StatusDetails {
		fmt.Printf("Status      : %s\n", status.Status)
		if status.Text != "" {
			fmt.Printf("Status Text : %s\n", status.Text)
		}
	}

	for _, status := range resp.Result.RGPStatuses {
		fmt.Printf("RGP Status  : %s\n", status)
	}

	printDomainDNSSEC(resp.Result.DNSSEC)
	printDomainFee(resp.Result.Fee)
	printDomainLaunch(resp.Result.Launch)

	if resp.Result.IDN.Table != "" {
		fmt.Printf("IDN Table   : %s\n", resp.Result.IDN.Table)
	}

	if resp.Result.CreatedDate != nil {
		fmt.Printf("Created Date: %s\n", resp.Result.CreatedDate.Format(time.RFC3339))
	}

	if resp.Result.UpdatedDate != nil {
		fmt.Printf("Updated Date: %s\n", resp.Result.UpdatedDate.Format(time.RFC3339))
	}

	if resp.Result.ExpiryDate != nil {
		fmt.Printf("Expiry Date : %s\n", resp.Result.ExpiryDate.Format(time.RFC3339))
	}

	if resp.Result.TransferDate != nil {
		fmt.Printf("Transfer Date: %s\n", resp.Result.TransferDate.Format(time.RFC3339))
	}

	return nil
}

func printDomainDNSSEC(info types.DomainDNSSECInfo) {
	if info.MaxSigLife != 0 {
		fmt.Printf("DNSSEC MaxSigLife: %d\n", info.MaxSigLife)
	}

	for _, ds := range info.DSData {
		fmt.Printf("DNSSEC DS   : keyTag=%d alg=%d digestType=%d digest=%s\n",
			ds.KeyTag,
			ds.Algorithm,
			ds.DigestType,
			ds.Digest,
		)
		if ds.KeyData != nil {
			fmt.Printf("DNSSEC Key  : flags=%d protocol=%d alg=%d pubKey=%s\n",
				ds.KeyData.Flags,
				ds.KeyData.Protocol,
				ds.KeyData.Algorithm,
				ds.KeyData.PublicKey,
			)
		}
	}

	for _, key := range info.KeyData {
		fmt.Printf("DNSSEC Key  : flags=%d protocol=%d alg=%d pubKey=%s\n",
			key.Flags,
			key.Protocol,
			key.Algorithm,
			key.PublicKey,
		)
	}
}

func printDomainFee(info types.DomainFeeInfo) {
	if info.Currency != "" {
		fmt.Printf("Fee Currency: %s\n", info.Currency)
	}

	for _, command := range info.Commands {
		fmt.Printf("Fee Command : %s phase=%s subphase=%s\n",
			command.Name,
			command.Phase,
			command.Subphase,
		)
	}

	for _, fee := range info.Fees {
		fmt.Printf("Fee         : %s description=%s refundable=%s grace=%s\n",
			fee.Amount,
			fee.Description,
			fee.Refundable,
			fee.GracePeriod,
		)
	}

	for _, credit := range info.Credits {
		fmt.Printf("Credit      : %s description=%s refundable=%s grace=%s\n",
			credit.Amount,
			credit.Description,
			credit.Refundable,
			credit.GracePeriod,
		)
	}

	if info.Balance != "" {
		fmt.Printf("Fee Balance : %s\n", info.Balance)
	}

	if info.CreditLimit != "" {
		fmt.Printf("Credit Limit: %s\n", info.CreditLimit)
	}
}

func printDomainLaunch(info types.DomainLaunchInfo) {
	if info.Phase != "" {
		fmt.Printf("Launch Phase: %s\n", info.Phase)
	}

	if info.ApplicationID != "" {
		fmt.Printf("Launch AppID: %s\n", info.ApplicationID)
	}

	if info.Status != "" {
		fmt.Printf("Launch Status: %s\n", info.Status)
		if info.StatusText != "" {
			fmt.Printf("Launch Status Text: %s\n", info.StatusText)
		}
	}
}
