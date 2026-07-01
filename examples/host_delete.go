package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func hostDeleteExample(client *epp.Client) (*types.HostDeleteResponse, error) {
	return client.HostDelete(types.HostDeleteRequest{
		HostName: "ns1.example.in",
	})
}
