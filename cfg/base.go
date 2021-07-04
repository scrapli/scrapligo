package cfg

import (
	"errors"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/network"
)

var ErrNoConfigSourcesProvided = errors.New("no configuration sources provided, cannot continue")

// Platform -- interface describing the methods the vendor specific platforms must implement, note
// that this is also the same api surface of the Cfg object that users see.
type Platform interface {
	GetVersion() *Response
	// GetConfig(source string) *Response
	// LoadConfig(config string, replace bool, options ...LoadOption) *Response
	// AbortConfig() *Response
	// CommitConfig(source string) *Response
	// DiffConfig(source string) *DiffResponse
}

// Cfg primary/base cfg platform struct.
type Cfg struct {
	ConfigSources       []string
	OnPrepare           func(*network.Driver) error
	DedicatedConnection bool
	IgnoreVersion       bool

	CandidateConfig string
	VersionString   string
	prepared        bool

	Platform Platform
	conn     *network.Driver
}

// NewCfg returns a new instance of Cfg.
func newCfg(
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	c := &Cfg{
		OnPrepare:           nil,
		DedicatedConnection: false,
		IgnoreVersion:       false,
		prepared:            false,
		conn:                conn,
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

// GetVersion get the version from the device.
func (d *Cfg) GetVersion() *Response {
	r := d.Platform.GetVersion()

	if r.Failed {
		logging.LogError(d.formatLogMessage("debug", "failed to fetch device version"))
	}

	if r.Result == "" {
		logging.LogDebug(d.formatLogMessage("warning", "failed to parse device version"))
	}

	return r
}

func (d *Cfg) formatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, d.conn.Host, d.conn.Port, msg)
}
