package epp

import (
	"encoding/xml"

	"github.com/hariom-pal/go-epp/constants"
)

// Login sends an EPP login command using the client's configured credentials.
func (c *Client) Login() error {

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
					ObjectURIs: []string{
						constants.DomainNamespace,
						constants.ContactNamespace,
						constants.HostNamespace,
					},
					ServiceExtension: loginServiceExtensionXML{
						ExtensionURIs: []string{
							constants.SecDNSNamespace,
							constants.RGPNamespace,
							constants.IDNNamespace,
							constants.FeeNamespace,
							constants.LaunchNamespace,
						},
					},
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

	response, err := c.Execute(loginXML)
	if err != nil {
		return err
	}

	return parseCommandResponse(response)
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
	ObjectURIs       []string                 `xml:"objURI"`
	ServiceExtension loginServiceExtensionXML `xml:"svcExtension"`
}

type loginServiceExtensionXML struct {
	ExtensionURIs []string `xml:"extURI"`
}
