package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testWithNetconfPreferredVersion(
	testName string,
	testCase *optionsStringTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithNetconfPreferredVersion(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*netconf.Driver)

		if !cmp.Equal(oo.PreferredVersion, testCase.s) {
			t.Fatalf(
				"%s: actual and preferred versions do not match\nactual: %s\nexpected:%s",
				testName,
				oo.PreferredVersion,
				testCase.s,
			)
		}
	}
}

func TestWithNetconfPreferredVersion(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-netconf-preferred-version-1.0": {
			description: "simple set option test",
			s:           "1.0",
			o:           &netconf.Driver{},
			isignored:   false,
		},
		"set-netconf-preferred-version-1.1": {
			description: "simple set option test",
			s:           "1.1",
			o:           &netconf.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "1.0",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithNetconfPreferredVersion(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithNetconfForceSelfClosingTags(
	testName string,
	testCase *optionsBoolTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithNetconfForceSelfClosingTags()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*netconf.Driver)

		if !cmp.Equal(oo.ForceSelfClosingTags, testCase.b) {
			t.Fatalf(
				"%s: actual and preferred versions do not match\nactual: %v\nexpected:%v",
				testName,
				oo.ForceSelfClosingTags,
				testCase.b,
			)
		}
	}
}

func TestWithNetconfForceSelfClosingTags(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-netconf-force-self-closing-tags": {
			description: "simple set option test",
			b:           true,
			o:           &netconf.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			b:           true,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithNetconfForceSelfClosingTags(testName, testCase)

		t.Run(testName, f)
	}
}
