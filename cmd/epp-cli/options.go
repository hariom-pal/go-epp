package main

import "flag"

type cliOptions struct {
	ConfigPath string

	CheckDomains string

	ContactCheckIDs []string

	InfoDomain string
	InfoHosts  string

	CreateDomain          string
	CreatePeriod          int
	CreateUnit            string
	CreateRegistrant      string
	CreateAuthInfo        string
	CreateAdminContacts   []string
	CreateTechContacts    []string
	CreateBillingContacts []string
	CreateNameServers     []string
}

func parseOptions() cliOptions {
	var options cliOptions

	var createAdminContacts stringList
	var createTechContacts stringList
	var createBillingContacts stringList
	var createNameServers stringList
	var contactCheckIDs stringList

	flag.StringVar(&options.ConfigPath, "config", "configs/config.yaml", "path to config YAML")
	flag.StringVar(&options.CheckDomains, "check", "", "comma-separated domains for domain check")
	flag.StringVar(&options.InfoDomain, "info", "", "domain for domain info")
	flag.StringVar(&options.InfoHosts, "hosts", "", "domain info hosts value: all, del, sub, or none")
	flag.StringVar(&options.CreateDomain, "create", "", "domain for domain create")
	flag.IntVar(&options.CreatePeriod, "period", 0, "registration period for domain create")
	flag.StringVar(&options.CreateUnit, "unit", "y", "registration period unit for domain create: y or m")
	flag.StringVar(&options.CreateRegistrant, "registrant", "", "registrant contact for domain create")
	flag.StringVar(&options.CreateAuthInfo, "authInfo", "", "authInfo password for domain create")

	flag.Var(&createAdminContacts, "admin", "admin contact for domain create; may be repeated")
	flag.Var(&createTechContacts, "tech", "tech contact for domain create; may be repeated")
	flag.Var(&createBillingContacts, "billing", "billing contact for domain create; may be repeated")
	flag.Var(&createNameServers, "ns", "hostObj name server for domain create; may be repeated")
	flag.Var(&contactCheckIDs, "contact-check", "contact ID for contact check; may be repeated")
	flag.Parse()

	options.ContactCheckIDs = []string(contactCheckIDs)
	options.CreateAdminContacts = []string(createAdminContacts)
	options.CreateTechContacts = []string(createTechContacts)
	options.CreateBillingContacts = []string(createBillingContacts)
	options.CreateNameServers = []string(createNameServers)

	return options
}
