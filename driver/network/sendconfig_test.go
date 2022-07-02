package network_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

type sendConfigTestCase struct {
	description string
	configs     string
	payloadFile string
	stripPrompt bool
	eager       bool
	privLevel   string
}

func testSendConfig(testName string, testCase *sendConfigTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendConfig(testCase.configs)
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

func TestSendConfig(t *testing.T) {
	cases := map[string]*sendConfigTestCase{
		"send-config-simple": {
			description: "simple send config test",
			configs:     "interface loopback1\nno interface loopback1",
			payloadFile: "send-config-simple.txt",
			stripPrompt: false,
			eager:       false,
		},
	}

	for testName, testCase := range cases {
		f := testSendConfig(testName, testCase)

		t.Run(testName, f)
	}
}

type sendConfigFunctionalTestCase struct {
	description string
}

func getTestSendConfigFunctionalConfig(t *testing.T, testName, platformName string) string {
	cSlice := getTestSendConfigsFunctionalConfig(t, testName, platformName)

	return strings.Join(cSlice, "\n")
}

func testSendConfigFunctional(
	testName, platformName, transportName string,
	testCase *sendConfigFunctionalTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		r, err := d.SendConfig(
			getTestSendConfigFunctionalConfig(t, testName, platformName),
		)
		if err != nil {
			t.Errorf(
				"%s: encountered error running network Driver SendConfig, error: %s",
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
			util.GetCleanFunc(platformName)(actualOut),
			cleanF(string(expectedOut)),
		) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestSendConfigFunctional(t *testing.T) {
	cases := map[string]*sendConfigFunctionalTestCase{
		"functional-send-config-simple": {
			description: "simple send config test",
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

				f := testSendConfigFunctional(testName, platformName, transportName, testCase)

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
