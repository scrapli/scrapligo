package driver

// PlatformName is an enum(ish) representing the name of a Platform.
type PlatformName string

// String (stringer) method for PlatformName for formatting and/or just to NokiaSrl.String() when
// passing to.
func (p PlatformName) String() string {
	return string(p)
}

const (
	AristaEos        PlatformName = "arista_eos"
	ArubaWlc         PlatformName = "aruba_wlc"
	CiscoIosxe       PlatformName = "cisco_iosxe"
	CiscoIosxr       PlatformName = "cisco_iosxr"
	CiscoNxos        PlatformName = "cisco_nxos"
	CumulusLinux     PlatformName = "cumulus_linux"
	CumulusVtysh     PlatformName = "cumulus_vtysh"
	HpComware        PlatformName = "hp_comware"
	HuaweiVrp        PlatformName = "huawei_vrp"
	IpInfusionOcnos  PlatformName = "ip_infusion_ocnos"
	JuniperJunos     PlatformName = "juniper_junos"
	NokiaSrl         PlatformName = "nokia_srl"
	NokiaSros        PlatformName = "nokia_sros"
	NokiaSrosClassic PlatformName = "nokia_sros_classic"
	PaloAltoPanos    PlatformName = "paloalto_panos"
	RuijieRgos       PlatformName = "rujie_rgos"
	VyattaVyos       PlatformName = "vyatta_vyos"
)

// GetPlatformNames is used to get the "core" (as in embedded in assets and used in testing)
// platform names. If your platform isn't listed here, or you need to have a tweaked platform
// definition you can always pass a definition filename instead of a platform name on Driver
// creation.
func GetPlatformNames() []string {
	return []string{
		string(AristaEos),
		string(ArubaWlc),
		string(CiscoIosxe),
		string(CiscoIosxr),
		string(CiscoNxos),
		string(CumulusLinux),
		string(CumulusVtysh),
		string(HpComware),
		string(HuaweiVrp),
		string(IpInfusionOcnos),
		string(JuniperJunos),
		string(NokiaSrl),
		string(NokiaSros),
		string(NokiaSrosClassic),
		string(PaloAltoPanos),
		string(RuijieRgos),
		string(VyattaVyos),
	}
}
