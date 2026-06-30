package types

//
// ============================================================
// CONTACT CHECK
// ============================================================
//

type ContactCheckRequest struct {
	IDs []string
}

type ContactCheckResponse struct {
	Response

	Results []ContactCheckResult
}

type ContactCheckResult struct {
	ID string

	Available bool
	Reason    string
}
