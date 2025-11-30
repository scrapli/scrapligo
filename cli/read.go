package cli

import (
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

const (
	defaultReadSize = 1_024
)

func newReadOptions(options ...Option) *readOptions {
	o := &readOptions{
		size: defaultReadSize,
	}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type readOptions struct {
	size uint64
}

// Read reads from the session -- this bypasses the driver operation loop in zig, use with caution.
// Also note that this does not accept a context because it is completely non-blocking -- this
// drains the read buffer, and if there is nothing to read we just return 0 bytes read.
func (c *Cli) Read(options ...Option) ([]byte, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	loadedOptions := newReadOptions(options...)

	buf := make([]byte, loadedOptions.size)

	var readSize uint64

	status := c.ffiMap.Session.Read(c.ptr, &buf, &readSize)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed executing read", nil)
	}

	return buf[0:readSize], nil
}
