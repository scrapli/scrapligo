package testhelper

import (
	"os"
	"testing"

	scrapligoconstants "github.com/scrapli/scrapligo/constants"
)

// ReadFile reads the file or fatals.
func ReadFile(t *testing.T, f string) []byte {
	t.Helper()

	b, err := os.ReadFile(f) //nolint: gosec
	if err != nil {
		t.Fatal(err)
	}

	return b
}

// WriteFile writes the content to the file or fatals.
// content.
func WriteFile(t *testing.T, f string, content []byte) {
	t.Helper()

	err := os.WriteFile(
		f,
		content,
		scrapligoconstants.PermissionsOwnerReadWriteEveryoneRead,
	)
	if err != nil {
		t.Fatal(err)
	}
}
