package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/rgp"
	"github.com/hariom-pal/go-epp/types"
)

func rgpRestoreRequestExample(client *epp.Client) (*types.DomainUpdateResponse, error) {
	return client.DomainUpdate(types.DomainUpdateRequest{
		Domain: "example.in",
		RGP: &rgp.UpdateRequest{
			Restore: &rgp.Restore{
				Operation: rgp.OperationRequest,
			},
		},
	})
}
