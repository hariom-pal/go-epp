package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func contactCreateExample(client *epp.Client) (*types.ContactCreateResponse, error) {
	return client.ContactCreate(types.ContactCreateRequest{
		ContactID: "CNT001",
		InternationalPostalInfo: &types.PostalInfo{
			Type:          "int",
			Name:          "Example User",
			Organization:  "Example Inc.",
			Street:        []string{"1 Example Street"},
			City:          "Dulles",
			StateProvince: "VA",
			PostalCode:    "20166",
			CountryCode:   "US",
		},
		Voice: types.Phone{
			Number: "+1.7035555555",
		},
		Email:    "user@example.test",
		AuthInfo: "change-me",
	})
}
