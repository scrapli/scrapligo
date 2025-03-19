package testhelper

import (
	"regexp"
	"testing"
)

const (
	cliPasswordPattern = "(?im)^\\s+password .*$"
)

// CleanCliOutput does what it says. For testing consistency.
func CleanCliOutput(t *testing.T, output string) []byte {
	t.Helper()

	pp := regexp.MustCompile(cliPasswordPattern)

	outputBytes := []byte(output)

	outputBytes = pp.ReplaceAll(outputBytes, []byte("__PASSWORD__"))

	return outputBytes
}
