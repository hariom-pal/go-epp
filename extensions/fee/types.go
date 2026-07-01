package fee

import "github.com/hariom-pal/go-epp/constants"

const (
	// CommandCreate requests or reports create-command fee data.
	CommandCreate = "create"

	// CommandRenew requests or reports renew-command fee data.
	CommandRenew = "renew"

	// CommandTransfer requests or reports transfer-command fee data.
	CommandTransfer = "transfer"

	// CommandRestore requests or reports restore-command fee data.
	CommandRestore = "restore"

	// AppliedImmediate indicates a fee is applied immediately.
	AppliedImmediate = "immediate"

	// AppliedDelayed indicates a fee is applied later.
	AppliedDelayed = "delayed"
)

// Namespace is the fee-0.7 XML namespace.
const Namespace = constants.FeeNamespace

// Command identifies the command, phase, and subphase for fee calculation.
type Command struct {
	Name     string
	Phase    string
	Subphase string
}

// Period contains the validity period used for fee calculation.
type Period struct {
	Value int
	Unit  string
}

// Amount contains a fee or credit amount with fee-0.7 attributes.
type Amount struct {
	Amount      string
	Description string
	Refundable  string
	GracePeriod string
	Applied     string
}

// Fee is an alias for a debit amount.
type Fee = Amount

// Credit is an alias for a credit amount.
type Credit = Amount

// Classification identifies the fee class of an object.
type Classification string

// CheckDomain requests fee data for one domain name.
type CheckDomain struct {
	Name     string
	Currency string
	Command  Command
	Period   *Period
}

// CheckRequest contains fee-0.7 data for a check command extension.
type CheckRequest struct {
	Domains []CheckDomain
}

// TransformRequest contains fee-0.7 data for transform command extensions.
type TransformRequest struct {
	Currency string
	Fees     []Fee
	Credits  []Credit
}

// CheckData contains parsed fee-0.7 check response data.
type CheckData struct {
	Results []CheckResult
}

// CheckResult contains parsed fee data for one checked domain.
type CheckResult struct {
	Name     string
	Currency string
	Command  Command
	Period   *Period

	Fees    []Fee
	Credits []Credit

	Class  Classification
	Reason string
}

// TransformData contains parsed fee-0.7 transform response data.
type TransformData struct {
	Currency string
	Period   *Period

	Fees    []Fee
	Credits []Credit

	Balance     string
	CreditLimit string
}

// InfoData contains parsed fee-0.7 info response data.
type InfoData struct {
	Currency string
	Command  Command
	Period   *Period

	Fees    []Fee
	Credits []Credit

	Class Classification
}
