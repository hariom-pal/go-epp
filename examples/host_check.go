package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func hostCheckExample(client *epp.Client) (*types.HostCheckResponse, error) {
	return client.HostCheck(types.HostCheckRequest{
		Hosts: []string{"ns1.example.in", "ns2.example.in"},
	})
}
