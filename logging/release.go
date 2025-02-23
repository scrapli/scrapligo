//go:build release
// +build release

package logging

import (
	"fmt"
	"os"
)

var (
	// Level is the level at which to emit log messages.
	Level = Warn //nolint: gochecknoglobals

	// Logger is the main logging function, used mostly for "global" non connection/device related
	// things like the ffi layer.
	Logger = func(level LogLevel, m string, a ...any) { //nolint: gochecknoglobals
		_, _ = fmt.Fprintln(os.Stderr, level, "|", fmt.Sprintf(m, a...))
	}
)
