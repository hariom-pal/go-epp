package rgp

import "github.com/hariom-pal/go-epp/constants"

// Namespace is the rgp-1.0 XML namespace.
const Namespace = constants.RGPNamespace

const (
	// OperationRequest requests restoration without a restore report.
	OperationRequest = "request"

	// OperationReport submits a restore report.
	OperationReport = "report"
)

const (
	StatusAddPeriod       = "addPeriod"
	StatusAutoRenewPeriod = "autoRenewPeriod"
	StatusRenewPeriod     = "renewPeriod"
	StatusTransferPeriod  = "transferPeriod"
	StatusPendingDelete   = "pendingDelete"
	StatusPendingRestore  = "pendingRestore"
	StatusRedemption      = "redemptionPeriod"
)

// UpdateRequest contains an RFC3915 restore update extension.
type UpdateRequest struct {
	Restore *Restore
}

// Restore contains the restore operation.
type Restore struct {
	Operation string
	Report    *RestoreReport
}

// RestoreReport contains the RFC3915 restore report data.
type RestoreReport struct {
	PreData       string
	PreDataXML    string
	PostData      string
	PostDataXML   string
	DeleteTime    string
	RestoreTime   string
	RestoreReason Text
	Statements    []Text
	Other         string
	OtherXML      string
}

// Text contains report text and an optional language tag.
type Text struct {
	Lang     string
	Value    string
	ValueXML string
}

// InfoData contains RGP statuses returned by a domain info response.
type InfoData struct {
	Statuses []Status
}

// UpdateData contains RGP statuses returned by a domain update response.
type UpdateData struct {
	Statuses []Status
}

// Status contains a grace period status and optional human-readable text.
type Status struct {
	Status string
	Lang   string
	Text   string
}
