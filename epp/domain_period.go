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

	// A unit on its own carries no meaning, so an optional period is absent
	// whenever no value was supplied. Requiring the unit to be empty too
	// would reject callers that simply default it, which is the normal way
	// to expose "y" as a default in a command line or API surface, and would
	// make optional-period commands such as transfer query unusable.
	if value == 0 && !required {
		return types.Period{}, nil
	}

	if value < 1 || value > domainCreateMaxPeriod {
		return types.Period{}, newValidationError(constants.ResultParameterError, "period must be between 1 and 99")
	}

	if unit == "" {
		unit = domainCreatePeriodUnitYears
	}

	if unit != domainCreatePeriodUnitYears &&
		unit != domainCreatePeriodUnitMonths {

		return types.Period{}, newValidationError(constants.ResultParameterError, "period unit must be y or m")
	}

	return types.Period{
		Value: value,
		Unit:  unit,
	}, nil
}
