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
