package generic_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type sendCommandsTestCase struct {
	description string
	commands    []string
	payloadFile string
	stripPrompt bool
	eager       bool
}

func testSendCommands(testName string, testCase *sendCommandsTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendCommands(testCase.commands)
		if err != nil {
			t.Errorf(
				"%s: encountered error running generic Driver GetPrompt, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := r.JoinedResult()
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

		if !cmp.Equal(actualOut, string(expectedOut)) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestSendCommands(t *testing.T) {
	cases := map[string]*sendCommandsTestCase{
		"send-commands-simple": {
			description: "simple send commands test",
			commands:    []string{"show run int vlan1", "show run int vlan1"},
			payloadFile: "send-commands-simple.txt",
			stripPrompt: false,
			eager:       false,
		},
	}

	for testName, testCase := range cases {
		f := testSendCommands(testName, testCase)

		t.Run(testName, f)
	}
}

type sendCommandsFromFileTestCase struct {
	description string
	f           string
	payloadFile string
	stripPrompt bool
	eager       bool
}

func testSendCommandsFromFile(
	testName string,
	testCase *sendCommandsFromFileTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendCommandsFromFile(resolveFile(t, testCase.f))
		if err != nil {
			t.Errorf(
				"%s: encountered error running generic Driver GetPrompt, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := r.JoinedResult()
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

		if !cmp.Equal(actualOut, string(expectedOut)) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestSendCommandsFromFile(t *testing.T) {
	cases := map[string]*sendCommandsFromFileTestCase{
		"send-commands-from-file-simple": {
			description: "simple send commands test",
			f:           "send-commands-from-file-simple-inputs.txt",
			payloadFile: "send-commands-from-file-simple.txt",
			stripPrompt: false,
			eager:       false,
		},
	}

	for testName, testCase := range cases {
		f := testSendCommandsFromFile(testName, testCase)

		t.Run(testName, f)
	}
}
