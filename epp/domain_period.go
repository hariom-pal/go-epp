package epp

import (
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/types"
)

func domainPeriod(
	value int,
	unit string,
	required bool,
) (types.Period, error) {

	unit = strings.ToLower(strings.TrimSpace(unit))
	if value == 0 && unit == "" && !required {
		return types.Period{}, nil
	}

	if value < 1 || value > domainCreateMaxPeriod {
		return types.Period{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "period must be between 1 and 99",
		}
	}

	if unit == "" {
		unit = domainCreatePeriodUnitYears
	}

	if unit != domainCreatePeriodUnitYears &&
		unit != domainCreatePeriodUnitMonths {

		return types.Period{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "period unit must be y or m",
		}
	}

	return types.Period{
		Value: value,
		Unit:  unit,
	}, nil
}
