package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/secdns"
	"github.com/hariom-pal/go-epp/types"
)

func dnssecDomainCreateExample(client *epp.Client) (*types.DomainCreateResponse, error) {
	return client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "REG123",
		AuthInfo:   "change-me",
		SecDNS: &secdns.CreateRequest{
			Data: secdns.Data{
				DSData: []secdns.DSData{
					{
						KeyTag:     12345,
						Algorithm:  8,
						DigestType: 2,
						Digest:     "49FD46E6C4B45C55D4AC",
					},
				},
			},
		},
	})
}
