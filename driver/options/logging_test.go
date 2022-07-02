package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

type withLoggerTestCase struct {
	description string
	l           *logging.Instance
	o           interface{}
	isignored   bool
}

func testWithLogger(testName string, testCase *withLoggerTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithLogger(testCase.l)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Driver)

		if !cmp.Equal(oo.Logger, testCase.l) {
			t.Fatalf(
				"%s: actual and expected loggers do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Logger,
				testCase.l,
			)
		}
	}
}

func TestWithLogger(t *testing.T) {
	cases := map[string]*withLoggerTestCase{
		"set-logger": {
			description: "simple set option test",
			l:           &logging.Instance{},
			o:           &generic.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			l:           &logging.Instance{},
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithLogger(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithDefaultLogger(testName string, testCase *optionsNoneTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithDefaultLogger()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Driver)

		if oo.Logger == nil {
			t.Fatalf(
				"%s: default logger not set",
				testName,
			)
		}
	}
}

func TestWithDefaultLogger(t *testing.T) {
	cases := map[string]*optionsNoneTestCase{
		"set-default-logger": {
			description: "simple set option test",
			o:           &generic.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithDefaultLogger(testName, testCase)

		t.Run(testName, f)
	}
}
