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

	Statuses    []string
	NameServers []string

	CreatedDate *time.Time
	UpdatedDate *time.Time
	ExpiryDate  *time.Time
}

//
// ============================================================
// DOMAIN CREATE
// ============================================================
//

type DomainCreateRequest struct {
	Domain string

	Period int

	Registrant string

	AdminContact   string
	TechContact    string
	BillingContact string

	NameServers []string

	AuthInfo string
}

type DomainCreateResponse struct {
	Response

	Result DomainCreateResult
}

type DomainCreateResult struct {
	Domain string
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
