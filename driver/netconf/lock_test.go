package netconf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/platform"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testLockUnlock(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.PayloadFile)

		lockr, err := d.Lock("running")
		if err != nil {
			t.Fatalf(
				"%s: encountered error running network Driver getConfig, error: %s",
				testName,
				err,
			)
		}

		if lockr.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		unlock, err := d.Unlock("running")
		if err != nil {
			t.Fatalf(
				"%s: encountered error running network Driver getConfig, error: %s",
				testName,
				err,
			)
		}

		if unlock.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := lockr.Result + "\n" + unlock.Result
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

func TestLockUnlock(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"lock-unlock-simple": {
			Description: "simple lock/unlock test",
			PayloadFile: "lock-unlock-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testLockUnlock(testName, testCase)
		t.Run(testName, f)
	}
}

func testLockUnlockFunctional( //nolint: funlen
	testName, platformName, transportName string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d := prepareFunctionalDriver(t, testName, platformName, transportName)

		var target string

		switch platformName {
		case platform.CiscoIosxe, platform.CiscoNxos, platform.AristaEos:
			target = "running"
		default:
			target = "candidate"
		}

		lockr, err := d.Lock(target)
		if err != nil {
			t.Fatalf(
				"%s: encountered error running netconf Driver lock, error: %s",
				testName,
				err,
			)
		}

		if lockr.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		unlockr, err := d.Unlock(target)
		if err != nil {
			t.Fatalf(
				"%s: encountered error running netconf Driver unlock, error: %s",
				testName,
				err,
			)
		}

		if unlockr.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		err = d.Close()
		if err != nil {
			t.Fatalf("%s: failed closing connection",
				testName)
		}

		actualLockOut := lockr.Result
		actualUnlockOut := unlockr.Result

		if *update {
			writeGoldenFunctional(
				t,
				fmt.Sprintf("%s-lock-%s-%s", testName, platformName, transportName),
				actualLockOut,
			)
			writeGoldenFunctional(
				t,
				fmt.Sprintf("%s-unlock-%s-%s", testName, platformName, transportName),
				actualUnlockOut,
			)
		}

		cleanF := util.GetCleanFunc(platformName)

		expectedLockOut := readFile(
			t,
			fmt.Sprintf("golden/%s-lock-%s-%s-out.txt", testName, platformName, transportName),
		)

		expectedUnlockOut := readFile(
			t,
			fmt.Sprintf("golden/%s-unlock-%s-%s-out.txt", testName, platformName, transportName),
		)

		if !cmp.Equal(
			util.GetCleanFunc(platformName)(actualLockOut),
			cleanF(string(expectedLockOut)),
		) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualLockOut,
				expectedLockOut,
			)
		}

		if !cmp.Equal(
			util.GetCleanFunc(platformName)(actualUnlockOut),
			cleanF(string(expectedUnlockOut)),
		) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualUnlockOut,
				expectedUnlockOut,
			)
		}
	}
}

func TestLockUnlockFunctional(t *testing.T) {
	cases := map[string]*struct {
		description string
	}{
		"functional-lock-unlock-simple": {
			description: "simple lock and unlock test",
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

				f := testLockUnlockFunctional(testName, platformName, transportName)

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
