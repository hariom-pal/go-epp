package idn

import "strings"

// ValidCreate reports whether an IDN create extension is usable. A nil
// request is valid because the extension is optional; a present one must
// carry a table tag, which is the only field a registry cannot derive.
func ValidCreate(req *CreateRequest) bool {
	if req == nil {
		return true
	}
	return strings.TrimSpace(req.Table) != ""
}

// NewCreate builds idn-1.0 create extension XML, or nil when no IDN data
// was supplied.
func NewCreate(req *CreateRequest) *DataXML {
	if req == nil || strings.TrimSpace(req.Table) == "" {
		return nil
	}
	return &DataXML{
		XMLNS: Namespace,
		Table: strings.TrimSpace(req.Table),
		UName: strings.TrimSpace(req.UName),
	}
}
