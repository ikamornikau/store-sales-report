package main

import (
	"fmt"

	"github.com/ikamornikau/store-sales-report/backend/internal/config"
)

type Config struct {
	Server   config.Server   `yaml:"server"`
	Postgres config.Postgres `yaml:"postgres"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("failed to validate server config: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("failed to validate postgres config: %w", err)
	}

	return nil
}
