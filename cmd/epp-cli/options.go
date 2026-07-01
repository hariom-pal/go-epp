package main

import "flag"

type cliOptions struct {
	ConfigPath string

	CheckDomains string

	HostCheckNames []string

	ContactCheckIDs []string
	ContactInfoID   string

	ContactCreateID             string
	ContactCreateName           string
	ContactCreateOrg            string
	ContactCreateStreets        []string
	ContactCreateCity           string
	ContactCreateState          string
	ContactCreatePostalCode     string
	ContactCreateCountryCode    string
	ContactCreateLocName        string
	ContactCreateLocOrg         string
	ContactCreateLocStreets     []string
	ContactCreateLocCity        string
	ContactCreateLocState       string
	ContactCreateLocPostalCode  string
	ContactCreateLocCountryCode string
	ContactCreateVoice          string
	ContactCreateVoiceExt       string
	ContactCreateFax            string
	ContactCreateFaxExt         string
	ContactCreateEmail          string
	ContactCreateAuthInfo       string

	ContactUpdateID             string
	ContactUpdateAddStatuses    []string
	ContactUpdateRemoveStatuses []string

	ContactDeleteID string

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
	var hostCheckNames stringList
	var contactCheckIDs stringList
	var contactCreateStreets stringList
	var contactCreateLocStreets stringList
	var contactUpdateAddStatuses stringList
	var contactUpdateRemoveStatuses stringList

	flag.StringVar(&options.ConfigPath, "config", "configs/config.yaml", "path to config YAML")
	flag.StringVar(&options.CheckDomains, "check", "", "comma-separated domains for domain check")
	flag.StringVar(&options.ContactInfoID, "contact-info", "", "contact ID for contact info")
	flag.StringVar(&options.ContactCreateID, "contact-create", "", "contact ID for contact create")
	flag.StringVar(&options.ContactCreateName, "contact-name", "", "international contact name for contact create")
	flag.StringVar(&options.ContactCreateOrg, "contact-org", "", "international contact organization for contact create")
	flag.StringVar(&options.ContactCreateCity, "contact-city", "", "international contact city for contact create")
	flag.StringVar(&options.ContactCreateState, "contact-state", "", "international contact state/province for contact create")
	flag.StringVar(&options.ContactCreatePostalCode, "contact-pc", "", "international contact postal code for contact create")
	flag.StringVar(&options.ContactCreateCountryCode, "contact-cc", "", "international contact country code for contact create")
	flag.StringVar(&options.ContactCreateLocName, "contact-loc-name", "", "localized contact name for contact create")
	flag.StringVar(&options.ContactCreateLocOrg, "contact-loc-org", "", "localized contact organization for contact create")
	flag.StringVar(&options.ContactCreateLocCity, "contact-loc-city", "", "localized contact city for contact create")
	flag.StringVar(&options.ContactCreateLocState, "contact-loc-state", "", "localized contact state/province for contact create")
	flag.StringVar(&options.ContactCreateLocPostalCode, "contact-loc-pc", "", "localized contact postal code for contact create")
	flag.StringVar(&options.ContactCreateLocCountryCode, "contact-loc-cc", "", "localized contact country code for contact create")
	flag.StringVar(&options.ContactCreateVoice, "contact-voice", "", "voice number for contact create")
	flag.StringVar(&options.ContactCreateVoiceExt, "contact-voice-ext", "", "voice extension for contact create")
	flag.StringVar(&options.ContactCreateFax, "contact-fax", "", "fax number for contact create")
	flag.StringVar(&options.ContactCreateFaxExt, "contact-fax-ext", "", "fax extension for contact create")
	flag.StringVar(&options.ContactCreateEmail, "contact-email", "", "email for contact create")
	flag.StringVar(&options.ContactCreateAuthInfo, "contact-authInfo", "", "authInfo password for contact create")
	flag.StringVar(&options.ContactUpdateID, "contact-update", "", "contact ID for contact update")
	flag.StringVar(&options.ContactDeleteID, "contact-delete", "", "contact ID for contact delete")
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
	flag.Var(&hostCheckNames, "host-check", "host name for host check; may be repeated")
	flag.Var(&contactCheckIDs, "contact-check", "contact ID for contact check; may be repeated")
	flag.Var(&contactCreateStreets, "contact-street", "international contact street for contact create; may be repeated")
	flag.Var(&contactCreateLocStreets, "contact-loc-street", "localized contact street for contact create; may be repeated")
	flag.Var(&contactUpdateAddStatuses, "contact-add-status", "contact status to add during contact update; may be repeated")
	flag.Var(&contactUpdateRemoveStatuses, "contact-rem-status", "contact status to remove during contact update; may be repeated")
	flag.Parse()

	options.HostCheckNames = []string(hostCheckNames)
	options.ContactCheckIDs = []string(contactCheckIDs)
	options.ContactCreateStreets = []string(contactCreateStreets)
	options.ContactCreateLocStreets = []string(contactCreateLocStreets)
	options.ContactUpdateAddStatuses = []string(contactUpdateAddStatuses)
	options.ContactUpdateRemoveStatuses = []string(contactUpdateRemoveStatuses)
	options.CreateAdminContacts = []string(createAdminContacts)
	options.CreateTechContacts = []string(createTechContacts)
	options.CreateBillingContacts = []string(createBillingContacts)
	options.CreateNameServers = []string(createNameServers)

	return options
}
