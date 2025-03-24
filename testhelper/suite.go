package testhelper

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/carlmontanari/difflibgo/difflibgo"
	scrapligocli "github.com/scrapli/scrapligo/cli"
)

const (
	cliUserAtHostPattern = "(?im)\\w+@[\\w.]+"
)

// AssertNotDefault asserts `v` is not a default value for that type.
func AssertNotDefault(t *testing.T, v any) {
	t.Helper()

	switch typedV := v.(type) {
	case int, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		if typedV == 0 {
			t.Fatal("expected non-default value")
		}
	case string:
		if typedV == "" {
			t.Fatal("expected non-default value")
		}
	case []byte:
		if typedV == nil {
			t.Fatal("expected non-default value")
		}
	case []*scrapligocli.Result:
		if len(typedV) == 0 {
			t.Fatal("expected non-default value")
		}
	default:
		t.Fatal("un assertable type")
	}
}

// AssertEqual asserts a and b are equal.
func AssertEqual[T comparable](t *testing.T, a, b T) {
	t.Helper()

	if a != b {
		t.Fatalf("expected '%v', got '%v'", a, b)
	}
}

// FailOutput is a simple func to nicely print actual vs expected output when a test fails.
func FailOutput(t *testing.T, actual, expected any) {
	t.Helper()

	switch actual.(type) {
	case string, []byte:
		diff := unifiedDiff(t, actual, expected)

		actualExpectedOut := fmt.Sprintf("\n\033[0;36m*** actual   >>>\033[0m"+
			"\n%s"+
			"\n\033[0;36m<<< actual   ***\033[0m"+
			"\n\033[0;35m*** expected >>>\033[0m"+
			"\n%s"+
			"\n\033[0;35m<<< expected ***\033[0m",
			actual, expected,
		)
		diffOut := fmt.Sprintf("\n\033[0;97m*** diff     >>>\033[0m"+
			"\n%s"+
			"\n\033[0;97m<<< diff     ***\033[0m", diff)

		t.Fatalf(
			"actual and expected outputs do not match...%s%s",
			actualExpectedOut,
			diffOut,
		)
	default:
		t.Fatalf(
			"actual and expected outputs do not match..."+
				"\n\033[0;36m*** actual   >>>\033[0m"+
				"\n%+v"+
				"\n\033[0;36m<<< actual   ***\033[0m"+
				"\n\033[0;35m*** expected >>>\033[0m"+
				"\n%+v"+
				"\n\033[0;35m<<< expected ***\033[0m",
			actual,
			expected,
		)
	}
}

func unifiedDiff(t *testing.T, actual, expected any) string {
	t.Helper()

	var aStr string

	var bStr string

	switch a := actual.(type) {
	case string:
		aStr = a
	case []byte:
		aStr = string(a)
	default:
		aBytes, err := json.MarshalIndent(a, "", "    ")
		if err != nil {
			t.Fatal(err)
		}

		aStr = string(aBytes)
	}

	switch bb := expected.(type) {
	case string:
		bStr = bb
	case []byte:
		bStr = string(bb)
	default:
		bBytes, err := json.MarshalIndent(expected, "", "    ")
		if err != nil {
			t.Fatal(err)
		}

		bStr = string(bBytes)
	}

	return difflibgo.UnifiedDiff(aStr, bStr)
}
