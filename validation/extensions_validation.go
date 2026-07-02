package validation

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/extensions/fee"
	"github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/extensions/rgp"
	"github.com/hariom-pal/go-epp/extensions/secdns"
)

// ValidateFeeCheck validates a fee check extension request.
func ValidateFeeCheck(req *fee.CheckRequest) error {
	if req == nil {
		return nil
	}
	for _, domain := range req.Domains {
		if strings.TrimSpace(domain.Name) == "" {
			return required("fee domain name")
		}
		if strings.TrimSpace(domain.Currency) == "" {
			return required("fee currency")
		}
		if strings.TrimSpace(domain.Command.Name) == "" {
			return required("fee command")
		}
		if domain.Period != nil {
			if err := validateFeePeriod(*domain.Period); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateFeeTransform validates a fee transform extension request.
func ValidateFeeTransform(req *fee.TransformRequest) error {
	if req == nil {
		return nil
	}
	if strings.TrimSpace(req.Currency) == "" {
		return required("fee currency")
	}
	return nil
}

// ValidateDNSSECCreate validates a secDNS create extension request.
func ValidateDNSSECCreate(req *secdns.CreateRequest) error {
	if !secdns.ValidCreate(req) {
		return fmt.Errorf("secDNS create request is invalid")
	}
	return nil
}

// ValidateDNSSECUpdate validates a secDNS update extension request.
func ValidateDNSSECUpdate(req *secdns.UpdateRequest) error {
	if !secdns.ValidUpdate(req) {
		return fmt.Errorf("secDNS update request is invalid")
	}
	return nil
}

// ValidateRGPUpdate validates an RGP update extension request.
func ValidateRGPUpdate(req *rgp.UpdateRequest) error {
	if !rgp.ValidUpdate(req) {
		return fmt.Errorf("RGP update request is invalid")
	}
	return nil
}

// ValidateLaunchCreate validates a launch create extension request.
func ValidateLaunchCreate(req *launch.CreateRequest) error {
	if !launch.ValidCreate(req) {
		return fmt.Errorf("launch create request is invalid")
	}
	return nil
}

// ValidateLaunchInfo validates a launch info extension request.
func ValidateLaunchInfo(req *launch.InfoRequest) error {
	if !launch.ValidInfo(req) {
		return fmt.Errorf("launch info request is invalid")
	}
	return nil
}

// ValidateLaunchUpdate validates a launch update extension request.
func ValidateLaunchUpdate(req *launch.UpdateRequest) error {
	if !launch.ValidUpdate(req) {
		return fmt.Errorf("launch update request is invalid")
	}
	return nil
}

// ValidateLaunchDelete validates a launch delete extension request.
func ValidateLaunchDelete(req *launch.DeleteRequest) error {
	if !launch.ValidDelete(req) {
		return fmt.Errorf("launch delete request is invalid")
	}
	return nil
}

func validateFeePeriod(period fee.Period) error {
	if period.Value < 1 || period.Value > maxPeriod {
		return fmt.Errorf("fee period must be between 1 and %d", maxPeriod)
	}
	unit := strings.ToLower(strings.TrimSpace(period.Unit))
	if unit != periodUnitYears && unit != periodUnitMonths {
		return fmt.Errorf("fee period unit must be %q or %q", periodUnitYears, periodUnitMonths)
	}
	return nil
}
