package netconf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testGetConfig(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.PayloadFile)

		r, err := d.GetConfig("running")
		if err != nil {
			t.Fatalf(
				"%s: encountered error running netconf Driver getConfig, error: %s",
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

func TestGetConfig(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"getconfig-simple": {
			Description: "simple getconfig test",
			PayloadFile: "getconfig-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testGetConfig(testName, testCase)
		t.Run(testName, f)
	}
}

func testGetConfigFunctional(
	testName, platformName, transportName string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		defer func() {
			err := d.Close()
			if err != nil {
				t.Fatalf("%s: failed closing connection",
					testName)
			}
		}()

		r, err := d.GetConfig("running")
		if err != nil {
			t.Logf(
				"%s: encountered error running netconf Driver getConfig, error: %s",
				testName,
				err,
			)

			t.Fail()

			return
		}

		if r.Failed != nil {
			t.Logf("%s: response object indicates failure",
				testName)

			t.Fail()

			return
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
			t.Logf(
				"%s: actual and expected outputs do *not* match, however we did not get "+
					"a failed operation and devices sometimes have empty elements and out of "+
					"order elements causing comparison to be difficult, this is "+
					"*probably* not an issue, but you should take a look at the output to confirm!",
				testName,
			)
		}
	}
}

func TestGetConfigFunctional(t *testing.T) {
	cases := map[string]*struct {
		description string
	}{
		"functional-getconfig-simple": {
			description: "simple get config test",
		},
	}

	if !*functional {
		t.Skip("skip: functional tests skipped without the '-functional' flag being passed")
	}

	for testName := range cases {
		for _, platformName := range getNetconfPlatformNames() {
			if !util.PlatformOK(platforms, platformName) {
				t.Logf("%s: skipping platform '%s'", testName, platformName)

				continue
			}

			for _, transportName := range getNetconfTransportNames() {
				if !util.TransportOK(transports, transportName) {
					t.Logf("%s: skipping transport '%s'", testName, transportName)

					continue
				}

				f := testGetConfigFunctional(testName, platformName, transportName)

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
