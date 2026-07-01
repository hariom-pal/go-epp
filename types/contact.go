package types

import "time"

//
// ============================================================
// CONTACT CHECK
// ============================================================
//

// ContactCheckRequest checks availability for one or more contact IDs.
type ContactCheckRequest struct {
	IDs []string
}

// ContactCheckResponse contains the response for a contact check command.
type ContactCheckResponse struct {
	Response

	Results []ContactCheckResult
}

// ContactCheckResult contains availability data for one contact ID.
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

// ContactInfoRequest requests detailed information for a contact.
type ContactInfoRequest struct {
	ContactID string
}

// ContactInfoResponse contains the response for a contact info command.
type ContactInfoResponse struct {
	Response

	Contact ContactInfo
}

// ContactInfo contains RFC5733 contact information.
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

// ContactCreateRequest creates a contact.
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

// ContactDisclosure contains contact disclose preferences.
type ContactDisclosure struct {
	Flag bool

	Name         ContactDisclosurePostal
	Organization ContactDisclosurePostal
	Address      ContactDisclosurePostal

	Voice bool
	Fax   bool
	Email bool
}

// ContactDisclosurePostal contains disclose preferences for postal fields.
type ContactDisclosurePostal struct {
	International bool
	Localized     bool
}

// ContactCreateResponse contains the response for a contact create command.
type ContactCreateResponse struct {
	Response

	Result ContactCreateResult
}

// ContactCreateResult contains contact create result data.
type ContactCreateResult struct {
	ContactID string

	CreatedDate time.Time
}

//
// ============================================================
// CONTACT UPDATE
// ============================================================
//

// ContactUpdateRequest updates contact status, postal, phone, email, authInfo, or disclosure data.
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

// ContactUpdateResponse contains the response for a contact update command.
type ContactUpdateResponse struct {
	Response
}

//
// ============================================================
// CONTACT DELETE
// ============================================================
//

// ContactDeleteRequest deletes a contact.
type ContactDeleteRequest struct {
	ContactID string
}

// ContactDeleteResponse contains the response for a contact delete command.
type ContactDeleteResponse struct {
	Response
}
