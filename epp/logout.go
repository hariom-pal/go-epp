package epp

import (
	"context"
	"fmt"

	"github.com/hariom-pal/go-epp/constants"
)

// Logout sends an EPP logout command.
func (c *Client) Logout() error {
	return c.LogoutContext(context.Background())
}

// LogoutContext sends an EPP logout command.
func (c *Client) LogoutContext(ctx context.Context) error {
	logoutXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="%s">
    <command>
        <logout/>
        <clTRID>%s</clTRID>
    </command>
</epp>`, constants.EPPNamespace, c.nextTRID("LOGOUT"))

	response, err := c.ExecuteContext(ctx, []byte(logoutXML))
	if err != nil {
		return err
	}

	if err := parseCommandResponse(response); err != nil {
		return err
	}
	c.setLoggedIn(false)
	c.emit(Event{Type: EventLogout})
	return nil
}
