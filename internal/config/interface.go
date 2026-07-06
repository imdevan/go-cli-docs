package config

import "go-cli-docs/internal/domain"

// Manager describes configuration persistence behavior.
type Manager interface {
	Load() (domain.Config, error)
	Save(config domain.Config) error
	Exists() (bool, error)
}
