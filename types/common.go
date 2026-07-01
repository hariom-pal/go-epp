package types

type Response struct {
	ResultCode int
	ResultMsg  string

	ClientTRID string
	ServerTRID string
}

type PostalInfo struct {
	Type string

	Name         string
	Organization string

	Street        []string
	City          string
	StateProvince string
	PostalCode    string
	CountryCode   string
}

type Phone struct {
	Number    string
	Extension string
}

type HostAddress struct {
	IPVersion string
	Address   string
}

type Period struct {
	Value int
	Unit  string
}
