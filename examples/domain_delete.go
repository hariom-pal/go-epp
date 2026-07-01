package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainDeleteExample(client *epp.Client) (*types.DomainDeleteResponse, error) {
	return client.DomainDelete(types.DomainDeleteRequest{
		DomainName: "example.in",
	})
}
