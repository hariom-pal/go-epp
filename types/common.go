package types

type Response struct {
	ResultCode int
	ResultMsg  string

	ClientTRID string
	ServerTRID string
}
