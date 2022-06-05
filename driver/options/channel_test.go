package options_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testWithPromptPattern(testName string, testCase *optionsRegexpTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithPromptPattern(testCase.p)(testCase.o)
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

		if !cmp.Equal(oo.PromptPattern.String(), testCase.p.String()) {
			t.Fatalf(
				"%s: actual and expected prompt patterns do not match\nactual: %s\nexpected:%s",
				testName,
				oo.PromptPattern,
				testCase.p,
			)
		}
	}
}

func TestWithPromptPattern(t *testing.T) {
	cases := map[string]*optionsRegexpTestCase{
		"set-prompt-pattern": {
			description: "simple set option test",
			p:           regexp.MustCompile("neatpattern"),
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			p:           regexp.MustCompile("neatpattern"),
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithPromptPattern(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithUsernamePattern(testName string, testCase *optionsRegexpTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithUsernamePattern(testCase.p)(testCase.o)
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

		if !cmp.Equal(oo.UsernamePattern.String(), testCase.p.String()) {
			t.Fatalf(
				"%s: actual and expected username patterns do not match\nactual: %s\nexpected:%s",
				testName,
				oo.UsernamePattern,
				testCase.p,
			)
		}
	}
}

func TestWithUsernamePattern(t *testing.T) {
	cases := map[string]*optionsRegexpTestCase{
		"set-username-pattern": {
			description: "simple set option test",
			p:           regexp.MustCompile("findthisusername!"),
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			p:           regexp.MustCompile("findthisusername!"),
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithUsernamePattern(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithPasswordPattern(testName string, testCase *optionsRegexpTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithPasswordPattern(testCase.p)(testCase.o)
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

		if !cmp.Equal(oo.PasswordPattern.String(), testCase.p.String()) {
			t.Fatalf(
				"%s: actual and expected password patterns do not match\nactual: %s\nexpected:%s",
				testName,
				oo.PasswordPattern,
				testCase.p,
			)
		}
	}
}

func TestWithPasswordPattern(t *testing.T) {
	cases := map[string]*optionsRegexpTestCase{
		"set-password-pattern": {
			description: "simple set option test",
			p:           regexp.MustCompile("helloisthatyoupasswordprompt!"),
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			p:           regexp.MustCompile("helloisthatyoupasswordprompt!"),
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithPasswordPattern(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithPassphrasePattern(
	testName string,
	testCase *optionsRegexpTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithPassphrasePattern(testCase.p)(testCase.o)
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

		if !cmp.Equal(oo.PassphrasePattern.String(), testCase.p.String()) {
			t.Fatalf(
				"%s: actual and expected passphrase patterns do not match\nactual: %s\nexpected:%s",
				testName,
				oo.PassphrasePattern,
				testCase.p,
			)
		}
	}
}

func TestWithPassphrasePattern(t *testing.T) {
	cases := map[string]*optionsRegexpTestCase{
		"set-pasphrase-pattern": {
			description: "simple set option test",
			p:           regexp.MustCompile("sneakypassphrase"),
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			p:           regexp.MustCompile("sneakypassphrase"),
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithPassphrasePattern(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithReturnChar(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithReturnChar(testCase.s)(testCase.o)
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

		if !cmp.Equal(oo.ReturnChar, []byte(testCase.s)) {
			t.Fatalf(
				"%s: actual and expected return char do not match\nactual: %s\nexpected:%s",
				testName,
				oo.ReturnChar,
				testCase.s,
			)
		}
	}
}

func TestWithReturnChar(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-return-char": {
			description: "simple set option test",
			s:           `\r\n`,
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           `\r\n`,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithReturnChar(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithTimeoutOps(testName string, testCase *optionsDurationTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithTimeoutOps(testCase.d)(testCase.o)
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

		if !cmp.Equal(oo.TimeoutOps, testCase.d) {
			t.Fatalf(
				"%s: actual and expected timeout ops do not match\nactual: %s\nexpected:%s",
				testName,
				oo.TimeoutOps,
				testCase.d,
			)
		}
	}
}

func TestWithTimeoutOps(t *testing.T) {
	cases := map[string]*optionsDurationTestCase{
		"set-timeout-ops": {
			description: "simple set option test",
			d:           123 * time.Second,
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			d:           123 * time.Second,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithTimeoutOps(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithReadDelay(testName string, testCase *optionsDurationTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithReadDelay(testCase.d)(testCase.o)
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

		if !cmp.Equal(oo.ReadDelay, testCase.d) {
			t.Fatalf(
				"%s: actual and expected read delay times do not match\nactual: %s\nexpected:%s",
				testName,
				oo.ReadDelay,
				testCase.d,
			)
		}
	}
}

func TestWithReadDelay(t *testing.T) {
	cases := map[string]*optionsDurationTestCase{
		"set-read-delay": {
			description: "simple set option test",
			d:           123 * time.Second,
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			d:           123 * time.Second,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithReadDelay(testName, testCase)

		t.Run(testName, f)
	}
}

type optionsWithChannelLog struct {
	description string
	log         io.Writer
	o           interface{}
	isignored   bool
}

func testWithChannelLog(testName string, testCase *optionsWithChannelLog) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := options.WithChannelLog(testCase.log)(testCase.o)
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

		if !cmp.Equal(fmt.Sprintf("%v", oo.ChannelLog), fmt.Sprintf("%v", testCase.log)) {
			t.Fatalf(
				"%s: actual and expected channel log objects do not match\nactual: %s\nexpected:%s",
				testName,
				oo.ChannelLog,
				testCase.log,
			)
		}
	}
}

func TestWithChannelLog(t *testing.T) {
	var logB bytes.Buffer

	cases := map[string]*optionsWithChannelLog{
		"set-channel-log": {
			description: "simple set option test",
			log:         &logB,
			o:           &channel.Channel{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			log:         &logB,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithChannelLog(testName, testCase)

		t.Run(testName, f)
	}
}
