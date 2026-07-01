package examples

import (
	"github.com/hariom-pal/go-epp/epp"
)

func loginExample(configPath string) (*epp.Client, error) {
	cfg, err := epp.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	client, err := epp.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := client.Login(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}
