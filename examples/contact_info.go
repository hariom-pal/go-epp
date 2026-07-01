package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func contactInfoExample(client *epp.Client) (*types.ContactInfoResponse, error) {
	return client.ContactInfo(types.ContactInfoRequest{
		ContactID: "CNT001",
	})
}
