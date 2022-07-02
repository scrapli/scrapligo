package util_test

import (
	"testing"

	"github.com/scrapli/scrapligo/util"
)

func testStringContainsAny(
	testName, s string,
	contains []string,
	expected bool,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		if util.StringContainsAny(s, contains) != expected {
			t.Fatalf(
				"%s: actual and expected results differ",
				testName,
			)
		}
	}
}

func TestStringContainsAny(t *testing.T) {
	cases := map[string]struct {
		description string
		s           string
		contains    []string
		expected    bool
	}{
		"string-contains-any-simple-true": {
			s:        "tacocat",
			contains: []string{"taco", "cat", "race", "car"},
			expected: true,
		},
		"sstring-contains-any-simple-false": {
			s:        "potato",
			contains: []string{"taco", "cat", "race", "car"},
			expected: false,
		},
	}

	for testName, testCase := range cases {
		f := testStringContainsAny(testName, testCase.s, testCase.contains, testCase.expected)

		t.Run(testName, f)
	}
}

func testStringContainsAnySubStrs(
	testName, s string,
	contains []string,
	expected string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		if util.StringContainsAnySubStrs(s, contains) != expected {
			t.Fatalf(
				"%s: actual and expected results differ",
				testName,
			)
		}
	}
}

func TestStringContainsAnySubStrs(t *testing.T) {
	cases := map[string]struct {
		description string
		s           string
		contains    []string
		expected    string
	}{
		"string-slice-contains-any-substr-match": {
			s:        "tacocat",
			contains: []string{"taco", "cat", "race", "car"},
			expected: "taco",
		},
		"string-slice-contains-any-substr-no-match": {
			s:        "potato",
			contains: []string{"taco", "cat", "race", "car"},
			expected: "",
		},
	}

	for testName, testCase := range cases {
		f := testStringContainsAnySubStrs(
			testName,
			testCase.s,
			testCase.contains,
			testCase.expected,
		)

		t.Run(testName, f)
	}
}
