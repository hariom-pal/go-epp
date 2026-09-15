package secdns

import (
	"strings"

	extcommon "github.com/hariom-pal/go-epp/extensions/common"
)

// RFC 5910 section 4 makes keyTag, alg, digestType and digest mandatory in a
// dsData record, and flags, protocol, alg and pubKey mandatory in a keyData
// record. An incomplete record cannot be serialised into schema-valid XML, so
// it is rejected here rather than dropped: silently omitting DNSSEC data would
// leave a caller believing a domain was signed when nothing was sent.
func validDSData(values []DSData) bool {
	for _, value := range values {
		if value.Algorithm <= 0 || value.DigestType <= 0 {
			return false
		}
		if !validHexDigest(value.Digest) {
			return false
		}
		if value.KeyData != nil && !validKeyData([]KeyData{*value.KeyData}) {
			return false
		}
	}
	return true
}

func validKeyData(values []KeyData) bool {
	for _, value := range values {
		if value.Algorithm <= 0 || value.Protocol <= 0 {
			return false
		}
		if strings.TrimSpace(value.PublicKey) == "" {
			return false
		}
	}
	return true
}

// validHexDigest reports whether value is non-empty hexBinary, as RFC 5910
// requires for a DS digest.
func validHexDigest(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || len(value)%2 != 0 {
		return false
	}
	for _, char := range value {
		switch {
		case char >= '0' && char <= '9':
		case char >= 'a' && char <= 'f':
		case char >= 'A' && char <= 'F':
		default:
			return false
		}
	}
	return true
}

// ValidCreate reports whether a create request follows RFC5910 interface rules.
func ValidCreate(req *CreateRequest) bool {
	if req == nil {
		return true
	}

	if mixedRequestDataInterfaces(req.DSData, req.KeyData) {
		return false
	}
	return validDSData(req.DSData) && validKeyData(req.KeyData)
}

// ValidUpdate reports whether an update request follows RFC5910 interface rules.
func ValidUpdate(req *UpdateRequest) bool {
	if req == nil {
		return true
	}

	if req.Add != nil {
		if mixedRequestDataInterfaces(req.Add.DSData, req.Add.KeyData) {
			return false
		}
		if !validDSData(req.Add.DSData) || !validKeyData(req.Add.KeyData) {
			return false
		}
	}

	hasDSData := req.Add != nil && len(req.Add.DSData) > 0
	hasKeyData := req.Add != nil && len(req.Add.KeyData) > 0

	if req.Remove != nil {
		if req.Remove.All != nil && (len(req.Remove.DSData) > 0 || len(req.Remove.KeyData) > 0) {
			return false
		}
		if mixedRequestDataInterfaces(req.Remove.DSData, req.Remove.KeyData) {
			return false
		}
		if !validDSData(req.Remove.DSData) || !validKeyData(req.Remove.KeyData) {
			return false
		}

		hasDSData = hasDSData || len(req.Remove.DSData) > 0
		hasKeyData = hasKeyData || len(req.Remove.KeyData) > 0
	}

	//return !(hasDSData && hasKeyData)
	return !hasDSData || !hasKeyData
}

// NewCreate returns a secDNS create extension, or nil when req is nil or empty.
func NewCreate(req *CreateRequest) *CreateXML {
	if req == nil {
		return nil
	}

	dsData := dsDataXMLs(req.DSData)
	keyData := keyDataXMLs(req.KeyData)
	if mixedDataInterfaces(dsData, keyData) {
		return nil
	}
	if req.MaxSigLife == 0 && len(dsData) == 0 && len(keyData) == 0 {
		return nil
	}

	return &CreateXML{
		XMLNS:      Namespace,
		MaxSigLife: req.MaxSigLife,
		DSData:     dsData,
		KeyData:    keyData,
	}
}

func mixedRequestDataInterfaces(
	dsData []DSData,
	keyData []KeyData,
) bool {

	return len(dsData) > 0 && len(keyData) > 0
}

// NewUpdate returns a secDNS update extension, or nil when req is nil or empty.
func NewUpdate(req *UpdateRequest) *UpdateXML {
	if req == nil {
		return nil
	}

	update := &UpdateXML{
		XMLNS: Namespace,
	}
	if req.Urgent {
		update.Urgent = extcommon.BoolString(req.Urgent)
	}

	update.Remove = updateRemoveXML(req.Remove)
	update.Add = updateAddXML(req.Add)
	update.Change = updateChangeXML(req.Change)

	if update.Remove == nil &&
		update.Add == nil &&
		update.Change == nil {

		return nil
	}

	return update
}

func updateAddXML(add *UpdateAdd) *UpdateAddXML {
	if add == nil {
		return nil
	}

	dsData := dsDataXMLs(add.DSData)
	keyData := keyDataXMLs(add.KeyData)
	if mixedDataInterfaces(dsData, keyData) {
		return nil
	}
	if add.MaxSigLife == 0 && len(dsData) == 0 && len(keyData) == 0 {
		return nil
	}

	return &UpdateAddXML{
		MaxSigLife: add.MaxSigLife,
		DSData:     dsData,
		KeyData:    keyData,
	}
}

func updateRemoveXML(remove *UpdateRemove) *UpdateRemoveXML {
	if remove == nil {
		return nil
	}

	dsData := dsDataXMLs(remove.DSData)
	keyData := keyDataXMLs(remove.KeyData)
	if remove.All != nil && (len(dsData) > 0 || len(keyData) > 0) {
		return nil
	}
	if mixedDataInterfaces(dsData, keyData) {
		return nil
	}
	if remove.All == nil && len(dsData) == 0 && len(keyData) == 0 {
		return nil
	}

	return &UpdateRemoveXML{
		All:     remove.All,
		DSData:  dsData,
		KeyData: keyData,
	}
}

func updateChangeXML(change *UpdateChange) *UpdateChangeXML {
	if change == nil || change.MaxSigLife == 0 {
		return nil
	}

	return &UpdateChangeXML{
		MaxSigLife: change.MaxSigLife,
	}
}

func mixedDataInterfaces(
	dsData []DSDataXML,
	keyData []KeyDataXML,
) bool {

	return len(dsData) > 0 && len(keyData) > 0
}

func dsDataXMLs(values []DSData) []DSDataXML {
	result := make([]DSDataXML, 0, len(values))
	for _, value := range values {
		digest := strings.TrimSpace(value.Digest)

		result = append(result, DSDataXML{
			KeyTag:     value.KeyTag,
			Algorithm:  value.Algorithm,
			DigestType: value.DigestType,
			Digest:     digest,
			KeyData:    keyDataXML(value.KeyData),
		})
	}

	return result
}

func keyDataXMLs(values []KeyData) []KeyDataXML {
	result := make([]KeyDataXML, 0, len(values))
	for _, value := range values {
		publicKey := strings.TrimSpace(value.PublicKey)

		result = append(result, KeyDataXML{
			Flags:     value.Flags,
			Protocol:  value.Protocol,
			Algorithm: value.Algorithm,
			PublicKey: publicKey,
		})
	}

	return result
}

func keyDataXML(value *KeyData) *KeyDataXML {
	if value == nil {
		return nil
	}

	publicKey := strings.TrimSpace(value.PublicKey)
	if publicKey == "" {
		return nil
	}

	return &KeyDataXML{
		Flags:     value.Flags,
		Protocol:  value.Protocol,
		Algorithm: value.Algorithm,
		PublicKey: publicKey,
	}
}
