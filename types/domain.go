package types

import "time"

//
// ============================================================
// DOMAIN CHECK
// ============================================================
//

type DomainCheckRequest struct {
	Domains []string
}

type DomainCheckResponse struct {
	Response

	Results []DomainCheckResult
}

type DomainCheckResult struct {
	Domain string
	ASCII  string

	Available bool
	Reason    string
}

//
// ============================================================
// DOMAIN INFO
// ============================================================
//

type DomainInfoRequest struct {
	Domain string
	Hosts  string
}

type DomainInfoResponse struct {
	Response

	Result DomainInfoResult
}

type DomainInfoResult struct {
	Domain string
	ASCII  string

	ROID       string
	Registrant string
	Registrar  string
	CreatedBy  string
	UpdatedBy  string

	Statuses       []string
	StatusDetails  []DomainStatus
	Contacts       []DomainContact
	NameServers    []string
	NameServerInfo []DomainNameServer
	AuthInfo       string
	AuthInfoROID   string
	RGPStatuses    []string
	DNSSEC         DomainDNSSECInfo
	Fee            DomainFeeInfo
	Launch         DomainLaunchInfo
	IDN            DomainIDNInfo

	CreatedDate  *time.Time
	UpdatedDate  *time.Time
	ExpiryDate   *time.Time
	TransferDate *time.Time
}

type DomainContact struct {
	Type string
	ID   string
}

type DomainStatus struct {
	Status string
	Lang   string
	Text   string
}

type DomainNameServer struct {
	HostName  string
	Addresses []DomainHostAddress
}

type DomainHostAddress struct {
	IP      string
	Version string
}

type DomainDNSSECInfo struct {
	MaxSigLife int
	DSData     []DomainDSData
	KeyData    []DomainKeyData
}

type DomainDSData struct {
	KeyTag     int
	Algorithm  int
	DigestType int
	Digest     string
	KeyData    *DomainKeyData
}

type DomainKeyData struct {
	Flags     int
	Protocol  int
	Algorithm int
	PublicKey string
}

type DomainFeeInfo struct {
	Currency    string
	Commands    []DomainFeeCommand
	Fees        []DomainFeeAmount
	Credits     []DomainFeeAmount
	Balance     string
	CreditLimit string
}

type DomainFeeCommand struct {
	Name     string
	Phase    string
	Subphase string
}

type DomainFeeAmount struct {
	Amount      string
	Description string
	Refundable  string
	GracePeriod string
}

type DomainLaunchInfo struct {
	Phase         string
	ApplicationID string
	Status        string
	StatusText    string
	StatusLang    string
}

type DomainIDNInfo struct {
	Table string
}

//
// ============================================================
// DOMAIN CREATE
// ============================================================
//

type DomainCreateRequest struct {
	Domain string

	Period int
	Unit   string

	Registrant string

	AdminContact   string
	TechContact    string
	BillingContact string

	AdminContacts   []string
	TechContacts    []string
	BillingContacts []string
	Contacts        []DomainContact

	NameServers    []string
	NameServerInfo []DomainNameServer

	AuthInfo string
}

type DomainCreateResponse struct {
	Response

	Result DomainCreateResult
}

type DomainCreateResult struct {
	Domain string

	CreatedDate time.Time
	ExpiryDate  time.Time
}

//
// ============================================================
// DOMAIN UPDATE
// ============================================================
//

type DomainUpdateRequest struct {
	Domain string

	NewRegistrant string

	AddStatuses    []string
	RemoveStatuses []string

	AddNameServers    []string
	RemoveNameServers []string

	AddAdminContacts    []string
	RemoveAdminContacts []string

	AddTechContacts    []string
	RemoveTechContacts []string

	AddBillingContacts    []string
	RemoveBillingContacts []string
}

type DomainUpdateResponse struct {
	Response

	Result DomainUpdateResult
}

type DomainUpdateResult struct {
	Domain string
}

//
// ============================================================
// DOMAIN DELETE
// ============================================================
//

type DomainDeleteRequest struct {
	Domain string
}

type DomainDeleteResponse struct {
	Response

	Result DomainDeleteResult
}

type DomainDeleteResult struct {
	Domain string
}

//
// ============================================================
// DOMAIN RENEW
// ============================================================
//

type DomainRenewRequest struct {
	Domain string

	CurrentExpiryDate time.Time

	Period int
}

type DomainRenewResponse struct {
	Response

	Result DomainRenewResult
}

type DomainRenewResult struct {
	Domain string

	NewExpiryDate *time.Time
}

//
// ============================================================
// DOMAIN TRANSFER
// ============================================================
//

type DomainTransferRequest struct {
	Domain string

	Operation string

	Period int

	AuthInfo string
}

type DomainTransferResponse struct {
	Response

	Result DomainTransferResult
}

type DomainTransferResult struct {
	Domain string

	Status string
}

//
// ============================================================
// DOMAIN RESTORE
// ============================================================
//

type DomainRestoreRequest struct {
	Domain string

	Reason string
}

type DomainRestoreResponse struct {
	Response

	Result DomainRestoreResult
}

type DomainRestoreResult struct {
	Domain string
}
