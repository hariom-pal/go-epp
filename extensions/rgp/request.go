package rgp

import "strings"

// ValidUpdate reports whether an update request follows RFC3915 restore rules.
func ValidUpdate(req *UpdateRequest) bool {
	if req == nil {
		return true
	}

	if req.Restore == nil {
		return false
	}

	switch req.Restore.Operation {
	case OperationRequest:
		return req.Restore.Report == nil
	case OperationReport:
		return validReport(req.Restore.Report)
	default:
		return false
	}
}

// NewUpdate returns an RGP update extension, or nil when req is nil.
func NewUpdate(req *UpdateRequest) *UpdateXML {
	if req == nil || !ValidUpdate(req) {
		return nil
	}

	return &UpdateXML{
		XMLNS: Namespace,
		Restore: &RestoreXML{
			Operation: req.Restore.Operation,
			Report:    reportXML(req.Restore.Report),
		},
	}
}

func validReport(report *RestoreReport) bool {
	if report == nil {
		return false
	}

	return hasMixedValue(report.PreData, report.PreDataXML) &&
		hasMixedValue(report.PostData, report.PostDataXML) &&
		strings.TrimSpace(report.DeleteTime) != "" &&
		strings.TrimSpace(report.RestoreTime) != "" &&
		hasTextValue(report.RestoreReason) &&
		len(report.Statements) == 2 &&
		validStatements(report.Statements)
}

func validStatements(statements []Text) bool {
	for _, statement := range statements {
		if !hasTextValue(statement) {
			return false
		}
	}

	return true
}

func reportXML(report *RestoreReport) *ReportXML {
	if report == nil {
		return nil
	}

	result := &ReportXML{
		PreData:       mixedXML(report.PreData, report.PreDataXML),
		PostData:      mixedXML(report.PostData, report.PostDataXML),
		DeleteTime:    strings.TrimSpace(report.DeleteTime),
		RestoreTime:   strings.TrimSpace(report.RestoreTime),
		RestoreReason: textXML(report.RestoreReason),
		Statements:    make([]TextXML, 0, len(report.Statements)),
	}

	for _, statement := range report.Statements {
		result.Statements = append(result.Statements, textXML(statement))
	}

	if hasMixedValue(report.Other, report.OtherXML) {
		other := mixedXML(report.Other, report.OtherXML)
		result.Other = &other
	}

	return result
}

func hasMixedValue(
	text string,
	rawXML string,
) bool {

	return strings.TrimSpace(text) != "" || strings.TrimSpace(rawXML) != ""
}

func mixedXML(
	text string,
	rawXML string,
) MixedXML {

	return MixedXML{
		Text: strings.TrimSpace(text),
		XML:  strings.TrimSpace(rawXML),
	}
}

func hasTextValue(value Text) bool {
	return strings.TrimSpace(value.Value) != "" || strings.TrimSpace(value.ValueXML) != ""
}

func textXML(value Text) TextXML {
	return TextXML{
		Lang: strings.TrimSpace(value.Lang),
		Text: strings.TrimSpace(value.Value),
		XML:  strings.TrimSpace(value.ValueXML),
	}
}
