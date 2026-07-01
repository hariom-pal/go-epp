package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainCreateExample(client *epp.Client) (*types.DomainCreateResponse, error) {
	return client.DomainCreate(types.DomainCreateRequest{
		Domain:          "example.in",
		Period:          1,
		Unit:            "y",
		Registrant:      "REG123",
		AdminContacts:   []string{"CNT001"},
		TechContacts:    []string{"CNT002"},
		BillingContacts: []string{"CNT003"},
		NameServers:     []string{"ns1.example.in", "ns2.example.in"},
		AuthInfo:        "change-me",
	})
}
