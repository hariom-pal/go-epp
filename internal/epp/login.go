package epp

import (
	"fmt"
)

func (c *Client) Login() error {

	trID := c.nextTRID("LOGIN")

	loginXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <login>
            <clID>%s</clID>
            <pw>%s</pw>

            <options>
                <version>1.0</version>
                <lang>en</lang>
            </options>

            <svcs>
                <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
                <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
                <objURI>urn:ietf:params:xml:ns:host-1.0</objURI>

                <svcExtension>
                    <extURI>urn:ietf:params:xml:ns:secDNS-1.1</extURI>
                    <extURI>urn:ietf:params:xml:ns:rgp-1.0</extURI>
                    <extURI>urn:ietf:params:xml:ns:idn-1.0</extURI>
                    <extURI>urn:ietf:params:xml:ns:fee-0.7</extURI>
                    <extURI>urn:ietf:params:xml:ns:launch-1.0</extURI>
                </svcExtension>

            </svcs>

        </login>

        <clTRID>%s</clTRID>

    </command>
</epp>`,
		c.config.Authentication.Username,
		c.config.Authentication.Password,
		trID,
	)

	response, err := c.Execute([]byte(loginXML))
	if err != nil {
		return err
	}

	fmt.Println("========== LOGIN RESPONSE ==========")
	fmt.Println(string(response))
	fmt.Println("====================================")

	return nil
}
