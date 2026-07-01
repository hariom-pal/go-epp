package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func helloExample(client *epp.Client) (*types.Greeting, error) {
	return client.Hello()
}
