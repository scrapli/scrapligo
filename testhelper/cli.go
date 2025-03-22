package testhelper

import (
	"regexp"
	"testing"
)

const (
	cliTimestampPattern = "(?im)((mon)|(tue)|(wed)|(thu)|(fri)|(sat)|(sun))\\s+((jan)|(feb)|(mar)|(apr)|(may)|(jun)|(jul)|(aug)|(sep)|(oct)|(nov)|(dec))\\s+\\d+\\s+\\d+:\\d+:\\d+ \\d+" //nolint:lll
	cliPasswordPattern  = "(?im)^\\s+password .*$"                                                                                                                                       //nolint:gosec
)

// CleanCliOutput does what it says. For testing consistency.
func CleanCliOutput(t *testing.T, output string) []byte {
	t.Helper()

	uhp := regexp.MustCompile(cliUserAtHostPattern)
	tp := regexp.MustCompile(cliTimestampPattern)
	pp := regexp.MustCompile(cliPasswordPattern)

	outputBytes := []byte(output)

	outputBytes = uhp.ReplaceAll(outputBytes, []byte("user@host"))
	outputBytes = tp.ReplaceAll(outputBytes, []byte("Mon Jan 1 00:00:00 2025"))
	outputBytes = pp.ReplaceAll(outputBytes, []byte("__PASSWORD__"))

	return outputBytes
}
