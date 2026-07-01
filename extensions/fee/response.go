package fee

import "strings"

// CheckDataFromXML converts fee-0.7 check XML into reusable response data.
func CheckDataFromXML(data CheckDataXML) CheckData {
	result := CheckData{
		Results: make([]CheckResult, 0, len(data.Results)),
	}

	for _, cd := range data.Results {
		result.Results = append(result.Results, CheckResult{
			Name:     strings.TrimSpace(cd.Name),
			Currency: strings.TrimSpace(cd.Currency),
			Command:  commandFromXML(cd.Command),
			Period:   periodFromXML(cd.Period),
			Fees:     amountsFromXML(cd.Fees),
			Credits:  amountsFromXML(cd.Credits),
			Class:    Classification(strings.TrimSpace(cd.Class)),
			Reason:   strings.TrimSpace(cd.Reason),
		})
	}

	return result
}

// InfoDataFromXML converts fee-0.7 info XML into reusable response data.
func InfoDataFromXML(data InfoDataXML) InfoData {
	return InfoData{
		Currency: strings.TrimSpace(data.Currency),
		Command:  commandFromXML(data.Command),
		Period:   periodFromXML(data.Period),
		Fees:     amountsFromXML(data.Fees),
		Credits:  amountsFromXML(data.Credits),
		Class:    Classification(strings.TrimSpace(data.Class)),
	}
}

// TransformDataFromXML converts fee-0.7 transform XML into reusable response data.
func TransformDataFromXML(data TransformDataXML) TransformData {
	return TransformData{
		Currency:    strings.TrimSpace(data.Currency),
		Period:      periodFromXML(data.Period),
		Fees:        amountsFromXML(data.Fees),
		Credits:     amountsFromXML(data.Credits),
		Balance:     strings.TrimSpace(data.Balance),
		CreditLimit: strings.TrimSpace(data.CreditLimit),
	}
}

func commandFromXML(command CommandXML) Command {
	return Command{
		Name:     strings.TrimSpace(command.Value),
		Phase:    strings.TrimSpace(command.Phase),
		Subphase: strings.TrimSpace(command.Subphase),
	}
}

func periodFromXML(period *PeriodXML) *Period {
	if period == nil {
		return nil
	}

	return &Period{
		Value: period.Value,
		Unit:  strings.TrimSpace(period.Unit),
	}
}

func amountsFromXML(amounts []AmountXML) []Amount {
	result := make([]Amount, 0, len(amounts))
	for _, amount := range amounts {
		value := strings.TrimSpace(amount.Value)
		if value == "" {
			continue
		}
		result = append(result, Amount{
			Amount:      value,
			Description: strings.TrimSpace(amount.Description),
			Refundable:  strings.TrimSpace(amount.Refundable),
			GracePeriod: strings.TrimSpace(amount.GracePeriod),
			Applied:     strings.TrimSpace(amount.Applied),
		})
	}

	return result
}
