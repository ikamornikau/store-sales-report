package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

func Load[T any](path string) (T, error) {
	var config T

	body, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(body, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return config, nil
}
