package launch

import "strings"

// IDDataFromXML converts launch create XML into reusable response data.
func IDDataFromXML(data IDDataXML) IDData {
	return IDData{
		Phase:         phaseFromXML(data.Phase),
		ApplicationID: strings.TrimSpace(data.ApplicationID),
	}
}

// InfoDataFromXML converts launch info XML into reusable response data.
func InfoDataFromXML(data InfoDataXML) InfoData {
	rawXML := strings.TrimSpace(data.RawXML)

	return InfoData{
		Phase:         phaseFromXML(data.Phase),
		ApplicationID: strings.TrimSpace(data.ApplicationID),
		Status: Status{
			Status: strings.TrimSpace(data.Status.Status),
			Name:   strings.TrimSpace(data.Status.Name),
			Lang:   strings.TrimSpace(data.Status.Lang),
			Text:   strings.TrimSpace(data.Status.Text),
		},
		MarkXML: rawXML,
		RawXML:  rawXML,
	}
}

func phaseFromXML(value PhaseXML) Phase {
	return Phase{
		Value: strings.TrimSpace(value.Value),
		Name:  strings.TrimSpace(value.Name),
	}
}
