package constants

const (
	// TransferQuery queries the current transfer state.
	TransferQuery = "query"

	// TransferRequest requests a transfer.
	TransferRequest = "request"

	// TransferApprove approves a pending transfer.
	TransferApprove = "approve"

	// TransferReject rejects a pending transfer.
	TransferReject = "reject"

	// TransferCancel cancels a pending transfer request.
	TransferCancel = "cancel"
)

// IsTransferOperation reports whether operation is a valid EPP transfer op value.
func IsTransferOperation(
	operation string,
) bool {

	switch operation {
	case TransferQuery,
		TransferRequest,
		TransferApprove,
		TransferReject,
		TransferCancel:

		return true
	default:
		return false
	}
}
