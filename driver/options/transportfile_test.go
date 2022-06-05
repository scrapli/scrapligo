package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testFileTransportFile(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithFileTransportFile(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.File)

		if !cmp.Equal(oo.F, testCase.s) {
			t.Fatalf(
				"%s: actual and expected transport files do not match\nactual: %s\nexpected:%s",
				testName,
				oo.F,
				testCase.s,
			)
		}
	}
}

func TestFileTransportFile(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-file-transport-file": {
			description: "simple set option test",
			s:           "some dumb file",
			o:           &transport.File{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "some dumb file",
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testFileTransportFile(testName, testCase)

		t.Run(testName, f)
	}
}
