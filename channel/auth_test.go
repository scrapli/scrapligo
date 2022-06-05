package channel_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testAuthenticateSSH(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		c, fileTransportObj := prepareChannel(t, testName, testCase.PayloadFile)

		actualOut, err := c.AuthenticateSSH([]byte("some-password"), nil)
		if err != nil {
			t.Errorf(
				"%s: encountered error running Channel AuthenticateSSH, error: %s",
				testName,
				err,
			)
		}

		actualIn := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualIn, actualOut)
		}

		expectedIn := readFile(t, fmt.Sprintf("golden/%s-in.txt", testName))
		expectedOut := string(readFile(t, fmt.Sprintf("golden/%s-out.txt", testName)))

		if !cmp.Equal(actualIn, expectedIn) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualIn,
				expectedIn,
			)
		}

		if !cmp.Equal(string(actualOut), expectedOut) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestAuthenticateSSH(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"auth-simple": {
			Description: "simple in channel auth test",
			PayloadFile: "auth-simple.txt",
		},
		"auth-two-attempts": {
			Description: "simple in channel auth where first attempt fails for some reason test",
			PayloadFile: "auth-two-attempts.txt",
		},
	}

	for testName, testCase := range cases {
		f := testAuthenticateSSH(testName, testCase)

		t.Run(testName, f)
	}
}

func testAuthenticateTelnet(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		c, fileTransportObj := prepareChannel(t, testName, testCase.PayloadFile)

		actualOut, err := c.AuthenticateTelnet([]byte("some-username"), []byte("some-password"))
		if err != nil {
			t.Errorf(
				"%s: encountered error running Channel AuthenticateTelnet, error: %s",
				testName,
				err,
			)
		}

		actualIn := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualIn, actualOut)
		}

		expectedIn := readFile(t, fmt.Sprintf("golden/%s-in.txt", testName))
		expectedOut := string(readFile(t, fmt.Sprintf("golden/%s-out.txt", testName)))

		if !cmp.Equal(actualIn, expectedIn) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualIn,
				expectedIn,
			)
		}

		if !cmp.Equal(string(actualOut), expectedOut) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestAuthenticateTelnet(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"auth-telnet-simple": {
			Description: "simple in channel auth test",
			PayloadFile: "auth-telnet-simple.txt",
		},
		"auth-telnet-two-attempts": {
			Description: "simple in channel auth where first attempt fails for some reason test",
			PayloadFile: "auth-telnet-two-attempts.txt",
		},
	}

	for testName, testCase := range cases {
		f := testAuthenticateTelnet(testName, testCase)

		t.Run(testName, f)
	}
}
