package logging

// LogLevel is an enum(ish) for log levels.
type LogLevel string

// String (stringer) method for LogLevel.
func (l LogLevel) String() string {
	return string(l)
}

const (
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

// LogLeveLAsInt is a uint8 that represents LogLevel values.
type LogLeveLAsInt uint8

// IntFromLevel returns the uint8 value of the given log level.
func IntFromLevel(level LogLevel) LogLeveLAsInt {
	switch level {
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
	// DebugAsInt is the debug log level.
	DebugAsInt LogLeveLAsInt = 0
	// InfoAsInt is the info(rmational) log level.
	InfoAsInt LogLeveLAsInt = 1
	// WarnAsInt is the warning log level.
	WarnAsInt LogLeveLAsInt = 2
	// CriticalAsInt is the critical log level.
	CriticalAsInt LogLeveLAsInt = 3
	// FatalAsInt is the fatal log level.
	FatalAsInt LogLeveLAsInt = 4
	// DisabledAsInt is the disabled (no logging) log level.
	DisabledAsInt LogLeveLAsInt = 5
)
