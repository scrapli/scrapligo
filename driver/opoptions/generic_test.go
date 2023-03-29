package opoptions_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/util"
)

func testWithStopOnFailed(testName string, testCase *optionsBoolTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithStopOnFailed()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error %s",
					testName, err,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.OperationOptions)

		if !cmp.Equal(oo.StopOnFailed, testCase.b) {
			t.Fatalf(
				"%s: actual and expected stop on faileds do not match\nactual: %v\nexpected:%v",
				testName,
				oo.StopOnFailed,
				testCase.b,
			)
		}
	}
}

func TestWithStopOnFailed(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-stop-on-failed": {
			description: "simple set option test",
			b:           true,
			o:           &generic.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			b:           true,
			o:           &channel.OperationOptions{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithStopOnFailed(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithFailedWhenContains(
	testName string,
	testCase *optionsStringSliceTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithFailedWhenContains(testCase.ss)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.OperationOptions)

		if !cmp.Equal(oo.FailedWhenContains, testCase.ss) {
			t.Fatalf(
				"%s: actual and expected failed when contains values do not match\nactual:"+
					" %s\nexpected:%s",
				testName,
				oo.FailedWhenContains,
				testCase.ss,
			)
		}
	}
}

func TestWithFailedWhenContains(t *testing.T) {
	cases := map[string]*optionsStringSliceTestCase{
		"set-failed-when-contains": {
			description: "simple set option test with an invalid transport type",
			ss:          []string{"failed", "when", "contains"},
			o:           &generic.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			ss:          []string{"failed", "when", "contains"},
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithFailedWhenContains(testName, testCase)

		t.Run(testName, f)
	}
}
