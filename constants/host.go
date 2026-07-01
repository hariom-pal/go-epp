package constants

const (
	// HostsAll requests all delegated and subordinate hosts in a domain info response.
	HostsAll = "all"

	// HostsDel requests delegated hosts in a domain info response.
	HostsDel = "del"

	// HostsSub requests subordinate hosts in a domain info response.
	HostsSub = "sub"

	// HostsNone requests no host data in a domain info response.
	HostsNone = "none"

	// HostsDefault is the default domain info hosts selector.
	HostsDefault = HostsAll
)

// IsHostsValue reports whether value is a valid RFC5731 domain info hosts selector.
func IsHostsValue(value string) bool {
	switch value {
	case HostsAll,
		HostsDel,
		HostsSub,
		HostsNone:
		return true
	default:
		return false
	}
}
