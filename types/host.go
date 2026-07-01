package types

//
// ============================================================
// HOST CHECK
// ============================================================
//

type HostCheckRequest struct {
	Hosts []string
}

type HostCheckResponse struct {
	Response

	Results []HostCheckResult
}

type HostCheckResult struct {
	HostName  string
	ASCIIName string

	Available bool
	Reason    string
}
