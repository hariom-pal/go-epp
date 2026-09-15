package epp

import (
	"context"
	"encoding/xml"
	"slices"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
)

// Login sends an EPP login command using the client's configured credentials.
func (c *Client) Login() error {
	return c.LoginContext(context.Background())
}

// LoginContext sends an EPP login command using the client's configured credentials.
func (c *Client) LoginContext(ctx context.Context) error {
	// An empty credential is always a caller mistake, and sending it costs a
	// round trip to be told so in terms that point at the schema rather than
	// at the configuration. RFC 5730 also bounds clID and pw lengths, but
	// registries commonly issue values outside those bounds, so only the
	// unambiguous case is enforced here.
	if strings.TrimSpace(c.config.Authentication.Username) == "" {
		return newValidationError(constants.ResultParameterError, "login client ID is required")
	}
	if strings.TrimSpace(c.config.Authentication.Password) == "" {
		return newValidationError(constants.ResultParameterError, "login password is required")
	}

	objects, extensions, err := c.loginServices()
	if err != nil {
		return err
	}
	request := loginRequestXML{
		XMLNS: constants.EPPNamespace,
		Command: loginCommandXML{
			Login: loginXML{
				ClientID: c.config.Authentication.Username,
				Password: c.config.Authentication.Password,
				Options: loginOptionsXML{
					Version: "1.0",
					Lang:    "en",
				},
				Services: loginServicesXML{
					ObjectURIs:       objects,
					ServiceExtension: serviceExtensionXML(extensions),
				},
			},
			ClientTRID: c.nextTRID("LOGIN"),
		},
	}

	loginXML, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		return err
	}

	loginXML = append([]byte(xml.Header), loginXML...)

	response, err := c.ExecuteContext(ctx, loginXML)
	if err != nil {
		return err
	}

	if err := parseCommandResponse(response); err != nil {
		return err
	}
	c.setLoggedIn(true)
	c.emit(Event{Type: EventLogin})
	return nil
}

func (c *Client) loginServices() ([]string, []string, error) {
	desiredObjects := append([]string(nil), c.config.Login.ObjectURIs...)
	if len(desiredObjects) == 0 {
		desiredObjects = []string{
			constants.DomainNamespace,
			constants.ContactNamespace,
			constants.HostNamespace,
		}
	}
	supportedObjects := map[string]bool{}
	if c.greetingInfo != nil {
		for _, uri := range c.greetingInfo.SupportedObjects {
			supportedObjects[strings.TrimSpace(uri)] = true
		}
	}

	objects := make([]string, 0, len(desiredObjects))
	for _, uri := range desiredObjects {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}
		if c.greetingInfo != nil && len(supportedObjects) > 0 && !supportedObjects[uri] {
			if c.config.Login.RequireSupportedObjects {
				return nil, nil, newSDKError(
					ErrorKindConfiguration,
					"requested login object is not advertised by greeting: "+uri,
					nil,
				)
			}
			continue
		}
		if !slices.Contains(objects, uri) {
			objects = append(objects, uri)
		}
	}
	if len(objects) == 0 {
		return nil, nil, newSDKError(ErrorKindConfiguration, "no login object services are available", nil)
	}

	desiredExtensions := append([]string(nil), c.config.Login.ExtensionURIs...)
	supportedExtensions := map[string]bool{}
	if c.greetingInfo != nil {
		for _, uri := range c.greetingInfo.SupportedExtensions {
			supportedExtensions[strings.TrimSpace(uri)] = true
		}
	}

	extensions := make([]string, 0, len(desiredExtensions))
	for _, uri := range desiredExtensions {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}
		if c.greetingInfo != nil && !supportedExtensions[uri] {
			if c.config.Login.RequireSupportedExtensions {
				return nil, nil, newSDKError(
					ErrorKindConfiguration,
					"requested login extension is not advertised by greeting: "+uri,
					nil,
				)
			}
			continue
		}
		if !slices.Contains(extensions, uri) {
			extensions = append(extensions, uri)
		}
	}

	return objects, extensions, nil
}

func serviceExtensionXML(extensions []string) *loginServiceExtensionXML {
	if len(extensions) == 0 {
		return nil
	}
	return &loginServiceExtensionXML{ExtensionURIs: extensions}
}

type loginRequestXML struct {
	XMLName xml.Name `xml:"epp"`
	XMLNS   string   `xml:"xmlns,attr"`

	Command loginCommandXML `xml:"command"`
}

type loginCommandXML struct {
	Login      loginXML `xml:"login"`
	ClientTRID string   `xml:"clTRID"`
}

type loginXML struct {
	ClientID string           `xml:"clID"`
	Password string           `xml:"pw"`
	Options  loginOptionsXML  `xml:"options"`
	Services loginServicesXML `xml:"svcs"`
}

type loginOptionsXML struct {
	Version string `xml:"version"`
	Lang    string `xml:"lang"`
}

type loginServicesXML struct {
	ObjectURIs       []string                  `xml:"objURI"`
	ServiceExtension *loginServiceExtensionXML `xml:"svcExtension,omitempty"`
}

type loginServiceExtensionXML struct {
	ExtensionURIs []string `xml:"extURI"`
}
