package epp

import (
	"fmt"

	"github.com/hariom-pal/go-epp/constants"
)

// Logout sends an EPP logout command.
func (c *Client) Logout() error {

	logoutXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="%s">
    <command>
        <logout/>
        <clTRID>%s</clTRID>
    </command>
</epp>`, constants.EPPNamespace, c.nextTRID("LOGOUT"))

	response, err := c.Execute([]byte(logoutXML))
	if err != nil {
		return err
	}

	return parseCommandResponse(response)
}
