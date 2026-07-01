package constants

const (
	// ContactStatusOK indicates a contact has no pending operations or prohibitions.
	ContactStatusOK = "ok"

	// ContactStatusLinked indicates a contact is linked to another object.
	ContactStatusLinked = "linked"

	// ContactStatusPendingCreate indicates a requested contact create is pending.
	ContactStatusPendingCreate = "pendingCreate"

	// ContactStatusPendingDelete indicates a requested contact delete is pending.
	ContactStatusPendingDelete = "pendingDelete"

	// ContactStatusPendingTransfer indicates a requested contact transfer is pending.
	ContactStatusPendingTransfer = "pendingTransfer"

	// ContactStatusPendingUpdate indicates a requested contact update is pending.
	ContactStatusPendingUpdate = "pendingUpdate"

	// ContactStatusClientDeleteProhibited indicates the client prohibits contact deletion.
	ContactStatusClientDeleteProhibited = "clientDeleteProhibited"

	// ContactStatusClientTransferProhibited indicates the client prohibits contact transfer.
	ContactStatusClientTransferProhibited = "clientTransferProhibited"

	// ContactStatusClientUpdateProhibited indicates the client prohibits contact update.
	ContactStatusClientUpdateProhibited = "clientUpdateProhibited"
)

// IsContactStatus reports whether status is a known RFC5733 contact status value.
func IsContactStatus(status string) bool {
	switch status {
	case ContactStatusOK,
		ContactStatusLinked,
		ContactStatusPendingCreate,
		ContactStatusPendingDelete,
		ContactStatusPendingTransfer,
		ContactStatusPendingUpdate,
		ContactStatusClientDeleteProhibited,
		ContactStatusClientTransferProhibited,
		ContactStatusClientUpdateProhibited:
		return true
	default:
		return false
	}
}
