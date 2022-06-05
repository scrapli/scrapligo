package opoptions_test

import (
	"flag"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/util"
)

var (
	update = flag.Bool( //nolint
		"update",
		false,
		"update the golden files",
	)
	functional = flag.Bool( //nolint
		"functional",
		false,
		"execute functional tests",
	)
	platforms = flag.String( //nolint
		"platforms",
		util.All,
		"comma sep list of platform(s) to target",
	)
	transports = flag.String( //nolint
		"transports",
		util.All,
		"comma sep list of transport(s) to target",
	)
)

type optionsBoolTestCase struct {
	description string
	b           bool
	o           interface{}
	isignored   bool
}

type optionsStringTestCase struct {
	description string
	s           string
	o           interface{}
	isignored   bool
	iserr       bool
}

type optionsStringSliceTestCase struct {
	description string
	ss          []string
	o           interface{}
	isignored   bool
}

type optionsDurationTestCase struct {
	description string
	d           time.Duration
	o           interface{}
	isignored   bool
}

type optionsRegexpTestCase struct {
	description string
	p           *regexp.Regexp
	o           interface{}
	isignored   bool
}
