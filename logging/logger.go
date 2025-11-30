package logging

import (
	"fmt"
	"log"
	"log/slog"
)

// AnyLogger is a logger that wraps any of the supported logger types (log.Logger, slog.Logger or
// the callback style) to give the Cli/Netconf objects an abstracted, simple logger object.
type AnyLogger struct {
	level    LogLeveLAsInt
	trace    func(s string)
	debug    func(s string)
	info     func(s string)
	warn     func(s string)
	critical func(s string)
}

// Trace logs at the trace level.
func (l *AnyLogger) Trace(s string) {
	if l.level > TraceAsInt {
		return
	}

	l.trace(s)
}

// Debug logs at the debug level.
func (l *AnyLogger) Debug(s string) {
	if l.level > DebugAsInt {
		return
	}

	l.debug(s)
}

// Info logs at the info level.
func (l *AnyLogger) Info(s string) {
	if l.level > InfoAsInt {
		return
	}

	l.info(s)
}

// Warn logs at the warn level.
func (l *AnyLogger) Warn(s string) {
	if l.level > WarnAsInt {
		return
	}

	l.warn(s)
}

// Critical logs at the critical level.
func (l *AnyLogger) Critical(s string) {
	if l.level > CriticalAsInt {
		return
	}

	l.critical(s)
}

// LoggerToAnyLogger wraps any of the supported logger flavors in an AnyLogger so the Cli/Netconf
// objects can easily log to them w/ a consistent/simple interface.
func LoggerToAnyLogger(logger any, logLevel LogLevel) *AnyLogger {
	al := &AnyLogger{
		level:    IntFromLevel(logLevel),
		trace:    func(_ string) {},
		debug:    func(_ string) {},
		info:     func(_ string) {},
		warn:     func(_ string) {},
		critical: func(_ string) {},
	}

	switch l := logger.(type) {
	case *log.Logger:
		al.trace = func(s string) { l.Printf("trace :: %s", s) }
		al.debug = func(s string) { l.Printf("debug :: %s", s) }
		al.info = func(s string) { l.Printf(" info :: %s", s) }
		al.warn = func(s string) { l.Printf(" warn :: %s", s) }
		al.critical = func(s string) { l.Printf(" crit ::: %s", s) }
	case *slog.Logger:
		al.trace = func(s string) { l.Debug(fmt.Sprintf("trace: %s", s)) }
		al.debug = func(s string) { l.Debug(s) }
		al.info = func(s string) { l.Info(s) }
		al.warn = func(s string) { l.Warn(s) }
		al.critical = func(s string) { l.Error(s) }
	case func(LogLevel, string):
		al.trace = func(s string) { l(Trace, s) }
		al.debug = func(s string) { l(Debug, s) }
		al.info = func(s string) { l(Info, s) }
		al.warn = func(s string) { l(Warn, s) }
		al.critical = func(s string) { l(Critical, s) }
	default:
	}

	return al
}
