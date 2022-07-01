package platform

import (
	"fmt"
	"strings"

	"github.com/scrapli/scrapligo/assets"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/scrapli/scrapligo/util"

	"gopkg.in/yaml.v3"
)

const (
	// AristaEos is a constant representing the platform string/name for Arista EOS devices.
	AristaEos = "arista_eos"
	// CiscoIosxe is a constant representing the platform string/name for Cisco IOSXE devices.
	CiscoIosxe = "cisco_iosxe"
	// CiscoIosxr is a constant representing the platform string/name for Cisco IOSXR devices.
	CiscoIosxr = "cisco_iosxr"
	// CiscoNxos is a constant representing the platform string/name for Cisco NXOS devices.
	CiscoNxos = "cisco_nxos"
	// JuniperJunos is a constant representing the platform string/name for Juniper JunOS devices.
	JuniperJunos = "juniper_junos"
	// NokiaSrl is a constant representing the platform string/name for Nokia SRL/SRLinux devices.
	NokiaSrl = "nokia_srl"
)

// GetPlatformNames is used to get the "core" (as in embedded in assets and used in testing)
// platform names.
func GetPlatformNames() []string {
	return []string{
		AristaEos, CiscoIosxe, CiscoIosxr, CiscoNxos, JuniperJunos, NokiaSrl,
	}
}

func loadPlatformDefinitionFromAssets(f string) ([]byte, error) {
	if !strings.HasSuffix(f, ".yaml") {
		f += ".yaml"
	}

	return assets.Assets.ReadFile(fmt.Sprintf("platforms/%s", f))
}

func loadPlatformDefinition(f string) (*Definition, error) {
	b, err := loadPlatformDefinitionFromAssets(f)
	if err != nil {
		b, err = util.ResolveAtFileOrURL(f) //nolint:gosec
		if err != nil {
			return nil, err
		}
	}

	pd := &Definition{}

	err = yaml.Unmarshal(b, pd)
	if err != nil {
		return nil, err
	}

	return pd, nil
}

func setDriver(host string, p *Platform, opts ...util.Option) error {
	finalOpts := p.AsOptions()
	finalOpts = append(finalOpts, opts...)

	var err error

	switch p.DriverType {
	case "generic":
		var d *generic.Driver

		d, err = generic.NewDriver(host, finalOpts...)
		if err != nil {
			return err
		}

		p.genericDriver = d
	case "network":
		var d *network.Driver

		d, err = network.NewDriver(host, finalOpts...)
		if err != nil {
			return err
		}

		p.networkDriver = d
	}

	return err
}

// NewPlatformVariant returns an instance of Platform from the platform definition file f. The
// provided variant data is merged back into the "base" platform definition. The host and
// any provided options are stored and will be applied when fetching the generic or network driver
// via the GetGenericDriver or GetNetworkDriver methods.
func NewPlatformVariant(f, variant, host string, opts ...util.Option) (*Platform, error) {
	pd, err := loadPlatformDefinition(f)
	if err != nil {
		return nil, err
	}

	p := pd.Default

	vp, ok := pd.Variants[variant]
	if !ok {
		return nil, fmt.Errorf("%w: no variant '%s' in platform", util.ErrPlatformError, variant)
	}

	p.mergeVariant(vp)

	err = setDriver(host, p, opts...)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// NewPlatform returns an instance of Platform from the platform definition file f. The host and
// any provided options are stored and will be applied when fetching the generic or network driver
// via the GetGenericDriver or GetNetworkDriver methods.
func NewPlatform(f, host string, opts ...util.Option) (*Platform, error) {
	pd, err := loadPlatformDefinition(f)
	if err != nil {
		return nil, err
	}

	p := pd.Default

	err = setDriver(host, p, opts...)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Definition is a struct representing a JSON or YAML platform definition file.
type Definition struct {
	Default  *Platform            `json:"default"  yaml:"default"`
	Variants map[string]*Platform `json:"variants" yaml:"variants"`
}

// Platform is a struct that contains JSON or YAML data that represent the attributes required to
// create a generic or network driver to connect to a given device type.
type Platform struct {
	// DriverType generic||network
	DriverType string `yaml:"driver-type"`

	FailedWhenContains []string       `json:"failed-when-contains" yaml:"failed-when-contains"`
	OnOpen             onXDefinitions `json:"on-open"              yaml:"on-open"`
	OnClose            onXDefinitions `json:"on-close"             yaml:"on-close"`

	PrivilegeLevels              network.PrivilegeLevels `json:"privilege-levels"                yaml:"privilege-levels"`
	DefaultDesiredPrivilegeLevel string                  `json:"default-desired-privilege-level" yaml:"default-desired-privilege-level"`
	NetworkOnOpen                onXDefinitions          `json:"network-on-open"                 yaml:"network-on-open"`
	NetworkOnClose               onXDefinitions          `json:"network-on-close"                yaml:"network-on-close"`

	Options optionDefinitions `json:"options" yaml:"options"`

	genericDriver *generic.Driver
	networkDriver *network.Driver
}

func (p *Platform) mergeVariant(v *Platform) {
	if v.DriverType != "" {
		p.DriverType = v.DriverType
	}

	if len(v.FailedWhenContains) > 0 {
		p.FailedWhenContains = v.FailedWhenContains
	}

	if v.OnOpen != nil {
		p.OnOpen = v.OnOpen
	}

	if v.OnClose != nil {
		p.OnClose = v.OnClose
	}

	if len(v.PrivilegeLevels) > 0 {
		p.PrivilegeLevels = v.PrivilegeLevels
	}

	if v.DefaultDesiredPrivilegeLevel != "" {
		p.DefaultDesiredPrivilegeLevel = v.DefaultDesiredPrivilegeLevel
	}

	if v.NetworkOnOpen != nil {
		p.NetworkOnOpen = v.NetworkOnOpen
	}

	if v.NetworkOnClose != nil {
		p.NetworkOnClose = v.NetworkOnClose
	}
}

// GetGenericDriver returns an instance of generic.Driver built from the Platform data. If the
// platform data (JSON/YAML) specifies a network driver type this will return an error.
func (p *Platform) GetGenericDriver() (*generic.Driver, error) {
	if p.genericDriver == nil {
		return nil, fmt.Errorf(
			"%w: requested generic driver, but generic driver is nil",
			util.ErrPlatformError,
		)
	}

	return p.genericDriver, nil
}

// GetNetworkDriver returns an instance of network.Driver built from the Platform data. If the
// platform data (JSON/YAML) specifies a generic driver type this will return an error.
func (p *Platform) GetNetworkDriver() (*network.Driver, error) {
	if p.networkDriver == nil {
		return nil, fmt.Errorf(
			"%w: requested network driver, but network driver is nil",
			util.ErrPlatformError,
		)
	}

	return p.networkDriver, nil
}

func (p *Platform) genericOptions() []util.Option {
	opts := make([]util.Option, 0)

	if len(p.FailedWhenContains) > 0 {
		opts = append(opts, options.WithFailedWhenContains(p.FailedWhenContains))
	}

	if len(p.OnOpen) > 0 {
		opts = append(opts, options.WithOnOpen(p.OnOpen.asGenericOnX()))
	}

	if len(p.OnClose) > 0 {
		opts = append(opts, options.WithOnClose(p.OnClose.asGenericOnX()))
	}

	return opts
}

// AsOptions returns a slice of options that the platform represents.
func (p *Platform) AsOptions() []util.Option {
	opts := p.genericOptions()

	opts = append(
		opts,
		options.WithPrivilegeLevels(p.PrivilegeLevels),
		options.WithDefaultDesiredPriv(p.DefaultDesiredPrivilegeLevel),
	)

	if len(p.NetworkOnOpen) > 0 {
		opts = append(opts, options.WithNetworkOnOpen(p.NetworkOnOpen.asNetworkOnX()))
	}

	if len(p.NetworkOnClose) > 0 {
		opts = append(opts, options.WithNetworkOnClose(p.NetworkOnClose.asNetworkOnX()))
	}

	opts = append(opts, p.Options.asOptions()...)

	return opts
}
