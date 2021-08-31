package testhelper

import (
	"regexp"
)

func cleanResponseMap() map[string]func(r string) string {
	return map[string]func(r string) string{
		"arista_eos":    aristaEosCleanResponse,
		"cisco_iosxr":   ciscoIosxrCleanResponse,
		"cisco_iosxe":   ciscoIosxeCleanResponse,
		"cisco_nxos":    ciscoNxosCleanResponse,
		"juniper_junos": juniperJunosCleanResponse,
	}
}

func GetCleanFunc(platform string) func(r string) string {
	cleanFuncs := cleanResponseMap()

	cleanFunc, ok := cleanFuncs[platform]
	if !ok {
		return cleanResponseNoop
	}

	return cleanFunc
}

func cleanResponseNoop(r string) string { return r }

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

func aristaEosCleanResponse(r string) string {
	replacePatterns := getAristaEosReplacePatterns()

	r = replacePatterns.datetimePattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.cryptoPattern.ReplaceAllString(r, "CRYPTO_REPLACED")

	return r
}

type ciscoIosxrReplacePatterns struct {
	datetimePattern         *regexp.Regexp
	cryptoPattern           *regexp.Regexp
	cfgByPattern            *regexp.Regexp
	commitInProgressPattern *regexp.Regexp
}

var ciscoIosxrReplacePatternsInstance *ciscoIosxrReplacePatterns //nolint:gochecknoglobals

func getCiscoIosxrReplacePatterns() *ciscoIosxrReplacePatterns {
	if ciscoIosxrReplacePatternsInstance == nil {
		ciscoIosxrReplacePatternsInstance = &ciscoIosxrReplacePatterns{
			datetimePattern: regexp.MustCompile(
				`(?im)(mon|tue|wed|thu|fri|sat|sun)` +
					`\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)` +
					`\s+\d+\s+\d+:\d+:\d+((\.\d+\s\w+)|\s\d+)`,
			),
			cryptoPattern: regexp.MustCompile(`(?im)^\ssecret\s5\s[\w$./]+$`),
			cfgByPattern: regexp.MustCompile(
				`(?im)^!! Last configuration change at TIME_STAMP_REPLACED by (\w+)$`,
			),
			commitInProgressPattern: regexp.MustCompile(`(?ims)System configuration.*`),
		}
	}

	return ciscoIosxrReplacePatternsInstance
}

func ciscoIosxrCleanResponse(r string) string {
	replacePatterns := getCiscoIosxrReplacePatterns()

	r = replacePatterns.datetimePattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.cryptoPattern.ReplaceAllString(r, "CRYPTO_REPLACED")
	r = replacePatterns.cfgByPattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.commitInProgressPattern.ReplaceAllString(r, "")

	return r
}

type ciscoIosxeReplacePatterns struct {
	configBytesPattern *regexp.Regexp
	datetimePattern    *regexp.Regexp
	cryptoPattern      *regexp.Regexp
	cfgByPattern       *regexp.Regexp
	callHomePattern    *regexp.Regexp
	certLicensePattern *regexp.Regexp
}

var ciscoIosxeReplacePatternsInstance *ciscoIosxeReplacePatterns //nolint:gochecknoglobals

func getCiscoIosxeReplacePatterns() *ciscoIosxeReplacePatterns {
	if ciscoIosxeReplacePatternsInstance == nil {
		ciscoIosxeReplacePatternsInstance = &ciscoIosxeReplacePatterns{
			configBytesPattern: regexp.MustCompile(`(?im)^Current configuration : \d+ bytes$`),
			datetimePattern: regexp.MustCompile(
				`(?im)\d+:\d+:\d+\d+\s+[a-z]{3}\s+(mon|tue|wed|thu|fri|sat|sun)` +
					`\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)\s+\d+\s+\d+`,
			),
			cryptoPattern: regexp.MustCompile(`(?im)^enable secret 5 (.*$)`),
			cfgByPattern: regexp.MustCompile(
				`(?im)^! Last configuration change at TIME_STAMP_REPLACED by (\w+)$`,
			),
			callHomePattern: regexp.MustCompile(
				`(?im)^! Call-home is enabled by Smart-Licensing.$`,
			),
			certLicensePattern: regexp.MustCompile(
				`(?ims)^crypto pki .*\nlicense udi pid CSR1000V sn \w+$`,
			),
		}
	}

	return ciscoIosxeReplacePatternsInstance
}

func ciscoIosxeCleanResponse(r string) string {
	replacePatterns := getCiscoIosxeReplacePatterns()

	r = replacePatterns.configBytesPattern.ReplaceAllString(r, "CONFIG_BYTES_REPLACED")
	r = replacePatterns.datetimePattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.cryptoPattern.ReplaceAllString(r, "CRYPTO_REPLACED")
	r = replacePatterns.cfgByPattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.callHomePattern.ReplaceAllString(r, "CALL_HOME_REPLACED")
	r = replacePatterns.certLicensePattern.ReplaceAllString(r, "CERT_LICENSE_REPLACED")

	return r
}

type ciscoNxosReplacePatterns struct {
	datetimePattern *regexp.Regexp
	cryptoPattern   *regexp.Regexp
	resourcePattern *regexp.Regexp
}

var ciscoNxosReplacePatternsInstance *ciscoNxosReplacePatterns //nolint:gochecknoglobals

func getCiscoNxosReplacePatterns() *ciscoNxosReplacePatterns {
	if ciscoNxosReplacePatternsInstance == nil {
		ciscoNxosReplacePatternsInstance = &ciscoNxosReplacePatterns{
			datetimePattern: regexp.MustCompile(
				`(?im)(mon|tue|wed|thu|fri|sat|sun)\s+` +
					`(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)\s+\d+\s+\d+:\d+:\d+\s\d+`,
			),
			cryptoPattern: regexp.MustCompile(`(?im)^(.*?\s(?:5|md5)\s)[\w$./]+.*$`),
			resourcePattern: regexp.MustCompile(
				`(?im)\d+\smaximum\s\d+$`,
			),
		}
	}

	return ciscoNxosReplacePatternsInstance
}

func ciscoNxosCleanResponse(r string) string {
	replacePatterns := getCiscoNxosReplacePatterns()

	r = replacePatterns.datetimePattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.cryptoPattern.ReplaceAllString(r, "CRYPTO_REPLACED")
	r = replacePatterns.resourcePattern.ReplaceAllString(r, "RESOURCES_REPLACED")

	return r
}

type juniperJunosReplacePatterns struct {
	datetimePattern *regexp.Regexp
	cryptoPattern   *regexp.Regexp
}

var juniperJunosReplacePatternsInstance *juniperJunosReplacePatterns //nolint:gochecknoglobals

func getJuniperJunosReplacePatterns() *juniperJunosReplacePatterns {
	if juniperJunosReplacePatternsInstance == nil {
		juniperJunosReplacePatternsInstance = &juniperJunosReplacePatterns{
			datetimePattern: regexp.MustCompile(
				`(?im)^## Last commit: \d+-\d+-\d+\s\d+:\d+:\d+\s\w+.*$`,
			),
			cryptoPattern: regexp.MustCompile(`(?im)^\s+encrypted-password\s"[\w$./]+";\s.*$`),
		}
	}

	return juniperJunosReplacePatternsInstance
}

func juniperJunosCleanResponse(r string) string {
	replacePatterns := getJuniperJunosReplacePatterns()

	r = replacePatterns.datetimePattern.ReplaceAllString(r, "TIME_STAMP_REPLACED")
	r = replacePatterns.cryptoPattern.ReplaceAllString(r, "CRYPTO_REPLACED")

	return r
}
