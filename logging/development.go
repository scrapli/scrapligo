//go:build !release
// +build !release

package logging

import (
	"fmt"
	"os"
)

var (
	// Level is the level at which to emit log messages.
	Level = Debug //nolint: gochecknoglobals

	// Logger is the main logging function, used mostly for "global" non connection/device related
	// things like the ffi layer.
	Logger = func(level LogLevel, m string, a ...any) { //nolint: gochecknoglobals
		_, _ = fmt.Fprintln(os.Stderr, level, "|", fmt.Sprintf(m, a...))
	}
)

// FfiLogger is a simple logger that can be passed as a logging callback to the zig layer.
func FfiLogger(level uint8, message *string) {
	switch level {
	case uint8(DebugAsInt):
		fmt.Println("debug :: ", *message) //nolint:forbidigo
	case uint8(InfoAsInt):
		fmt.Println(" info :: ", *message) //nolint:forbidigo
	case uint8(WarnAsInt):
		fmt.Println(" warn :: ", *message) //nolint:forbidigo
	case uint8(CriticalAsInt):
		fmt.Println(" crit :: ", *message) //nolint:forbidigo
	case uint8(FatalAsInt):
		fmt.Println("fatal :: ", *message) //nolint:forbidigo
	case uint8(DisabledAsInt):
	}
}
