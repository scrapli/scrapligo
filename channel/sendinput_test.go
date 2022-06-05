package channel_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

type sendInputTestcase struct {
	description   string
	input         string
	payloadFile   string
	noStripPrompt bool
	eager         bool
}

func testSendInput(testName string, testCase *sendInputTestcase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		c, fileTransportObj := prepareChannel(t, testName, testCase.payloadFile)

		var opts []util.Option

		if testCase.eager {
			opts = append(opts, opoptions.WithEager())
		}

		if testCase.noStripPrompt {
			opts = append(opts, opoptions.WithNoStripPrompt())
		}

		actualOut, err := c.SendInput(
			testCase.input,
			opts...,
		)
		if err != nil {
			t.Errorf("%s: encountered error running Channel SendInput, error: %s", testName, err)
		}

		actualIn := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualIn, actualOut)
		}

		expectedIn := readFile(t, fmt.Sprintf("golden/%s-in.txt", testName))
		expectedOut := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !cmp.Equal(actualIn, expectedIn) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualIn,
				expectedIn,
			)
		}

		if !cmp.Equal(actualOut, expectedOut) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestSendInput(t *testing.T) {
	cases := map[string]*sendInputTestcase{
		"send-input-simple": {
			description:   "simple send input test",
			input:         "show run int vlan1",
			payloadFile:   "send-input-simple.txt",
			noStripPrompt: true,
			eager:         false,
		},
		"send-input-simple-strip-prompt": {
			description:   "simple send input test",
			input:         "show run int vlan1",
			payloadFile:   "send-input-simple.txt",
			noStripPrompt: false,
			eager:         false,
		},
	}

	for testName, testCase := range cases {
		f := testSendInput(testName, testCase)

		t.Run(testName, f)
	}
}
