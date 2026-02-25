package cli

import (
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
)

// Write writes the input to the session -- this bypasses the driver operation loop in zig, use
// with caution.
func (c *Cli) Write(input string) error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	status := c.ffiMap.Session.Write(c.ptr, input, false)
	if status != 0 {
		return scrapligoerrors.NewFfiError("failed executing write", nil)
	}

	return nil
}

// WriteAndReturn writes the given input and then sends a return character -- this bypasses the
// driver operation loop in zig, use with caution.
func (c *Cli) WriteAndReturn(input string) error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	status := c.ffiMap.Session.WriteAndReturn(c.ptr, input, false)
	if status != 0 {
		return scrapligoerrors.NewFfiError("failed executing writeAndReturn", nil)
	}

	return nil
}

// WriteReturn writes a return character -- this bypasses the driver operation loop in zig, use
// with caution.
func (c *Cli) WriteReturn() error {
	if c.ptr == 0 {
		return scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	status := c.ffiMap.Session.WriteReturn(c.ptr)
	if status != 0 {
		return scrapligoerrors.NewFfiError("failed executing writeReturn", nil)
	}

	return nil
}
