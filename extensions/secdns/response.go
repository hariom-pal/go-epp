package secdns

import "strings"

// InfoDataFromXML converts secDNS info XML into reusable response data.
func InfoDataFromXML(data InfoDataXML) InfoData {
	return InfoData{
		Data: Data{
			MaxSigLife: data.MaxSigLife,
			DSData:     dsDataFromXML(data.DSData),
			KeyData:    keyDataFromXML(data.KeyData),
		},
	}
}

func dsDataFromXML(values []InfoDSDataXML) []DSData {
	result := make([]DSData, 0, len(values))
	for _, value := range values {
		digest := strings.TrimSpace(value.Digest)
		if digest == "" {
			continue
		}

		result = append(result, DSData{
			KeyTag:     value.KeyTag,
			Algorithm:  value.Algorithm,
			DigestType: value.DigestType,
			Digest:     digest,
			KeyData:    keyDataPointerFromXML(value.KeyData),
		})
	}

	return result
}

func keyDataFromXML(values []InfoKeyDataXML) []KeyData {
	result := make([]KeyData, 0, len(values))
	for _, value := range values {
		publicKey := strings.TrimSpace(value.PublicKey)
		if publicKey == "" {
			continue
		}

		result = append(result, KeyData{
			Flags:     value.Flags,
			Protocol:  value.Protocol,
			Algorithm: value.Algorithm,
			PublicKey: publicKey,
		})
	}

	return result
}

func keyDataPointerFromXML(value *InfoKeyDataXML) *KeyData {
	if value == nil {
		return nil
	}

	publicKey := strings.TrimSpace(value.PublicKey)
	if publicKey == "" {
		return nil
	}

	return &KeyData{
		Flags:     value.Flags,
		Protocol:  value.Protocol,
		Algorithm: value.Algorithm,
		PublicKey: publicKey,
	}
}
