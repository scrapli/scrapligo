package opoptions_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testWithCallbackContains(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackContains(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.Contains, testCase.s) {
			t.Fatalf(
				"%s: actual and expected contains do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Contains,
				testCase.s,
			)
		}
	}
}

func TestWithCallbackContains(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-callback-contains": {
			description: "simple set option test",
			s:           "potato",
			o:           &generic.Callback{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "potato",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithCallbackContains(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackNotContains(
	testName string,
	testCase *optionsStringTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackNotContains(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.NotContains, testCase.s) {
			t.Fatalf(
				"%s: actual and expected not contains do not match\nactual: %v\nexpected:%v",
				testName,
				oo.NotContains,
				testCase.s,
			)
		}
	}
}

func TestWithCallbackNotContains(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-callback-not-contains": {
			description: "simple set option test",
			s:           "potato",
			o:           &generic.Callback{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "potato",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithCallbackNotContains(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackContainsRe(
	testName string,
	testCase *optionsRegexpTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackContainsRe(testCase.p)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.ContainsRe.String(), testCase.p.String()) {
			t.Fatalf(
				"%s: actual and expected contains regex do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Contains,
				testCase.p,
			)
		}
	}
}

func TestWithCallbackContainsRe(t *testing.T) {
	cases := map[string]*optionsRegexpTestCase{
		"set-callback-contains-re": {
			description: "simple set option test",
			p:           regexp.MustCompile("potato"),
			o:           &generic.Callback{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			p:           regexp.MustCompile("potato"),
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithCallbackContainsRe(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackInsensitive(
	testName string,
	testCase *optionsBoolTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackInsensitive(testCase.b)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.Insensitive, testCase.b) {
			t.Fatalf(
				"%s: actual and expected insensitive do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Insensitive,
				testCase.b,
			)
		}
	}
}

func TestWithCallbackInsensitive(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-callback-insensitive": {
			description: "simple set option test",
			b:           true,
			o:           &generic.Callback{},
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
		f := testWithCallbackInsensitive(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackResetOutput(
	testName string,
	testCase *optionsBoolTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackResetOutput()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.ResetOutput, testCase.b) {
			t.Fatalf(
				"%s: actual and expected reset output do not match\nactual: %v\nexpected:%v",
				testName,
				oo.ResetOutput,
				testCase.b,
			)
		}
	}
}

func TestWithCallbackResetOutput(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-callback-reset-output": {
			description: "simple set option test",
			b:           true,
			o:           &generic.Callback{},
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
		f := testWithCallbackResetOutput(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackOnce(
	testName string,
	testCase *optionsBoolTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackOnce()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.Once, testCase.b) {
			t.Fatalf(
				"%s: actual and expected once do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Once,
				testCase.b,
			)
		}
	}
}

func TestWithCallbackOnce(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-callback-once": {
			description: "simple set option test",
			b:           true,
			o:           &generic.Callback{},
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
		f := testWithCallbackOnce(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackComplete(
	testName string,
	testCase *optionsBoolTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackComplete()(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.Complete, testCase.b) {
			t.Fatalf(
				"%s: actual and expected complete do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Complete,
				testCase.b,
			)
		}
	}
}

func TestWithCallbackComplete(t *testing.T) {
	cases := map[string]*optionsBoolTestCase{
		"set-callback-complete": {
			description: "simple set option test",
			b:           true,
			o:           &generic.Callback{},
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
		f := testWithCallbackComplete(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackName(testName string, testCase *optionsStringTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackName(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.Name, testCase.s) {
			t.Fatalf(
				"%s: actual and expected name do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Name,
				testCase.s,
			)
		}
	}
}

func TestWithCallbackName(t *testing.T) {
	cases := map[string]*optionsStringTestCase{
		"set-callback-name": {
			description: "simple set option test",
			s:           "potato",
			o:           &generic.Callback{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "potato",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithCallbackName(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithCallbackNextTimeout(
	testName string,
	testCase *optionsDurationTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithCallbackNextTimeout(testCase.d)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*generic.Callback)

		if !cmp.Equal(oo.NextTimeout, testCase.d) {
			t.Fatalf(
				"%s: actual and expected next timeout do not match\nactual: %v\nexpected:%v",
				testName,
				oo.NextTimeout,
				testCase.d,
			)
		}
	}
}

func TestWithCallbackNextTimeout(t *testing.T) {
	cases := map[string]*optionsDurationTestCase{
		"set-callback-duration": {
			description: "simple set option test",
			d:           99 * time.Second,
			o:           &generic.Callback{},
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
		f := testWithCallbackNextTimeout(testName, testCase)

		t.Run(testName, f)
	}
}
