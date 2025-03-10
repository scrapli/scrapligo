package driver

// PlatformName is an enum(ish) representing the name of a Platform.
type PlatformName string

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
func GetPlatformNames() []PlatformName {
	return []PlatformName{
		AristaEos,
		ArubaWlc,
		CiscoIosxe,
		CiscoIosxr,
		CiscoNxos,
		CumulusLinux,
		CumulusVtysh,
		HpComware,
		HuaweiVrp,
		IpInfusionOcnos,
		JuniperJunos,
		NokiaSrl,
		NokiaSros,
		NokiaSrosClassic,
		PaloAltoPanos,
		RuijieRgos,
		VyattaVyos,
	}
}
