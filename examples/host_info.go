package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func hostInfoExample(client *epp.Client) (*types.HostInfoResponse, error) {
	return client.HostInfo(types.HostInfoRequest{
		HostName: "ns1.example.in",
	})
}
