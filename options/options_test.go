package options_test

import (
	"fmt"
	"os"
	"testing"

	scrapligoffi "github.com/kentik/scrapligo/v2/ffi"
	scrapligotesthelper "github.com/kentik/scrapligo/v2/testhelper"
)

func TestMain(m *testing.M) {
	scrapligotesthelper.Flags()

	exitCode := m.Run()

	if scrapligoffi.AssertNoLeaks() != nil {
		_, _ = fmt.Fprintln(os.Stderr, "memory leak(s) detected!")

		os.Exit(127)
	}

	_, _ = fmt.Fprintln(os.Stderr, "no memory leak(s) detected!")

	os.Exit(exitCode)
}
