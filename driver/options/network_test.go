package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"
)

func testWithNetworkOnOpen(
	testName string,
	testCase *optionsNetworkDriverOnXTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithNetworkOnOpen(testCase.f)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*network.Driver)

		if oo.OnOpen == nil {
			t.Fatalf(
				"%s: on open function is *not* set",
				testName,
			)
		}
	}
}

func TestWithNetworkOnOpen(t *testing.T) {
	cases := map[string]*optionsNetworkDriverOnXTestCase{
		"set-on-open": {
			description: "simple set option test",
			f:           func(d *network.Driver) error { return nil },
			o:           &network.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			f:           func(d *network.Driver) error { return nil },
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithNetworkOnOpen(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithNetworkOnClose(
	testName string,
	testCase *optionsNetworkDriverOnXTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithNetworkOnClose(testCase.f)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*network.Driver)

		if oo.OnClose == nil {
			t.Fatalf(
				"%s: on close function is *not* set",
				testName,
			)
		}
	}
}

func TestWithNetworkOnClose(t *testing.T) {
	cases := map[string]*optionsNetworkDriverOnXTestCase{
		"set-on-close": {
			description: "simple set option test",
			f:           func(d *network.Driver) error { return nil },
			o:           &network.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			f:           func(d *network.Driver) error { return nil },
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithNetworkOnClose(testName, testCase)

		t.Run(testName, f)
	}
}
