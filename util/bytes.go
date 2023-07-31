package util

import (
	"bytes"
	"regexp"
	"sync"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?" +
	"\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var (
	ansiPattern     *regexp.Regexp //nolint: gochecknoglobals
	ansiPatternOnce sync.Once      //nolint: gochecknoglobals
)

func getAnsiPattern() *regexp.Regexp {
	ansiPatternOnce.Do(func() {
		ansiPattern = regexp.MustCompile(ansi)
	})

	return ansiPattern
}

// StripANSI removes ANSI escape codes from the given byte slice b.
func StripANSI(b []byte) []byte {
	return getAnsiPattern().ReplaceAll(b, []byte{})
}

// ByteIsAny checks if byte b is contained in byte slice l.
func ByteIsAny(b byte, l []byte) bool {
	for _, ss := range l {
		if b == ss {
			return true
		}
	}

	return false
}

// ByteContainsAny checks if byte slice b is contained in any byte slice in the slice of byte
// slices l.
func ByteContainsAny(b []byte, l [][]byte) bool {
	for _, ss := range l {
		if bytes.Contains(b, ss) {
			return true
		}
	}

	return false
}
