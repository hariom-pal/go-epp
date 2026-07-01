package epp

import (
	"fmt"
)

// Logout sends an EPP logout command.
func (c *Client) Logout() error {

	trID := c.nextTRID("LOGOUT")

	logoutXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <command>
        <logout/>
        <clTRID>%s</clTRID>
    </command>
</epp>`, trID)

	response, err := c.Execute([]byte(logoutXML))
	if err != nil {
		return err
	}

	fmt.Println("========== LOGOUT RESPONSE ==========")
	fmt.Println(string(response))
	fmt.Println("=====================================")

	return nil
}
