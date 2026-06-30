package epp

import (
	"fmt"
	"time"
)

func (c *Client) nextTRID(prefix string) string {

	id := c.sequence.Add(1)

	return fmt.Sprintf(
		"%s-%s-%06d",
		prefix,
		time.Now().UTC().Format("20060102T150405"),
		id,
	)
}
