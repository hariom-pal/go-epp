package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runContactCreate(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.ContactCreateID) == "" {
		return nil
	}

	req := types.ContactCreateRequest{
		ContactID: options.ContactCreateID,
		Voice: types.Phone{
			Number:    options.ContactCreateVoice,
			Extension: options.ContactCreateVoiceExt,
		},
		Fax: types.Phone{
			Number:    options.ContactCreateFax,
			Extension: options.ContactCreateFaxExt,
		},
		Email:    options.ContactCreateEmail,
		AuthInfo: options.ContactCreateAuthInfo,
	}

	if hasPostalInfo(
		options.ContactCreateName,
		options.ContactCreateOrg,
		options.ContactCreateStreets,
		options.ContactCreateCity,
		options.ContactCreateState,
		options.ContactCreatePostalCode,
		options.ContactCreateCountryCode,
	) {

		req.InternationalPostalInfo = &types.PostalInfo{
			Type:          "int",
			Name:          options.ContactCreateName,
			Organization:  options.ContactCreateOrg,
			Street:        options.ContactCreateStreets,
			City:          options.ContactCreateCity,
			StateProvince: options.ContactCreateState,
			PostalCode:    options.ContactCreatePostalCode,
			CountryCode:   options.ContactCreateCountryCode,
		}
	}

	if hasPostalInfo(
		options.ContactCreateLocName,
		options.ContactCreateLocOrg,
		options.ContactCreateLocStreets,
		options.ContactCreateLocCity,
		options.ContactCreateLocState,
		options.ContactCreateLocPostalCode,
		options.ContactCreateLocCountryCode,
	) {

		req.LocalizedPostalInfo = &types.PostalInfo{
			Type:          "loc",
			Name:          options.ContactCreateLocName,
			Organization:  options.ContactCreateLocOrg,
			Street:        options.ContactCreateLocStreets,
			City:          options.ContactCreateLocCity,
			StateProvince: options.ContactCreateLocState,
			PostalCode:    options.ContactCreateLocPostalCode,
			CountryCode:   options.ContactCreateLocCountryCode,
		}
	}

	resp, err := client.ContactCreate(req)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== CONTACT CREATE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("------------------------------------")

	fmt.Printf("Contact ID  : %s\n", resp.Result.ContactID)

	if !resp.Result.CreatedDate.IsZero() {
		fmt.Printf("Created Date: %s\n", resp.Result.CreatedDate.Format(time.RFC3339))
	}

	return nil
}

func hasPostalInfo(
	name string,
	organization string,
	streets []string,
	city string,
	state string,
	postalCode string,
	countryCode string,
) bool {

	if strings.TrimSpace(name) != "" ||
		strings.TrimSpace(organization) != "" ||
		strings.TrimSpace(city) != "" ||
		strings.TrimSpace(state) != "" ||
		strings.TrimSpace(postalCode) != "" ||
		strings.TrimSpace(countryCode) != "" {

		return true
	}

	for _, street := range streets {
		if strings.TrimSpace(street) != "" {
			return true
		}
	}

	return false
}
