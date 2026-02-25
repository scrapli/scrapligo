package util_test

import (
	"testing"

	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

func TestGetEnvStrOrDefault(t *testing.T) {
	cases := []struct {
		name     string
		k        string
		setV     string
		defaultV string
		expected string
	}{
		{
			name:     "simple-default",
			k:        "SOME_ENV_VAR",
			setV:     "",
			defaultV: "foo",
			expected: "foo",
		},
		{
			name:     "simple-already-set",
			k:        "SOME_ENV_VAR",
			setV:     "foo",
			defaultV: "taco",
			expected: "foo",
		},
	}

	for _, testCase := range cases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Logf("%s: starting", testCase.name)

				t.Setenv(testCase.k, testCase.setV)

				actual := scrapligoutil.GetEnvStrOrDefault(testCase.k, testCase.defaultV)

				if actual != testCase.expected {
					t.Logf(
						"actual and expected outputs do not match, got %q want %q",
						actual,
						testCase.expected,
					)
				}
			})
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	cases := []struct {
		name     string
		k        string
		setV     string
		defaultV int
		expected int
	}{
		{
			name:     "simple-default",
			k:        "SOME_ENV_VAR",
			setV:     "",
			defaultV: 1,
			expected: 1,
		},
		{
			name:     "simple-already-set",
			k:        "SOME_ENV_VAR",
			setV:     "1",
			defaultV: 2,
			expected: 1,
		},
	}

	for _, testCase := range cases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Logf("%s: starting", testCase.name)

				t.Setenv(testCase.k, testCase.setV)

				actual := scrapligoutil.GetEnvIntOrDefault(testCase.k, testCase.defaultV)

				if actual != testCase.expected {
					t.Logf(
						"actual and expected outputs do not match, got %q want %q",
						actual,
						testCase.expected,
					)
				}
			})
	}
}
