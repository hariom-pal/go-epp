package rgp

import "strings"

// InfoDataFromXML converts RGP info XML into reusable response data.
func InfoDataFromXML(data InfoDataXML) InfoData {
	return InfoData{
		Statuses: statusesFromXML(data.Statuses),
	}
}

// UpdateDataFromXML converts RGP update XML into reusable response data.
func UpdateDataFromXML(data UpdateDataXML) UpdateData {
	return UpdateData{
		Statuses: statusesFromXML(data.Statuses),
	}
}

func statusesFromXML(values []StatusXML) []Status {
	result := make([]Status, 0, len(values))
	for _, value := range values {
		status := strings.TrimSpace(value.Status)
		if status == "" {
			continue
		}

		result = append(result, Status{
			Status: status,
			Lang:   strings.TrimSpace(value.Lang),
			Text:   strings.TrimSpace(value.Text),
		})
	}

	return result
}
