package options_test

import (
	"flag"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"

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

type optionsIntTestCase struct {
	description string
	i           int
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

type optionsRegexpTestCase struct {
	description string
	p           *regexp.Regexp
	o           interface{}
	isignored   bool
}

type optionsDurationTestCase struct {
	description string
	d           time.Duration
	o           interface{}
	isignored   bool
}

type optionsGenericDriverOnXTestCase struct {
	description string
	f           func(d *generic.Driver) error
	o           interface{}
	isignored   bool
}

type optionsNetworkDriverOnXTestCase struct {
	description string
	f           func(d *network.Driver) error
	o           interface{}
	isignored   bool
}

type optionsNoneTestCase struct {
	description string
	o           interface{}
	isignored   bool
}
