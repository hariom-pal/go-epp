package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func hostCreateExample(client *epp.Client) (*types.HostCreateResponse, error) {
	return client.HostCreate(types.HostCreateRequest{
		HostName: "ns1.example.in",
		Addresses: []types.HostAddress{
			{IPVersion: "v4", Address: "192.0.2.1"},
			{IPVersion: "v6", Address: "2001:db8::1"},
		},
	})
}
