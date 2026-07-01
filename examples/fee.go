package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/fee"
	"github.com/hariom-pal/go-epp/types"
)

func feeDomainCreateExample(client *epp.Client) (*types.DomainCreateResponse, error) {
	return client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "REG123",
		AuthInfo:   "change-me",
		Fee: &fee.TransformRequest{
			Currency: "USD",
			Fees: []fee.Amount{
				{Amount: "5.00"},
			},
		},
	})
}
