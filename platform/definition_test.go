package platform_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/platform"
)

type testPlatformTestCase struct {
	description string
	f           string
}

func testNewPlatform(testName string, testCase *testPlatformTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		p, err := platform.NewPlatform(testCase.f, "localhost")
		if err != nil {
			t.Fatalf("failed creating platform instance, error: %s", err)
		}

		actual, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("failed marshaling platform, error: %s", err)
		}

		if *update {
			writeGolden(t, testName, actual)
		}

		expected := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !cmp.Equal(actual, expected) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actual,
				expected,
			)
		}
	}
}

func TestNewPlatform(t *testing.T) {
	cases := map[string]*testPlatformTestCase{
		"new-platform-from-assets": {
			description: "simple test to generate platform from embedded assets",
			f:           "cisco_iosxe",
		},
		"new-platform-from-file": {
			description: "simple test to generate platform from user provided file",
			f:           "./test-fixtures/explicit_cisco_iosxe.yaml",
		},
	}

	for testName, testCase := range cases {
		f := testNewPlatform(testName, testCase)

		t.Run(testName, f)
	}
}

type testPlatformVariantTestCase struct {
	description string
	f           string
	variant     string
}

func testNewPlatformVariant(
	testName string,
	testCase *testPlatformVariantTestCase,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		p, err := platform.NewPlatformVariant(testCase.f, testCase.variant, "localhost")
		if err != nil {
			t.Fatalf("failed creating platform instance, error: %s", err)
		}

		actual, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("failed marshaling platform, error: %s", err)
		}

		if *update {
			writeGolden(t, testName, actual)
		}

		expected := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !cmp.Equal(actual, expected) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actual,
				expected,
			)
		}
	}
}

func TestNewPlatformVariant(t *testing.T) {
	cases := map[string]*testPlatformVariantTestCase{
		"new-platform-variant-from-file": {
			description: "simple test to generate platform variant",
			f:           "./test-fixtures/explicit_cisco_iosxe.yaml",
			variant:     "testing1",
		},
	}

	for testName, testCase := range cases {
		f := testNewPlatformVariant(testName, testCase)

		t.Run(testName, f)
	}
}
