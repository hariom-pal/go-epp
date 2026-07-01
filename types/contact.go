package types

import "time"

//
// ============================================================
// CONTACT CHECK
// ============================================================
//

type ContactCheckRequest struct {
	IDs []string
}

type ContactCheckResponse struct {
	Response

	Results []ContactCheckResult
}

type ContactCheckResult struct {
	ContactID string
	ID        string

	Available bool
	Reason    string
}

//
// ============================================================
// CONTACT INFO
// ============================================================
//

type ContactInfoRequest struct {
	ContactID string
}

type ContactInfoResponse struct {
	Response

	Contact ContactInfo
}

type ContactInfo struct {
	ContactID string
	ROID      string

	Statuses []string

	InternationalPostalInfo *PostalInfo
	LocalizedPostalInfo     *PostalInfo

	Voice Phone
	Fax   Phone
	Email string

	ClientID  string
	CreatedBy string
	UpdatedBy string

	CreatedDate  time.Time
	UpdatedDate  time.Time
	TransferDate time.Time

	AuthInfo string
}

//
// ============================================================
// CONTACT CREATE
// ============================================================
//

type ContactCreateRequest struct {
	ContactID string

	InternationalPostalInfo *PostalInfo
	LocalizedPostalInfo     *PostalInfo

	Voice Phone
	Fax   Phone
	Email string

	AuthInfo string

	Disclosure *ContactDisclosure
}

type ContactDisclosure struct {
	Flag bool

	Name         ContactDisclosurePostal
	Organization ContactDisclosurePostal
	Address      ContactDisclosurePostal

	Voice bool
	Fax   bool
	Email bool
}

type ContactDisclosurePostal struct {
	International bool
	Localized     bool
}

type ContactCreateResponse struct {
	Response

	Result ContactCreateResult
}

type ContactCreateResult struct {
	ContactID string

	CreatedDate time.Time
}

//
// ============================================================
// CONTACT UPDATE
// ============================================================
//

type ContactUpdateRequest struct {
	ContactID string

	AddStatuses    []string
	RemoveStatuses []string

	InternationalPostalInfo *PostalInfo
	LocalizedPostalInfo     *PostalInfo

	Voice *Phone
	Fax   *Phone
	Email string

	AuthInfo string

	Disclosure *ContactDisclosure
}

type ContactUpdateResponse struct {
	Response
}
