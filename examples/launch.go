package examples

import (
	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/extensions/launch"
	"github.com/hariom-pal/go-epp/types"
)

func launchApplicationCreateExample(client *epp.Client) (*types.DomainCreateResponse, error) {
	return client.DomainCreate(types.DomainCreateRequest{
		Domain:     "example.in",
		Period:     1,
		Unit:       "y",
		Registrant: "REG123",
		AuthInfo:   "change-me",
		Launch: &launch.CreateRequest{
			Type: launch.ObjectApplication,
			Phase: launch.Phase{
				Value: launch.PhaseSunrise,
			},
			CodeMarks: []launch.CodeMark{
				{
					Code: &launch.Code{
						Value:       "49FD46E6C4B45C55D4AC",
						ValidatorID: "tmch",
					},
				},
			},
		},
	})
}
