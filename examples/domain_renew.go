package examples

import (
	"time"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

func domainRenewExample(client *epp.Client) (*types.DomainRenewResponse, error) {
	return client.DomainRenew(types.DomainRenewRequest{
		DomainName:        "example.in",
		CurrentExpiryDate: time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC),
		PeriodInfo: types.Period{
			Value: 1,
			Unit:  "y",
		},
	})
}
