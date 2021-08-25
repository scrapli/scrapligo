package util

import (
	"regexp"
	"strings"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?" +
	"\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

type ansiPattern struct {
	pattern *regexp.Regexp
}

var ansiPatternInstance *ansiPattern //nolint:gochecknoglobals

func getAnsiPattern() *ansiPattern {
	if ansiPatternInstance == nil {
		ansiPatternInstance = &ansiPattern{
			pattern: regexp.MustCompile(ansi),
		}
	}

	return ansiPatternInstance
}

// StrContainsAnySubStr checks string `s` for any occurrences of substrings in `l`, returns first
// found substring if any, otherwise an empty string.
func StrContainsAnySubStr(s string, l []string) string {
	for _, ss := range l {
		if strings.Contains(s, ss) {
			return s
		}
	}

	return ""
}

func ByteInSlice(b byte, s []byte) bool {
	for _, i := range s {
		if b == i {
			return true
		}
	}

	return false
}

// StripAnsi strips ansi characters from a byte slice.
func StripAnsi(b []byte) []byte {
	ap := getAnsiPattern()
	return ap.pattern.ReplaceAll(b, []byte{})
}
