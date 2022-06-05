package transport

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/util"
)

const (
	// DefaultTransport is the default transport constant for scrapligo, this defaults to the
	// "system" transport.
	DefaultTransport = "system"

	defaultPort                    = 22
	defaultTimeoutSocketSeconds    = 30
	defaultTimeoutTransportSeconds = 30
	defaultReadSize                = 65_535
	defaultTermHeight              = 255
	defaultTermWidth               = 80
	defaultSSHStrictKey            = true
	tcp                            = "tcp"
)

// GetTransportNames is returns a slice of available transport type names.
func GetTransportNames() []string {
	return []string{SystemTransport, StandardTransport, TelnetTransport}
}

// GetNetconfTransportNames returns a slice of available NETCONF transport type names.
func GetNetconfTransportNames() []string {
	return []string{SystemTransport, StandardTransport}
}

// NewArgs returns an instance of Args with the logging instance, host, and any provided args
// set. Users should *generally* not need to call this function as this is called during Transport
// creation (which is called by the Driver creation).
func NewArgs(l *logging.Instance, host string, options ...util.Option) (*Args, error) {
	a := &Args{
		l:                l,
		Host:             host,
		Port:             defaultPort,
		TimeoutSocket:    defaultTimeoutSocketSeconds * time.Second,
		TimeoutTransport: defaultTimeoutTransportSeconds * time.Second,
		ReadSize:         defaultReadSize,
		TermHeight:       defaultTermHeight,
		TermWidth:        defaultTermWidth,
	}

	for _, option := range options {
		err := option(a)
		if err != nil {
			if !errors.Is(err, util.ErrIgnoredOption) {
				return nil, err
			}
		}
	}

	return a, nil
}

// Args is a struct representing common transport arguments.
type Args struct {
	l                *logging.Instance
	Host             string
	Port             int
	User             string
	Password         string
	TimeoutSocket    time.Duration
	TimeoutTransport time.Duration
	ReadSize         int
	TermHeight       int
	TermWidth        int
}

// NewSSHArgs returns an instance of SSH arguments with provided options set. Just like NewArgs,
// this should generally not be called by users directly.
func NewSSHArgs(options ...util.Option) (*SSHArgs, error) {
	a := &SSHArgs{
		StrictKey: defaultSSHStrictKey,
	}

	for _, option := range options {
		err := option(a)
		if err != nil {
			if !errors.Is(err, util.ErrIgnoredOption) {
				return nil, err
			}
		}
	}

	return a, nil
}

// SSHArgs is a struct representing common transport SSH specific arguments.
type SSHArgs struct {
	StrictKey            bool
	PrivateKeyPath       string
	PrivateKeyPassPhrase string
	ConfigFile           string
	KnownHostsFile       string
	NetconfConnection    bool
}

// NewTelnetArgs returns an instance of TelnetArgs with any provided options set. This should,
// just like the other NewXArgs functions, not be called directly by users.
func NewTelnetArgs(options ...util.Option) (*TelnetArgs, error) {
	a := &TelnetArgs{}

	for _, option := range options {
		err := option(a)
		if err != nil {
			if !errors.Is(err, util.ErrIgnoredOption) {
				return nil, err
			}
		}
	}

	return a, nil
}

// TelnetArgs is a struct representing common transport Telnet specific arguments.
type TelnetArgs struct{}

type transportImpl interface {
	Open(a *Args) error
	Close() error
	IsAlive() bool
	Read(n int) ([]byte, error)
	Write(b []byte) error
}

// Transport is a struct which wraps a transportImpl object and provides a unified interface to any
// type of transport selected by the user.
type Transport struct {
	Args        *Args
	Impl        transportImpl
	timeoutLock *sync.Mutex
}

// SetTimeoutTransport is a convenience function that acquires a lock on the timeoutLock and updates
// the timeout transport value. Without this we may encounter data races when reading the timeout
// value in the read loop (which is ultimately in a goroutine in the Channel read loop), and trying
// to update the timeout value from elsewhere.
func (t *Transport) SetTimeoutTransport(timeout time.Duration) {
	t.timeoutLock.Lock()
	defer t.timeoutLock.Unlock()

	t.Args.TimeoutTransport = timeout
}

// GetTimeoutTransport fetches the current TimeoutTransport value for the Transport. See also
// SetTimeoutTransport.
func (t *Transport) GetTimeoutTransport() time.Duration {
	t.timeoutLock.Lock()
	defer t.timeoutLock.Unlock()

	return t.Args.TimeoutTransport
}

// Open opens the underlying transportImpl transport object.
func (t *Transport) Open() error {
	return t.Impl.Open(t.Args)
}

// Close closes the underlying transportImpl transport object.
func (t *Transport) Close() error {
	return t.Impl.Close()
}

// IsAlive returns true if the underlying transportImpl reports liveness, otherwise false.
func (t *Transport) IsAlive() bool {
	return t.Impl.IsAlive()
}

func (t *Transport) read(n int) ([]byte, error) {
	c := make(chan *readResult)

	go func() {
		b, err := t.Impl.Read(n)

		c <- &readResult{
			r:   b,
			err: err,
		}
	}()

	timeout := t.GetTimeoutTransport()

	if timeout <= 0 {
		t.Args.l.Debug("transport timeout is 0, using max timeout")

		timeout = util.MaxTimeout * time.Second
	}

	timer := time.NewTimer(timeout)

	select {
	case r := <-c:
		return r.r, r.err
	case <-timer.C:
		t.Args.l.Critical("timed out reading from transport")

		return nil, fmt.Errorf("%w: timed out reading from transport", util.ErrTimeoutError)
	}
}

// Read reads the Transport Args ReadSize bytes from the transportImpl.
func (t *Transport) Read() ([]byte, error) {
	return t.read(t.Args.ReadSize)
}

// ReadN reads n bytes from the transportImpl.
func (t *Transport) ReadN(n int) ([]byte, error) {
	return t.read(n)
}

// Write writes bytes b to the transportImpl.
func (t *Transport) Write(b []byte) error {
	return t.Impl.Write(b)
}

// GetHost is a convenience method to return the Transport Args Host value.
func (t *Transport) GetHost() string {
	return t.Args.Host
}

// GetPort is a convenience method to return the Transport Args Port value.
func (t *Transport) GetPort() int {
	return t.Args.Port
}

type readResult struct {
	r   []byte
	err error
}
