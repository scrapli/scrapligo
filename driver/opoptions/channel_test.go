package opoptions_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/channel"
	"github.com/scrapli/scrapligo/util"
)

func testWithStripPrompt(testName string, testCase *optionsBoolTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithNoStripPrompt()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*channel.OperationOptions)

		if !cmp.Equal(oo.StripPrompt, testCase.b) {
			t.Fatalf(
				"%s: actual and expected strip prompts do not match\nactual: %v\nexpected:%v",
				testName,
				oo.StripPrompt,
				testCase.b,
			)
		}
	}
}

func TestWithStripPrompt(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-strip-prompt": {
			description: "simple set option test",
			b:           false,
			o:           &channel.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			b:           false,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithStripPrompt(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithEager(testName string, testCase *optionsBoolTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithEager()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*channel.OperationOptions)

		if !cmp.Equal(oo.Eager, testCase.b) {
			t.Fatalf(
				"%s: actual and expected eagers do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Eager,
				testCase.b,
			)
		}
	}
}

func TestWithEager(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-eager": {
			description: "simple set option test",
			b:           true,
			o:           &channel.OperationOptions{},
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
		f := testWithEager(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithTimeout(testName string, testCase *optionsDurationTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithTimeoutOps(testCase.d)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		var timeout time.Duration

		switch oo := testCase.o.(type) {
		case *channel.OperationOptions:
			timeout = oo.Timeout
		case *netconf.OperationOptions:
			timeout = oo.Timeout
		}

		if !cmp.Equal(timeout, testCase.d) {
			t.Fatalf(
				"%s: actual and expected timeouts do not match\nactual: %v\nexpected:%v",
				testName,
				timeout,
				testCase.d,
			)
		}
	}
}

func TestWithTimeout(t *testing.T) {
	cases := map[string]*optionsDurationTestCase{
		"set-timeout": {
			description: "simple set option test",
			d:           99 * time.Second,
			o:           &channel.OperationOptions{},
			isignored:   false,
		},
		"set-timeout-netconf": {
			description: "simple set option test",
			d:           99 * time.Second,
			o:           &netconf.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			d:           99 * time.Second,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithTimeout(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCompletePatterns(testName string, testCase *struct {
	description string
	p           []*regexp.Regexp
	o           interface{}
	isignored   bool
},
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCompletePatterns(testCase.p)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*channel.OperationOptions)

		if !cmp.Equal(oo.CompletePatterns[0].String(), testCase.p[0].String()) {
			t.Fatalf(
				"%s: actual and expected complete patterns do not match\nactual: %v\nexpected:%v",
				testName,
				oo.CompletePatterns[0].String(),
				testCase.p[0].String(),
			)
		}
	}
}

func TestWithCompletePatterns(t *testing.T) {
	cases := map[string]*struct {
		description string
		p           []*regexp.Regexp
		o           interface{}
		isignored   bool
	}{
		"set-complete-patterns": {
			description: "simple set option test",
			p:           []*regexp.Regexp{regexp.MustCompile(`pattern!`)},
			o:           &channel.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			p:           nil,
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithCompletePatterns(testName, testCase)

		t.Run(testName, f)
	}
}
