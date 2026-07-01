package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func contactCheckExample(client *epp.Client) (*types.ContactCheckResponse, error) {
	return client.ContactCheck(types.ContactCheckRequest{
		IDs: []string{"CNT001", "CNT002"},
	})
}
