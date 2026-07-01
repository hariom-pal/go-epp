package types

import "time"

//
// ============================================================
// HOST CHECK
// ============================================================
//

// HostCheckRequest checks availability for one or more host names.
type HostCheckRequest struct {
	Hosts []string
}

// HostCheckResponse contains the response for a host check command.
type HostCheckResponse struct {
	Response

	Results []HostCheckResult
}

// HostCheckResult contains availability data for one host.
type HostCheckResult struct {
	HostName  string
	ASCIIName string

	Available bool
	Reason    string
}

//
// ============================================================
// HOST INFO
// ============================================================
//

// HostInfoRequest requests detailed information for a host.
type HostInfoRequest struct {
	HostName string
}

// HostInfoResponse contains the response for a host info command.
type HostInfoResponse struct {
	Response

	Host HostInfo
}

// HostInfo contains RFC5732 host information.
type HostInfo struct {
	HostName  string
	ASCIIName string
	ROID      string

	Statuses  []string
	Addresses []HostAddress

	ClientID  string
	CreatedBy string
	UpdatedBy string

	CreatedDate  time.Time
	UpdatedDate  time.Time
	TransferDate time.Time
}

//
// ============================================================
// HOST CREATE
// ============================================================
//

// HostCreateRequest creates a host.
type HostCreateRequest struct {
	HostName string

	Addresses []HostAddress
}

// HostCreateResponse contains the response for a host create command.
type HostCreateResponse struct {
	Response

	Result HostCreateResult
}

// HostCreateResult contains host create result data.
type HostCreateResult struct {
	HostName string

	CreatedDate time.Time
}

//
// ============================================================
// HOST UPDATE
// ============================================================
//

// HostUpdateRequest updates host addresses, statuses, or the host name.
type HostUpdateRequest struct {
	HostName string

	AddAddresses    []HostAddress
	RemoveAddresses []HostAddress

	AddStatuses    []string
	RemoveStatuses []string

	NewHostName string
}

// HostUpdateResponse contains the response for a host update command.
type HostUpdateResponse struct {
	Response
}

//
// ============================================================
// HOST DELETE
// ============================================================
//

// HostDeleteRequest deletes a host.
type HostDeleteRequest struct {
	HostName string
}

// HostDeleteResponse contains the response for a host delete command.
type HostDeleteResponse struct {
	Response
}
