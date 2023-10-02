package options_test

import (
	"errors"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testWithTransportReadSize(testName string, testCase *optionsIntTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithTransportReadSize(testCase.i)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.Args)

		if !cmp.Equal(oo.ReadSize, testCase.i) {
			t.Fatalf(
				"%s: actual and expected transport read sizes do not match\nactual:"+
					" %d\nexpected:%d",
				testName,
				oo.ReadSize,
				testCase.i,
			)
		}
	}
}

func TestWithTransportReadSize(t *testing.T) {
	cases := map[string]*optionsIntTestCase{
		"set-transport-read-size": {
			description: "simple set option test",
			i:           99,
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			i:           99,
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithTransportReadSize(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithPort(testName string, testCase *optionsIntTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithPort(testCase.i)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.Args)

		if !cmp.Equal(oo.Port, testCase.i) {
			t.Fatalf(
				"%s: actual and expected ports do not match\nactual: %d\nexpected:%d",
				testName,
				oo.Port,
				testCase.i,
			)
		}
	}
}

func TestWithPort(t *testing.T) {
	cases := map[string]*optionsIntTestCase{
		"set-port": {
			description: "simple set option test",
			i:           99,
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			i:           99,
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithPort(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithTermHeight(
	testName string,
	testCase *optionsIntTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithTermHeight(testCase.i)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.Args)

		if !cmp.Equal(oo.TermHeight, testCase.i) {
			t.Fatalf(
				"%s: actual and expected term heights do not match\nactual: %d\nexpected:%d",
				testName,
				oo.TermHeight,
				testCase.i,
			)
		}
	}
}

func TestWithTermHeight(t *testing.T) {
	cases := map[string]*optionsIntTestCase{
		"set-port": {
			description: "simple set option test",
			i:           99,
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			i:           99,
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithTermHeight(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithTermWidth(
	testName string,
	testCase *optionsIntTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithTermWidth(testCase.i)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.Args)

		if !cmp.Equal(oo.TermWidth, testCase.i) {
			t.Fatalf(
				"%s: actual and expected term widths do not match\nactual: %d\nexpected:%d",
				testName,
				oo.TermWidth,
				testCase.i,
			)
		}
	}
}

func TestWithTermWidth(t *testing.T) {
	cases := map[string]*optionsIntTestCase{
		"set-term-width": {
			description: "simple set option test",
			i:           99,
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			i:           99,
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithTermWidth(testName, testCase)

		t.Run(testName, f)
	}
}

func testTransportTimeoutSocket(testName string, testCase *optionsDurationTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithTimeoutSocket(testCase.d)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.Args)

		if !cmp.Equal(oo.TimeoutSocket, testCase.d) {
			t.Fatalf(
				"%s: actual and expected dial timeout do not match\nactual: %s\nexpected:%s",
				testName,
				oo.TimeoutSocket,
				testCase.d,
			)
		}
	}
}

func TestTransportTimeoutSocket(t *testing.T) {
	cases := map[string]*optionsDurationTestCase{
		"set-timeout": {
			description: "simple set option test",
			d:           123 * time.Second,
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			d:           123 * time.Second,
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testTransportTimeoutSocket(testName, testCase)

		t.Run(testName, f)
	}
}
