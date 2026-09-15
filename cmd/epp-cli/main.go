package main

import (
	"fmt"
	"log"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/internal/config"
)

func main() {
	options := parseOptions()
	plan := planOperations(options)

	cfg, err := config.LoadFromFile(options.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	environment, err := validateUATSafety(cfg, options, plan)
	if err != nil {
		log.Fatal(err)
	}
	if options.UAT {
		printTargetSummary(cfg, environment, plan)
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	capture, err := newCaptureLogger(options.CaptureDir)
	if err != nil {
		log.Fatal(err)
	}
	defer capture.Close()
	if capture != nil {
		client.SetLogger(capture)
		capture.EPPEvent(epp.Event{Type: epp.EventConnect})
	}

	fmt.Println("======================================")
	fmt.Println("TLS Connected Successfully")
	fmt.Println("======================================")

	fmt.Println("========== SERVER GREETING ==========")
	fmt.Println(string(client.Greeting()))
	fmt.Println("=====================================")
	if options.ConnectOnly {
		return
	}

	if err := runHello(client, options.Hello); err != nil {
		log.Fatal(err)
	}
	if options.Hello && !needsLogin(options) {
		return
	}

	// --------------------------------------------------
	// LOGIN
	// --------------------------------------------------

	if err := client.Login(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Login Successful")
	if options.LoginOnly {
		if err := client.Logout(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Logout Successful")
		return
	}

	if err := runPoll(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainCheck(client, options.CheckDomains); err != nil {
		log.Fatal(err)
	}

	if err := runHostCheck(client, options.HostCheckNames); err != nil {
		log.Fatal(err)
	}

	if err := runHostInfo(client, options.HostInfoName); err != nil {
		log.Fatal(err)
	}

	if err := runHostCreate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runHostUpdate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runHostDelete(client, options.HostDeleteName); err != nil {
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

	if err := runContactTransfer(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainCreate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainUpdate(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainRenew(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainTransfer(client, options); err != nil {
		log.Fatal(err)
	}

	if err := runDomainDelete(client, options.DomainDeleteName); err != nil {
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

func needsLogin(options cliOptions) bool {
	return options.LoginOnly ||
		options.Poll ||
		options.PollAckID != "" ||
		options.CheckDomains != "" ||
		len(options.HostCheckNames) > 0 ||
		options.HostInfoName != "" ||
		options.HostCreateName != "" ||
		options.HostUpdateName != "" ||
		options.HostDeleteName != "" ||
		len(options.ContactCheckIDs) > 0 ||
		options.ContactInfoID != "" ||
		options.ContactCreateID != "" ||
		options.ContactUpdateID != "" ||
		options.ContactDeleteID != "" ||
		options.ContactTransferID != "" ||
		options.CreateDomain != "" ||
		options.DomainUpdateName != "" ||
		options.DomainRenewName != "" ||
		options.DomainTransferName != "" ||
		options.DomainDeleteName != "" ||
		options.InfoDomain != ""
}
