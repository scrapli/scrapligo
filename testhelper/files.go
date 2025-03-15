package testhelper

import (
	"os"
	"testing"

	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
)

// ReadFile reads the file or fatals.
func ReadFile(t *testing.T, f string) []byte {
	b, err := os.ReadFile(f) //nolint: gosec
	if err != nil {
		t.Fatal(err)
	}

	return b
}

// WriteFile writes the content to the file or fatals. Also strips any ascii/ansi bits out the
// content.
func WriteFile(t *testing.T, f string, content []byte) {
	sContent, err := scrapligoffi.StripASCIIAndAnsiControlCharsInPlace(string(content))
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		f,
		[]byte(sContent),
		scrapligoconstants.PermissionsOwnerReadWriteEveryoneRead,
	)
	if err != nil {
		t.Fatal(err)
	}
}
