package header

import (
	"encoding/json"
	"os"
)

func LoadHeaderConfig(path string) (*HeaderConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg HeaderConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
