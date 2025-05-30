package logging

import (
	"fmt"
	"os"

	scrapligoconstants "github.com/scrapli/scrapligo/constants"
)

var (
	// Level is the level at which to emit log messages. When the ScrapligoDebug env var is set
	// we use that value (if it == one of the log levels, otherwise it defaults to debug), for all
	// other cases this defaults to warn.
	Level LogLevel //nolint: gochecknoglobals

	// Logger is the main logging function, used mostly for "global" non connection/device related
	// things like the ffi layer.
	Logger = func(level LogLevel, m string, a ...any) { //nolint: gochecknoglobals
		if IntFromLevel(Level) <= IntFromLevel(level) {
			_, _ = fmt.Fprintln(os.Stderr, level, "::", fmt.Sprintf(m, a...))
		}
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

// normally i *really* dislike inits but... meh?
func init() { //nolint: gochecknoinits
	v := os.Getenv(scrapligoconstants.ScrapligoDebug)
	if v != "" {
		switch v {
		case Debug.String():
			Level = Debug
		case Info.String():
			Level = Info
		case Warn.String():
			Level = Warn
		case Critical.String():
			Level = Critical
		case Fatal.String():
			Level = Fatal
		case Disabled.String():
			Level = Disabled
		default:
			Level = Debug
		}
	} else {
		Level = Warn
	}
}
