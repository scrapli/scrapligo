package cfg

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/scrapli/scrapligo/driver/base"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/network"
)

const (
	runningConfig = "running"
)

var ErrNoConfigSourcesProvided = errors.New("no configuration sources provided, cannot continue")
var ErrInvalidConfigTarget = errors.New("provided config source is not valid")

var ErrConnectionNotOpen = errors.New(
	"underlying scrapli connection is not open and `DedicatedConnection` is false, cannot continue")
var ErrVersionError = errors.New("failed getting or parsing device version")

var ErrPrepareNotCalled = errors.New("the Prepare method has not been called, cannot continue")

var ErrConfigSessionAlreadyExists = errors.New(
	"configuration session name already exists, cannot create it")

var ErrInvalidSource = errors.New("invalid config source/target provided")

var ErrFailedToDetermineDeviceState = errors.New(
	"failed to determine device state, ex: failed to determine available space")
var ErrInsufficientSpaceAvailable = errors.New("insufficient space available on device")

var ErrGetConfigFailed = errors.New("get config operation failed")
var ErrLoadConfigFailed = errors.New("load config operation failed")
var ErrAbortConfigFailed = errors.New("abort config operation failed")
var ErrCommitConfigFailed = errors.New("commit config operation failed")
var ErrDiffConfigFailed = errors.New("diff config operation failed")

// Platform -- interface describing the methods the vendor specific platforms must implement, note
// that this is also similar (but not the same!) to the same api surface of the Cfg object that
// users see.
type Platform interface {
	GetVersion() (string, []*base.Response, error)
	GetConfig(source string) (string, []*base.Response, error)
	LoadConfig(config string, replace bool, options *OperationOptions) ([]*base.Response, error)
	AbortConfig() ([]*base.Response, error)
	CommitConfig(source string) ([]*base.Response, error)
	DiffConfig(source, candidateConfig string) ([]*base.Response, string, string, string, error)
	ClearConfigSession()
}

func FormatLogMessage(conn *network.Driver, level, msg string) string {
	return logging.FormatLogMessage(level, conn.Host, conn.Transport.BaseTransportArgs.Port, msg)
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

func parseOperationOptions(o []OperationOption) *OperationOptions {
	opts := &OperationOptions{Source: "running", DiffColorize: true, AutoClean: true}

	if len(o) > 0 && o[0] != nil {
		for _, option := range o {
			option(opts)
		}
	}

	return opts
}

func determineCandidateConfigFilename(c string) string {
	candidateConfigFilename := c

	if c == "" {
		candidateConfigFilename = fmt.Sprintf("scrapli_cfg_%d", time.Now().Unix())
	}

	return candidateConfigFilename
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

func (d *Cfg) invalidConfigSource(err error) error {
	logging.LogError(
		FormatLogMessage(
			d.conn,
			"error",
			"invalid configuration source",
		),
	)

	return err
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

func (d *Cfg) validateAndSetVersion(versionResponse *Response) error {
	if versionResponse.Failed != nil {
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
		logging.LogDebug(FormatLogMessage(d.conn, "debug", "OnPrepare provided, executing now"))

		prepareErr := d.OnPrepare(d.conn)
		if prepareErr != nil {
			return prepareErr
		}
	}

	d.prepared = true

	return nil
}

func (d *Cfg) clearConfigSession() {
	logging.LogDebug(
		FormatLogMessage(
			d.conn,
			"debug",
			"resetting config session data",
		),
	)

	d.CandidateConfig = ""
	d.Platform.ClearConfigSession()
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

// Cleanup cleans up the cfg session.
func (d *Cfg) Cleanup() error {
	err := d.close()
	if err != nil {
		return err
	}

	d.VersionString = ""
	d.prepared = false

	d.clearConfigSession()

	return nil
}

// RenderSubstitutedConfig renders a config with provided substitutions.
func (d *Cfg) RenderSubstitutedConfig() (string, error) {
	return "", nil
}

func (d *Cfg) configSourceValid(source string) bool {
	for _, configSource := range d.ConfigSources {
		if configSource == source {
			return true
		}
	}

	return false
}

// GetVersion gets the version from the device.
func (d *Cfg) GetVersion() (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(d.conn, "info", "get version requested"),
	)

	r := NewResponse(d.conn.Host, "GetVersion", nil)

	versionString, scrapliResponses, err := d.Platform.GetVersion()

	r.Record(scrapliResponses, versionString)

	if r.Failed != nil {
		logging.LogDebug(FormatLogMessage(d.conn, "warning", "failed to fetch device version"))
	}

	if r.Result == "" {
		logging.LogDebug(FormatLogMessage(d.conn, "warning", "failed to parse device version"))
	}

	return r, err
}

// GetConfig gets the configuration of a source datastore from the device.
func (d *Cfg) GetConfig(source string) (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(
			d.conn,
			"info",
			fmt.Sprintf("get config requested for config source '%s'", source),
		),
	)

	r := NewResponse(d.conn.Host, "GetConfig", ErrGetConfigFailed)

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	if !d.configSourceValid(source) {
		return r, d.invalidConfigSource(ErrInvalidSource)
	}

	cfgString, scrapliResponses, err := d.Platform.GetConfig(source)

	r.Record(scrapliResponses, cfgString)

	if r.Failed != nil {
		logging.LogError(FormatLogMessage(d.conn, "debug", "failed to fetch config from device"))
	}

	return r, err
}

// LoadConfig loads a candidate configuration.
func (d *Cfg) LoadConfig(
	config string,
	replace bool,
	options ...OperationOption,
) (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(d.conn, "info", "load config requested"),
	)

	opts := parseOperationOptions(options)

	d.CandidateConfig = config
	r := NewResponse(d.conn.Host, "LoadConfig", ErrLoadConfigFailed)

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	scrapliResponses, err := d.Platform.LoadConfig(config, replace, opts)

	r.Record(scrapliResponses, "")

	if r.Failed != nil {
		logging.LogError(
			FormatLogMessage(d.conn, "error", "failed to load candidate configuration"),
		)
	}

	return r, err
}

// LoadConfigFromFile loads a candidate configuration from a provided file.
func (d *Cfg) LoadConfigFromFile(
	f string,
	replace bool,
	options ...OperationOption,
) (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(d.conn, "info", "load config from file requested"),
	)

	c, err := util.LoadFileLines(f)
	if err != nil {
		return nil, err
	}

	return d.LoadConfig(strings.Join(c, "\n"), replace, options...)
}

// AbortConfig aborts the loaded candidate configuration.
func (d *Cfg) AbortConfig() (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(d.conn, "info", "abort config requested"),
	)

	r := NewResponse(d.conn.Host, "AbortConfig", ErrAbortConfigFailed)

	if d.CandidateConfig == "" {
		logging.LogError(
			FormatLogMessage(
				d.conn,
				"error",
				"no candidate configuration exists, you must load a config in order to abort it!",
			),
		)

		return r, ErrAbortConfigFailed
	}

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	scrapliResponses, err := d.Platform.AbortConfig()

	r.Record(scrapliResponses, "")

	if r.Failed != nil {
		logging.LogError(
			FormatLogMessage(d.conn, "error", "failed to abort candidate configuration"),
		)
	}

	d.clearConfigSession()

	return r, err
}

// CommitConfig commits the loaded candidate configuration.
func (d *Cfg) CommitConfig(options ...OperationOption) (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(d.conn, "info", "commit config requested"),
	)

	opts := parseOperationOptions(options)

	r := NewResponse(d.conn.Host, "CommitConfig", ErrCommitConfigFailed)

	if d.CandidateConfig == "" {
		logging.LogError(
			FormatLogMessage(
				d.conn,
				"error",
				"no candidate configuration exists, you must load a config in order to commit it!",
			),
		)

		return r, ErrCommitConfigFailed
	}

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	if !d.configSourceValid(opts.Source) {
		return r, d.invalidConfigSource(ErrInvalidSource)
	}

	scrapliResponses, err := d.Platform.CommitConfig(opts.Source)

	r.Record(scrapliResponses, "")

	if r.Failed != nil {
		logging.LogError(
			FormatLogMessage(d.conn, "error", "failed to commit candidate configuration"),
		)
	}

	d.clearConfigSession()

	return r, err
}

// DiffConfig diffs the candidate configuration against a source config.
func (d *Cfg) DiffConfig(options ...OperationOption) (*DiffResponse, error) {
	logging.LogDebug(
		FormatLogMessage(d.conn, "info", "diff config requested"),
	)

	opts := parseOperationOptions(options)

	r := NewDiffResponse(d.conn.Host, opts.Source, opts.DiffColorize, opts.DiffSideBySideWidth)

	operationOkErr := d.operationOk()
	if operationOkErr != nil {
		return r, operationOkErr
	}

	if d.CandidateConfig == "" {
		logging.LogError(
			FormatLogMessage(
				d.conn,
				"error",
				"no candidate configuration exists, you must load a config in order to diff it!",
			),
		)

		return r, ErrDiffConfigFailed
	}

	if !d.configSourceValid(opts.Source) {
		return r, d.invalidConfigSource(ErrDiffConfigFailed)
	}

	scrapliResponses, sourceConfig, candidateConfig, deviceDiff, err := d.Platform.DiffConfig(
		opts.Source,
		d.CandidateConfig,
	)

	r.Record(scrapliResponses, "")
	r.RecordDiff(sourceConfig, candidateConfig, deviceDiff)

	if r.Failed != nil {
		logging.LogError(
			FormatLogMessage(d.conn, "error", "failed to diff configuration"),
		)
	}

	return r, err
}
