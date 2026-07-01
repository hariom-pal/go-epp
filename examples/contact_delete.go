package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func contactDeleteExample(client *epp.Client) (*types.ContactDeleteResponse, error) {
	return client.ContactDelete(types.ContactDeleteRequest{
		ContactID: "CNT001",
	})
}
