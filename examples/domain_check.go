package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainCheckExample(client *epp.Client) (*types.DomainCheckResponse, error) {
	return client.DomainCheck(types.DomainCheckRequest{
		Domains: []string{"example.in", "भारत.भारत"},
	})
}
