package examples

import "github.com/hariom-pal/go-epp/epp"

func logoutExample(client *epp.Client) error {
	return client.Logout()
}
