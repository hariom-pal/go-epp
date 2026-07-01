package launch

import (
	"encoding/xml"
	"strings"

	extcommon "github.com/hariom-pal/go-epp/extensions/common"
)

const (
	signedMarkNamespace = "urn:ietf:params:xml:ns:signedMark-1.0"
	markNamespace       = "urn:ietf:params:xml:ns:mark-1.0"
)

// ValidCreate reports whether a create request contains the required launch phase.
func ValidCreate(req *CreateRequest) bool {
	if req == nil {
		return true
	}

	return validPhase(req.Phase)
}

// ValidInfo reports whether an info request contains the required launch phase.
func ValidInfo(req *InfoRequest) bool {
	if req == nil {
		return true
	}

	return validPhase(req.Phase)
}

// ValidUpdate reports whether an update request contains required launch data.
func ValidUpdate(req *UpdateRequest) bool {
	if req == nil {
		return true
	}

	return validPhase(req.Phase) && strings.TrimSpace(req.ApplicationID) != ""
}

// ValidDelete reports whether a delete request contains required launch data.
func ValidDelete(req *DeleteRequest) bool {
	if req == nil {
		return true
	}

	return validPhase(req.Phase) && strings.TrimSpace(req.ApplicationID) != ""
}

// NewCreate returns a launch create extension, or nil when req is nil.
func NewCreate(req *CreateRequest) *CreateXML {
	if req == nil || !ValidCreate(req) {
		return nil
	}

	return &CreateXML{
		XMLNS:     Namespace,
		XMLNSSMD:  namespaceIf(len(req.SignedMarkXML) > 0 || len(req.EncodedSignedMarkXML) > 0, signedMarkNamespace),
		XMLNSMark: namespaceIf(hasMarkXML(req.CodeMarks), markNamespace),
		Type:      strings.TrimSpace(req.Type),
		Phase:     phaseXML(req.Phase),
		ChoiceXML: RawXML{Value: createChoiceXML(req)},
		Notices:   noticeXMLs(req.Notices),
	}
}

// NewInfo returns a launch info extension, or nil when req is nil.
func NewInfo(req *InfoRequest) *InfoXML {
	if req == nil || !ValidInfo(req) {
		return nil
	}

	result := &InfoXML{
		XMLNS:         Namespace,
		Phase:         phaseXML(req.Phase),
		ApplicationID: strings.TrimSpace(req.ApplicationID),
	}
	if req.IncludeMark {
		result.IncludeMark = extcommon.BoolString(req.IncludeMark)
	}

	return result
}

// NewUpdate returns a launch update extension, or nil when req is nil.
func NewUpdate(req *UpdateRequest) *UpdateXML {
	if req == nil || !ValidUpdate(req) {
		return nil
	}

	return &UpdateXML{
		XMLNS:         Namespace,
		Phase:         phaseXML(req.Phase),
		ApplicationID: strings.TrimSpace(req.ApplicationID),
	}
}

// NewDelete returns a launch delete extension, or nil when req is nil.
func NewDelete(req *DeleteRequest) *DeleteXML {
	if req == nil || !ValidDelete(req) {
		return nil
	}

	return &DeleteXML{
		XMLNS:         Namespace,
		Phase:         phaseXML(req.Phase),
		ApplicationID: strings.TrimSpace(req.ApplicationID),
	}
}

func validPhase(phase Phase) bool {
	return strings.TrimSpace(phase.Value) != ""
}

func phaseXML(phase Phase) PhaseXML {
	return PhaseXML{
		Name:  strings.TrimSpace(phase.Name),
		Value: strings.TrimSpace(phase.Value),
	}
}

func createChoiceXML(req *CreateRequest) string {
	var builder strings.Builder

	for _, codeMark := range req.CodeMarks {
		value := codeMarkXML(codeMark)
		if value.Code == nil && strings.TrimSpace(value.MarkXML.Value) == "" {
			continue
		}

		out, err := xml.Marshal(value)
		if err == nil {
			builder.Write(out)
		}
	}

	for _, raw := range req.SignedMarkXML {
		builder.WriteString(strings.TrimSpace(raw))
	}

	for _, raw := range req.EncodedSignedMarkXML {
		builder.WriteString(strings.TrimSpace(raw))
	}

	return builder.String()
}

func codeMarkXML(value CodeMark) CodeMarkXML {
	return CodeMarkXML{
		Code:    codeXML(value.Code),
		MarkXML: RawXML{Value: strings.TrimSpace(value.MarkXML)},
	}
}

func codeXML(value *Code) *CodeXML {
	if value == nil {
		return nil
	}

	code := strings.TrimSpace(value.Value)
	if code == "" {
		return nil
	}

	return &CodeXML{
		ValidatorID: strings.TrimSpace(value.ValidatorID),
		Value:       code,
	}
}

func noticeXMLs(values []Notice) []NoticeXML {
	result := make([]NoticeXML, 0, len(values))
	for _, value := range values {
		id := strings.TrimSpace(value.ID)
		notAfter := strings.TrimSpace(value.NotAfter)
		acceptedDate := strings.TrimSpace(value.AcceptedDate)
		if id == "" || notAfter == "" || acceptedDate == "" {
			continue
		}

		result = append(result, NoticeXML{
			NoticeID: NoticeIDXML{
				ValidatorID: strings.TrimSpace(value.ValidatorID),
				Value:       id,
			},
			NotAfter:     notAfter,
			AcceptedDate: acceptedDate,
		})
	}

	return result
}

func hasMarkXML(values []CodeMark) bool {
	for _, value := range values {
		if strings.TrimSpace(value.MarkXML) != "" {
			return true
		}
	}

	return false
}

func namespaceIf(include bool, namespace string) string {
	if include {
		return namespace
	}

	return ""
}
