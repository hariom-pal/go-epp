package validation

import (
	"fmt"
	"net"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

const (
	periodUnitYears  = "y"
	periodUnitMonths = "m"
	maxPeriod        = 99
)

// ValidateDomainCheck validates a domain check request.
func ValidateDomainCheck(req types.DomainCheckRequest) error {
	if len(req.Domains) == 0 {
		return required("domain check name")
	}
	if err := validateDomainNames("domain check name", req.Domains); err != nil {
		return err
	}
	return ValidateFeeCheck(req.Fee)
}

// ValidateDomainInfo validates a domain info request.
func ValidateDomainInfo(req types.DomainInfoRequest) error {
	if err := validateDomainName("domain", req.Domain); err != nil {
		return err
	}
	hosts := strings.TrimSpace(req.Hosts)
	if hosts != "" && !constants.IsHostsValue(hosts) {
		return fmt.Errorf("domain info hosts value is invalid")
	}
	return ValidateLaunchInfo(req.Launch)
}

// ValidateDomainCreate validates a domain create request.
func ValidateDomainCreate(req types.DomainCreateRequest) error {
	if err := validateDomainName("domain", req.Domain); err != nil {
		return err
	}
	if err := validatePeriod(req.Period, req.Unit, true); err != nil {
		return err
	}
	if strings.TrimSpace(req.Registrant) == "" {
		return required("registrant contact")
	}
	if strings.TrimSpace(req.AuthInfo) == "" {
		return required("authInfo")
	}
	if err := validateDomainContacts(req.Contacts); err != nil {
		return err
	}
	if err := validateLegacyContactIDs(req.AdminContact, req.TechContact, req.BillingContact); err != nil {
		return err
	}
	if err := validateContactIDs("admin contact", req.AdminContacts); err != nil {
		return err
	}
	if err := validateContactIDs("tech contact", req.TechContacts); err != nil {
		return err
	}
	if err := validateContactIDs("billing contact", req.BillingContacts); err != nil {
		return err
	}
	if err := validateNameServers(req.NameServers, req.NameServerInfo); err != nil {
		return err
	}
	if err := ValidateFeeTransform(req.Fee); err != nil {
		return err
	}
	if err := ValidateDNSSECCreate(req.SecDNS); err != nil {
		return err
	}
	return ValidateLaunchCreate(req.Launch)
}

// ValidateDomainUpdate validates a domain update request.
func ValidateDomainUpdate(req types.DomainUpdateRequest) error {
	if err := validateDomainName("domain", firstNonEmpty(req.DomainName, req.Domain)); err != nil {
		return err
	}
	if err := validateNameServers(req.AddNameServers, req.AddNameServerInfo); err != nil {
		return err
	}
	if err := validateNameServers(req.RemoveNameServers, req.RemoveNameServerInfo); err != nil {
		return err
	}
	if err := validateDomainContacts(req.AddContacts); err != nil {
		return err
	}
	if err := validateDomainContacts(req.RemoveContacts); err != nil {
		return err
	}
	for label, ids := range map[string][]string{
		"add admin contact":      req.AddAdminContacts,
		"remove admin contact":   req.RemoveAdminContacts,
		"add tech contact":       req.AddTechContacts,
		"remove tech contact":    req.RemoveTechContacts,
		"add billing contact":    req.AddBillingContacts,
		"remove billing contact": req.RemoveBillingContacts,
	} {
		if err := validateContactIDs(label, ids); err != nil {
			return err
		}
	}
	if err := validateDomainStatuses(req.AddStatuses, req.AddStatusInfo); err != nil {
		return err
	}
	if err := validateDomainStatuses(req.RemoveStatuses, req.RemoveStatusInfo); err != nil {
		return err
	}
	if err := ValidateDNSSECUpdate(req.SecDNS); err != nil {
		return err
	}
	if err := ValidateRGPUpdate(req.RGP); err != nil {
		return err
	}
	if err := ValidateLaunchUpdate(req.Launch); err != nil {
		return err
	}
	if !domainUpdateHasOperation(req) {
		return fmt.Errorf("at least one domain update operation is required")
	}
	return nil
}

// ValidateDomainRenew validates a domain renew request.
func ValidateDomainRenew(req types.DomainRenewRequest) error {
	if err := validateDomainName("domain", firstNonEmpty(req.DomainName, req.Domain)); err != nil {
		return err
	}
	if req.CurrentExpiryDate.IsZero() {
		return required("current expiry date")
	}
	if err := validatePeriodValue(domainPeriodValue(req.PeriodInfo, req.Period, req.Unit), true); err != nil {
		return err
	}
	return ValidateFeeTransform(req.Fee)
}

// ValidateDomainTransferRequest validates a domain transfer request operation.
func ValidateDomainTransferRequest(req types.DomainTransferRequest) error {
	req.Operation = constants.TransferRequest
	return ValidateDomainTransfer(req)
}

// ValidateDomainTransferQuery validates a domain transfer query operation.
func ValidateDomainTransferQuery(req types.DomainTransferRequest) error {
	req.Operation = constants.TransferQuery
	return ValidateDomainTransfer(req)
}

// ValidateDomainTransferApprove validates a domain transfer approve operation.
func ValidateDomainTransferApprove(req types.DomainTransferRequest) error {
	req.Operation = constants.TransferApprove
	return ValidateDomainTransfer(req)
}

// ValidateDomainTransferReject validates a domain transfer reject operation.
func ValidateDomainTransferReject(req types.DomainTransferRequest) error {
	req.Operation = constants.TransferReject
	return ValidateDomainTransfer(req)
}

// ValidateDomainTransferCancel validates a domain transfer cancel operation.
func ValidateDomainTransferCancel(req types.DomainTransferRequest) error {
	req.Operation = constants.TransferCancel
	return ValidateDomainTransfer(req)
}

// ValidateDomainTransfer validates a domain transfer request.
func ValidateDomainTransfer(req types.DomainTransferRequest) error {
	if err := validateDomainName("domain", firstNonEmpty(req.DomainName, req.Domain)); err != nil {
		return err
	}
	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	if !constants.IsTransferOperation(operation) {
		return fmt.Errorf("transfer operation is invalid")
	}
	if operation == constants.TransferRequest && strings.TrimSpace(req.AuthInfo) == "" {
		return required("authInfo")
	}
	if err := validatePeriodValue(domainPeriodValue(req.PeriodInfo, req.Period, req.Unit), false); err != nil {
		return err
	}
	if operation == constants.TransferRequest {
		return ValidateFeeTransform(req.Fee)
	}
	return nil
}

// ValidateDomainDelete validates a domain delete request.
func ValidateDomainDelete(req types.DomainDeleteRequest) error {
	if err := validateDomainName("domain", firstNonEmpty(req.DomainName, req.Domain)); err != nil {
		return err
	}
	return ValidateLaunchDelete(req.Launch)
}

func validateDomainName(field string, value string) error {
	value = strings.TrimSuffix(strings.TrimSpace(value), ".")
	if value == "" {
		return required(field)
	}
	if _, err := idn.ToASCII(value); err != nil {
		return fmt.Errorf("%s is invalid: %w", field, err)
	}
	return nil
}

func validateDomainNames(field string, values []string) error {
	valid := 0
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if err := validateDomainName(field, value); err != nil {
			return err
		}
		valid++
	}
	if valid == 0 {
		return required(field)
	}
	return nil
}

func validatePeriod(value int, unit string, required bool) error {
	return validatePeriodValue(types.Period{Value: value, Unit: unit}, required)
}

func validatePeriodValue(period types.Period, required bool) error {
	unit := strings.ToLower(strings.TrimSpace(period.Unit))
	if period.Value == 0 && unit == "" && !required {
		return nil
	}
	if period.Value < 1 || period.Value > maxPeriod {
		return fmt.Errorf("period must be between 1 and %d", maxPeriod)
	}
	if unit == "" {
		return nil
	}
	if unit != periodUnitYears && unit != periodUnitMonths {
		return fmt.Errorf("period unit must be %q or %q", periodUnitYears, periodUnitMonths)
	}
	return nil
}

func domainPeriodValue(period types.Period, legacyValue int, legacyUnit string) types.Period {
	if period.Value == 0 && period.Unit == "" {
		return types.Period{Value: legacyValue, Unit: legacyUnit}
	}
	return period
}

func validateDomainContacts(contacts []types.DomainContact) error {
	for _, contact := range contacts {
		if strings.TrimSpace(contact.Type) == "" || strings.TrimSpace(contact.ID) == "" {
			return fmt.Errorf("domain contact type and ID are required")
		}
	}
	return nil
}

func validateLegacyContactIDs(values ...string) error {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
	}
	return nil
}

func validateContactIDs(field string, values []string) error {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return required(field)
		}
	}
	return nil
}

func validateNameServers(hostObjects []string, hostAttrs []types.DomainNameServer) error {
	if err := validateDomainNamesIfPresent("name server", hostObjects); err != nil {
		return err
	}
	for _, host := range hostAttrs {
		if err := validateDomainName("name server", host.HostName); err != nil {
			return err
		}
		for _, address := range host.Addresses {
			if err := validateDomainHostAddress(address); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateDomainNamesIfPresent(field string, values []string) error {
	for _, value := range values {
		if err := validateDomainName(field, value); err != nil {
			return err
		}
	}
	return nil
}

func validateDomainHostAddress(address types.DomainHostAddress) error {
	ip := strings.TrimSpace(address.IP)
	if ip == "" {
		return required("host address")
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return fmt.Errorf("host address is invalid")
	}
	version := strings.TrimSpace(address.Version)
	switch version {
	case "", "v4", "v6":
	default:
		return fmt.Errorf("host address version must be %q or %q", "v4", "v6")
	}
	if version == "v4" && parsed.To4() == nil {
		return fmt.Errorf("host address must be IPv4")
	}
	if version == "v6" && parsed.To4() != nil {
		return fmt.Errorf("host address must be IPv6")
	}
	return nil
}

func validateDomainStatuses(statuses []string, statusInfo []types.DomainStatus) error {
	for _, status := range statuses {
		if err := validateDomainStatus(status); err != nil {
			return err
		}
	}
	for _, status := range statusInfo {
		if err := validateDomainStatus(status.Status); err != nil {
			return err
		}
	}
	return nil
}

func validateDomainStatus(status string) error {
	status = strings.TrimSpace(status)
	if status == "" {
		return required("domain status")
	}
	if !constants.IsDomainStatus(status) {
		return fmt.Errorf("domain status is invalid")
	}
	return nil
}

func domainUpdateHasOperation(req types.DomainUpdateRequest) bool {
	return len(req.AddStatuses) > 0 ||
		len(req.RemoveStatuses) > 0 ||
		len(req.AddStatusInfo) > 0 ||
		len(req.RemoveStatusInfo) > 0 ||
		len(req.AddNameServers) > 0 ||
		len(req.RemoveNameServers) > 0 ||
		len(req.AddNameServerInfo) > 0 ||
		len(req.RemoveNameServerInfo) > 0 ||
		len(req.AddContacts) > 0 ||
		len(req.RemoveContacts) > 0 ||
		len(req.AddAdminContacts) > 0 ||
		len(req.RemoveAdminContacts) > 0 ||
		len(req.AddTechContacts) > 0 ||
		len(req.RemoveTechContacts) > 0 ||
		len(req.AddBillingContacts) > 0 ||
		len(req.RemoveBillingContacts) > 0 ||
		strings.TrimSpace(req.Registrant) != "" ||
		strings.TrimSpace(req.NewRegistrant) != "" ||
		strings.TrimSpace(req.AuthInfo) != "" ||
		req.SecDNS != nil ||
		req.RGP != nil ||
		req.Launch != nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
