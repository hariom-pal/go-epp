package validation

import (
	"fmt"
	"net"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ValidateHostCheck validates a host check request.
func ValidateHostCheck(req types.HostCheckRequest) error {
	if len(req.Hosts) == 0 {
		return required("host name")
	}
	return validateDomainNames("host name", req.Hosts)
}

// ValidateHostInfo validates a host info request.
func ValidateHostInfo(req types.HostInfoRequest) error {
	return validateDomainName("host name", req.HostName)
}

// ValidateHostCreate validates a host create request.
func ValidateHostCreate(req types.HostCreateRequest) error {
	if err := validateDomainName("host name", req.HostName); err != nil {
		return err
	}
	return validateHostAddresses(req.Addresses)
}

// ValidateHostUpdate validates a host update request.
func ValidateHostUpdate(req types.HostUpdateRequest) error {
	if err := validateDomainName("host name", req.HostName); err != nil {
		return err
	}
	if err := validateHostAddresses(req.AddAddresses); err != nil {
		return err
	}
	if err := validateHostAddresses(req.RemoveAddresses); err != nil {
		return err
	}
	if err := validateHostStatuses(req.AddStatuses); err != nil {
		return err
	}
	if err := validateHostStatuses(req.RemoveStatuses); err != nil {
		return err
	}
	if strings.TrimSpace(req.NewHostName) != "" {
		if err := validateDomainName("new host name", req.NewHostName); err != nil {
			return err
		}
	}
	if len(req.AddAddresses) == 0 &&
		len(req.RemoveAddresses) == 0 &&
		len(req.AddStatuses) == 0 &&
		len(req.RemoveStatuses) == 0 &&
		strings.TrimSpace(req.NewHostName) == "" {
		return fmt.Errorf("at least one host update operation is required")
	}
	return nil
}

// ValidateHostDelete validates a host delete request.
func ValidateHostDelete(req types.HostDeleteRequest) error {
	return validateDomainName("host name", req.HostName)
}

func validateHostAddresses(addresses []types.HostAddress) error {
	for _, address := range addresses {
		version := strings.ToLower(strings.TrimSpace(address.IPVersion))
		if version != "v4" && version != "v6" {
			return fmt.Errorf("host address IP version must be %q or %q", "v4", "v6")
		}
		ip := strings.TrimSpace(address.Address)
		parsed := net.ParseIP(ip)
		if parsed == nil {
			return fmt.Errorf("host address is invalid")
		}
		if version == "v4" && parsed.To4() == nil {
			return fmt.Errorf("host address must be IPv4")
		}
		if version == "v6" && parsed.To4() != nil {
			return fmt.Errorf("host address must be IPv6")
		}
	}
	return nil
}

func validateHostStatuses(statuses []string) error {
	for _, status := range statuses {
		status = strings.TrimSpace(status)
		if status == "" {
			return required("host status")
		}
		if !constants.IsHostStatus(status) {
			return fmt.Errorf("host status is invalid")
		}
	}
	return nil
}
