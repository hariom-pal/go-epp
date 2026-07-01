package main

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runHostUpdate(
	client *epp.Client,
	options cliOptions,
) error {

	if strings.TrimSpace(options.HostUpdateName) == "" {
		return nil
	}

	req := types.HostUpdateRequest{
		HostName:     options.HostUpdateName,
		AddAddresses: hostUpdateAddresses(options.HostUpdateAddIPv4, options.HostUpdateAddIPv6),
		RemoveAddresses: hostUpdateAddresses(
			options.HostUpdateRemoveIPv4,
			options.HostUpdateRemoveIPv6,
		),
		AddStatuses:    options.HostUpdateAddStatuses,
		RemoveStatuses: options.HostUpdateRemoveStatuses,
		NewHostName:    options.HostUpdateNewName,
	}

	resp, err := client.HostUpdate(req)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== HOST UPDATE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	return nil
}

func hostUpdateAddresses(
	ipv4 []string,
	ipv6 []string,
) []types.HostAddress {

	addresses := make([]types.HostAddress, 0, len(ipv4)+len(ipv6))

	for _, address := range ipv4 {
		addresses = append(addresses, types.HostAddress{
			IPVersion: "v4",
			Address:   address,
		})
	}

	for _, address := range ipv6 {
		addresses = append(addresses, types.HostAddress{
			IPVersion: "v6",
			Address:   address,
		})
	}

	return addresses
}
