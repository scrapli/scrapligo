package logging

import "fmt"

// Logger accepts logging interface to set as library logger(s).
type Logger func(...interface{})

// debugLog default DebugLog -- defaults to nil.
var debugLog Logger //nolint:gochecknoglobals

// errorLog default ErrorLog -- defaults to nil.
var errorLog Logger

// SetDebugLogger function to set debug logger to something that implements `Logger`.
func SetDebugLogger(logger Logger) {
	debugLog = logger
}

// SetErrorLogger function to set error logger to something that implements `Logger`.
func SetErrorLogger(logger Logger) {
	errorLog = logger
}

// LogDebug writes debug message to the debug log.
func LogDebug(msg string) {
	if debugLog != nil {
		debugLog(msg)
	}
}

// LogError writes error message to the error log.
func LogError(msg string) {
	if errorLog != nil {
		errorLog(msg)
	}
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func FormatLogMessage(level, host string, port int, msg string) string {
	return fmt.Sprintf("%s::%s::%d::%s", level, host, port, msg)
}
