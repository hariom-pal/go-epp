package main

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runDomainUpdate(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.DomainUpdateName) == "" {
		return nil
	}

	addContacts, err := domainUpdateContactsFromValues(options.DomainUpdateAddContacts)
	if err != nil {
		return err
	}

	removeContacts, err := domainUpdateContactsFromValues(options.DomainUpdateRemoveContacts)
	if err != nil {
		return err
	}

	resp, err := client.DomainUpdate(types.DomainUpdateRequest{
		DomainName:        options.DomainUpdateName,
		AddNameServers:    options.DomainUpdateAddNameServers,
		RemoveNameServers: options.DomainUpdateRemoveNameServers,
		AddContacts:       addContacts,
		RemoveContacts:    removeContacts,
		AddStatuses:       options.DomainUpdateAddStatuses,
		RemoveStatuses:    options.DomainUpdateRemoveStatuses,
		Registrant:        options.DomainUpdateRegistrant,
		AuthInfo:          options.DomainUpdateAuthInfo,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== DOMAIN UPDATE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	return nil
}

func domainUpdateContactsFromValues(
	values []string,
) ([]types.DomainContact, error) {

	contacts := make([]types.DomainContact, 0, len(values))

	for _, value := range values {
		contactType, id, ok := strings.Cut(value, ":")
		if !ok {
			contactType, id, ok = strings.Cut(value, "=")
		}

		contactType = strings.TrimSpace(contactType)
		id = strings.TrimSpace(id)

		if !ok || contactType == "" || id == "" {
			return nil, fmt.Errorf("domain contact must be type:id")
		}

		contacts = append(contacts, types.DomainContact{
			Type: contactType,
			ID:   id,
		})
	}

	return contacts, nil
}
