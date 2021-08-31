package network

import (
	"errors"
	"regexp"
	"strings"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/base"
)

type privilegeAction string

const (
	deescalateAction privilegeAction = "deescalateAction"
	escalateAction   privilegeAction = "escalateAction"
	noAction         privilegeAction = "noAction"
)

// ErrInvalidDesiredPriv error raised when user attempts to acquire an invalid privilege level.
var ErrInvalidDesiredPriv = errors.New("invalid desired priv name")

// ErrCouldNotDeterminePriv error raised when unable to determine the current privilege level.
var ErrCouldNotDeterminePriv = errors.New("could not determine current privilege level")

// Driver driver for the "network" layer -- adds privilege levels, on open/close, and augments to
// the base driver it extends.
type Driver struct {
	base.Driver
	OnOpen      func(*Driver) error
	OnClose     func(*Driver) error
	privGraph   map[string]map[string]bool
	CurrentPriv string
	Augments    map[string]func(d *Driver) (*base.Response, error)
}

// NewNetworkDriver returns a new driver of the network flavor.
func NewNetworkDriver(
	host string,
	privilegeLevels map[string]*base.PrivilegeLevel,
	defaultDesiredPriv string,
	failedWhenContains []string,
	onOpen func(d *Driver) error,
	onClose func(d *Driver) error,
	options ...base.Option,
) (*Driver, error) {
	newDriver, err := base.NewDriver(host, options...)

	if err != nil {
		return nil, err
	}

	d := &Driver{
		Driver:   *newDriver,
		OnOpen:   onOpen,
		OnClose:  onClose,
		Augments: map[string]func(d *Driver) (*base.Response, error){},
	}

	if len(d.FailedWhenContains) == 0 {
		d.FailedWhenContains = failedWhenContains
	}

	if len(d.PrivilegeLevels) == 0 {
		d.PrivilegeLevels = privilegeLevels
	}

	if d.DefaultDesiredPriv == "" {
		d.DefaultDesiredPriv = defaultDesiredPriv
	}

	d.buildPrivGraph()
	d.generateJoinedCommsPromptPattern()

	return d, nil
}

// Open opens a connection; calls the base driver `open` method, but additionally executes the
// `OnOpen` callable.
func (d *Driver) Open() error {
	err := d.Driver.Open()
	if err != nil {
		return err
	}

	err = d.OnOpen(d)

	return err
}

// Close closes a connection; calls the base driver `close` method, but additionally executes the
// `OnClose` callable.
func (d *Driver) Close() error {
	if d.OnClose != nil {
		err := d.OnClose(d)
		if err != nil {
			logging.LogError(
				d.FormatLogMessage(
					"error",
					"encountered error on OnClose - continuing to close transport...",
				),
			)
		}
	}

	err := d.Driver.Close()

	return err
}

func (d *Driver) generateJoinedCommsPromptPattern() {
	// handle setting up the "joined" priv pattern
	allPatterns := make([]string, 0, len(d.PrivilegeLevels))
	for _, pLevel := range d.PrivilegeLevels {
		allPatterns = append(allPatterns, pLevel.Pattern)
	}

	joinedPattern := strings.Join(allPatterns, "|")

	d.Channel.CommsPromptPattern = regexp.MustCompile(joinedPattern)
}
