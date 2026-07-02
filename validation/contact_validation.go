package validation

import (
	"fmt"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

// ValidateContactCheck validates a contact check request.
func ValidateContactCheck(req types.ContactCheckRequest) error {
	if len(req.IDs) == 0 {
		return required("contact ID")
	}
	return validateContactIDs("contact ID", req.IDs)
}

// ValidateContactInfo validates a contact info request.
func ValidateContactInfo(req types.ContactInfoRequest) error {
	return validateContactID(req.ContactID)
}

// ValidateContactCreate validates a contact create request.
func ValidateContactCreate(req types.ContactCreateRequest) error {
	if err := validateContactID(req.ContactID); err != nil {
		return err
	}
	if err := validateOptionalPostalInfo(req.InternationalPostalInfo, req.LocalizedPostalInfo); err != nil {
		return err
	}
	if strings.TrimSpace(req.Voice.Number) == "" {
		return required("voice")
	}
	if strings.TrimSpace(req.Email) == "" {
		return required("email")
	}
	if strings.TrimSpace(req.AuthInfo) == "" {
		return required("authInfo")
	}
	return nil
}

// ValidateContactUpdate validates a contact update request.
func ValidateContactUpdate(req types.ContactUpdateRequest) error {
	if err := validateContactID(req.ContactID); err != nil {
		return err
	}
	if err := validateContactStatuses(req.AddStatuses); err != nil {
		return err
	}
	if err := validateContactStatuses(req.RemoveStatuses); err != nil {
		return err
	}
	if err := validatePostalInfo(req.InternationalPostalInfo, req.LocalizedPostalInfo); err != nil {
		return err
	}
	if req.Voice != nil && strings.TrimSpace(req.Voice.Number) == "" {
		return required("voice")
	}
	if req.Fax != nil && strings.TrimSpace(req.Fax.Number) == "" {
		return required("fax")
	}
	if !contactUpdateHasOperation(req) {
		return fmt.Errorf("at least one contact update operation is required")
	}
	return nil
}

// ValidateContactDelete validates a contact delete request.
func ValidateContactDelete(req types.ContactDeleteRequest) error {
	return validateContactID(req.ContactID)
}

func validateContactID(value string) error {
	if strings.TrimSpace(value) == "" {
		return required("contact ID")
	}
	return nil
}

func validatePostalInfo(values ...*types.PostalInfo) error {
	seen := false
	for _, value := range values {
		if value == nil {
			continue
		}
		seen = true
		if strings.TrimSpace(value.Name) == "" {
			return required("postalInfo name")
		}
		if strings.TrimSpace(value.City) == "" {
			return required("postalInfo city")
		}
		if strings.TrimSpace(value.CountryCode) == "" {
			return required("postalInfo country code")
		}
	}
	if !seen {
		return required("postalInfo")
	}
	return nil
}

func validateOptionalPostalInfo(values ...*types.PostalInfo) error {
	seen := false
	for _, value := range values {
		if value != nil {
			seen = true
			break
		}
	}
	if !seen {
		return nil
	}
	return validatePostalInfo(values...)
}

func validateContactStatuses(statuses []string) error {
	for _, status := range statuses {
		status = strings.TrimSpace(status)
		if status == "" {
			return required("contact status")
		}
		if !constants.IsContactStatus(status) {
			return fmt.Errorf("contact status is invalid")
		}
	}
	return nil
}

func contactUpdateHasOperation(req types.ContactUpdateRequest) bool {
	return len(req.AddStatuses) > 0 ||
		len(req.RemoveStatuses) > 0 ||
		req.InternationalPostalInfo != nil ||
		req.LocalizedPostalInfo != nil ||
		req.Voice != nil ||
		req.Fax != nil ||
		strings.TrimSpace(req.Email) != "" ||
		strings.TrimSpace(req.AuthInfo) != "" ||
		req.Disclosure != nil
}
