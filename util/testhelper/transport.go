package testhelper

import (
	"os"
	"testing"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
)

// TestingTransport patched transport for testing.
type TestingTransport struct {
	BaseTransportArgs *transport.BaseTransportArgs
	FakeSession       *os.File
}

// Open do nothing!
func (t *TestingTransport) Open() error {
	return nil
}

// OpenNetconf do nothing!
func (t *TestingTransport) OpenNetconf() error {
	return nil
}

// Close do nothing!
func (t *TestingTransport) Close() error {
	return nil
}

// Read read from the fake session.
func (t *TestingTransport) Read() ([]byte, error) {
	// need to read one byte at a time so we dont auto read past prompts and commands and such
	// its sorta strange that 65535 works for scrapli IRL but i guess its just consuming a byte at
	// a time out of a stream rather than just reading an already present file?
	b := make([]byte, 1)
	_, err := t.FakeSession.Read(b)

	return b, err
}

// Write do nothing!
func (t *TestingTransport) Write(channelInput []byte) error {
	// for now just ignoring this as the inputs are reflected already in the file that is loaded
	_ = channelInput
	return nil
}

// IsAlive do nothing!
func (t *TestingTransport) IsAlive() bool {
	return true
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (t *TestingTransport) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, t.BaseTransportArgs.Host, t.BaseTransportArgs.Port, msg)
}

// WithPatchedTransport option to use to patch a driver transport.
func WithPatchedTransport(sessionFile string, t *testing.T) base.Option {
	return func(d *base.Driver) error {
		f, err := os.Open(sessionFile)
		if err != nil {
			t.Fatalf("failed opening transport session file '%s' err: %v", sessionFile, err)
		}

		d.Transport = &TestingTransport{
			FakeSession: f,
		}

		return nil
	}
}
