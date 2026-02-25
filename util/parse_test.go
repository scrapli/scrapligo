package util_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	scrapligotesthelper "github.com/scrapli/scrapligo/v2/testhelper"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

func TestTextFsmParse(t *testing.T) {
	cases := map[string]struct {
		description string
	}{
		"textfsm-parse-simple": {},
	}

	for testName := range cases {
		t.Run(testName, func(t *testing.T) {
			t.Logf("%s: starting", testName)

			testFixturePath, err := filepath.Abs(fmt.Sprintf("./fixtures/%s", testName))
			if err != nil {
				t.Fatal(err)
			}

			testTemplateFixturePath, err := filepath.Abs(
				"./fixtures/cisco_ios_show_version.textfsm",
			)
			if err != nil {
				t.Fatal(err)
			}

			testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
			if err != nil {
				t.Fatal(err)
			}

			actualOut, err := scrapligoutil.TextFsmParse(
				t.Context(),
				string(scrapligotesthelper.ReadFile(t, testFixturePath)),
				testTemplateFixturePath,
			)
			if err != nil {
				t.Fatalf("%s: encountered error parsing with textfsm, error: %s", testName, err)
			}

			actualOutJSON, err := json.Marshal(actualOut)
			if err != nil {
				t.Fatalf(
					"%s: encountered error dumping textfsm output to json, error: %s",
					testName,
					err,
				)
			}

			if *scrapligotesthelper.Update {
				scrapligotesthelper.WriteFile(
					t,
					testGoldenPath,
					actualOutJSON,
				)
			} else {
				testGoldenContent := scrapligotesthelper.ReadFile(t, testGoldenPath)

				if !bytes.Equal(actualOutJSON, testGoldenContent) {
					scrapligotesthelper.FailOutput(t, actualOutJSON, testGoldenContent)
				}
			}
		})
	}
}
