package network_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

type sendConfigsTestCase struct {
	description string
	configs     []string
	payloadFile string
	stripPrompt bool
	eager       bool
	privLevel   string
}

func testSendConfigs(testName string, testCase *sendConfigsTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendConfigs(testCase.configs)
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

func TestSendConfigs(t *testing.T) {
	cases := map[string]*sendConfigsTestCase{
		"send-configs-simple": {
			description: "simple send config test",
			configs:     []string{"interface loopback1", "no interface loopback1"},
			payloadFile: "send-configs-simple.txt",
			stripPrompt: false,
			eager:       false,
		},
	}

	for testName, testCase := range cases {
		f := testSendConfigs(testName, testCase)

		t.Run(testName, f)
	}
}

type sendConfigsFromFileTestCase struct {
	description string
	f           string
	payloadFile string
	stripPrompt bool
	eager       bool
}

func testSendConfigsFromFile(
	testName string,
	testCase *sendConfigsFromFileTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		r, err := d.SendConfigsFromFile(resolveFile(t, testCase.f))
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

func TestSendConfigsFromFile(t *testing.T) {
	cases := map[string]*sendConfigsFromFileTestCase{
		"send-configs-from-file-simple": {
			description: "simple send commands test",
			f:           "send-configs-from-file-simple-inputs.txt",
			payloadFile: "send-configs-from-file-simple.txt",
			stripPrompt: false,
			eager:       false,
		},
	}

	for testName, testCase := range cases {
		f := testSendConfigsFromFile(testName, testCase)

		t.Run(testName, f)
	}
}

type sendConfigsFunctionalTestcase struct {
	description string
	stripPrompt bool
	eager       bool
}

func getTestSendConfigsFunctionalConfig(t *testing.T, testName, platformName string) []string {
	commands := map[string][]string{
		platform.CiscoIosxe: {"interface loopback0", "no interface loopback0"},
		platform.CiscoIosxr: {"interface loopback0", "no interface loopback0", "commit"},
		platform.CiscoNxos:  {"interface loopback0", "no interface loopback0"},
		platform.AristaEos:  {"interface loopback0", "no interface loopback0"},
		platform.JuniperJunos: {
			"set interfaces fxp0 unit 0 description \"scrapli was here\"",
			"delete interfaces fxp0 unit 0 description",
			"commit",
		},
		platform.NokiaSrl: {
			"interface ethernet-1/50",
			"subinterface 0",
			"ipv4 address 1.1.1.1/30",
			"discard now",
		},
	}

	c, ok := commands[platformName]
	if !ok {
		t.Skipf("%s: skipping platform '%s', no command in commands map", testName, platformName)
	}

	return c
}

func testSendConfigsFunctional(
	testName, platformName, transportName string,
	testCase *sendConfigsFunctionalTestcase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		r, err := d.SendConfigs(
			getTestSendConfigsFunctionalConfig(t, testName, platformName),
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

		actualOut := r.JoinedResult()

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

func TestSendConfigsFunctional(t *testing.T) {
	cases := map[string]*sendConfigsFunctionalTestcase{
		"functional-send-configs-simple": {
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

				f := testSendConfigsFunctional(testName, platformName, transportName, testCase)

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
