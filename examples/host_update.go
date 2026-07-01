package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func hostUpdateExample(client *epp.Client) (*types.HostUpdateResponse, error) {
	return client.HostUpdate(types.HostUpdateRequest{
		HostName: "ns1.example.in",
		AddAddresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.10"},
		},
		RemoveAddresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.1"},
		},
		AddStatuses: []string{
			"clientUpdateProhibited",
		},
		RemoveStatuses: []string{
			"clientDeleteProhibited",
		},
		NewHostName: "ns2.example.in",
	})
}
