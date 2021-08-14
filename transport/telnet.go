package transport

import (
	"fmt"
	"net"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

const (
	IAC  = byte(255)
	DONT = byte(254)
	DO   = byte(253)
	WONT = byte(252)
	WILL = byte(251)
	SGA  = byte(3)
)

// Telnet the telnet transport option for scrapligo.
type Telnet struct {
	BaseTransportArgs   *BaseTransportArgs
	TelnetTransportArgs *TelnetTransportArgs
	Conn                net.Conn
	initialBuf          []byte
}

// TelnetTransportArgs struct representing attributes required for the Telnet transport.
type TelnetTransportArgs struct {
}

func byteInSlice(b byte, s []byte) bool {
	for _, i := range s {
		if b == i {
			return true
		}
	}

	return false
}

func (t *Telnet) handleControlCharResponse(ctrlBuf []byte, c byte) ([]byte, error) {
	if len(ctrlBuf) == 0 { //nolint:nestif
		if c != IAC {
			t.initialBuf = append(t.initialBuf, c)
		} else {
			ctrlBuf = append(ctrlBuf, c)
		}
	} else if len(ctrlBuf) == 1 && byteInSlice(c, []byte{DO, DONT, WILL, WONT}) {
		ctrlBuf = append(ctrlBuf, c)
	} else if len(ctrlBuf) == 2 { //nolint:gomnd
		cmd := ctrlBuf[1:2][0]
		ctrlBuf = make([]byte, 0)

		var writeErr error

		if cmd == DO && c == SGA {
			_, writeErr = t.Conn.Write([]byte{IAC, WILL, c})
		} else if byteInSlice(cmd, []byte{DO, DONT}) {
			_, writeErr = t.Conn.Write([]byte{IAC, WONT, c})
		} else if cmd == WILL {
			_, writeErr = t.Conn.Write([]byte{IAC, DO, c})
		} else if cmd == WONT {
			_, writeErr = t.Conn.Write([]byte{IAC, DONT, c})
		}

		if writeErr != nil {
			return nil, writeErr
		}
	}

	return ctrlBuf, nil
}

func (t *Telnet) handleControlChars() error {
	socketTimeout := t.BaseTransportArgs.TimeoutSocket
	d := *socketTimeout / 4

	var handleErr error

	ctrlBuf := make([]byte, 0)

	for {
		setDeadlineErr := t.Conn.SetReadDeadline(time.Now().Add(d))
		if setDeadlineErr != nil {
			return setDeadlineErr
		}

		// speed up timeout after initial read
		d = *socketTimeout / 10

		charBuf := make([]byte, 1)

		_, readErr := t.Conn.Read(charBuf)
		if readErr != nil { //nolint:nestif
			if opErr, ok := readErr.(*net.OpError); ok {
				if opErr.Timeout() {
					// timeout is good -- we want to be done reading control chars, so cancel the
					// deadline by setting it to "zero"
					cancelDeadlineErr := t.Conn.SetReadDeadline(time.Time{})
					if cancelDeadlineErr != nil {
						return cancelDeadlineErr
					}

					return nil
				}

				return opErr
			}

			return readErr
		}

		ctrlBuf, handleErr = t.handleControlCharResponse(ctrlBuf, charBuf[0])
		if handleErr != nil {
			return handleErr
		}
	}
}

// Open opens a telnet connection.
func (t *Telnet) Open() error {
	var dialErr error

	t.Conn, dialErr = net.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", t.BaseTransportArgs.Host, t.BaseTransportArgs.Port),
	)
	if dialErr != nil {
		return dialErr
	}

	logging.LogDebug(t.FormatLogMessage("debug", "tcp socket to host opened"))

	controlCharErr := t.handleControlChars()
	if controlCharErr != nil {
		return controlCharErr
	}

	logging.LogDebug(t.FormatLogMessage("debug", "telnet control characters exchanged"))

	return nil
}

// OpenNetconf returns an error, netconf does not support telnet... duh.
func (t *Telnet) OpenNetconf() error {
	return ErrUnsupportedOperation
}

// Close closes the transport connection to the device.
func (t *Telnet) Close() error {
	err := t.Conn.Close()

	t.Conn = nil
	logging.LogDebug(t.FormatLogMessage("debug", "transport connection to host closed"))

	return err
}

func (t *Telnet) read(n int) *transportResult {
	if len(t.initialBuf) > 0 {
		b := t.initialBuf
		t.initialBuf = []byte{}

		return &transportResult{
			result: b,
			error:  nil,
		}
	}

	b := make([]byte, n)
	_, err := t.Conn.Read(b)

	if err != nil {
		return &transportResult{
			result: nil,
			error:  ErrTransportFailure,
		}
	}

	return &transportResult{
		result: b,
		error:  nil,
	}
}

// Read reads bytes from the transport.
func (t *Telnet) Read() ([]byte, error) {
	b, err := transportTimeout(
		*t.BaseTransportArgs.TimeoutTransport,
		t.read,
		ReadSize,
	)

	if err != nil {
		logging.LogError(t.FormatLogMessage("error", "timed out reading from transport"))
		return b, err
	}

	return b, nil
}

// ReadN reads N bytes from the transport.
func (t *Telnet) ReadN(n int) ([]byte, error) {
	b, err := transportTimeout(
		*t.BaseTransportArgs.TimeoutTransport,
		t.read,
		n,
	)

	if err != nil {
		logging.LogError(t.FormatLogMessage("error", "timed out reading from transport"))
		return b, err
	}

	return b, nil
}

// Write writes bytes to the transport.
func (t *Telnet) Write(channelInput []byte) error {
	_, err := t.Conn.Write(channelInput)
	if err != nil {
		return err
	}

	return nil
}

// IsAlive indicates if the transport is alive or not.
func (t *Telnet) IsAlive() bool {
	return t.Conn != nil
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (t *Telnet) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, t.BaseTransportArgs.Host, t.BaseTransportArgs.Port, msg)
}
