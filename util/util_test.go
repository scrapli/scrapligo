package util_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	scrapligoconstants "github.com/scrapli/scrapligo/constants"
)

var update = flag.Bool( //nolint: gochecknoglobals
	"update",
	false,
	"update the golden files",
)

func readFile(t *testing.T, f string) []byte {
	b, err := os.ReadFile(fmt.Sprintf("./test-fixtures/%s", f))
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func writeGolden(t *testing.T, testName string, actualOut []byte) {
	goldenOut := filepath.Join("test-fixtures", "golden", testName+"-out.txt")

	err := os.WriteFile(
		goldenOut,
		actualOut,
		scrapligoconstants.PermissionsOwnerReadWriteEveryoneRead,
	)
	if err != nil {
		t.Fatal(err)
	}
}
