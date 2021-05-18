package base

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/scrapli/scrapligo/logging"
)

// ErrFileNotFound error for being unable to find requested file.
var ErrFileNotFound = errors.New("file not found")

func resolveFilePath(f string) (string, error) {
	if _, err := os.Stat(f); err == nil {
		return f, nil
	}

	// if didnt stat a fully qualified file, strip user dir (if exists) and then check there
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
	resolvedFile, err := resolveFilePath(f)

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
