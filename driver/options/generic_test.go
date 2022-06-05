package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testWithTransportType(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithTransportType(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			if errors.Is(err, util.ErrIgnoredOption) && testCase.isignored {
				return
			}

			if !errors.Is(err, util.ErrIgnoredOption) && testCase.iserr {
				return
			}

			t.Fatalf("%s: encountered an error but we should not", testName)
		}

		oo, _ := testCase.o.(*generic.Driver)

		if !cmp.Equal(oo.TransportType, testCase.s) {
			t.Fatalf(
				"%s: actual and expected transport types do not match\nactual: %s\nexpected:%s",
				testName,
				oo.TransportType,
				testCase.s,
			)
		}
	}
}

func TestWithTransportType(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-transport-type": {
			description: "simple set option test",
			s:           transport.StandardTransport,
			o:           &generic.Driver{},
			isignored:   false,
		},
		"set-transport-type-invalid": {
			description: "simple set option test with an invalid transport type",
			s:           "potato",
			o:           &generic.Driver{},
			isignored:   false,
			iserr:       true,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           transport.StandardTransport,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithTransportType(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithFailedWhenContains(
	testName string,
	testCase *optionsStringSliceTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithFailedWhenContains(testCase.ss)(testCase.o)
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

		if !cmp.Equal(oo.FailedWhenContains, testCase.ss) {
			t.Fatalf(
				"%s: actual and expected failed when contains values do not match\nactual: %s\nexpected:%s",
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
			description: "simple set option test",
			ss:          []string{"failed", "when", "contains"},
			o:           &generic.Driver{},
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

func testWithOnOpen(testName string, testCase *optionsGenericDriverOnXTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithOnOpen(testCase.f)(testCase.o)
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

		if oo.OnOpen == nil {
			t.Fatalf(
				"%s: on open function is *not* set",
				testName,
			)
		}
	}
}

func TestWithOnOpen(t *testing.T) {
	cases := map[string]*optionsGenericDriverOnXTestCase{
		"set-on-open": {
			description: "simple set option test",
			f:           func(d *generic.Driver) error { return nil },
			o:           &generic.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			f:           func(d *generic.Driver) error { return nil },
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithOnOpen(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithOnClose(
	testName string,
	testCase *optionsGenericDriverOnXTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithOnClose(testCase.f)(testCase.o)
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

		if oo.OnClose == nil {
			t.Fatalf(
				"%s: on close function is *not* set",
				testName,
			)
		}
	}
}

func TestWithOnClose(t *testing.T) {
	cases := map[string]*optionsGenericDriverOnXTestCase{
		"set-on-close": {
			description: "simple set option test",
			f:           func(d *generic.Driver) error { return nil },
			o:           &generic.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			f:           func(d *generic.Driver) error { return nil },
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithOnClose(testName, testCase)

		t.Run(testName, f)
	}
}
