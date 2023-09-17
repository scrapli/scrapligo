package util_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testStripANSI(testName string, in, expected []byte) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		if !cmp.Equal(util.StripANSI(in), expected) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				in,
				expected,
			)
		}
	}
}

func TestStripANSI(t *testing.T) {
	cases := map[string]struct {
		in       []byte
		expected []byte
	}{
		"strip-ansi-simple": {
			in:       []byte("[admin@CoolDevice.Sea1: \x1b[1m/\x1b[0;0m]$"),
			expected: []byte("[admin@CoolDevice.Sea1: /]$"),
		},
	}

	for testName, testCase := range cases {
		f := testStripANSI(testName, testCase.in, testCase.expected)

		t.Run(testName, f)
	}
}

func testByteIsAny(testName string, b byte, l []byte, expected bool) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		if !cmp.Equal(util.ByteIsAny(b, l), expected) {
			t.Fatalf(
				"%s: actual and expected outputs do not match",
				testName,
			)
		}
	}
}

func TestByteIsAny(t *testing.T) {
	cases := map[string]struct {
		b        byte
		l        []byte
		expected bool
	}{
		"byte-is-any-simple": {
			b:        byte(1),
			l:        []byte{3, 2, 1},
			expected: true,
		},
		"byte-is-any-simple-false": {
			b:        byte(0),
			l:        []byte{3, 2, 1},
			expected: false,
		},
	}

	for testName, testCase := range cases {
		f := testByteIsAny(testName, testCase.b, testCase.l, testCase.expected)

		t.Run(testName, f)
	}
}

func testByteContainsAny(testName string, b []byte, l [][]byte, expected bool) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf(
			"%s: starting", testName)

		if !cmp.Equal(util.ByteContainsAny(b, l), expected) {
			t.Fatalf(
				"%s: actual and expected outputs do not match",
				testName,
			)
		}
	}
}

func TestByteContainsAny(t *testing.T) {
	cases := map[string]struct {
		b        []byte
		l        [][]byte
		expected bool
	}{
		"byte-contains-any-simple": {
			b:        []byte("one"),
			l:        [][]byte{[]byte("potato"), []byte("one")},
			expected: true,
		},
		"byte-contains-any-simple-false": {
			b:        []byte("one"),
			l:        [][]byte{[]byte("potato"), []byte("two")},
			expected: false,
		},
	}

	for testName, testCase := range cases {
		f := testByteContainsAny(testName, testCase.b, testCase.l, testCase.expected)

		t.Run(testName, f)
	}
}

func testBytesRoughlyContains(
	testName string,
	input, output []byte,
	expected bool,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf(
			"%s: starting", testName)

		if !cmp.Equal(util.BytesRoughlyContains(input, output), expected) {
			t.Fatalf(
				"%s: actual and expected outputs do not match",
				testName,
			)
		}
	}
}

//nolint:lll
func TestBytesRoughlyContains(t *testing.T) {
	bannerInput := `set system banner motd-banner "................................................................
:                  Welcome to Nokia SR Linux!                  :
:              Open Network OS for the NetOps era.             :
:                                                              :
:    This is a freely distributed official container image.    :
:                      Use it - Share it                       :
:                                                              :
: Get started: https://learn.srlinux.dev                       :
: Container:   https://go.srlinux.dev/container-image          :
: Docs:        https://doc.srlinux.dev/22-11                   :
: Rel. notes:  https://doc.srlinux.dev/rn22-11-2               :
: YANG:        https://yang.srlinux.dev/v22.11.2               :
: Discord:     https://go.srlinux.dev/discord                  :
: Contact:     https://go.srlinux.dev/contact-sales            :
................................................................"`

	bannerOutput := `set system banner motd-banner "................................................................
...:                  Welcome to Nokia SR Linux!                  :
...:              Open Network OS for the NetOps era.             :
...:                                                              :
...:    This is a freely distributed official container image.    :
...:                      Use it - Share it                       :
...:                                                              :
...: Get started: https://learn.srlinux.dev                       :
...: Container:   https://go.srlinux.dev/container-image          :
...: Docs:        https://doc.srlinux.dev/22-11                   :
...: Rel. notes:  https://doc.srlinux.dev/rn22-11-2               :
...: YANG:        https://yang.srlinux.dev/v22.11.2               :
...: Discord:     https://go.srlinux.dev/discord                  :
...: Contact:     https://go.srlinux.dev/contact-sales            :
..................................................................."`

	cases := map[string]struct {
		input    []byte
		output   []byte
		expected bool
	}{
		"exactly-equal": {
			[]byte("potato"),
			[]byte("potato"),
			true,
		},
		"output-has-prefix": {
			[]byte("potato"),
			[]byte("...potato"),
			true,
		},
		"output-has-suffix": {
			[]byte("potato"),
			[]byte("potato..."),
			true,
		},
		"output-has-prefix-and-suffix": {
			[]byte("potato"),
			[]byte("...potato..."),
			true,
		},
		"output-has-newlines": {
			[]byte("potato"),
			[]byte("po\nta\nto"),
			true,
		},
		"output-is-reversed": {
			[]byte("potato"),
			[]byte("otatop"),
			false,
		},
		"nokia_srl_banner": {
			[]byte(bannerInput),
			[]byte(bannerOutput),
			true,
		},
	}

	for testName, testCase := range cases {
		f := testBytesRoughlyContains(testName, testCase.input, testCase.output, testCase.expected)

		t.Run(testName, f)
	}
}
