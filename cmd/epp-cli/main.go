package main

import (
	"fmt"
	"log"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
	"github.com/hariom-pal/go-epp/types"
)

func main() {

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("======================================")
	fmt.Println("TLS Connected Successfully")
	fmt.Println("======================================")

	fmt.Println("========== SERVER GREETING ==========")
	fmt.Println(string(client.Greeting()))
	fmt.Println("=====================================")

	// --------------------------------------------------
	// LOGIN
	// --------------------------------------------------

	if err := client.Login(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Login Successful")

	// --------------------------------------------------
	// DOMAIN CHECK
	// --------------------------------------------------

	req := types.DomainCheckRequest{
		Domains: []string{
			"google.in",
			"भारत.भारत",
		},
	}

	resp, err := client.DomainCheck(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("========== DOMAIN CHECK ==========")

	fmt.Printf("Result Code : %d\n", resp.ResultCode)
	fmt.Printf("Result Msg  : %s\n", resp.ResultMsg)
	fmt.Printf("ClientTRID  : %s\n", resp.ClientTRID)
	fmt.Printf("ServerTRID  : %s\n", resp.ServerTRID)

	fmt.Println("----------------------------------")

	for _, result := range resp.Results {

		fmt.Printf("Domain      : %s\n", result.Domain)
		fmt.Printf("ASCII       : %s\n", result.ASCII)
		fmt.Printf("Available   : %t\n", result.Available)

		if result.Reason != "" {
			fmt.Printf("Reason      : %s\n", result.Reason)
		}

		fmt.Println("----------------------------------")
	}

	// --------------------------------------------------
	// LOGOUT
	// --------------------------------------------------

	if err := client.Logout(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Logout Successful")
}
