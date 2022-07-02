package options_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testWithAuthUsername(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthUsername(testCase.s)(testCase.o)
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

		if !cmp.Equal(oo.User, testCase.s) {
			t.Fatalf(
				"%s: actual and expected usernames do not match\nactual: %s\nexpected:%s",
				testName,
				oo.User,
				testCase.s,
			)
		}
	}
}

func TestWithAuthUsername(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-auth-username": {
			description: "simple set option test",
			s:           "userperson",
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "userperson",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthUsername(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithAuthPassword(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthPassword(testCase.s)(testCase.o)
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

		if !cmp.Equal(oo.Password, testCase.s) {
			t.Fatalf(
				"%s: actual and expected passwords do not match\nactual: %s\nexpected:%s",
				testName,
				oo.Password,
				testCase.s,
			)
		}
	}
}

func TestWithAuthPassword(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-auth-password": {
			description: "simple set option test",
			s:           "supergoodpassword",
			o:           &transport.Args{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "userperson",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthPassword(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithAuthSecondary(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthSecondary(testCase.s)(testCase.o)
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

		if !cmp.Equal(oo.AuthSecondary, testCase.s) {
			t.Fatalf(
				"%s: actual and expected secondary passwords do not match\nactual: %s\nexpected:%s",
				testName,
				oo.AuthSecondary,
				testCase.s,
			)
		}
	}
}

func TestWithAuthSecondary(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-auth-secondary": {
			description: "simple set option test",
			s:           "supergoodsecondarypassword",
			o:           &network.Driver{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "userperson",
			o:           &transport.Args{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthSecondary(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithAuthPassphrase(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthPassphrase(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*transport.SSHArgs)

		if !cmp.Equal(oo.PrivateKeyPassPhrase, testCase.s) {
			t.Fatalf(
				"%s: actual and expected passphrases do not match\nactual: %s\nexpected:%s",
				testName,
				oo.PrivateKeyPassPhrase,
				testCase.s,
			)
		}
	}
}

func TestWithAuthPassphrase(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-auth-passphrase": {
			description: "simple set option test",
			s:           "verysecurepassphrase",
			o:           &transport.SSHArgs{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "userperson",
			o:           &transport.Args{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthPassphrase(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithAuthBypass(testName string, testCase *optionsBoolTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthBypass()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*channel.Channel)

		if !cmp.Equal(oo.AuthBypass, testCase.b) {
			t.Fatalf(
				"%s: actual and expected auth bypass do not match\nactual: %v\nexpected:%v",
				testName,
				oo.AuthBypass,
				testCase.b,
			)
		}
	}
}

func TestWithAuthBypass(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-auth-bypass": {
			description: "simple set option test",
			b:           true,
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			b:           true,
			o:           &transport.Args{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthBypass(testName, testCase)

		t.Run(testName, f)
	}
}
