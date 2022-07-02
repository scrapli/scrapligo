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

func testSystemTransportOpenArgs(
	testName string,
	testCase *optionsStringSliceTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithSystemTransportOpenArgs(testCase.ss)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.System)

		if !cmp.Equal(oo.ExtraArgs, testCase.ss) {
			t.Fatalf(
				"%s: actual and expected transport extra args do not match\nactual: %v\nexpected:%v",
				testName,
				oo.ExtraArgs,
				testCase.ss,
			)
		}
	}
}

func TestSystemTransportOpenArgs(t *testing.T) {
	cases := map[string]*optionsStringSliceTestCase{
		"set-system-transport-open-args": {
			description: "simple set option test",
			ss:          []string{"some", "neat", "args"},
			o:           &transport.System{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			ss:          []string{"some", "neat", "args"},
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testSystemTransportOpenArgs(testName, testCase)

		t.Run(testName, f)
	}
}
