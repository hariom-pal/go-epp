package types

import (
	"time"

	"github.com/hariom-pal/go-epp/extensions/fee"
)

//
// ============================================================
// DOMAIN CHECK
// ============================================================
//

// DomainCheckRequest checks availability for one or more domain names.
type DomainCheckRequest struct {
	Domains []string

	Fee *fee.CheckRequest
}

// DomainCheckResponse contains the response for a domain check command.
type DomainCheckResponse struct {
	Response

	Results []DomainCheckResult
	Fee     fee.CheckData
}

// DomainCheckResult contains availability data for one domain.
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

// DomainInfoRequest requests detailed information for a domain.
type DomainInfoRequest struct {
	Domain string
	Hosts  string
}

// DomainInfoResponse contains the response for a domain info command.
type DomainInfoResponse struct {
	Response

	Result DomainInfoResult
}

// DomainInfoResult contains RFC5731 domain information and supported extension data.
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

// DomainContact identifies a domain contact and its contact type.
type DomainContact struct {
	Type string
	ID   string
}

// DomainStatus contains a domain status value and optional display metadata.
type DomainStatus struct {
	Status string
	Lang   string
	Text   string
}

// DomainNameServer contains a domain name server and optional host attributes.
type DomainNameServer struct {
	HostName  string
	Addresses []DomainHostAddress
}

// DomainHostAddress contains an address for a domain hostAttr name server.
type DomainHostAddress struct {
	IP      string
	Version string
}

// DomainDNSSECInfo contains secDNS data returned for a domain.
type DomainDNSSECInfo struct {
	MaxSigLife int
	DSData     []DomainDSData
	KeyData    []DomainKeyData
}

// DomainDSData contains DNSSEC DS data for a domain.
type DomainDSData struct {
	KeyTag     int
	Algorithm  int
	DigestType int
	Digest     string
	KeyData    *DomainKeyData
}

// DomainKeyData contains DNSSEC key data for a domain.
type DomainKeyData struct {
	Flags     int
	Protocol  int
	Algorithm int
	PublicKey string
}

// DomainFeeInfo contains fee extension data returned for a domain.
type DomainFeeInfo struct {
	Currency    string
	Commands    []DomainFeeCommand
	Fees        []DomainFeeAmount
	Credits     []DomainFeeAmount
	Balance     string
	CreditLimit string
}

// DomainFeeCommand contains fee information for a single command.
type DomainFeeCommand struct {
	Name     string
	Phase    string
	Subphase string
}

// DomainFeeAmount contains a fee or credit amount with registry metadata.
type DomainFeeAmount struct {
	Amount      string
	Description string
	Refundable  string
	GracePeriod string
}

// DomainLaunchInfo contains launch extension data returned for a domain.
type DomainLaunchInfo struct {
	Phase         string
	ApplicationID string
	Status        string
	StatusText    string
	StatusLang    string
}

// DomainIDNInfo contains IDN extension data returned for a domain.
type DomainIDNInfo struct {
	Table string
}

//
// ============================================================
// DOMAIN CREATE
// ============================================================
//

// DomainCreateRequest creates a domain.
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

	Fee *fee.TransformRequest
}

// DomainCreateResponse contains the response for a domain create command.
type DomainCreateResponse struct {
	Response

	Result DomainCreateResult
}

// DomainCreateResult contains domain create result data.
type DomainCreateResult struct {
	Domain string

	CreatedDate time.Time
	ExpiryDate  time.Time

	Fee fee.TransformData
}

//
// ============================================================
// DOMAIN UPDATE
// ============================================================
//

// DomainUpdateRequest updates a domain using add, remove, and change sections.
type DomainUpdateRequest struct {
	Domain     string
	DomainName string

	Registrant    string
	NewRegistrant string
	AuthInfo      string

	AddStatuses      []string
	RemoveStatuses   []string
	AddStatusInfo    []DomainStatus
	RemoveStatusInfo []DomainStatus

	AddNameServers       []string
	RemoveNameServers    []string
	AddNameServerInfo    []DomainNameServer
	RemoveNameServerInfo []DomainNameServer

	AddContacts    []DomainContact
	RemoveContacts []DomainContact

	AddAdminContacts    []string
	RemoveAdminContacts []string

	AddTechContacts    []string
	RemoveTechContacts []string

	AddBillingContacts    []string
	RemoveBillingContacts []string
}

// DomainUpdateResponse contains the response for a domain update command.
type DomainUpdateResponse struct {
	Response

	Result DomainUpdateResult
}

// DomainUpdateResult contains domain update result data.
type DomainUpdateResult struct {
	Domain string
}

//
// ============================================================
// DOMAIN DELETE
// ============================================================
//

// DomainDeleteRequest deletes a domain.
type DomainDeleteRequest struct {
	Domain     string
	DomainName string
}

// DomainDeleteResponse contains the response for a domain delete command.
type DomainDeleteResponse struct {
	Response

	Result DomainDeleteResult
}

// DomainDeleteResult contains domain delete result data.
type DomainDeleteResult struct {
	Domain     string
	DomainName string
}

//
// ============================================================
// DOMAIN RENEW
// ============================================================
//

// DomainRenewRequest renews a domain for an additional registration period.
type DomainRenewRequest struct {
	Domain     string
	DomainName string

	CurrentExpiryDate time.Time

	Period     int
	Unit       string
	PeriodInfo Period

	Fee *fee.TransformRequest
}

// DomainRenewResponse contains the response for a domain renew command.
type DomainRenewResponse struct {
	Response

	Result DomainRenewResult
}

// DomainRenewResult contains domain renew result data.
type DomainRenewResult struct {
	Domain     string
	DomainName string

	NewExpiryDate time.Time

	Fee fee.TransformData
}

//
// ============================================================
// DOMAIN TRANSFER
// ============================================================
//

// DomainTransferRequest performs a domain transfer operation.
type DomainTransferRequest struct {
	Domain     string
	DomainName string

	Operation string

	Period     int
	Unit       string
	PeriodInfo Period

	AuthInfo string

	Fee *fee.TransformRequest
}

// DomainTransferResponse contains the response for a domain transfer command.
type DomainTransferResponse struct {
	Response

	TransferData DomainTransferData
	Result       DomainTransferResult
}

// DomainTransferData contains domain-specific transfer response data.
type DomainTransferData struct {
	TransferData

	DomainName string
	Fee        fee.TransformData
}

// DomainTransferResult contains compatibility fields for domain transfer responses.
type DomainTransferResult struct {
	DomainTransferData

	Domain string
	Status string
}

//
// ============================================================
// DOMAIN RESTORE
// ============================================================
//

// DomainRestoreRequest contains fields for a future domain restore command.
type DomainRestoreRequest struct {
	Domain string

	Reason string
}

// DomainRestoreResponse contains fields for a future domain restore response.
type DomainRestoreResponse struct {
	Response

	Result DomainRestoreResult
}

// DomainRestoreResult contains future domain restore result data.
type DomainRestoreResult struct {
	Domain string
}
