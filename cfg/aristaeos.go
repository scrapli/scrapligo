package cfg

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

type eosPatterns struct {
	globalCommentLinePattern *regexp.Regexp
	bannerPattern            *regexp.Regexp
	endPattern               *regexp.Regexp
}

var (
	eosPatternsInstance     *eosPatterns //nolint:gochecknoglobals
	eosPatternsInstanceOnce sync.Once    //nolint:gochecknoglobals
)

func getEOSPatterns() *eosPatterns {
	eosPatternsInstanceOnce.Do(func() {
		eosPatternsInstance = &eosPatterns{
			globalCommentLinePattern: regexp.MustCompile(`(?im)^! .*$`),
			bannerPattern:            regexp.MustCompile(`(?ims)^banner.*EOF$`),
			endPattern:               regexp.MustCompile(`end$`),
		}
	})

	return eosPatternsInstance
}

type EOSCfg struct {
	conn              *network.Driver
	VersionPattern    *regexp.Regexp
	configCommandMap  map[string]string
	configSessionName string
}

// NewEOSCfg return a cfg instance setup for an Arista EOS device.
func NewEOSCfg(
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	options = append([]Option{WithConfigSources([]string{"running", "startup"})}, options...)

	c, err := newCfg(conn, options...)
	if err != nil {
		return nil, err
	}

	c.Platform = &EOSCfg{
		conn:           conn,
		VersionPattern: regexp.MustCompile(`(?i)\d+\.\d+\.[a-z0-9\-]+(\.\d+[a-z]?)?`),
		configCommandMap: map[string]string{
			"running": "show running-config",
			"startup": "show startup-config",
		},
	}

	err = setPlatformOptions(c.Platform, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *EOSCfg) ClearConfigSession() {
	p.configSessionName = ""
}

// GetVersion get the version from the device.
func (p *EOSCfg) GetVersion() (string, []*base.Response, error) {
	versionResult, err := p.conn.SendCommand("show version | i Software image version")
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
func (p *EOSCfg) GetConfig(source string) (string, []*base.Response, error) {
	cmd, err := getConfigCommand(p.configCommandMap, source)
	if err != nil {
		return "", nil, err
	}

	configResult, err := p.conn.SendCommand(cmd)

	if err != nil {
		return "", nil, err
	}

	return configResult.Result, []*base.Response{configResult}, nil
}

func (p *EOSCfg) prepareConfigPayloads(config string) (stdConfig, eagerConfig string) {
	patterns := getEOSPatterns()

	// remove comment lines
	config = patterns.globalCommentLinePattern.ReplaceAllString(config, "!")

	// remove "end" at the end of the config - if its present it will drop scrapli out
	// of the config session which we do not want
	config = patterns.endPattern.ReplaceAllString(config, "!")

	// find all sections that need to be "eagerly" sent; remove those sections from the "normal"
	// config, then join all the eager sections into a single string
	eagerSections := patterns.bannerPattern.FindStringSubmatch(config)
	eagerConfig = strings.Join(eagerSections, "\n")

	for _, section := range eagerSections {
		config = strings.Replace(config, section, "!", -1)
	}

	return config, eagerConfig
}

// RegisterConfigSession register a configuration session in EOS.
func (p *EOSCfg) RegisterConfigSession(sessionName string) error {
	_, ok := p.conn.PrivilegeLevels[sessionName]

	if ok {
		return ErrConfigSessionAlreadyExists
	}

	sessionPrompt := regexp.QuoteMeta(sessionName[:6])
	sessionPromptPattern := fmt.Sprintf(
		`(?im)^[\w.\-@()/:\s]{1,63}\(config\-s\-%s[\w.\-@_/:]{0,32}\)#\s?$`,
		sessionPrompt,
	)

	sessionPrivilegeLevel := &base.PrivilegeLevel{
		Pattern:        sessionPromptPattern,
		Name:           sessionName,
		PreviousPriv:   execPrivLevel,
		Deescalate:     "end",
		Escalate:       fmt.Sprintf("configure session %s", sessionName),
		EscalateAuth:   false,
		EscalatePrompt: "",
	}

	p.conn.PrivilegeLevels[sessionName] = sessionPrivilegeLevel
	p.conn.UpdatePrivilegeLevels()

	return nil
}

func (p *EOSCfg) loadConfig(
	stdConfig, eagerConfig string,
	replace bool,
) ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	if replace {
		rollbackCleanConfigResult, rollbackErr := p.conn.SendConfig("rollback clean-config",
			base.WithDesiredPrivilegeLevel(p.configSessionName))
		if rollbackErr != nil {
			return scrapliResponses, rollbackErr
		}

		scrapliResponses = append(scrapliResponses, rollbackCleanConfigResult)
	}

	configResult, stdConfigErr := p.conn.SendConfig(
		stdConfig,
		base.WithDesiredPrivilegeLevel(p.configSessionName),
	)
	if stdConfigErr != nil || configResult.Failed != nil {
		return scrapliResponses, stdConfigErr
	}

	scrapliResponses = append(scrapliResponses, configResult)

	eagerResult, eagerConfigErr := p.conn.SendConfig(
		eagerConfig,
		base.WithDesiredPrivilegeLevel(p.configSessionName),
		base.WithSendEager(true),
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

// LoadConfig load a candidate configuration.
func (p *EOSCfg) LoadConfig(
	config string,
	replace bool,
	options *OperationOptions,
) ([]*base.Response, error) {
	// options are unused for eos load config
	_ = options

	stdConfig, eagerConfig := p.prepareConfigPayloads(config)

	if p.configSessionName == "" {
		p.configSessionName = fmt.Sprintf("scrapli_cfg_%d", time.Now().Unix())

		logging.LogDebug(
			FormatLogMessage(
				p.conn,
				"debug",
				fmt.Sprintf("configuration session name will be %s", p.configSessionName),
			),
		)

		err := p.RegisterConfigSession(p.configSessionName)
		if err != nil {
			return nil, err
		}
	}

	return p.loadConfig(stdConfig, eagerConfig, replace)
}

// AbortConfig abort the loaded candidate configuration.
func (p *EOSCfg) AbortConfig() ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	err := p.conn.AcquirePriv(p.configSessionName)
	if err != nil {
		return scrapliResponses, err
	}

	_, err = p.conn.Channel.SendInput("abort", false, false, p.conn.Channel.TimeoutOps)
	if err != nil {
		return scrapliResponses, err
	}

	p.conn.CurrentPriv = execPrivLevel

	return scrapliResponses, nil
}

// CommitConfig commit the loaded candidate configuration.
func (p *EOSCfg) CommitConfig(source string) ([]*base.Response, error) {
	if source != runningConfig {
		logging.LogDebug(
			FormatLogMessage(
				p.conn,
				"warning",
				"eos only supports committing to running config, running config is automatically copied to "+
					"startup during commit operation",
			),
		)
	}

	var scrapliResponses []*base.Response

	commands := []string{
		fmt.Sprintf("configure session %s commit", p.configSessionName),
		"copy running-config startup-config",
	}

	m, err := p.conn.SendCommands(commands)
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, m.Responses...)

	return scrapliResponses, err
}

func (p *EOSCfg) normalizeSourceAndCandidateConfigs(
	sourceConfig, candidateConfig string,
) (normalizedSourceConfig, normalizedCandidateConfig string) {
	patterns := getEOSPatterns()

	// Remove all comment lines from both the source and candidate configs -- this is only done
	// here pre-diff, so we dont modify the user provided candidate config which can totally have
	// those comment lines - we only remove "global" (top level) comments though... user comments
	// attached to interfaces and the stuff will remain
	normalizedSourceConfig = patterns.globalCommentLinePattern.ReplaceAllString(sourceConfig, "")
	normalizedSourceConfig = strings.Replace(normalizedSourceConfig, "\n\n", "\n", -1)

	normalizedCandidateConfig = patterns.globalCommentLinePattern.ReplaceAllString(
		candidateConfig,
		"",
	)
	normalizedCandidateConfig = strings.Replace(normalizedCandidateConfig, "\n\n", "\n", -1)

	return normalizedSourceConfig, normalizedCandidateConfig
}

// DiffConfig diff the candidate configuration against a source config.
func (p *EOSCfg) DiffConfig(
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
		"show session-config diffs",
		base.WithDesiredPrivilegeLevel(p.configSessionName),
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
