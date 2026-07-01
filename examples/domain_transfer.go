package examples

import (
	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainTransferExample(client *epp.Client) (*types.DomainTransferResponse, error) {
	return client.DomainTransfer(types.DomainTransferRequest{
		DomainName: "example.in",
		Operation:  constants.TransferRequest,
		AuthInfo:   "change-me",
	})
}
