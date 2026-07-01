package examples

import (
	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainUpdateExample(client *epp.Client) (*types.DomainUpdateResponse, error) {
	return client.DomainUpdate(types.DomainUpdateRequest{
		DomainName:        "example.in",
		AddNameServers:    []string{"ns3.example.in"},
		RemoveNameServers: []string{"ns1.example.in"},
		AddContacts: []types.DomainContact{
			{Type: "admin", ID: "CNT004"},
		},
		RemoveContacts: []types.DomainContact{
			{Type: "tech", ID: "CNT002"},
		},
		AddStatuses:    []string{constants.StatusClientTransferProhibited},
		RemoveStatuses: []string{constants.StatusClientHold},
		Registrant:     "REG456",
		AuthInfo:       "change-me",
	})
}
