package testhelper

import (
	"flag"
	"regexp"
	"testing"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

var Functional = flag.Bool( //nolint:gochecknoglobals
	"functional", false, "perform functional tests")

type FunctionalTestHostConnData struct {
	Host       string
	Port       int
	TelnetPort int
}

func FunctionalTestHosts() map[string]*FunctionalTestHostConnData {
	return map[string]*FunctionalTestHostConnData{
		"cisco_iosxe": {
			Host:       "localhost",
			Port:       21022,
			TelnetPort: 21023,
		},
		"arista_eos": {
			Host:       "localhost",
			Port:       24022,
			TelnetPort: 24023,
		},
	}
}

func NewFunctionalTestDriver(
	t *testing.T,
	host, platform, transportName string,
	port int,
) *network.Driver {
	d, driverErr := core.NewCoreDriver(
		host,
		platform,
		base.WithAuthUsername("boxen"),
		base.WithAuthPassword("b0x3N-b0x3N"),
		base.WithAuthSecondary("b0x3N-b0x3N"),
		base.WithPort(port),
		base.WithTransportType(transportName),
		base.WithAuthStrictKey(false),
	)

	if driverErr != nil {
		t.Fatalf("failed creating test device: %v", driverErr)
	}

	return d
}

type aristaEosReplacePatterns struct {
	datetimePattern *regexp.Regexp
	cryptoPattern   *regexp.Regexp
}

var aristaEosReplacePatternsInstance *aristaEosReplacePatterns //nolint:gochecknoglobals

func getAristaEosReplacePatterns() *aristaEosReplacePatterns {
	if aristaEosReplacePatternsInstance == nil {
		aristaEosReplacePatternsInstance = &aristaEosReplacePatterns{
			datetimePattern: regexp.MustCompile(
				`(?im)(mon|tue|wed|thu|fri|sat|sun)` +
					`\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)` +
					`\s+\d+\s+\d+:\d+:\d+\s+\d+$`,
			),
			cryptoPattern: regexp.MustCompile(`(?im)secret\ssha512\s[\w$./]+$`),
		}
	}

	return aristaEosReplacePatternsInstance
}

func CleanResponseNoop(r string) string { return r }

func CleanResponseMap() map[string]func(r string) string {
	return map[string]func(r string) string{
		"arista_eos": AristaEosCleanResponse,
	}
}

func AristaEosCleanResponse(r string) string {
	replacePatterns := getAristaEosReplacePatterns()

	r = replacePatterns.datetimePattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.cryptoPattern.ReplaceAllString(r, "CRYPTO_REPLACED")

	return r
}
