package cfg

import (
	"errors"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/network"
)

var ErrNoConfigSourcesProvided = errors.New("no configuration sources provided, cannot continue")
var ErrInvalidConfigTarget = errors.New("provided config source is not valid")

var ErrConnectionNotOpen = errors.New(
	"underlying scrapli connection is not open and `DedicatedConnection` is false, cannot continue")
var ErrVersionError = errors.New("failed getting or parsing device version")

var ErrPrepareNotCalled = errors.New("the Prepare method has not been called, cannot continue")

var ErrConfigSessionAlreadyExists = errors.New(
	"configuration session name already exists, cannot create it")

// Platform -- interface describing the methods the vendor specific platforms must implement, note
// that this is also the same api surface of the Cfg object that users see.
type Platform interface {
	GetVersion() (string, []*base.Response, error)
	GetConfig(source string) (string, []*base.Response, error)
	LoadConfig(config string, replace bool, options ...LoadOption) ([]*base.Response, error)
	// AbortConfig() ([]*base.Response, error)
	// CommitConfig(source string) ([]*base.Response, error)
	// DiffConfig(source string) *DiffResponse
}

// simple interface to allow for checking if platform implements ClearConfigSession.
type platformWithClearConfigSession interface {
	ClearConfigSession()
}

func FormatLogMessage(conn *network.Driver, level, msg string) string {
	return logging.FormatLogMessage(level, conn.Host, conn.Port, msg)
}

func setPlatformOptions(p Platform, options ...Option) error {
	for _, option := range options {
		err := option(p)

		if err != nil {
			if errors.Is(err, ErrIgnoredOption) {
				continue
			} else {
				return err
			}
		}
	}

	return nil
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

// newCfg returns a new instance of Cfg; private because users should be calling the platform
// specific new functions (or using the factory).
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
		err := option(c)

		if err != nil {
			if errors.Is(err, ErrIgnoredOption) {
				continue
			} else {
				return nil, err
			}
		}
	}

	if len(c.ConfigSources) == 0 {
		// if for some reason we dont have config sources we cant really do anything... this should
		// be set by the specific platform so this *shouldn't* happen but... who knows!
		return nil, ErrNoConfigSourcesProvided
	}

	return c, nil
}

func (d *Cfg) prepareOk() error {
	if d.OnPrepare != nil && !d.prepared {
		logging.LogError(
			FormatLogMessage(
				d.conn,
				"error",
				"OnPrepare provided, but prepare method not called. call prepare method prior "+
					"to using the Cfg object",
			),
		)

		return ErrPrepareNotCalled
	}

	return nil
}

func (d *Cfg) versionOk() error {
	if !d.IgnoreVersion && d.VersionString == "" {
		logging.LogError(
			FormatLogMessage(
				d.conn,
				"error",
				"IgnoreVersion is false, but version has not yet been fetched. call prepare method prior "+
					"to using the Cfg object to ensure version is properly gathered",
			),
		)

		return ErrPrepareNotCalled
	}

	return nil
}

func (d *Cfg) operationOk() error {
	prepareErr := d.prepareOk()
	if prepareErr != nil {
		return prepareErr
	}

	versionErr := d.versionOk()

	if versionErr != nil {
		return versionErr
	}

	return nil
}

func (d *Cfg) open() error {
	if d.conn.Transport.IsAlive() {
		// nothing to do, connection is already open!
		return nil
	}

	if d.DedicatedConnection {
		err := d.conn.Open()
		return err
	}

	return ErrConnectionNotOpen
}

func (d *Cfg) validateAndSetVersion(versionResponse *Response) error {
	if versionResponse.Failed {
		logging.LogError(FormatLogMessage(d.conn, "error", "failed getting version from device"))
		return ErrVersionError
	}

	if versionResponse.Result == "" {
		logging.LogError(
			FormatLogMessage(d.conn, "error", "failed parsing version string from device output"),
		)

		return ErrVersionError
	}

	d.VersionString = versionResponse.Result

	return nil
}

// Prepare the connection.
func (d *Cfg) Prepare() error {
	logging.LogDebug(FormatLogMessage(d.conn, "info", "preparing cfg connection"))

	err := d.open()
	if err != nil {
		return err
	}

	if !d.IgnoreVersion {
		logging.LogDebug(
			FormatLogMessage(d.conn, "debug", "IgnoreVersion is false, fetching device version"),
		)

		versionResponse, getVersionErr := d.GetVersion()

		if getVersionErr != nil {
			return getVersionErr
		}

		validateVersionErr := d.validateAndSetVersion(versionResponse)
		if validateVersionErr != nil {
			return validateVersionErr
		}
	}

	if d.OnPrepare != nil {
		logging.LogDebug(FormatLogMessage(d.conn, "debug", "OnPrepare provided, exucting now"))

		prepareErr := d.OnPrepare(d.conn)
		if prepareErr != nil {
			return prepareErr
		}
	}

	d.prepared = true

	return nil
}

func (d *Cfg) close() error {
	if d.DedicatedConnection && d.conn.Transport.IsAlive() {
		logging.LogDebug(
			FormatLogMessage(
				d.conn,
				"info",
				"DedicatedConnection is true, closing scrapli connection",
			),
		)

		err := d.conn.Close()

		return err
	}

	return nil
}

// Cleanup the connection.
func (d *Cfg) Cleanup() error {
	err := d.close()
	if err != nil {
		return err
	}

	d.VersionString = ""
	d.CandidateConfig = ""
	d.prepared = false

	platform, ok := interface{}(d.Platform).(platformWithClearConfigSession)
	if ok {
		platform.ClearConfigSession()
	}

	return nil
}

// RenderSubstitutedConfig renders a config with provided substitutions.
func (d *Cfg) RenderSubstitutedConfig() (string, error) {
	return "", nil
}

// GetVersion get the version from the device.
func (d *Cfg) GetVersion() (*Response, error) {
	r := NewResponse(d.conn.Host, nil)

	versionString, scrapliResponses, err := d.Platform.GetVersion()

	r.Record(scrapliResponses, versionString)

	if r.Failed {
		logging.LogDebug(FormatLogMessage(d.conn, "warning", "failed to fetch device version"))
	}

	if r.Result == "" {
		logging.LogDebug(FormatLogMessage(d.conn, "warning", "failed to parse device version"))
	}

	return r, err
}

// GetConfig get the configuration of a source datastore from the device.
func (d *Cfg) GetConfig(source string) (*Response, error) {
	r := NewResponse(d.conn.Host, nil)

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	cfgString, scrapliResponses, err := d.Platform.GetConfig(source)

	r.Record(scrapliResponses, cfgString)

	if r.Failed {
		logging.LogError(FormatLogMessage(d.conn, "debug", "failed to fetch config from device"))
	}

	return r, err
}

// LoadConfig load a candidate configuration.
func (d *Cfg) LoadConfig(config string, replace bool, options ...LoadOption) (*Response, error) {
	d.CandidateConfig = config
	r := NewResponse(d.conn.Host, nil)

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	scrapliResponses, err := d.Platform.LoadConfig(config, replace, options...)

	r.Record(scrapliResponses, "")

	if r.Failed {
		logging.LogError(
			FormatLogMessage(d.conn, "error", "failed loading candidate configuration"),
		)
	}

	return r, err
}
