package network_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/transport"

	"github.com/google/go-cmp/cmp"
)

type sendCommandTestCase struct {
	description string
	command     string
	payloadFile string
	stripPrompt bool
	eager       bool
}

func testSendCommand(testName string, testCase *sendCommandTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendCommand(testCase.command)
		if err != nil {
			t.Errorf(
				"%s: encountered error running network Driver SendCommand, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := r.Result
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

func TestSendCommand(t *testing.T) {
	cases := map[string]*sendCommandTestCase{
		"send-command-simple": {
			description: "simple send command test",
			command:     "show run int vlan1",
			payloadFile: "send-command-simple.txt",
			stripPrompt: false,
			eager:       false,
		},
		"send-command-acquire-priv": {
			description: "simple send command test plus acquire priv",
			command:     "show run int vlan1",
			payloadFile: "send-command-acquire-priv.txt",
			stripPrompt: false,
			eager:       false,
		},
	}

	for testName, testCase := range cases {
		f := testSendCommand(testName, testCase)
		t.Run(testName, f)
	}
}

type sendCommandFunctionalTestcase struct {
	description string
	stripPrompt bool
	eager       bool
}

func getTestSendCommandFunctionalCommand(t *testing.T, testName, platformName string) string {
	commands := map[string]string{
		platform.CiscoIosxe:   "show run",
		platform.CiscoIosxr:   "show run",
		platform.CiscoNxos:    "show run",
		platform.AristaEos:    "show run",
		platform.JuniperJunos: "show configuration",
		platform.NokiaSrl:     "show interface all",
	}

	c, ok := commands[platformName]
	if !ok {
		t.Skipf("%s: skipping platform '%s', no command in commands map", testName, platformName)
	}

	return c
}

func testSendCommandFunctional(
	testName, platformName, transportName string,
	_ *sendCommandFunctionalTestcase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		r, err := d.SendCommand(
			getTestSendCommandFunctionalCommand(t, testName, platformName),
		)
		if err != nil {
			t.Errorf(
				"%s: encountered error running network Driver SendCommand, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		err = d.Close()
		if err != nil {
			t.Fatalf("%s: failed closing connection",
				testName)
		}

		actualOut := r.Result

		if *update {
			writeGoldenFunctional(
				t,
				fmt.Sprintf("%s-%s-%s", testName, platformName, transportName),
				actualOut,
			)
		}

		cleanF := util.GetCleanFunc(platformName)

		expectedOut := readFile(
			t,
			fmt.Sprintf("golden/%s-%s-%s-out.txt", testName, platformName, transportName),
		)

		if !cmp.Equal(
			cleanF(actualOut),
			cleanF(string(expectedOut)),
		) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				cleanF(actualOut),
				cleanF(string(expectedOut)),
			)
		}
	}
}

func TestSendCommandFunctional(t *testing.T) {
	cases := map[string]*sendCommandFunctionalTestcase{
		"functional-send-command-simple": {
			description: "simple send command test",
			stripPrompt: false,
			eager:       false,
		},
	}

	if !*functional {
		t.Skip("skip: functional tests skipped without the '-functional' flag being passed")
	}

	for testName, testCase := range cases {
		for _, platformName := range platform.GetPlatformNames() {
			if !util.PlatformOK(platforms, platformName) {
				t.Logf("%s: skipping platform '%s'", testName, platformName)

				continue
			}

			for _, transportName := range transport.GetTransportNames() {
				if !util.TransportOK(transports, transportName) {
					t.Logf("%s: skipping transport '%s'", testName, transportName)

					continue
				}

				f := testSendCommandFunctional(testName, platformName, transportName, testCase)

				t.Run(
					fmt.Sprintf(
						"%s;platform=%s;transport=%s",
						testName,
						platformName,
						transportName,
					),
					f,
				)

				interTestSleep()
			}
		}
	}
}

func TestSendCommandFunctionalGiant(t *testing.T) {
	testName := "send-command-giant"
	platformName := platform.NokiaSrl

	if !*functional {
		t.Skip("skip: functional tests skipped without the '-functional' flag being passed")
	}

	t.Logf("%s: starting", testName)

	for _, transportName := range []string{transport.SystemTransport, transport.StandardTransport} {
		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		r, err := d.SendCommand("info from running")
		if err != nil {
			t.Errorf(
				"%s: encountered error running network Driver SendCommand (giant), error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		err = d.Close()
		if err != nil {
			t.Fatalf("%s: failed closing connection",
				testName)
		}

		actualOut := r.Result

		if *update {
			writeGoldenFunctional(
				t,
				fmt.Sprintf("%s-%s-%s", testName, platformName, transportName),
				actualOut,
			)
		}

		cleanF := util.GetCleanFunc(platformName)

		expectedOut := readFile(
			t,
			fmt.Sprintf("golden/%s-%s-%s-out.txt", testName, platformName, transportName),
		)

		if !cmp.Equal(
			cleanF(actualOut),
			cleanF(string(expectedOut)),
		) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				cleanF(actualOut),
				cleanF(string(expectedOut)),
			)
		}

		// in seconds; realistically this looks like it finishes in ~10-12s pretty consistently in
		// local tests, but we need to make sure it doesn't go kaboom in ci as well. also note that
		// this is usually ran w/ race flag so that slows things down even more.
		if r.ElapsedTime > 30 {
			t.Fatalf("%s: test completed but was greater than maximum expected duration, took %fs",
				testName, r.ElapsedTime)
		}
	}
}
