package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runContactInfo(
	client *epp.Client,
	contactID string,
) error {

	if strings.TrimSpace(contactID) == "" {
		return nil
	}

	resp, err := client.ContactInfo(types.ContactInfoRequest{
		ContactID: contactID,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== CONTACT INFO ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("----------------------------------")

	fmt.Printf("Contact ID  : %s\n", resp.Contact.ContactID)
	fmt.Printf("ROID        : %s\n", resp.Contact.ROID)

	for _, status := range resp.Contact.Statuses {
		fmt.Printf("Status      : %s\n", status)
	}

	printPostalInfo("International Postal Info", resp.Contact.InternationalPostalInfo)
	printPostalInfo("Localized Postal Info", resp.Contact.LocalizedPostalInfo)

	fmt.Println("----------------------------------")

	if resp.Contact.Voice.Number != "" {
		fmt.Printf("Voice       : %s\n", formatPhone(resp.Contact.Voice))
	}

	if resp.Contact.Fax.Number != "" {
		fmt.Printf("Fax         : %s\n", formatPhone(resp.Contact.Fax))
	}

	fmt.Printf("Email       : %s\n", resp.Contact.Email)

	fmt.Println("----------------------------------")

	fmt.Printf("Created By  : %s\n", resp.Contact.CreatedBy)
	fmt.Printf("Updated By  : %s\n", resp.Contact.UpdatedBy)
	fmt.Printf("Client ID   : %s\n", resp.Contact.ClientID)

	fmt.Println("----------------------------------")

	if !resp.Contact.CreatedDate.IsZero() {
		fmt.Printf("Created Date: %s\n", resp.Contact.CreatedDate.Format(time.RFC3339))
	}

	if !resp.Contact.UpdatedDate.IsZero() {
		fmt.Printf("Updated Date: %s\n", resp.Contact.UpdatedDate.Format(time.RFC3339))
	}

	if !resp.Contact.TransferDate.IsZero() {
		fmt.Printf("Transfer Date: %s\n", resp.Contact.TransferDate.Format(time.RFC3339))
	}

	return nil
}

func printPostalInfo(
	title string,
	info *types.PostalInfo,
) {

	if info == nil {
		return
	}

	fmt.Println("----------------------------------")
	fmt.Println(title)

	fmt.Printf("Name        : %s\n", info.Name)
	fmt.Printf("Organization: %s\n", info.Organization)

	for _, street := range info.Street {
		fmt.Printf("Street      : %s\n", street)
	}

	fmt.Printf("City        : %s\n", info.City)
	fmt.Printf("State       : %s\n", info.StateProvince)
	fmt.Printf("Postal Code : %s\n", info.PostalCode)
	fmt.Printf("Country     : %s\n", info.CountryCode)
}

func formatPhone(phone types.Phone) string {
	if phone.Extension == "" {
		return phone.Number
	}

	return fmt.Sprintf("%s x%s", phone.Number, phone.Extension)
}
