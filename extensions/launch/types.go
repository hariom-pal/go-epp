package launch

import "github.com/hariom-pal/go-epp/constants"

// Namespace is the launch-1.0 XML namespace.
const Namespace = constants.LaunchNamespace

const (
	// PhaseSunrise identifies the trademark-holder sunrise launch phase.
	PhaseSunrise = "sunrise"

	// PhaseLandrush identifies a high-volume pre-open launch phase.
	PhaseLandrush = "landrush"

	// PhaseClaims identifies the trademark claims launch phase.
	PhaseClaims = "claims"

	// PhaseOpen identifies the steady-state/open launch phase.
	PhaseOpen = "open"

	// PhaseCustom identifies a server-defined launch phase using Phase.Name.
	PhaseCustom = "custom"
)

const (
	// ObjectApplication requests creation of a launch application.
	ObjectApplication = "application"

	// ObjectRegistration requests creation of a launch registration.
	ObjectRegistration = "registration"
)

const (
	// StatusPendingValidation indicates that application validation is pending.
	StatusPendingValidation = "pendingValidation"

	// StatusValidated indicates that an application was validated.
	StatusValidated = "validated"

	// StatusInvalid indicates that an application failed validation.
	StatusInvalid = "invalid"

	// StatusPendingAllocation indicates that allocation is pending.
	StatusPendingAllocation = "pendingAllocation"

	// StatusAllocated indicates that an application was allocated.
	StatusAllocated = "allocated"

	// StatusRejected indicates that an application was rejected.
	StatusRejected = "rejected"

	// StatusCustom identifies a server-defined launch status using Status.Name.
	StatusCustom = "custom"
)

// Phase identifies a launch phase and optional sub-phase/name.
type Phase struct {
	Value string
	Name  string
}

// Status contains launch application status metadata.
type Status struct {
	Status string
	Name   string
	Lang   string
	Text   string
}

// Validator identifies a mark or claims validator.
type Validator struct {
	ID string
}

// Code contains a launch validation code.
type Code struct {
	Value       string
	ValidatorID string
}

// CodeMark contains a launch code and optional mark XML.
type CodeMark struct {
	Code    *Code
	MarkXML string
}

// Notice contains claims notice acknowledgement data.
type Notice struct {
	ID           string
	ValidatorID  string
	NotAfter     string
	AcceptedDate string
}

// CreateRequest contains launch create extension data.
type CreateRequest struct {
	Phase                Phase
	Type                 string
	CodeMarks            []CodeMark
	SignedMarkXML        []string
	EncodedSignedMarkXML []string
	Notices              []Notice
}

// InfoRequest contains launch info extension data.
type InfoRequest struct {
	Phase         Phase
	ApplicationID string
	IncludeMark   bool
}

// UpdateRequest contains launch update extension data.
type UpdateRequest struct {
	Phase         Phase
	ApplicationID string
}

// DeleteRequest contains launch delete extension data.
type DeleteRequest struct {
	Phase         Phase
	ApplicationID string
}

// IDData contains launch phase and application identifier response data.
type IDData struct {
	Phase         Phase
	ApplicationID string
}

// InfoData contains launch info response data.
type InfoData struct {
	Phase         Phase
	ApplicationID string
	Status        Status
	MarkXML       string
	RawXML        string
}
