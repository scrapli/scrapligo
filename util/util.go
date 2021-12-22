package util

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/scrapli/scrapligo/logging"
)

// ErrFileNotFound error for being unable to find requested file.
var ErrFileNotFound = errors.New("file not found")

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

// StrInSlice checks for any occurrence of `s` in slice of strings `l`. Returns true if `s` found,
// otherwise false.
func StrInSlice(s string, l []string) bool {
	for _, i := range l {
		if s == i {
			return true
		}
	}

	return false
}

// BytesContainsAnySubBytes checks byte `b` for any occurrences of substrings in `l`, returns first
// found substring if any, otherwise an empty string.
func BytesContainsAnySubBytes(b []byte, l [][]byte) []byte {
	for _, ss := range l {
		if bytes.Contains(b, ss) {
			return b
		}
	}

	return []byte{}
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

func ResolveFilePath(f string) (string, error) {
	if _, err := os.Stat(f); err == nil {
		return f, nil
	}

	// if didn't stat a fully qualified file, strip user dir (if exists) and then check there
	f = strings.TrimPrefix(f, "~/")
	homeDir, err := os.UserHomeDir()

	if err != nil {
		logging.LogError(fmt.Sprintf("couldnt determine users home directory: %v", err))

		return "", err
	}

	f = fmt.Sprintf("%s/%s", homeDir, f)

	if _, err = os.Stat(f); err == nil {
		return f, nil
	}

	return "", ErrFileNotFound
}

// LoadFileLines convenience function to load a file and return slice of strings of lines in that
// file.
func LoadFileLines(f string) ([]string, error) {
	resolvedFile, err := ResolveFilePath(f)

	if err != nil {
		return []string{}, err
	}

	file, err := os.Open(resolvedFile)
	if err != nil {
		return []string{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
