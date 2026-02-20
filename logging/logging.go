package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ebitengine/purego"
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

// normally i *really* dislike inits but... meh?
func init() { //nolint: gochecknoinits
	v := os.Getenv(scrapligoconstants.ScrapligoDebug)

	if v != "" {
		switch v {
		case Trace.String():
			Level = Trace
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
