package logging

import "fmt"

// Logger accept logging interface to set as library logger(s).
type Logger func(...interface{})

// DebugLog default DebugLog -- defaults to nil.
var DebugLog Logger

// ErrorLog default ErrorLog -- defaults to nil.
var ErrorLog Logger

// SetDebugLogger function to set debug logger to something that implements `Logger`.
func SetDebugLogger(logger Logger) {
	DebugLog = logger
}

// SetErrorLogger function to set error logger to something that implements `Logger`.
func SetErrorLogger(logger Logger) {
	ErrorLog = logger
}

// LogDebug write debug message to the debug log.
func LogDebug(msg string) {
	if DebugLog != nil {
		DebugLog(msg)
	}
}

// LogError write error message to the error log.
func LogError(msg string) {
	if ErrorLog != nil {
		ErrorLog(msg)
	}
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func FormatLogMessage(level, host string, port int, msg string) string {
	return fmt.Sprintf("%s::%s::%d::%s", level, host, port, msg)
}
