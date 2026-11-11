package config

import (
	"errors"
	"time"
)

type Server struct {
	Address           string        `yaml:"address"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
}

func (c Server) Validate() error {
	if c.Address == "" {
		return errors.New("address should be set")
	}

	if c.ReadHeaderTimeout <= 0 {
		return errors.New("read_header_timeout should be greater than 0")
	}

	if c.ReadTimeout <= 0 {
		return errors.New("read_timeout should be greater than 0")
	}

	if c.WriteTimeout <= 0 {
		return errors.New("write_timeout should be greater than 0")
	}

	if c.IdleTimeout <= 0 {
		return errors.New("idle_timeout should be greater than 0")
	}

	return nil
}
