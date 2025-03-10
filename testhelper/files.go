package testhelper

import (
	"os"
	"testing"

	scrapligoconstants "github.com/scrapli/scrapligo/constants"
)

// ReadFile reads the file or fatals.
func ReadFile(t *testing.T, f string) []byte {
	b, err := os.ReadFile(f) //nolint: gosec
	if err != nil {
		t.Fatal(err)
	}

	return b
}

// WriteFile writes the conteent to the file or fatals.
func WriteFile(t *testing.T, f string, content []byte) {
	err := os.WriteFile(
		f,
		content,
		scrapligoconstants.PermissionsOwnerReadWriteEveryoneRead,
	)
	if err != nil {
		t.Fatal(err)
	}
}
