package main

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runContactUpdate(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.ContactUpdateID) == "" {
		return nil
	}

	req := types.ContactUpdateRequest{
		ContactID:      options.ContactUpdateID,
		AddStatuses:    options.ContactUpdateAddStatuses,
		RemoveStatuses: options.ContactUpdateRemoveStatuses,
		Email:          options.ContactCreateEmail,
		AuthInfo:       options.ContactCreateAuthInfo,
	}

	if strings.TrimSpace(options.ContactCreateVoice) != "" {
		req.Voice = &types.Phone{
			Number:    options.ContactCreateVoice,
			Extension: options.ContactCreateVoiceExt,
		}
	}

	if strings.TrimSpace(options.ContactCreateFax) != "" {
		req.Fax = &types.Phone{
			Number:    options.ContactCreateFax,
			Extension: options.ContactCreateFaxExt,
		}
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

	resp, err := client.ContactUpdate(req)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== CONTACT UPDATE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	return nil
}
