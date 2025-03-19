package testhelper

import (
	"regexp"
	"testing"
)

const (
	netconfTimestampPattern = "\\d{4}-\\d{2}-\\d{2}T\\d+:\\d+:\\d+.\\d+Z"
	netconfSessionIDPattern = "<session-id>\\d+</session-id>"
	netconfPasswordPattern  = "<password>.*</password>"
)

// CleanNetconfOutput does what it says. For testing consistency.
func CleanNetconfOutput(t *testing.T, output string) []byte {
	t.Helper()

	tp := regexp.MustCompile(netconfTimestampPattern)
	sp := regexp.MustCompile(netconfSessionIDPattern)
	pp := regexp.MustCompile(netconfPasswordPattern)

	outputBytes := []byte(output)

	outputBytes = tp.ReplaceAll(outputBytes, []byte("__TIMESTAMP__"))
	outputBytes = sp.ReplaceAll(outputBytes, []byte("__SESSIONID__"))
	outputBytes = pp.ReplaceAll(outputBytes, []byte("__PASSWORD__"))

	return outputBytes
}
