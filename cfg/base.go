package cfg

import (
	"errors"

	"github.com/scrapli/scrapligo/driver/network"
)

var ErrNoConfigSourcesProvided = errors.New("no configuration sources provided, cannot continue")

// Cfg primary/base cfg platform struct.
type Cfg struct {
	Conn                *network.Driver
	ConfigSources       []string
	OnPrepare           func(*network.Driver) error
	DedicatedConnection bool
	IgnoreVersion       bool

	candidateConfig   string
	getVersionCommand string
	versionString     string
	prepared          bool
}

// NewCfg returns a new instance of Cfg.
func NewCfg(
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	c := &Cfg{
		Conn:                conn,
		OnPrepare:           nil,
		DedicatedConnection: false,
		IgnoreVersion:       false,
		prepared:            false,
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if len(c.ConfigSources) == 0 {
		// if for some reason we dont have config sources we cant really do anything... this should
		// be set by the specific platform so this *shouldn't* happen but... who knows!
		return nil, ErrNoConfigSourcesProvided
	}

	return c, nil
}

// Prepare the connection.
func (d *Cfg) Prepare() error {
	return nil
}

// Cleanup the connection.
func (d *Cfg) Cleanup() error {
	return nil
}

// RenderSubstitutedConfig renders a config with provided substitutions.
func (d *Cfg) RenderSubstitutedConfig() (string, error) {
	return "", nil
}

type Platform interface {
	GetVersion() *Response
	GetConfig(source string) *Response
	LoadConfig(config string, replace bool, options ...LoadOption) *Response
	AbortConfig() *Response
	CommitConfig(source string) *Response
	DiffConfig(source string) *DiffResponse
}
