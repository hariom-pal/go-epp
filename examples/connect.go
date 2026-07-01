package examples

import "github.com/hariom-pal/go-epp/epp"

func connectExample(configPath string) (*epp.Client, error) {
	cfg, err := epp.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return epp.Connect(cfg)
}
