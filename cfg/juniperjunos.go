package cfg

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

type junosPatterns struct {
	outputHeaderPattern *regexp.Regexp
	editPattern         *regexp.Regexp
}

var (
	junosPatternsInstance     *junosPatterns //nolint:gochecknoglobals
	junosPatternsInstanceOnce sync.Once      //nolint:gochecknoglobals
)

func getJUNOSPatterns() *junosPatterns {
	junosPatternsInstanceOnce.Do(func() {
		junosPatternsInstance = &junosPatterns{
			outputHeaderPattern: regexp.MustCompile(`(?im)^## last commit.*$\nversion.*$`),
			editPattern:         regexp.MustCompile(`(?m)^\[edit\]$`),
		}
	})

	return junosPatternsInstance
}

type JUNOSCfg struct {
	conn                    *network.Driver
	VersionPattern          *regexp.Regexp
	Filesystem              string
	replaceConfig           bool
	configInProgress        bool
	CandidateConfigFilename string
	candidateConfigFilename string
	configSetStyle          bool
}

// NewJUNOSCfg return a cfg instance setup for a Juniper JunOS device.
func NewJUNOSCfg(
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	options = append(
		[]Option{
			WithConfigSources([]string{"running"}),
			WithFilesystem("/config/"),
		},
		options...)

	c, err := newCfg(conn, options...)
	if err != nil {
		return nil, err
	}

	c.Platform = &JUNOSCfg{
		conn:                    conn,
		VersionPattern:          regexp.MustCompile(`\d+\.[\w-]+\.\w+`),
		candidateConfigFilename: "",
	}

	err = setPlatformOptions(c.Platform, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *JUNOSCfg) ClearConfigSession() {
	p.candidateConfigFilename = ""
	p.configInProgress = false
	p.configSetStyle = false
}

// GetVersion get the version from the device.
func (p *JUNOSCfg) GetVersion() (string, []*base.Response, error) {
	var versionResult *base.Response

	var err error

	if !p.configInProgress {
		versionResult, err = p.conn.SendCommand("show version")
	} else {
		versionResult, err = p.conn.SendConfig("run show version")
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
func (p *JUNOSCfg) GetConfig(source string) (string, []*base.Response, error) {
	var configResult *base.Response

	var err error

	if !p.configInProgress {
		configResult, err = p.conn.SendCommand("show configuration")
	} else {
		configResult, err = p.conn.SendConfig("run show configuration")
	}

	if err != nil {
		return "", nil, err
	}

	return configResult.Result, []*base.Response{configResult}, nil
}

func (p *JUNOSCfg) prepareConfigPayload(config string) string {
	finalConfigs := make([]string, 0)

	for _, configLine := range strings.Split(config, "\n") {
		finalConfigs = append(
			finalConfigs,
			fmt.Sprintf("echo >> %s%s '%s'", p.Filesystem, p.candidateConfigFilename, configLine),
		)
	}

	return strings.Join(
		finalConfigs,
		"\n",
	)
}

// LoadConfig load a candidate configuration.
func (p *JUNOSCfg) LoadConfig(
	config string,
	replace bool,
	options *OperationOptions,
) ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	p.replaceConfig = replace

	// the actual value is irrelevant, if there is a key "set" w/ any value we assume user is
	// loading a "set" style config
	_, ok := options.Kwargs["set"]
	if ok {
		p.configSetStyle = true
	}

	if p.candidateConfigFilename == "" {
		p.candidateConfigFilename = determineCandidateConfigFilename(p.CandidateConfigFilename)

		logging.LogDebug(
			FormatLogMessage(
				p.conn,
				"debug",
				fmt.Sprintf(
					"candidate configuration filename name will be %s",
					p.candidateConfigFilename,
				),
			),
		)
	}

	config = p.prepareConfigPayload(config)

	r, err := p.conn.SendConfig(config, base.WithDesiredPrivilegeLevel("root_shell"))
	if err != nil {
		return nil, err
	}

	p.configInProgress = true

	scrapliResponses = append(scrapliResponses, r)

	loadCommand := fmt.Sprintf("load override %s%s", p.Filesystem, p.candidateConfigFilename)
	if !p.replaceConfig {
		loadCommand = fmt.Sprintf("load merge %s%s", p.Filesystem, p.candidateConfigFilename)
		if p.configSetStyle {
			loadCommand = fmt.Sprintf("load set %s%s", p.Filesystem, p.candidateConfigFilename)
		}
	}

	loadResult, err := p.conn.SendConfig(loadCommand)
	if err != nil {
		return nil, err
	}

	scrapliResponses = append(scrapliResponses, loadResult)

	return scrapliResponses, nil
}

func (p *JUNOSCfg) deleteCandidateConfigFile() (*base.Response, error) {
	deleteCommand := fmt.Sprintf("rm %s%s", p.Filesystem, p.candidateConfigFilename)

	return p.conn.SendConfig(deleteCommand, base.WithDesiredPrivilegeLevel("root_shell"))
}

// AbortConfig abort the loaded candidate configuration.
func (p *JUNOSCfg) AbortConfig() ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	rollbackResponse, err := p.conn.SendConfig("rollback 0")
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, rollbackResponse)

	deleteResponse, err := p.deleteCandidateConfigFile()
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, deleteResponse)

	return scrapliResponses, nil
}

// CommitConfig commit the loaded candidate configuration.
func (p *JUNOSCfg) CommitConfig(source string) ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	commitResult, err := p.conn.SendConfig("commit")

	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, commitResult)

	cleanupResult, err := p.deleteCandidateConfigFile()
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, cleanupResult)

	return scrapliResponses, nil
}

func (p *JUNOSCfg) normalizeSourceAndCandidateConfigs(
	sourceConfig, candidateConfig string,
) (normalizedSourceConfig, normalizedCandidateConfig string) {
	patterns := getJUNOSPatterns()

	normalizedSourceConfig = patterns.outputHeaderPattern.ReplaceAllString(sourceConfig, "")
	normalizedSourceConfig = patterns.editPattern.ReplaceAllString(normalizedSourceConfig, "")
	normalizedSourceConfig = strings.Replace(normalizedSourceConfig, "\n\n", "\n", -1)

	normalizedCandidateConfig = patterns.outputHeaderPattern.ReplaceAllString(
		candidateConfig,
		"",
	)
	normalizedCandidateConfig = strings.Replace(normalizedCandidateConfig, "\n\n", "\n", -1)
	normalizedCandidateConfig = patterns.editPattern.ReplaceAllString(normalizedCandidateConfig, "")

	return normalizedSourceConfig, normalizedCandidateConfig
}

// DiffConfig diff the candidate configuration against a source config.
func (p *JUNOSCfg) DiffConfig(
	source, candidateConfig string,
) (responses []*base.Response,
	normalizedSourceConfig,
	normalizedCandidateConfig,
	deviceDiff string, err error) {
	var scrapliResponses []*base.Response

	diffResult, diffErr := p.conn.SendConfig("show | compare")
	if diffErr != nil {
		return scrapliResponses, "", "", "", diffErr
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
