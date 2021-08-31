package transport

import (
	"fmt"
	"net"
	"time"

	"github.com/scrapli/scrapligo/util"

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
	TelnetTransportArgs *TelnetTransportArgs
	Conn                net.Conn
	initialBuf          []byte
}

// TelnetTransportArgs struct representing attributes required for the Telnet transport.
type TelnetTransportArgs struct {
}

func (t *Telnet) handleControlCharResponse(ctrlBuf []byte, c byte) ([]byte, error) {
	if len(ctrlBuf) == 0 { //nolint:nestif
		if c != IAC {
			t.initialBuf = append(t.initialBuf, c)
		} else {
			ctrlBuf = append(ctrlBuf, c)
		}
	} else if len(ctrlBuf) == 1 && util.ByteInSlice(c, []byte{DO, DONT, WILL, WONT}) {
		ctrlBuf = append(ctrlBuf, c)
	} else if len(ctrlBuf) == 2 { //nolint:gomnd
		cmd := ctrlBuf[1:2][0]
		ctrlBuf = make([]byte, 0)

		var writeErr error

		if cmd == DO && c == SGA {
			_, writeErr = t.Conn.Write([]byte{IAC, WILL, c})
		} else if util.ByteInSlice(cmd, []byte{DO, DONT}) {
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

func (t *Telnet) handleControlChars(baseArgs *BaseTransportArgs) error {
	d := baseArgs.TimeoutSocket / 4

	var handleErr error

	ctrlBuf := make([]byte, 0)

	for {
		setDeadlineErr := t.Conn.SetReadDeadline(time.Now().Add(d))
		if setDeadlineErr != nil {
			return setDeadlineErr
		}

		// speed up timeout after initial Read
		d = baseArgs.TimeoutSocket / 10

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

func (t *Telnet) Open(baseArgs *BaseTransportArgs) error {
	var dialErr error

	t.Conn, dialErr = net.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", baseArgs.Host, baseArgs.Port),
	)
	if dialErr != nil {
		return dialErr
	}

	logging.LogDebug(FormatLogMessage(baseArgs, "debug", "tcp socket to host opened"))

	controlCharErr := t.handleControlChars(baseArgs)
	if controlCharErr != nil {
		return controlCharErr
	}

	logging.LogDebug(
		FormatLogMessage(baseArgs, "debug", "telnet control characters exchanged"),
	)

	return nil
}

func (t *Telnet) OpenNetconf(baseArgs *BaseTransportArgs) error {
	_ = baseArgs
	return ErrUnsupportedOperation
}

func (t *Telnet) Close() error {
	err := t.Conn.Close()

	t.Conn = nil

	return err
}

func (t *Telnet) IsAlive() bool {
	return t.Conn != nil
}

func (t *Telnet) Read(n int) *ReadResult {
	if len(t.initialBuf) > 0 {
		b := t.initialBuf
		t.initialBuf = []byte{}

		return &ReadResult{
			Result: b,
			Error:  nil,
		}
	}

	b := make([]byte, n)
	_, err := t.Conn.Read(b)

	if err != nil {
		return &ReadResult{
			Result: nil,
			Error:  ErrTransportFailure,
		}
	}

	return &ReadResult{
		Result: b,
		Error:  nil,
	}
}

func (t *Telnet) Write(channelInput []byte) error {
	_, err := t.Conn.Write(channelInput)

	return err
}
