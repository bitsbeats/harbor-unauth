package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bitsbeats/harbor-unauth/core"
)

func LoadConfig() (*core.Config, error) {
	path := "/etc/harbor-unauth.json"
	if overwrite, ok := os.LookupEnv("UNAUTH_CONFIG"); ok {
		path = overwrite
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open config: %w", err)
	}

	config := &core.Config{}
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}
	return config, nil
}
