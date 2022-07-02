package util_test

import (
	"testing"

	"github.com/scrapli/scrapligo/util"
)

func testStringSliceContains(testName, s string, ss []string, expected bool) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		if util.StringSliceContains(ss, s) != expected {
			t.Fatalf(
				"%s: actual and expected results differ",
				testName,
			)
		}
	}
}

func TestStringSliceContains(t *testing.T) {
	cases := map[string]struct {
		description string
		ss          []string
		s           string
		expected    bool
	}{
		"string-slice-contains-simple-true": {
			ss:       []string{"taco", "cat", "race", "car"},
			s:        "taco",
			expected: true,
		},
		"string-slice-contains-simple-false": {
			ss:       []string{"taco", "cat", "race", "car"},
			s:        "potato",
			expected: false,
		},
	}

	for testName, testCase := range cases {
		f := testStringSliceContains(testName, testCase.s, testCase.ss, testCase.expected)

		t.Run(testName, f)
	}
}
