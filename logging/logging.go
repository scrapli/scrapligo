package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ebitengine/purego"
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

// LoggerToLoggerCallback wraps a given supported logger type in a callback to be passed to the
// underlying libscrapli bits.
func LoggerToLoggerCallback(logger any, logLevel uint8) uintptr { //nolint: gocyclo
	var loggerCallback uintptr

	switch l := logger.(type) {
	case *log.Logger:
		loggerCallback = purego.NewCallback(func(level uint8, message *string) {
			if logLevel > level {
				return
			}

			switch level {
			case uint8(TraceAsInt):
				l.Printf("trace :: %s", *message)
			case uint8(DebugAsInt):
				l.Printf("debug :: %s", *message)
			case uint8(InfoAsInt):
				l.Printf(" info :: %s", *message)
			case uint8(WarnAsInt):
				l.Printf(" warn :: %s", *message)
			case uint8(CriticalAsInt):
				l.Printf(" crit :: %s", *message)
			case uint8(FatalAsInt):
				l.Printf("fatal :: %s", *message)
			case uint8(DisabledAsInt):
			}
		})
	case *slog.Logger:
		loggerCallback = purego.NewCallback(func(level uint8, message *string) {
			if logLevel > level {
				return
			}

			// ignoring context things since we (currently?) expose no means to actually pass
			// a context with things here anyway
			switch level {
			case uint8(TraceAsInt):
				// no "trace" level, so... just debug it and add the trace prefix for clarity
				l.Debug(fmt.Sprintf("trace: %s", *message))
			case uint8(DebugAsInt):
				l.Debug(*message)
			case uint8(InfoAsInt):
				l.Info(*message)
			case uint8(WarnAsInt):
				l.Warn(*message)
			case uint8(CriticalAsInt):
				l.Error(*message)
			case uint8(FatalAsInt):
				l.Error(*message)
			case uint8(DisabledAsInt):
			}
		})
	case func(LogLevel, string):
		loggerCallback = purego.NewCallback(func(level uint8, message *string) {
			if logLevel > level {
				return
			}

			l(LevelFromInt(level), *message)
		})
	default:
	}

	return loggerCallback
}
