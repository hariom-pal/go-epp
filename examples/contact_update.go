package examples

import (
	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func contactUpdateExample(client *epp.Client) (*types.ContactUpdateResponse, error) {
	voice := types.Phone{Number: "+1.7035550101"}

	return client.ContactUpdate(types.ContactUpdateRequest{
		ContactID:      "CNT001",
		AddStatuses:    []string{constants.ContactStatusClientUpdateProhibited},
		RemoveStatuses: []string{constants.ContactStatusClientDeleteProhibited},
		Voice:          &voice,
		Email:          "updated@example.test",
		AuthInfo:       "change-me",
	})
}
