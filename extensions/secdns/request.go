package secdns

import (
	"strings"

	extcommon "github.com/hariom-pal/go-epp/extensions/common"
)

// ValidCreate reports whether a create request follows RFC5910 interface rules.
func ValidCreate(req *CreateRequest) bool {
	if req == nil {
		return true
	}

	return !mixedRequestDataInterfaces(req.DSData, req.KeyData)
}

// ValidUpdate reports whether an update request follows RFC5910 interface rules.
func ValidUpdate(req *UpdateRequest) bool {
	if req == nil {
		return true
	}

	if req.Add != nil && mixedRequestDataInterfaces(req.Add.DSData, req.Add.KeyData) {
		return false
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
		if digest == "" {
			continue
		}

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
		if publicKey == "" {
			continue
		}

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
