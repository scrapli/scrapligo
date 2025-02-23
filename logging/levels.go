package logging

// LogLevel is an enum(ish) for log levels.
type LogLevel string

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
