package transport

import (
	"os"
)

const (
	// FileTransport transport name.
	FileTransport = "file"
)

// NewFileTransport returns an instance of File transport. This is for testing purposes only.
func NewFileTransport() (*File, error) {
	t := &File{
		fd: nil,
	}

	return t, nil
}

// File transport is a transport object that "connects" to a file rather than a device, it probably
// has no use outside of testing.
type File struct {
	F  string
	fd *os.File

	Writes [][]byte
}

// Open opens the File transport.
func (t *File) Open(a *Args) error {
	_ = a

	f, err := os.Open(t.F)
	if err != nil {
		return err
	}

	t.fd = f

	return nil
}

// Close is a noop for the File transport.
func (t *File) Close() error {
	return nil
}

// IsAlive always returns true for File transport.
func (t *File) IsAlive() bool {
	return true
}

// Read reads n bytes from the transport. File transport ignores EOF errors, see comment below.
func (t *File) Read(n int) ([]byte, error) {
	b := make([]byte, n)

	// we don't care about errors here, only one we really would get is EOF and since this should
	// only ever be used for testing we can ignore that. in some test situations we could read too
	// fast and not d-queue things fast enough, so we basically hit the EOF w/out actually "finding"
	// the prompt we are looking for which of course causes issues. again, shouldn't matter for
	// anything "real" since this is just for testing!
	_, _ = t.fd.Read(b)

	return b, nil
}

// Write writes bytes b to the transport.
func (t *File) Write(b []byte) error {
	t.Writes = append(t.Writes, b)

	return nil
}
