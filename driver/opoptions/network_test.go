package opoptions_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testWithPrivilegeLevel(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithPrivilegeLevel(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*network.OperationOptions)

		if !cmp.Equal(oo.PrivilegeLevel, testCase.s) {
			t.Fatalf(
				"%s: actual and expected privilege levels do not match\nactual: %s\nexpected:%s",
				testName,
				oo.PrivilegeLevel,
				testCase.s,
			)
		}
	}
}

func TestWithPrivilegeLevel(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-privilege-level": {
			description: "simple set option test",
			s:           "notconfig",
			o:           &network.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "notconfig",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithPrivilegeLevel(testName, testCase)

		t.Run(testName, f)
	}
}
