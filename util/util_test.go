package util_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scrapli/scrapligo/util"
)

var (
	update = flag.Bool( //nolint
		"update",
		false,
		"update the golden files",
	)
	functional = flag.Bool( //nolint
		"functional",
		false,
		"execute functional tests",
	)
	platforms = flag.String( //nolint
		"platforms",
		util.All,
		"comma sep list of platform(s) to target",
	)
	transports = flag.String( //nolint
		"transports",
		util.All,
		"comma sep list of transport(s) to target",
	)
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

	err := os.WriteFile(goldenOut, actualOut, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
}
