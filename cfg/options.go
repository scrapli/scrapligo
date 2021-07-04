package cfg

import (
	"github.com/scrapli/scrapligo/driver/network"
)

// Option function to set cfg platform options.
type Option func(*Cfg) error

// WithConfigSources modify the default config sources for your platform.
func WithConfigSources(sources []string) Option {
	return func(c *Cfg) error {
		c.ConfigSources = sources
		return nil
	}
}

// WithOnPrepare provide an OnPrepare callable for the Cfg instance.
func WithOnPrepare(onPrepare func(*network.Driver) error) Option {
	return func(c *Cfg) error {
		c.OnPrepare = onPrepare
		return nil
	}
}

// WithDedicatedConnection set dedicated connection for Cfg instance.
func WithDedicatedConnection(dedicatedConnection bool) Option {
	return func(c *Cfg) error {
		c.DedicatedConnection = dedicatedConnection
		return nil
	}
}

// WithIgnoreVersion set ignore version for Cfg instance.
func WithIgnoreVersion(ignoreVersion bool) Option {
	return func(c *Cfg) error {
		c.IgnoreVersion = ignoreVersion
		return nil
	}
}

// LoadOptions struct for LoadConfig options.
type LoadOptions struct {
}

// LoadOption function to set options for cfg LoadConfig operations.
type LoadOption func(*LoadOptions) error
