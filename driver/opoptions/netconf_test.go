package opoptions_test

import (
	"errors"
	"testing"

	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/opoptions"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testWithFilterType(testName string, testCase *struct {
	description string
	s           string
	o           interface{}
	isignored   bool
},
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithFilterType(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*netconf.OperationOptions)

		if !cmp.Equal(oo.FilterType, testCase.s) {
			t.Fatalf(
				"%s: actual and expected filter types do not match\nactual: %v\nexpected:%v",
				testName,
				oo.FilterType,
				testCase.s,
			)
		}
	}
}

func TestWithFilterType(t *testing.T) {
	cases := map[string]*struct {
		description string
		s           string
		o           interface{}
		isignored   bool
	}{
		"set-filter-type": {
			description: "simple set option test",
			s:           "xpath",
			o:           &netconf.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "xpath",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithFilterType(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithDefaultType(testName string, testCase *struct {
	description string
	s           string
	o           interface{}
	isignored   bool
},
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithDefaultType(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*netconf.OperationOptions)

		if !cmp.Equal(oo.DefaultType, testCase.s) {
			t.Fatalf(
				"%s: actual and expected filter types do not match\nactual: %v\nexpected:%v",
				testName,
				oo.DefaultType,
				testCase.s,
			)
		}
	}
}

func TestWithDefaultType(t *testing.T) {
	cases := map[string]*struct {
		description string
		s           string
		o           interface{}
		isignored   bool
	}{
		"set-default-type": {
			description: "simple set option test",
			s:           "potato",
			o:           &netconf.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "xpath",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithDefaultType(testName, testCase)

		t.Run(testName, f)
	}
}

func testWithFilter(testName string, testCase *struct {
	description string
	s           string
	o           interface{}
	isignored   bool
},
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		err := opoptions.WithFilter(testCase.s)(testCase.o)
		if err != nil {
			if errors.Is(err, util.ErrIgnoredOption) && !testCase.isignored {
				t.Fatalf(
					"%s: option should be ignored, but returned different error",
					testName,
				)
			}

			return
		}

		oo, _ := testCase.o.(*netconf.OperationOptions)

		if !cmp.Equal(oo.Filter, testCase.s) {
			t.Fatalf(
				"%s: actual and expected filter do not match\nactual: %v\nexpected:%v",
				testName,
				oo.Filter,
				testCase.s,
			)
		}
	}
}

func TestWithFilter(t *testing.T) {
	cases := map[string]*struct {
		description string
		s           string
		o           interface{}
		isignored   bool
	}{
		"set-filter": {
			description: "simple set option test",
			s:           "potato",
			o:           &netconf.OperationOptions{},
			isignored:   false,
		},
		"ignored": {
			description: "skipped due to ignored type",
			s:           "xpath",
			o:           &network.Driver{},
			isignored:   true,
		},
	}

	for testName, testCase := range cases {
		f := testWithFilter(testName, testCase)

		t.Run(testName, f)
	}
}
