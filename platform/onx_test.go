package platform_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testOnXTestCase struct {
	description  string
	payloadFile  string
	platformFile string
}

func testOnXGeneric(testName string, testCase *testOnXTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(
			t,
			testName,
			testCase.platformFile,
			testCase.payloadFile,
		)

		err := d.Driver.OnOpen(d.Driver)
		if err != nil {
			t.Fatalf("%s: response object indicates failure", testName)
		}

		actualOut := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualOut)
		}

		expectedOut := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

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

func TestOnXGeneric(t *testing.T) {
	cases := map[string]*testOnXTestCase{
		"generic-on-x-simple": {
			description:  "simple generic on-x operation test",
			platformFile: "test-platform.yaml",
			payloadFile:  "generic-on-x-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testOnXGeneric(testName, testCase)

		t.Run(testName, f)
	}
}

func testOnXNetwork(testName string, testCase *testOnXTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(
			t,
			testName,
			testCase.platformFile,
			testCase.payloadFile,
		)

		err := d.OnOpen(d)
		if err != nil {
			t.Fatalf("%s: response object indicates failure", testName)
		}

		actualOut := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualOut)
		}

		expectedOut := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

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

func TestOnXNetwork(t *testing.T) {
	cases := map[string]*testOnXTestCase{
		"network-on-x-simple": {
			description:  "simple network on-x operation test",
			platformFile: "test-platform.yaml",
			payloadFile:  "network-on-x-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testOnXNetwork(testName, testCase)

		t.Run(testName, f)
	}
}
