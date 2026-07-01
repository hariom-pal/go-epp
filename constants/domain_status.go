package constants

const (
	// StatusOK indicates a domain has no pending operations or prohibitions.
	StatusOK = "ok"

	// StatusInactive indicates a domain has no delegated hosts.
	StatusInactive = "inactive"

	// StatusPendingCreate indicates a requested domain create is pending.
	StatusPendingCreate = "pendingCreate"

	// StatusPendingDelete indicates a requested domain delete is pending.
	StatusPendingDelete = "pendingDelete"

	// StatusPendingRenew indicates a requested domain renew is pending.
	StatusPendingRenew = "pendingRenew"

	// StatusPendingTransfer indicates a requested domain transfer is pending.
	StatusPendingTransfer = "pendingTransfer"

	// StatusPendingUpdate indicates a requested domain update is pending.
	StatusPendingUpdate = "pendingUpdate"

	// StatusClientHold indicates the client has placed the domain on hold.
	StatusClientHold = "clientHold"

	// StatusClientDeleteProhibited indicates the client prohibits domain deletion.
	StatusClientDeleteProhibited = "clientDeleteProhibited"

	// StatusClientRenewProhibited indicates the client prohibits domain renewal.
	StatusClientRenewProhibited = "clientRenewProhibited"

	// StatusClientTransferProhibited indicates the client prohibits domain transfer.
	StatusClientTransferProhibited = "clientTransferProhibited"

	// StatusClientUpdateProhibited indicates the client prohibits domain update.
	StatusClientUpdateProhibited = "clientUpdateProhibited"

	// StatusServerHold indicates the server has placed the domain on hold.
	StatusServerHold = "serverHold"

	// StatusServerDeleteProhibited indicates the server prohibits domain deletion.
	StatusServerDeleteProhibited = "serverDeleteProhibited"

	// StatusServerRenewProhibited indicates the server prohibits domain renewal.
	StatusServerRenewProhibited = "serverRenewProhibited"

	// StatusServerTransferProhibited indicates the server prohibits domain transfer.
	StatusServerTransferProhibited = "serverTransferProhibited"

	// StatusServerUpdateProhibited indicates the server prohibits domain update.
	StatusServerUpdateProhibited = "serverUpdateProhibited"
)

// IsDomainStatus reports whether status is a known RFC5731 domain status value.
func IsDomainStatus(status string) bool {
	switch status {
	case StatusOK,
		StatusInactive,
		StatusPendingCreate,
		StatusPendingDelete,
		StatusPendingRenew,
		StatusPendingTransfer,
		StatusPendingUpdate,
		StatusClientHold,
		StatusClientDeleteProhibited,
		StatusClientRenewProhibited,
		StatusClientTransferProhibited,
		StatusClientUpdateProhibited,
		StatusServerHold,
		StatusServerDeleteProhibited,
		StatusServerRenewProhibited,
		StatusServerTransferProhibited,
		StatusServerUpdateProhibited:
		return true
	default:
		return false
	}
}
