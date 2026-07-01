package main

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func runHostDelete(
	client *epp.Client,
	hostName string,
) error {

	if strings.TrimSpace(hostName) == "" {
		return nil
	}

	resp, err := client.HostDelete(types.HostDeleteRequest{
		HostName: hostName,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("========== HOST DELETE ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	return nil
}
