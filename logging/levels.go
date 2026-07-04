package logging

// LogLevel is an enum(ish) for log levels.
type LogLevel string

// String (stringer) method for LogLevel.
func (l LogLevel) String() string {
	return string(l)
}

const (
	// Trace is the trace log level.
	Trace LogLevel = "trace"
	// Debug is the debug log level.
	Debug LogLevel = "debug"
	// Info is the info(rmational) log level.
	Info LogLevel = "info"
	// Warn is the warning log level.
	Warn LogLevel = "warn"
	// Critical is the critical log level.
	Critical LogLevel = "critical"
	// Fatal is the fatal log level.
	Fatal LogLevel = "fatal"
	// Disabled is the disabled (no logging) log level.
	Disabled LogLevel = "disabled"
)

// LogLevelAsInt is a uint8 that represents LogLevel values.
type LogLevelAsInt uint8

// LogLeveLAsInt is kept for backward compatibility.
type LogLeveLAsInt = LogLevelAsInt

// IntFromLevel returns the uint8 value of the given log level.
func IntFromLevel(level LogLevel) LogLevelAsInt {
	switch level {
	case Trace:
		return TraceAsInt
	case Debug:
		return DebugAsInt
	case Info:
		return InfoAsInt
	case Warn:
		return WarnAsInt
	case Critical:
		return CriticalAsInt
	case Fatal:
		return FatalAsInt
	case Disabled:
		return DisabledAsInt
	default:
		return DisabledAsInt
	}
}

// LevelFromInt returns the LogLevel value of the given uint8 level.
func LevelFromInt(level uint8) LogLevel {
	switch level {
	case uint8(TraceAsInt):
		return Trace
	case uint8(DebugAsInt):
		return Debug
	case uint8(InfoAsInt):
		return Info
	case uint8(WarnAsInt):
		return Warn
	case uint8(CriticalAsInt):
		return Critical
	case uint8(FatalAsInt):
		return Fatal
	case uint8(DisabledAsInt):
		return Disabled
	default:
		return Disabled
	}
}

const (
	// TraceAsInt is the debug log level.
	TraceAsInt LogLevelAsInt = 0
	// DebugAsInt is the debug log level.
	DebugAsInt LogLevelAsInt = 1
	// InfoAsInt is the info(rmational) log level.
	InfoAsInt LogLevelAsInt = 2
	// WarnAsInt is the warning log level.
	WarnAsInt LogLevelAsInt = 3
	// CriticalAsInt is the critical log level.
	CriticalAsInt LogLevelAsInt = 4
	// FatalAsInt is the fatal log level.
	FatalAsInt LogLevelAsInt = 5
	// DisabledAsInt is the disabled (no logging) log level.
	DisabledAsInt LogLevelAsInt = 6
)
