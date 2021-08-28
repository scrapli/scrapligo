package testhelper

import (
	"flag"
	"strings"

	"github.com/scrapli/scrapligo/util"
)

const all = "all"

var Functional = flag.Bool( //nolint:gochecknoglobals
	"functional", false, "perform functional tests")

var FunctionalPlatform = flag.String( //nolint:gochecknoglobals
	"platform", "all", "list comma sep platform(s) to target")

var FunctionalTransport = flag.String( //nolint:gochecknoglobals
	"transport", "all", "list comma sep transport(s) to target")

func RunPlatform(p string) bool {
	if *FunctionalPlatform == all {
		return true
	}

	platformTargetSplit := strings.Split(*FunctionalPlatform, ",")

	return util.StrInSlice(p, platformTargetSplit)
}

func RunTransport(t string) bool {
	if *FunctionalTransport == all {
		return true
	}

	platformTargetSplit := strings.Split(*FunctionalTransport, ",")

	return util.StrInSlice(t, platformTargetSplit)
}
