package constants

const (
	TransferQuery = "query"

	TransferRequest = "request"

	TransferApprove = "approve"

	TransferReject = "reject"

	TransferCancel = "cancel"
)

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
