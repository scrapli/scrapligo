package generic_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/scrapli/scrapligo/driver/generic"

	"github.com/google/go-cmp/cmp"
)

type sendWithCallbacksTestCase struct {
	description  string
	payloadFile  string
	initialInput string
	callbacks    []*generic.Callback
	failed       bool
}

func testSendWithCallbacks(
	testName string,
	testCase *sendWithCallbacksTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendWithCallbacks(testCase.initialInput, testCase.callbacks, 1*time.Second)
		if err != nil {
			t.Fatalf(
				"%s: encountered error running generic Driver SendWithCallbacks, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil && !testCase.failed {
			t.Fatalf(
				"%s: response object indicates failure, this shouldn't happen",
				testName,
			)
		}

		actualIn := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualIn, "")
		}

		expectedIn := readFile(t, fmt.Sprintf("golden/%s-in.txt", testName))

		if !cmp.Equal(actualIn, expectedIn) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualIn,
				expectedIn,
			)
		}
	}
}

func TestSendWithCallbacks(t *testing.T) {
	cases := map[string]*sendWithCallbacksTestCase{
		"send-with-callbacks-simple": {
			description:  "simple send with callbacks test",
			payloadFile:  "send-with-callbacks-simple.txt",
			initialInput: "",
			callbacks: []*generic.Callback{
				{
					Callback: func(d *generic.Driver, _ string) error {
						_, err := d.Channel.SendInput("configure terminal")

						return err
					},
					Contains:    "C3560CX#",
					Name:        "callback-one",
					ResetOutput: true,
				},
				{
					Callback: func(d *generic.Driver, _ string) error {
						return d.Channel.WriteAndReturn([]byte("show version"), false)
					},
					ContainsRe: regexp.MustCompile(`(?im)^c3560cx\(config\)#`),
					Name:       "callback-two",
					Complete:   true,
				},
			},
		},
		"send-with-callbacks-failed-when-contains": {
			description:  "simple send with callbacks test w/ failed when contains output",
			payloadFile:  "send-with-callbacks-failed-when-contains.txt",
			initialInput: "",
			callbacks: []*generic.Callback{
				{
					Callback: func(d *generic.Driver, _ string) error {
						_, err := d.Channel.SendInput("configure terminal")

						return err
					},
					Contains:    "C3560CX#",
					Name:        "callback-one",
					ResetOutput: true,
				},
				{
					Callback: func(d *generic.Driver, _ string) error {
						return d.Channel.WriteAndReturn([]byte("show version"), false)
					},
					ContainsRe: regexp.MustCompile(`(?im)^c3560cx\(config\)#`),
					Name:       "callback-two",
					Complete:   true,
				},
			},
			failed: true,
		},
	}

	for testName, testCase := range cases {
		f := testSendWithCallbacks(testName, testCase)

		t.Run(testName, f)
	}
}
