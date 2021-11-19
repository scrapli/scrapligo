package cfg

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/scrapli/scrapligo/channel"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
)

const (
	filePromptNoisy = "noisy"
	filePromptQuiet = "quiet"
	filePromptAlert = "alert"
)

type iosxePatterns struct {
	bytesFreePattern      *regexp.Regexp
	filePromptModePattern *regexp.Regexp
	outputHeaderPattern   *regexp.Regexp
}

var (
	iosxePatternsInstance     *iosxePatterns //nolint:gochecknoglobals
	iosxePatternsInstanceOnce sync.Once      //nolint:gochecknoglobals
)

func getIOSXEPatterns() *iosxePatterns {
	iosxePatternsInstanceOnce.Do(func() {
		iosxePatternsInstance = &iosxePatterns{
			bytesFreePattern: regexp.MustCompile(
				`(?i)(?P<bytes_available>\d+)(?: bytes free)`,
			),
			filePromptModePattern: regexp.MustCompile(`(?i)(?:file prompt )(?P<prompt_mode>\w+)`),
			// sort of a bad name, but it matches python version --  used to find the version
			// string in the config so we can remove anything in front of it
			outputHeaderPattern: regexp.MustCompile(`(?im)(^version \d+\.\d+$)`),
		}
	})

	return iosxePatternsInstance
}

type IOSXECfg struct {
	conn                           *network.Driver
	VersionPattern                 *regexp.Regexp
	Filesystem                     string
	filesystemSpaceAvailBufferPerc float32
	configCommandMap               map[string]string
	CandidateConfigFilename        string
	candidateConfigFilename        string
	replaceConfig                  bool
}

// NewIOSXECfg return a cfg instance setup for a Cisco IOSXE device.
func NewIOSXECfg( //nolint:dupl
	conn *network.Driver,
	options ...Option,
) (*Cfg, error) {
	options = append(
		[]Option{WithConfigSources([]string{"running", "startup"}), WithFilesystem("flash:")},
		options...)

	c, err := newCfg(conn, options...)
	if err != nil {
		return nil, err
	}

	c.Platform = &IOSXECfg{
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

func (p *IOSXECfg) ClearConfigSession() {
	p.candidateConfigFilename = ""
}

// GetVersion get the version from the device.
func (p *IOSXECfg) GetVersion() (string, []*base.Response, error) {
	versionResult, err := p.conn.SendCommand("show version | i Version")
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
func (p *IOSXECfg) GetConfig(source string) (string, []*base.Response, error) {
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

func (p *IOSXECfg) cleanConfig(config string) string {
	patterns := getIOSXEPatterns()

	configSectionIndices := patterns.outputHeaderPattern.FindStringIndex(config)

	if len(configSectionIndices) == 0 {
		// didnt find the header pattern
		return config
	}

	if len(configSectionIndices) == 2 { //nolint:gomnd
		return config[configSectionIndices[0]:]
	}

	panic("stripping config header failed, this is a bug, provided config is wonky, or both...")
}

func (p *IOSXECfg) prepareConfigPayload(config string) string {
	tcslhStartFile := fmt.Sprintf(
		`puts [open "%s%s" w+] {`,
		p.Filesystem,
		p.candidateConfigFilename,
	)
	tclshEndFile := "}"

	return strings.Join([]string{tcslhStartFile, config, tclshEndFile}, "\n")
}

func (p *IOSXECfg) getFilesystemSpaceAvail() (int, error) {
	patterns := getIOSXEPatterns()

	filesystemSizeResult, err := p.conn.SendCommand(fmt.Sprintf("dir %s | i bytes", p.Filesystem))
	if err != nil {
		return -1, ErrFailedToDetermineDeviceState
	}

	return parseSpaceAvail(patterns.bytesFreePattern, filesystemSizeResult)
}

// LoadConfig load a candidate configuration.
func (p *IOSXECfg) LoadConfig(
	config string,
	replace bool,
	options *OperationOptions,
) ([]*base.Response, error) {
	p.replaceConfig = replace

	var scrapliResponses []*base.Response

	if options.AutoClean {
		config = p.cleanConfig(config)
	}

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

	originalReturnChar := p.conn.Channel.CommsReturnChar
	tclCommsReturnChar := "\r"

	err = p.conn.AcquirePriv("tclsh")
	if err != nil {
		return nil, err
	}

	p.conn.Channel.CommsReturnChar = tclCommsReturnChar

	r, err := p.conn.SendConfig(config, base.WithDesiredPrivilegeLevel("tclsh"))
	if err != nil {
		return nil, err
	}

	scrapliResponses = append(scrapliResponses, r)

	err = p.conn.AcquirePriv(p.conn.DefaultDesiredPriv)
	if err != nil {
		return scrapliResponses, err
	}

	p.conn.Channel.CommsReturnChar = originalReturnChar

	return scrapliResponses, nil
}

func (p *IOSXECfg) determineFilePromptMode() (string, error) {
	r, err := p.conn.SendCommand("show run | i file prompt")
	if err != nil {
		return "", err
	}

	patterns := getIOSXEPatterns()

	filePromptMatch := patterns.filePromptModePattern.FindString(r.Result)

	if filePromptMatch == "" {
		return filePromptAlert, nil
	}

	if strings.Contains(filePromptMatch, filePromptNoisy) {
		return filePromptNoisy, nil
	}

	return filePromptQuiet, nil
}

// AbortConfig abort the loaded candidate configuration.
func (p *IOSXECfg) AbortConfig() ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	r, err := p.deleteCandidateConfigFile()

	scrapliResponses = append(scrapliResponses, r)

	return scrapliResponses, err
}

func (p *IOSXECfg) commitConfigMerge() (*base.Response, error) {
	filePromptMode, err := p.determineFilePromptMode()
	if err != nil {
		return nil, err
	}

	var mergeEvents []*channel.SendInteractiveEvent

	if filePromptMode == filePromptAlert {
		mergeEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput: fmt.Sprintf(
					"copy %s%s running-config",
					p.Filesystem,
					p.candidateConfigFilename,
				),
				ChannelResponse: "Destination filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	} else if filePromptMode == filePromptNoisy {
		mergeEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput: fmt.Sprintf(
					"copy %s%s running-config", p.Filesystem, p.candidateConfigFilename),
				ChannelResponse: "Source filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "Destination filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	} else {
		mergeEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput: fmt.Sprintf(
					"copy %s%s running-config", p.Filesystem, p.candidateConfigFilename),
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	}

	return p.conn.SendInteractive(mergeEvents)
}

// SaveConfig writes running config to startup config.
func (p *IOSXECfg) SaveConfig() (*base.Response, error) {
	filePromptMode, err := p.determineFilePromptMode()
	if err != nil {
		return nil, err
	}

	var saveEvents []*channel.SendInteractiveEvent

	if filePromptMode == filePromptAlert {
		saveEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput:    "copy running-config startup-config",
				ChannelResponse: "Destination filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	} else if filePromptMode == filePromptNoisy {
		saveEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput:    "copy running-config startup-config",
				ChannelResponse: "Source filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "Destination filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	} else {
		saveEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput:    "copy running-config startup-config",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	}

	return p.conn.SendInteractive(saveEvents)
}

func (p *IOSXECfg) deleteCandidateConfigFile() (*base.Response, error) {
	filePromptMode, err := p.determineFilePromptMode()
	if err != nil {
		return nil, err
	}

	var saveEvents []*channel.SendInteractiveEvent

	if filePromptMode == filePromptAlert || filePromptMode == filePromptNoisy {
		saveEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput: fmt.Sprintf(
					"delete %s%s",
					p.Filesystem,
					p.candidateConfigFilename,
				),
				ChannelResponse: "Delete filename",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "[confirm]",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	} else {
		saveEvents = []*channel.SendInteractiveEvent{
			{
				ChannelInput:    fmt.Sprintf("delete %s%s", p.Filesystem, p.candidateConfigFilename),
				ChannelResponse: "[confirm]",
				HideInput:       false,
			},
			{
				ChannelInput:    "",
				ChannelResponse: "",
				HideInput:       false,
			},
		}
	}

	return p.conn.SendInteractive(saveEvents)
}

// CommitConfig commit the loaded candidate configuration.
func (p *IOSXECfg) CommitConfig(source string) ([]*base.Response, error) {
	var scrapliResponses []*base.Response

	var commitResult *base.Response

	var err error

	if p.replaceConfig {
		commitResult, err = p.conn.SendCommand(
			fmt.Sprintf("configure replace %s%s force", p.Filesystem, p.candidateConfigFilename),
		)
	} else {
		commitResult, err = p.commitConfigMerge()
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

	scrapliResponses = append(scrapliResponses, cleanupResult)

	return scrapliResponses, nil
}

func (p *IOSXECfg) getDiffCommand(source string) string {
	if p.replaceConfig {
		return fmt.Sprintf(
			"show archive config differences system:%s-config %s%s",
			source,
			p.Filesystem,
			p.candidateConfigFilename,
		)
	}

	return fmt.Sprintf(
		"show archive config incremental-diffs %s%s ignorecase",
		p.Filesystem,
		p.candidateConfigFilename,
	)
}

func (p *IOSXECfg) normalizeSourceAndCandidateConfigs(
	sourceConfig, candidateConfig string,
) (normalizedSourceConfig, normalizedCandidateConfig string) {
	// remove any of the leading timestamp/building config/config size/last change lines in both the
	// source and candidate configs so they dont need to be compared
	normalizedSourceConfig = p.cleanConfig(sourceConfig)
	normalizedCandidateConfig = p.cleanConfig(candidateConfig)

	return normalizedSourceConfig, normalizedCandidateConfig
}

// DiffConfig diff the candidate configuration against a source config.
func (p *IOSXECfg) DiffConfig(
	source, candidateConfig string,
) (responses []*base.Response,
	normalizedSourceConfig,
	normalizedCandidateConfig,
	deviceDiff string, err error) {
	var scrapliResponses []*base.Response

	diffResult, err := p.conn.SendCommand(p.getDiffCommand(source))
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
