package v1

import (
	access "go.indent.com/apis/pkg/access/v1"
)

type ProviderConfig struct{}

type Provider interface {
	Provide(cfg ProviderConfig) access.Resources
}
