package testhelper

import (
	"os"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/transport"
)

// TestingTransport patched transport for testing.
type TestingTransport struct {
	FakeSession    *os.File
	CapturedWrites [][]byte
	ReadSize       *int
}

// SetOpenCmd implements SystemTransport SetOpenCmd here for testing purposes.
func (t *TestingTransport) SetOpenCmd(openCmd []string) {
	_ = openCmd
}

// SetExecCmd implements SystemTransport SetOpenCmd here for testing purposes.
func (t *TestingTransport) SetExecCmd(execCmd string) {
	_ = execCmd
}

// Open do nothing!
func (t *TestingTransport) Open(baseArgs *transport.BaseTransportArgs) error {
	_ = baseArgs
	return nil
}

// OpenNetconf do nothing!
func (t *TestingTransport) OpenNetconf(baseArgs *transport.BaseTransportArgs) error {
	_ = baseArgs
	return nil
}

// Close do nothing!
func (t *TestingTransport) Close() error {
	return nil
}

func (t *TestingTransport) Read(n int) *transport.ReadResult {
	_ = n
	readSize := 1

	if t.ReadSize != nil {
		readSize = *t.ReadSize
	}

	b := make([]byte, readSize)
	_, err := t.FakeSession.Read(b)

	return &transport.ReadResult{Result: b, Error: err}
}

func (t *TestingTransport) Write(channelInput []byte) error {
	t.CapturedWrites = append(t.CapturedWrites, channelInput)
	return nil
}

// IsAlive do nothing!
func (t *TestingTransport) IsAlive() bool {
	return true
}

// WithPatchedTransport option to use to patch a driver transport.
func WithPatchedTransport(sessionFile string) base.Option {
	return func(o interface{}) error {
		f, err := os.Open(sessionFile)
		if err != nil {
			return err
		}

		d, ok := o.(*base.Driver)

		if ok {
			d.Transport = &transport.Transport{
				Impl: &TestingTransport{
					FakeSession: f,
				},
				BaseTransportArgs: &transport.BaseTransportArgs{
					Host: "localhost",
					Port: 22,
				},
			}
		}

		return base.ErrIgnoredOption
	}
}
