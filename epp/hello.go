package epp

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/types"
)

// Hello sends an RFC5730 hello command and returns the server greeting.
func (c *Client) Hello() (*types.Greeting, error) {
	responseXML, err := c.Execute([]byte(helloRequestXML))
	if err != nil {
		return nil, err
	}

	return parseGreetingXML(responseXML)
}

// GreetingInfo parses the raw server greeting received when the connection opened.
func (c *Client) GreetingInfo() (*types.Greeting, error) {
	return parseGreetingXML(c.greeting)
}

func parseGreetingXML(
	responseXML []byte,
) (*types.Greeting, error) {

	var response greetingEnvelopeXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, err
	}

	greeting := response.Greeting
	result := &types.Greeting{
		ServerID:            strings.TrimSpace(greeting.ServerID),
		Versions:            trimStrings(greeting.Service.Versions),
		Languages:           trimStrings(greeting.Service.Languages),
		SupportedObjects:    trimStrings(greeting.Service.Objects),
		SupportedExtensions: trimStrings(greeting.Service.Extensions.URIs),
		DCP:                 greetingDCPFromXML(greeting.DCP),
	}

	if serverDate := parseEPPDateTime(greeting.ServerDate); serverDate != nil {
		result.ServerDate = serverDate
	}

	return result, nil
}

func greetingDCPFromXML(
	dcp greetingDCPXML,
) types.GreetingDCP {

	result := types.GreetingDCP{
		Access:     firstElementLocalName(dcp.Access.Elements),
		Statements: make([]types.GreetingDCPStatement, 0, len(dcp.Statements)),
	}

	for _, statement := range dcp.Statements {
		next := types.GreetingDCPStatement{
			Purposes:   elementLocalNames(statement.Purpose.Elements),
			Recipients: elementLocalNames(statement.Recipient.Elements),
			Retentions: elementLocalNames(statement.Retention.Elements),

			ExpiryRelative: strings.TrimSpace(statement.Expiry.Relative),
		}

		if absolute := parseGreetingDateTime(statement.Expiry.Absolute); absolute != nil {
			next.ExpiryAbsolute = absolute
		}

		result.Statements = append(result.Statements, next)
	}

	return result
}

func parseGreetingDateTime(
	value string,
) *time.Time {

	if parsed := parseEPPDateTime(value); parsed != nil {
		return parsed
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil
	}

	return &parsed
}

func elementLocalNames(
	elements []greetingAnyElementXML,
) []string {

	result := make([]string, 0, len(elements))
	for _, element := range elements {
		name := elementLocalName(element.XMLName)
		if name == "" {
			continue
		}
		result = append(result, name)
	}

	return result
}

func firstElementLocalName(
	elements []greetingAnyElementXML,
) string {

	for _, element := range elements {
		name := elementLocalName(element.XMLName)
		if name != "" {
			return name
		}
	}

	return ""
}

func elementLocalName(
	name xml.Name,
) string {

	return strings.TrimSpace(name.Local)
}

func trimStrings(
	values []string,
) []string {

	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		result = append(result, value)
	}

	return result
}
