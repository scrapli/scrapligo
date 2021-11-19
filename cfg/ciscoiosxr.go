package cfg

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

type iosxrPatterns struct {
	bannerDelimPattern   *regexp.Regexp
	timestampPattern     *regexp.Regexp
	buildConfigPattern   *regexp.Regexp
	configVersionPattern *regexp.Regexp
	configChangePattern  *regexp.Regexp
	outputHeaderPattern  *regexp.Regexp
	endPattern           *regexp.Regexp
}

var (
	iosxrPatternsInstance     *iosxrPatterns //nolint:gochecknoglobals
	iosxrPatternsInstanceOnce sync.Once      //nolint:gochecknoglobals
)

func getIOSXRPatterns() *iosxrPatterns {
	iosxrPatternsInstanceOnce.Do(func() {
		iosxrPatternsInstance = &iosxrPatterns{
			bannerDelimPattern: regexp.MustCompile(
				`(?im)(^banner\s(?:exec|incoming|login|motd|prompt-timeout|slip-ppp)\s(.))`,
			),
			timestampPattern: regexp.MustCompile(
				`(?im)^(mon|tue|wed|thur|fri|sat|sun)\s+` +
					`(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)` +
					`\s+\d+\s+\d+:\d+:\d+((\.\d+\s\w+)|\s\d+)$`,
			),
			buildConfigPattern:   regexp.MustCompile(`(?im)(^building configuration\.{3}$)`),
			configVersionPattern: regexp.MustCompile(`(?im)(^!! ios xr.*$)`),
			configChangePattern:  regexp.MustCompile(`(?im)(^!! last config.*$)`),
			endPattern:           regexp.MustCompile(`end$`),
		}

		iosxrPatternsInstance.outputHeaderPattern = regexp.MustCompile(
			fmt.Sprintf(
				`(?im)%s|%s|%s|%s`,
				iosxrPatternsInstance.timestampPattern.String(),
				iosxrPatternsInstance.buildConfigPattern.String(),
				iosxrPatternsInstance.configVersionPattern.String(),
				iosxrPatternsInstance.configChangePattern.String(),
			),
		)
	})

	return iosxrPatternsInstance
}

type IOSXRCfg struct {
	conn             *network.Driver
	VersionPattern   *regexp.Regexp
	configCommandMap map[string]string
	replaceConfig    bool
	configInProgress bool
	configPrivLevel  string
}

// NewIOSXRCfg return a cfg instance setup for an Cisco IOSXR device.
func NewIOSXRCfg(
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	options = append([]Option{WithConfigSources([]string{"running"})}, options...)

	c, err := newCfg(conn, options...)
	if err != nil {
		return nil, err
	}

	c.Platform = &IOSXRCfg{
		conn:           conn,
		VersionPattern: regexp.MustCompile(`(?i)\d+\.\d+\.\d+`),
		configCommandMap: map[string]string{
			"running": "show running-config",
		},
		configPrivLevel: "configuration",
	}

	err = setPlatformOptions(c.Platform, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *IOSXRCfg) ClearConfigSession() {
	p.configInProgress = false
	p.configPrivLevel = "configuration"
}

// GetVersion get the version from the device.
func (p *IOSXRCfg) GetVersion() (string, []*base.Response, error) {
	var versionResult *base.Response

	var err error

	if !p.configInProgress {
		versionResult, err = p.conn.SendCommand("show version | i Version")
	} else {
		versionResult, err = p.conn.SendConfig("do show version | i Version")
	}

	if err != nil {
		return "", nil, err
	}

	return p.VersionPattern.FindString(
			versionResult.Result,
		), []*base.Response{
			versionResult,
		}, nil
}

// GetConfig get the configuration of a source datastore from the device.
func (p *IOSXRCfg) GetConfig(source string) (string, []*base.Response, error) {
	cmd, err := getConfigCommand(p.configCommandMap, source)
	if err != nil {
		return "", nil, err
	}

	var configResult *base.Response

	if !p.configInProgress {
		configResult, err = p.conn.SendCommand(cmd)
	} else {
		configResult, err = p.conn.SendConfig(cmd, base.WithDesiredPrivilegeLevel(p.configPrivLevel))
	}

	if err != nil {
		return "", nil, err
	}

	return configResult.Result, []*base.Response{configResult}, nil
}

func (p *IOSXRCfg) prepareConfigPayloads(config string) (stdConfig, eagerConfig string) {
	patterns := getIOSXRPatterns()

	// remove comment lines
	config = patterns.outputHeaderPattern.ReplaceAllString(config, "!")

	// remove "end" at the end of the config - if its present it will drop scrapli out
	// of the config session which we do not want
	config = patterns.endPattern.ReplaceAllString(config, "!")

	// find all sections that need to be "eagerly" sent; remove those sections from the "normal"
	// config, then join all the eager sections into a single string
	eagerSections := make([]string, 0)
	bannerSections := patterns.bannerDelimPattern.FindAllString(config, -1)

	for _, bannerHeader := range bannerSections {
		bannerDelim := bannerHeader[len(bannerHeader)-1:]

		currentBannerPattern := regexp.MustCompile(
			fmt.Sprintf(
				`(?ims)^%s.*?%s$`,
				regexp.QuoteMeta(bannerHeader),
				regexp.QuoteMeta(bannerDelim),
			),
		)

		currentBanner := currentBannerPattern.FindString(config)
		eagerSections = append(eagerSections, currentBanner)

		config = strings.Replace(config, currentBanner, "!", 1)
	}

	return config, strings.Join(eagerSections, "\n")
}

// LoadConfig load a candidate configuration.
func (p *IOSXRCfg) LoadConfig(
	config string,
	replace bool,
	options *OperationOptions,
) ([]*base.Response, error) {
	p.replaceConfig = replace
	p.configInProgress = true

	// the actual value is irrelevant, if there is a key "exclusive" w/ any value we assume user is
	// wanting to use configuration_exclusive config mode
	_, ok := options.Kwargs["exclusive"]
	if ok {
		p.configPrivLevel = configExclusivePrivLevel
	}

	var scrapliResponses []*base.Response

	stdConfig, eagerConfig := p.prepareConfigPayloads(config)

	configResult, stdConfigErr := p.conn.SendConfig(
		stdConfig, base.WithDesiredPrivilegeLevel(p.configPrivLevel),
	)
	if stdConfigErr != nil || configResult.Failed != nil {
		return scrapliResponses, stdConfigErr
	}

	scrapliResponses = append(scrapliResponses, configResult)

	eagerResult, eagerConfigErr := p.conn.SendConfig(
		eagerConfig,
		base.WithSendEager(true),
		base.WithDesiredPrivilegeLevel(p.configPrivLevel),
	)

	if eagerConfigErr != nil {
		return scrapliResponses, eagerConfigErr
	}

	if eagerResult.Failed != nil {
		return scrapliResponses, eagerConfigErr
	}

	scrapliResponses = append(scrapliResponses, eagerResult)

	return scrapliResponses, nil
}

// AbortConfig abort the loaded candidate configuration.
func (p *IOSXRCfg) AbortConfig() ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	_, err := p.conn.Channel.SendInput("abort", false, false, p.conn.Channel.TimeoutOps)
	if err != nil {
		return scrapliResponses, err
	}

	p.conn.CurrentPriv = "privilege_exec"
	p.configInProgress = false

	return scrapliResponses, nil
}

// CommitConfig commit the loaded candidate configuration.
func (p *IOSXRCfg) CommitConfig(source string) ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	var commitResult *base.Response

	var err error

	if p.replaceConfig {
		replaceEvents := []*channel.SendInteractiveEvent{
			{ChannelInput: "commit replace", ChannelResponse: "proceed?", HideInput: false},
			{ChannelInput: "yes", ChannelResponse: "", HideInput: false},
		}
		commitResult, err = p.conn.SendInteractive(
			replaceEvents,
			base.WithDesiredPrivilegeLevel(p.configPrivLevel),
		)
	} else {
		commitResult, err = p.conn.SendConfig("commit", base.WithDesiredPrivilegeLevel(p.configPrivLevel))
	}

	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, commitResult)

	p.configInProgress = false

	return scrapliResponses, nil
}

func (p *IOSXRCfg) getDiffCommand() string {
	if p.replaceConfig {
		return "show configuration changes diff"
	}

	return "show commit changes diff"
}

func (p *IOSXRCfg) normalizeSourceAndCandidateConfigs(
	sourceConfig, candidateConfig string,
) (normalizedSourceConfig, normalizedCandidateConfig string) {
	patterns := getIOSXRPatterns()

	normalizedSourceConfig = patterns.outputHeaderPattern.ReplaceAllString(sourceConfig, "")
	normalizedSourceConfig = strings.Replace(normalizedSourceConfig, "\n\n", "\n", -1)

	normalizedCandidateConfig = patterns.outputHeaderPattern.ReplaceAllString(
		candidateConfig,
		"",
	)
	normalizedCandidateConfig = strings.Replace(normalizedCandidateConfig, "\n\n", "\n", -1)

	return normalizedSourceConfig, normalizedCandidateConfig
}

// DiffConfig diff the candidate configuration against a source config.
func (p *IOSXRCfg) DiffConfig(
	source, candidateConfig string,
) (responses []*base.Response,
	normalizedSourceConfig,
	normalizedCandidateConfig,
	deviceDiff string, err error) {
	if source != runningConfig {
		logging.LogDebug(
			FormatLogMessage(
				p.conn,
				"warning",
				"eos only supports diffing against the running config",
			),
		)
	}

	var scrapliResponses []*base.Response

	diffResult, err := p.conn.SendConfig(
		p.getDiffCommand(), base.WithDesiredPrivilegeLevel(p.configPrivLevel),
	)

	if err != nil {
		return scrapliResponses, "", "", "", err
	}

	scrapliResponses = append(scrapliResponses, diffResult)

	if diffResult.Failed != nil {
		logging.LogError(
			FormatLogMessage(
				p.conn,
				"error",
				"failed generating diff for config session",
			),
		)

		return scrapliResponses, "", "", "", nil
	}

	deviceDiff = diffResult.Result

	sourceConfig, getConfigR, err := p.GetConfig(source)
	if err != nil {
		return scrapliResponses, "", "", "", nil
	}

	scrapliResponses = append(scrapliResponses, getConfigR[0])

	if getConfigR[0].Failed != nil {
		logging.LogError(
			FormatLogMessage(
				p.conn,
				"error",
				"failed fetching source config for diff comparison",
			),
		)

		return scrapliResponses, "", "", "", nil
	}

	normalizedSourceConfig, normalizedCandidateConfig = p.normalizeSourceAndCandidateConfigs(
		sourceConfig,
		candidateConfig,
	)

	return scrapliResponses, normalizedSourceConfig, normalizedCandidateConfig, deviceDiff, nil
}
