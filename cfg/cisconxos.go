package cfg

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

var ErrGetCheckpointFailed = errors.New("get checkpoint operation failed")

type nxosPatterns struct {
	bytesFreePattern      *regexp.Regexp
	buildConfigPattern    *regexp.Regexp
	configVersionPattern  *regexp.Regexp
	configChangePattern   *regexp.Regexp
	outputHeaderPattern   *regexp.Regexp
	checkpointLinePattern *regexp.Regexp
}

var (
	nxosPatternsInstance     *nxosPatterns //nolint:gochecknoglobals
	nxosPatternsInstanceOnce sync.Once     //nolint:gochecknoglobals
)

func getNXOSPatterns() *nxosPatterns {
	nxosPatternsInstanceOnce.Do(func() {
		nxosPatternsInstance = &nxosPatterns{
			bytesFreePattern: regexp.MustCompile(
				`(?i)(?P<bytes_available>\d+)(?: bytes free)`,
			),
			buildConfigPattern: regexp.MustCompile(`(?im)(^!command:.*$)`),
			configVersionPattern: regexp.MustCompile(
				`(?im)(^!running configuration last done.*$)`,
			),
			configChangePattern:   regexp.MustCompile(`(?im)(^!time.*$)`),
			checkpointLinePattern: regexp.MustCompile(`(?m)^\s*!#.*$`),
		}

		nxosPatternsInstance.outputHeaderPattern = regexp.MustCompile(
			fmt.Sprintf(
				`(?im)%s|%s|%s`,
				nxosPatternsInstance.buildConfigPattern.String(),
				nxosPatternsInstance.configVersionPattern.String(),
				nxosPatternsInstance.configChangePattern.String(),
			),
		)
	})

	return nxosPatternsInstance
}

type NXOSCfg struct {
	conn                           *network.Driver
	VersionPattern                 *regexp.Regexp
	Filesystem                     string
	filesystemSpaceAvailBufferPerc float32
	configCommandMap               map[string]string
	replaceConfig                  bool
	CandidateConfigFilename        string
	candidateConfigFilename        string
}

// NewNXOSCfg return a cfg instance setup for an Cisco NXOS device.
func NewNXOSCfg( //nolint:dupl
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	options = append(
		[]Option{
			WithConfigSources([]string{"running", "startup"}),
			WithFilesystem("bootflash:"),
		},
		options...)

	c, err := newCfg(conn, options...)
	if err != nil {
		return nil, err
	}

	c.Platform = &NXOSCfg{
		conn:           conn,
		VersionPattern: regexp.MustCompile(`(?i)\d+\.[a-z0-9\(\).]+`),
		configCommandMap: map[string]string{
			"running": "show running-config",
			"startup": "show startup-config",
		},
		filesystemSpaceAvailBufferPerc: 10.0,
	}

	err = setPlatformOptions(c.Platform, options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (p *NXOSCfg) ClearConfigSession() {
	p.candidateConfigFilename = ""
}

// GetVersion get the version from the device.
func (p *NXOSCfg) GetVersion() (string, []*base.Response, error) {
	versionResult, err := p.conn.SendCommand("show version | i \"NXOS: version\"")
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
func (p *NXOSCfg) GetConfig(source string) (string, []*base.Response, error) {
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

func (p *NXOSCfg) getFilesystemSpaceAvail() (int, error) {
	patterns := getNXOSPatterns()

	filesystemSizeResult, err := p.conn.SendCommand(fmt.Sprintf("dir %s | i bytes", p.Filesystem))
	if err != nil {
		return -1, ErrFailedToDetermineDeviceState
	}

	return parseSpaceAvail(patterns.bytesFreePattern, filesystemSizeResult)
}

func (p *NXOSCfg) prepareConfigPayload(config string) string {
	tclshFilesystem := fmt.Sprintf("/%s/", strings.TrimSuffix(p.Filesystem, ":"))
	tcslhStartFile := fmt.Sprintf(
		`set fl [open "%s%s" wb+]`,
		tclshFilesystem,
		p.candidateConfigFilename,
	)

	splitConfig := strings.Split(config, "\n")

	tclshConfig := make([]string, 0)

	for _, configLine := range splitConfig {
		tclshConfig = append(tclshConfig, fmt.Sprintf("puts -nonewline $fl {%s\n}", configLine))
	}

	tclshEndFile := "close $fl"

	return strings.Join(
		[]string{tcslhStartFile, strings.Join(tclshConfig, "\n"), tclshEndFile},
		"\n",
	)
}

// LoadConfig load a candidate configuration.
func (p *NXOSCfg) LoadConfig(
	config string,
	replace bool,
	options *OperationOptions,
) ([]*base.Response, error) {
	p.replaceConfig = replace

	var scrapliResponses []*base.Response

	filesystemBytesAvail, err := p.getFilesystemSpaceAvail()
	if err != nil {
		return nil, err
	}

	spaceSufficient := isSpaceSufficient(
		filesystemBytesAvail,
		p.filesystemSpaceAvailBufferPerc,
		config,
	)

	if !spaceSufficient {
		return nil, ErrInsufficientSpaceAvailable
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

	err = p.conn.AcquirePriv("tclsh")
	if err != nil {
		return nil, err
	}

	r, err := p.conn.SendConfig(config, base.WithDesiredPrivilegeLevel("tclsh"))
	if err != nil {
		return nil, err
	}

	err = p.conn.AcquirePriv(p.conn.DefaultDesiredPriv)
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, r)

	return scrapliResponses, nil
}

func (p *NXOSCfg) deleteCandidateConfigFile() (*base.MultiResponse, error) {
	deleteCommands := []string{
		"terminal dont-ask",
		fmt.Sprintf("delete %s%s", p.Filesystem, p.candidateConfigFilename),
	}

	return p.conn.SendCommands(deleteCommands)
}

// AbortConfig abort the loaded candidate configuration.
func (p *NXOSCfg) AbortConfig() ([]*base.Response, error) {
	scrapliResponses := make([]*base.Response, 0)

	r, err := p.deleteCandidateConfigFile()
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, r.Responses...)

	return scrapliResponses, err
}

func (p *NXOSCfg) SaveConfig() (*base.Response, error) {
	return p.conn.SendCommand("copy running-config startup-config")
}

// CommitConfig commit the loaded candidate configuration.
func (p *NXOSCfg) CommitConfig(source string) ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	var commitResult *base.Response

	var err error

	if p.replaceConfig {
		commitResult, err = p.conn.SendCommand(
			fmt.Sprintf(
				"rollback running-config file %s%s",
				p.Filesystem,
				p.candidateConfigFilename,
			),
		)
	} else {
		commitResult, err = p.conn.SendCommand(
			fmt.Sprintf("copy %s%s running-config", p.Filesystem, p.candidateConfigFilename),
		)
	}

	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, commitResult)

	saveResult, err := p.SaveConfig()
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, saveResult)

	cleanupResult, err := p.deleteCandidateConfigFile()
	if err != nil {
		return scrapliResponses, err
	}

	scrapliResponses = append(scrapliResponses, cleanupResult.Responses...)

	return scrapliResponses, nil
}

func (p *NXOSCfg) getDiffCommand(source string) string {
	if p.replaceConfig {
		return fmt.Sprintf(
			"show diff rollback-patch %s-config file %s%s",
			source,
			p.Filesystem,
			p.candidateConfigFilename,
		)
	}

	return ""
}

func (p *NXOSCfg) normalizeSourceAndCandidateConfigs(
	sourceConfig, candidateConfig string,
) (normalizedSourceConfig, normalizedCandidateConfig string) {
	patterns := getNXOSPatterns()

	normalizedSourceConfig = patterns.outputHeaderPattern.ReplaceAllString(sourceConfig, "")
	normalizedSourceConfig = strings.Replace(normalizedSourceConfig, "\n\n", "\n", -1)

	normalizedCandidateConfig = patterns.checkpointLinePattern.ReplaceAllString(candidateConfig, "")
	normalizedCandidateConfig = patterns.outputHeaderPattern.ReplaceAllString(
		normalizedCandidateConfig,
		"",
	)
	normalizedCandidateConfig = strings.Replace(normalizedCandidateConfig, "\n\n", "\n", -1)

	return normalizedSourceConfig, normalizedCandidateConfig
}

// DiffConfig diff the candidate configuration against a source config.
func (p *NXOSCfg) DiffConfig(
	source, candidateConfig string,
) (responses []*base.Response,
	normalizedSourceConfig,
	normalizedCandidateConfig,
	deviceDiff string, err error) {
	var scrapliResponses []*base.Response

	deviceDiff = ""

	diffCmd := p.getDiffCommand(source)

	if diffCmd != "" {
		diffResult, diffErr := p.conn.SendCommand(p.getDiffCommand(source))
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

		deviceDiff = diffResult.Result
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

// GetCheckpoint return a checkpoint file from the target device.
func (p *NXOSCfg) GetCheckpoint() (*Response, error) {
	logging.LogDebug(
		FormatLogMessage(p.conn, "info", "get checkpoint requested"),
	)

	r := NewResponse(p.conn.Host, "GetCheckpoint", ErrGetCheckpointFailed)

	timestamp := time.Now().Unix()
	checkpointCommands := []string{
		"terminal dont-ask",
		fmt.Sprintf("checkpoint file %sscrapli_cfg_tmp_%d", p.Filesystem, timestamp),
		fmt.Sprintf("show file %sscrapli_cfg_tmp_%d", p.Filesystem, timestamp),
		fmt.Sprintf("delete %sscrapli_cfg_tmp_%d", p.Filesystem, timestamp),
	}

	scrapliResponses, err := p.conn.SendCommands(checkpointCommands)
	if err != nil {
		return r, err
	}

	r.Record(scrapliResponses.Responses, "")

	return r, nil
}
