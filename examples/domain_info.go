package examples

import (
	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainInfoExample(client *epp.Client) (*types.DomainInfoResponse, error) {
	return client.DomainInfo(types.DomainInfoRequest{
		Domain: "example.in",
		Hosts:  constants.HostsAll,
	})
}
