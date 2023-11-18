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

type optionsWithTwoStringTestCase struct {
	description string
	ks          string
	ps          string
	o           interface{}
	isignored   bool
}

func testWithAuthPrivateKey(
	testName string,
	testCase *optionsWithTwoStringTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthPrivateKey(testCase.ks, testCase.ps)(testCase.o)
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

		if !cmp.Equal(oo.PrivateKeyPath, testCase.ks) {
			t.Fatalf(
				"%s: actual and expected private key paths do not match\nactual: %v\nexpected:%v",
				testName,
				oo.PrivateKeyPath,
				testCase.ks,
			)
		}

		if !cmp.Equal(oo.PrivateKeyPassPhrase, testCase.ps) {
			t.Fatalf(
				"%s: actual and expected private key paths do not match\nactual: %v\nexpected:%v",
				testName,
				oo.PrivateKeyPath,
				testCase.ks,
			)
		}
	}
}

func TestWithAuthPrivateKey(t *testing.T) {
	cases := map[string]*optionsWithTwoStringTestCase{
		"set-transport-auth-private-key": {
			description: "simple set option test",
			ks:          "/key/path/thingy",
			ps:          "keypassphrase",
			o:           &transport.SSHArgs{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			ks:          "/key/path/thingy",
			ps:          "keypassphrase",
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthPrivateKey(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithAuthNoStrictKey(testName string, testCase *optionsBoolTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithAuthNoStrictKey()(testCase.o)
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

		if !cmp.Equal(oo.StrictKey, testCase.b) {
			t.Fatalf(
				"%s: actual and expected strict keys do not match\nactual: %v\nexpected:%v",
				testName,
				oo.StrictKey,
				testCase.b,
			)
		}
	}
}

func TestWithAuthNoStrictKey(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-transport-auth-strict-key": {
			description: "simple set option test",
			b:           false,
			o:           &transport.SSHArgs{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			b:           true,
			o:           &generic.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithAuthNoStrictKey(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithSSHConfigFile(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithSSHConfigFile(testCase.s)(testCase.o)
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

		oo, _ := testCase.o.(*transport.SSHArgs)

		if !cmp.Equal(oo.ConfigFile, testCase.s) {
			t.Fatalf(
				"%s: actual and expected resolved ssh config files do not match\nactual:"+
					" %v\nexpected:%v", //nolint:goconst
				testName,
				oo.ConfigFile,
				testCase.s,
			)
		}
	}
}

func TestWithSSHConfigFile(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-ssh-config-file": {
			description: "simple set option test",
			s:           "transportssh_test.go", // option resolves the file, so make it a real one
			o:           &transport.SSHArgs{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "transportssh_test.go", // option resolves the file, so make it a real one
			o:           &generic.Driver{},
			isignored:   true,
		},
		"set-ssh-config-file-cant-resolve": {
			description: "option fails to resolves provided file path",
			s:           "not/a/real/file",
			o:           &transport.SSHArgs{},
			isignored:   false,
			iserr:       true,
		},
	}

	for testName, testCase := range cases {
		f := testWithSSHConfigFile(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithSSHKnownHostsFile(
	testName string,
	testCase *optionsStringTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithSSHKnownHostsFile(testCase.s)(testCase.o)
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

		oo, _ := testCase.o.(*transport.SSHArgs)

		if !cmp.Equal(oo.KnownHostsFile, testCase.s) {
			t.Fatalf(
				"%s: actual and expected resolved known hosts files do not match\nactual:"+
					" %v\nexpected:%v",
				testName,
				oo.KnownHostsFile,
				testCase.s,
			)
		}
	}
}

func TestWithSSHKnownHostsFile(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-ssh-known-hosts-file": {
			description: "simple set option test",
			s:           "transportssh_test.go", // option resolves the file, so make it a real one
			o:           &transport.SSHArgs{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "transportssh_test.go", // option resolves the file, so make it a real one
			o:           &generic.Driver{},
			isignored:   true,
		},
		"set-ssh-known-hosts-file-cant-resolve": {
			description: "option fails to resolves provided file path",
			s:           "not/a/real/file",
			o:           &transport.SSHArgs{},
			isignored:   false,
			iserr:       true,
		},
	}

	for testName, testCase := range cases {
		f := testWithSSHKnownHostsFile(testName, testCase)

		t.Run(testName, f)
	}
}

// not testing system ssh config file option yet due to not wanting to mess with patching
// file paths somehow (does something like pyfakefs exist in go?)
