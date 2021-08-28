package transport

import (
	"errors"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

const (
	ReadSize              = 65_535
	SystemTransportName   = "system"
	StandardTransportName = "standard"
	TelnetTransportName   = "telnet"
	// MaxTimeout maximum allowable timeout value -- one day.
	MaxTimeout = 86_400
)

// SupportedTransports pseudo constant providing slice of supported transport types.
func SupportedTransports() []string {
	return []string{SystemTransportName, StandardTransportName, TelnetTransportName}
}

// SupportedNetconfTransports pseudo constant providing slice of supported netconf transport types.
func SupportedNetconfTransports() []string {
	return []string{SystemTransportName, StandardTransportName}
}

// ErrTransportFailure error for EOF/failure reading from the transport.
var ErrTransportFailure = errors.New("error reading from transport, cannot continue")

// ErrUnknownTransport error for when user provides an unknown/unsupported transport name.
var ErrUnknownTransport = errors.New("unknown transport provided")

// ErrTransportTimeout error for transport operations timing out.
var ErrTransportTimeout = errors.New("transport operation timed out")

// ErrKeyVerificationFailed ssh key verification failure.
var ErrKeyVerificationFailed = errors.New("ssh key verification failed")

// ErrUnsupportedOperation error for things like trying to use telnet transport with netconf.
var ErrUnsupportedOperation = errors.New("unsupported operation for this transport type")

// BaseTransportArgs struct for attributes that are required for any transport type.
type BaseTransportArgs struct {
	Host             string
	Port             int
	AuthUsername     string
	TimeoutSocket    time.Duration
	TimeoutTransport time.Duration
	PtyHeight        int
	PtyWidth         int
}

// ReadResult is an object used to return from the read goroutine.
type ReadResult struct {
	Result []byte
	Error  error
}

// Implementation defines an interface that a transport plugins must implement.
type Implementation interface {
	Open(baseArgs *BaseTransportArgs) error
	OpenNetconf(baseArgs *BaseTransportArgs) error
	Close() error
	IsAlive() bool
	Read(n int) *ReadResult
	Write([]byte) error
}

// Transport interface defining required methods for any transport type.
type Transport struct {
	Impl              Implementation
	BaseTransportArgs *BaseTransportArgs
}

// Open opens the transport in "normal" mode (usually telnet/ssh).
func (t *Transport) Open() error {
	err := t.Impl.Open(t.BaseTransportArgs)

	if err != nil {
		logging.LogError(
			FormatLogMessage(
				t.BaseTransportArgs,
				"error",
				"failed opening transport connection to host",
			),
		)
	} else {
		logging.LogDebug(
			FormatLogMessage(t.BaseTransportArgs, "debug", "transport connection to host opened"),
		)
	}

	return err
}

// OpenNetconf opens a netconf connection.
func (t *Transport) OpenNetconf() error {
	err := t.Impl.OpenNetconf(t.BaseTransportArgs)

	if err != nil {
		logging.LogError(
			FormatLogMessage(
				t.BaseTransportArgs,
				"error",
				"failed opening netconf transport connection to host",
			),
		)
	} else {
		logging.LogDebug(
			FormatLogMessage(t.BaseTransportArgs, "debug", "netconf transport connection to host opened"),
		)
	}

	return err
}

// Close closes the transport connection.
func (t *Transport) Close() error {
	err := t.Impl.Close()

	logging.LogDebug(
		FormatLogMessage(t.BaseTransportArgs, "debug", "transport connection to host closed"),
	)

	return err
}

// IsAlive indicates if the transport is alive or not.
func (t *Transport) IsAlive() bool {
	return t.Impl.IsAlive()
}

// Read reads bytes from the transport.
func (t *Transport) Read() ([]byte, error) {
	b, err := t.transportTimeout(
		t.Impl.Read,
		ReadSize,
	)

	if err != nil {
		logging.LogError(
			FormatLogMessage(t.BaseTransportArgs, "error", "timed out reading from transport"),
		)

		return b, err
	}

	return b, nil
}

// ReadN reads N bytes from the transport.
func (t *Transport) ReadN(n int) ([]byte, error) {
	b, err := t.transportTimeout(
		t.Impl.Read,
		n,
	)

	if err != nil {
		logging.LogError(
			FormatLogMessage(t.BaseTransportArgs, "error", "timed out reading from transport"),
		)

		return b, err
	}

	return b, nil
}

// Write writes bytes to the transport.
func (t *Transport) Write(channelInput []byte) error {
	return t.Impl.Write(channelInput)
}

func (t *Transport) transportTimeout(
	f func(int) *ReadResult,
	n int,
) ([]byte, error) {
	c := make(chan *ReadResult)

	go func() {
		r := f(n)
		c <- r
		close(c)
	}()

	timeout := t.BaseTransportArgs.TimeoutTransport
	if t.BaseTransportArgs.TimeoutTransport <= 0 {
		timeout = MaxTimeout * time.Second
	}

	timer := time.NewTimer(timeout)

	select {
	case r := <-c:
		return r.Result, r.Error
	case <-timer.C:
		return []byte{}, ErrTransportTimeout
	}
}

// FormatLogMessage formats log messages for transport level logging.
func FormatLogMessage(b *BaseTransportArgs, level, msg string) string {
	return logging.FormatLogMessage(level, b.Host, b.Port, msg)
}
