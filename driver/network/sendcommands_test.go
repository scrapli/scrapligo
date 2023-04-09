package network_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/transport"

	"github.com/google/go-cmp/cmp"
)

type sendCommandsTestCase struct {
	description    string
	commands       []string
	payloadFile    string
	stripPrompt    bool
	eager          bool
	interimPrompts []*regexp.Regexp
}

func testSendCommands(testName string, testCase *sendCommandsTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payloadFile)

		var opts []util.Option

		if len(testCase.interimPrompts) > 0 {
			opts = append(opts, opoptions.WithInterimPromptPattern(testCase.interimPrompts))
		}

		r, err := d.SendCommands(testCase.commands, opts...)
		if err != nil {
			t.Fatalf(
				"%s: encountered error running generic Driver SendCommands, error: %s",
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
		"send-commands-interim-prompt": {
			description: "simple send commands test with interim prompt patterns",
			commands: []string{
				"some command that starts interim prompt thing",
				"subcommand1",
				"subcommand2",
			},
			payloadFile:    "send-commands-interim-prompt.txt",
			stripPrompt:    false,
			eager:          false,
			interimPrompts: []*regexp.Regexp{regexp.MustCompile(`(?m)^\.{3}`)},
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

type sendCommandsFunctionalTestcase struct {
	description string
	stripPrompt bool
	eager       bool
}

func getTestSendCommandsFunctionalCommand(t *testing.T, testName, platformName string) []string {
	commands := map[string][]string{
		platform.CiscoIosxe:   {"show run", "show ip interface brief"},
		platform.CiscoIosxr:   {"show run", "show interfaces brief"},
		platform.CiscoNxos:    {"show run", "show interface status"},
		platform.AristaEos:    {"show run", "show ip interface brief"},
		platform.JuniperJunos: {"show configuration", "show version"},
		platform.NokiaSrl:     {"show interface all", "show platform environment"},
	}

	c, ok := commands[platformName]
	if !ok {
		t.Skipf("%s: skipping platform '%s', no command in commands map", testName, platformName)
	}

	return c
}

func testSendCommandsFunctional(
	testName, platformName, transportName string,
	_ *sendCommandsFunctionalTestcase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		r, err := d.SendCommands(
			getTestSendCommandsFunctionalCommand(t, testName, platformName),
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

func TestSendCommandsFunctional(t *testing.T) {
	cases := map[string]*sendCommandsFunctionalTestcase{
		"functional-send-commands-simple": {
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

				f := testSendCommandsFunctional(testName, platformName, transportName, testCase)

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
