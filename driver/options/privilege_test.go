package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

type optionsPrivilegeLevelsTestCase struct {
	description string
	privs       network.PrivilegeLevels
	o           interface{}
	isignored   bool
}

func testWithPrivilegeLevels(
	testName string,
	testCase *optionsPrivilegeLevelsTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithPrivilegeLevels(testCase.privs)(testCase.o)
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

		if !cmp.Equal(oo.PrivilegeLevels["exec"].Pattern, testCase.privs["exec"].Pattern) {
			t.Fatalf(
				"%s: actual and expected privilege levels do not match\nactual: %v\nexpected:%v",
				testName,
				oo.PrivilegeLevels,
				testCase.privs,
			)
		}
	}
}

func TestWithPrivilegeLevels(t *testing.T) {
	cases := map[string]*optionsPrivilegeLevelsTestCase{
		"set-privilege-levels": {
			description: "simple set option test",
			privs: map[string]*network.PrivilegeLevel{
				"exec": {
					Pattern:        `(?im)^[\w.\-@/:]{1,63}>$`,
					Name:           "exec",
					PreviousPriv:   "",
					Deescalate:     "",
					Escalate:       "",
					EscalateAuth:   false,
					EscalatePrompt: "",
				},
			},
			o:         &network.Driver{},
			isignored: false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			privs: map[string]*network.PrivilegeLevel{
				"exec": {
					Pattern:        `(?im)^[\w.\-@/:]{1,63}>$`,
					Name:           "exec",
					PreviousPriv:   "",
					Deescalate:     "",
					Escalate:       "",
					EscalateAuth:   false,
					EscalatePrompt: "",
				},
			},
			o:         &generic.Driver{},
			isignored: true,
		},
	}

	for testName, testCase := range cases {
		f := testWithPrivilegeLevels(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithDefaultDesiredPriv(
	testName string,
	testCase *optionsStringTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithDefaultDesiredPriv(testCase.s)(testCase.o)
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

		if !cmp.Equal(oo.DefaultDesiredPriv, testCase.s) {
			t.Fatalf(
				"%s: actual and expected default desired privilege levels do not "+
					"match\nactual: %s\nexpected:%s",
				testName,
				oo.DefaultDesiredPriv,
				testCase.s,
			)
		}
	}
}

func TestWithDefaultDesiredPriv(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-default-desired-priv": {
			description: "simple set option test",
			s:           "exec",
			o:           &network.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "exec",
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithDefaultDesiredPriv(testName, testCase)

		t.Run(testName, f)
	}
}
