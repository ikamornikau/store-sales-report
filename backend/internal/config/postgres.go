package config

import "errors"

type Postgres struct {
	URL string `yaml:"url"`
}

func (p Postgres) Validate() error {
	if p.URL == "" {
		return errors.New("postgres URL should be set")
	}

	return nil
}
