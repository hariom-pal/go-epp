package constants

const (
	// HostStatusOK indicates a host has no pending operations or prohibitions.
	HostStatusOK = "ok"

	// HostStatusLinked indicates a host is linked to another object.
	HostStatusLinked = "linked"

	// HostStatusPendingCreate indicates a requested host create is pending.
	HostStatusPendingCreate = "pendingCreate"

	// HostStatusPendingDelete indicates a requested host delete is pending.
	HostStatusPendingDelete = "pendingDelete"

	// HostStatusPendingUpdate indicates a requested host update is pending.
	HostStatusPendingUpdate = "pendingUpdate"

	// HostStatusClientDeleteProhibited indicates the client prohibits host deletion.
	HostStatusClientDeleteProhibited = "clientDeleteProhibited"

	// HostStatusClientUpdateProhibited indicates the client prohibits host update.
	HostStatusClientUpdateProhibited = "clientUpdateProhibited"
)

// IsHostStatus reports whether status is a known RFC5732 host status value.
func IsHostStatus(status string) bool {
	switch status {
	case HostStatusOK,
		HostStatusLinked,
		HostStatusPendingCreate,
		HostStatusPendingDelete,
		HostStatusPendingUpdate,
		HostStatusClientDeleteProhibited,
		HostStatusClientUpdateProhibited:
		return true
	default:
		return false
	}
}
