package constants

const (
	HostsDefault = HostsAll
)

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
