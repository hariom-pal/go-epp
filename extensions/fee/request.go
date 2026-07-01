package fee

import "strings"

// NewCheckExtension returns a fee-0.7 check extension, or nil when req is nil or empty.
func NewCheckExtension(req *CheckRequest) *CheckExtensionXML {
	if req == nil || len(req.Domains) == 0 {
		return nil
	}

	domains := make([]CheckDomainXML, 0, len(req.Domains))
	for _, domain := range req.Domains {
		name := strings.TrimSpace(domain.Name)
		command := commandXML(domain.Command)
		if name == "" || command.Value == "" {
			continue
		}

		domains = append(domains, CheckDomainXML{
			Name:     name,
			Currency: strings.TrimSpace(domain.Currency),
			Command:  command,
			Period:   periodXML(domain.Period),
		})
	}

	if len(domains) == 0 {
		return nil
	}

	return &CheckExtensionXML{
		Check: CheckXML{
			XMLNS:   Namespace,
			Domains: domains,
		},
	}
}

// NewTransformExtension returns a fee-0.7 transform extension for operation.
func NewTransformExtension(
	operation string,
	req *TransformRequest,
) *TransformExtensionXML {

	transform := transformXML(req)
	if transform == nil {
		return nil
	}

	extension := &TransformExtensionXML{}
	switch operation {
	case CommandCreate:
		extension.Create = transform
	case CommandRenew:
		extension.Renew = transform
	case CommandTransfer:
		extension.Transfer = transform
	default:
		return nil
	}

	return extension
}

func transformXML(
	req *TransformRequest,
) *TransformXML {

	if req == nil {
		return nil
	}

	fees := amountXMLs(req.Fees)
	credits := amountXMLs(req.Credits)
	if len(fees) == 0 {
		return nil
	}

	return &TransformXML{
		XMLNS:    Namespace,
		Currency: strings.TrimSpace(req.Currency),
		Fees:     fees,
		Credits:  credits,
	}
}

func commandXML(command Command) CommandXML {
	return CommandXML{
		Phase:    strings.TrimSpace(command.Phase),
		Subphase: strings.TrimSpace(command.Subphase),
		Value:    strings.TrimSpace(command.Name),
	}
}

func periodXML(period *Period) *PeriodXML {
	if period == nil {
		return nil
	}

	return &PeriodXML{
		Unit:  strings.TrimSpace(period.Unit),
		Value: period.Value,
	}
}

func amountXMLs(amounts []Amount) []AmountXML {

	result := make([]AmountXML, 0, len(amounts))
	for _, amount := range amounts {
		value := strings.TrimSpace(amount.Amount)
		if value == "" {
			continue
		}
		result = append(result, AmountXML{
			Description: strings.TrimSpace(amount.Description),
			Refundable:  strings.TrimSpace(amount.Refundable),
			GracePeriod: strings.TrimSpace(amount.GracePeriod),
			Applied:     strings.TrimSpace(amount.Applied),
			Value:       value,
		})
	}

	return result
}
