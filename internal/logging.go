package internal

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"

	"github.com/ebitengine/purego"
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligologging "github.com/scrapli/scrapligo/v2/logging"
)

var (
	loggingDispatcherInst     *loggerDispatcher //nolint: gochecknoglobals
	loggingDispatcherInstOnce sync.Once         //nolint: gochecknoglobals
)

type loggerCallbackF func(level uint8, message string)

// LoggerDispatcher is the interface returned when fetching the singleton loggerDispatcher. The
// dispatcher exists because in purego we have a max of 2000 callbacks -- every logger func was
// a callback historically, so on a long running program or a program with just a lot of connections
// it was possible to quickly eat up those 2000 slots -- moreover they are *never* freed up -- so
// once you use it you use it. So this interface lets us have a single logger callback passed to
// libscrapli, and it properly dispatches logs based on the id which is the pointer to the Cli or
// Netconf object.
type LoggerDispatcher interface {
	Register(userData uintptr, logger any, logLevel scrapligologging.LogLevel) error
	Deregister(userData uintptr)

	GetLoggerCallback() uintptr
}

// GetLoggerDispatcher returns the LoggerDispatcher singleton.
func GetLoggerDispatcher() LoggerDispatcher { //nolint: ireturn
	loggingDispatcherInstOnce.Do(
		func() {
			loggingDispatcherInst = &loggerDispatcher{
				lock:    sync.RWMutex{},
				loggers: map[uintptr]loggerCallbackF{},
			}

			cb := purego.NewCallback(loggingDispatcherInst.log)

			loggingDispatcherInst.cb = cb
		},
	)

	return loggingDispatcherInst
}

type loggerDispatcher struct {
	lock    sync.RWMutex
	loggers map[uintptr]loggerCallbackF
	cb      uintptr
}

func (l *loggerDispatcher) Register( //nolint: gocyclo
	userData uintptr,
	logger any,
	logLevel scrapligologging.LogLevel,
) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	configuredLevel := uint8(scrapligologging.IntFromLevel(logLevel))

	switch tl := logger.(type) {
	case *log.Logger:
		l.loggers[userData] = func(level uint8, message string) {
			if configuredLevel > level {
				return
			}

			switch level {
			case uint8(scrapligologging.TraceAsInt):
				tl.Printf("trace :: %s", message)
			case uint8(scrapligologging.DebugAsInt):
				tl.Printf("debug :: %s", message)
			case uint8(scrapligologging.InfoAsInt):
				tl.Printf(" info :: %s", message)
			case uint8(scrapligologging.WarnAsInt):
				tl.Printf(" warn :: %s", message)
			case uint8(scrapligologging.CriticalAsInt):
				tl.Printf(" crit :: %s", message)
			case uint8(scrapligologging.FatalAsInt):
				tl.Printf("fatal :: %s", message)
			case uint8(scrapligologging.DisabledAsInt):
			}
		}
	case *slog.Logger:
		l.loggers[userData] = func(level uint8, message string) {
			if configuredLevel > level {
				return
			}

			// ignoring context things since we (currently?) expose no means to actually pass
			// a context with things here anyway
			switch level {
			case uint8(scrapligologging.TraceAsInt):
				// no "trace" level, so... just debug it and add the trace prefix for clarity
				tl.Debug(fmt.Sprintf("trace: %s", message))
			case uint8(scrapligologging.DebugAsInt):
				tl.Debug(message)
			case uint8(scrapligologging.InfoAsInt):
				tl.Info(message)
			case uint8(scrapligologging.WarnAsInt):
				tl.Warn(message)
			case uint8(scrapligologging.CriticalAsInt):
				tl.Error(message)
			case uint8(scrapligologging.FatalAsInt):
				tl.Error(message)
			case uint8(scrapligologging.DisabledAsInt):
			}
		}
	case func(scrapligologging.LogLevel, string):
		l.loggers[userData] = func(level uint8, message string) {
			if configuredLevel > level {
				return
			}

			tl(scrapligologging.LevelFromInt(level), message)
		}
	case nil:
		return nil
	default:
		return scrapligoerrors.NewFfiError(fmt.Sprintf("invalid logger type %T", tl), nil)
	}

	return nil
}

func (l *loggerDispatcher) Deregister(userData uintptr) {
	l.lock.Lock()
	defer l.lock.Unlock()

	delete(l.loggers, userData)
}

func (l *loggerDispatcher) GetLoggerCallback() uintptr {
	return l.cb
}

func (l *loggerDispatcher) log(userData uintptr, level uint8, message *string) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	lf, ok := l.loggers[userData]
	if !ok {
		return
	}

	msg := strings.Clone(*message)

	lf(level, msg)
}
