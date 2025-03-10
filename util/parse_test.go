package util_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	scrapligotesthelper "github.com/scrapli/scrapligo/testhelper"
	scrapligoutil "github.com/scrapli/scrapligo/util"
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

			testGoldenPath, err := filepath.Abs(fmt.Sprintf("./golden/%s", testName))
			if err != nil {
				t.Fatal(err)
			}

			s := string(scrapligotesthelper.ReadFile(t, fmt.Sprintf("%s.txt", testName)))

			actualOut, err := scrapligoutil.TextFsmParse(
				s,
				"test-fixtures/cisco_ios_show_version.textfsm",
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
				scrapligotesthelper.WriteFile(t, testGoldenPath, actualOutJSON)
			} else {
				testGoldenContent := scrapligotesthelper.ReadFile(
					t,
					fmt.Sprintf("golden/%s-out.txt", testName),
				)

				if !bytes.Equal(actualOutJSON, testGoldenContent) {
					t.Fatalf(
						"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
						testName,
						actualOutJSON,
						testGoldenContent,
					)
				}
			}
		})
	}
}
