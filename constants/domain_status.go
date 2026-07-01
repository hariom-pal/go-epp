package constants

const (
	StatusOK = "ok"

	StatusInactive = "inactive"

	StatusPendingCreate   = "pendingCreate"
	StatusPendingDelete   = "pendingDelete"
	StatusPendingRenew    = "pendingRenew"
	StatusPendingTransfer = "pendingTransfer"
	StatusPendingUpdate   = "pendingUpdate"

	StatusClientHold = "clientHold"

	StatusClientDeleteProhibited   = "clientDeleteProhibited"
	StatusClientRenewProhibited    = "clientRenewProhibited"
	StatusClientTransferProhibited = "clientTransferProhibited"
	StatusClientUpdateProhibited   = "clientUpdateProhibited"

	StatusServerHold = "serverHold"

	StatusServerDeleteProhibited   = "serverDeleteProhibited"
	StatusServerRenewProhibited    = "serverRenewProhibited"
	StatusServerTransferProhibited = "serverTransferProhibited"
	StatusServerUpdateProhibited   = "serverUpdateProhibited"
)
