package testhelper

import (
	"regexp"
	"testing"
)

const (
	timestampPattern = "\\d{4}-\\d{2}-\\d{2}T\\d+:\\d+:\\d+.\\d+Z"
	sessionIDPattern = "<session-id>\\d+</session-id>"
)

// CleanNetconfOutput does what it says. For testing consistency.
func CleanNetconfOutput(t *testing.T, output string) []byte {
	t.Helper()

	tp := regexp.MustCompile(timestampPattern)
	sp := regexp.MustCompile(sessionIDPattern)

	outputBytes := []byte(output)

	outputBytes = tp.ReplaceAll(outputBytes, []byte("__TIMESTAMP__"))
	outputBytes = sp.ReplaceAll(outputBytes, []byte("__SESSIONID__"))

	return outputBytes
}
