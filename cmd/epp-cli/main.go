package main

import (
	"fmt"
	"log"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
)

func main() {
	options := parseOptions()

	cfg, err := config.Load(options.ConfigPath)
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

	if err := runDomainCheck(client, options.CheckDomains); err != nil {
		log.Fatal(err)
	}

	if err := runHostCheck(client, options.HostCheckNames); err != nil {
		log.Fatal(err)
	}

	if err := runContactCheck(client, options.ContactCheckIDs); err != nil {
		log.Fatal(err)
	}

	if err := runContactInfo(client, options.ContactInfoID); err != nil {
		log.Fatal(err)
	}

	if err := runContactCreate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runContactUpdate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runContactDelete(client, options.ContactDeleteID); err != nil {
		log.Fatal(err)
	}

	if err := runDomainCreate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainInfo(client, options.InfoDomain, options.InfoHosts); err != nil {
		log.Fatal(err)
	}

	// --------------------------------------------------
	// LOGOUT
	// --------------------------------------------------

	if err := client.Logout(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Logout Successful")
}
