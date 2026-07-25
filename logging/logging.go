package logging

import (
	"fmt"
	"os"

	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
)

var (
	// Level is the level at which to emit log messages.
	Level = defaultLevel() //nolint:gochecknoglobals

	// Logger is the main logging function, used mostly for "global" non connection/device related
	// things like the ffi layer.
	Logger = func(level LogLevel, m string, a ...any) { //nolint: gochecknoglobals
		if IntFromLevel(Level) <= IntFromLevel(level) {
			_, _ = fmt.Fprintln(os.Stderr, level, "::", fmt.Sprintf(m, a...))
		}
	}
)

func defaultLevel() LogLevel {
	v := os.Getenv(scrapligoconstants.ScrapligoDebug)
	if v == "" {
		return Warn
	}

	switch v {
	case Trace.String():
		return Trace
	case Debug.String():
		return Debug
	case Info.String():
		return Info
	case Warn.String():
		return Warn
	case Critical.String():
		return Critical
	case Fatal.String():
		return Fatal
	case Disabled.String():
		return Disabled
	default:
		return Debug
	}
}
