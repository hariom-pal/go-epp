package types

import "time"

//
// ============================================================
// HOST CHECK
// ============================================================
//

type HostCheckRequest struct {
	Hosts []string
}

type HostCheckResponse struct {
	Response

	Results []HostCheckResult
}

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

type HostInfoRequest struct {
	HostName string
}

type HostInfoResponse struct {
	Response

	Host HostInfo
}

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

type HostCreateRequest struct {
	HostName string

	Addresses []HostAddress
}

type HostCreateResponse struct {
	Response

	Result HostCreateResult
}

type HostCreateResult struct {
	HostName string

	CreatedDate time.Time
}
