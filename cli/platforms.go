// DO NOT EDIT, GENERATED FILE
package cli

// PlatformName is an enum(ish) representing the name of a Platform.
type PlatformName string

// String (stringer) method for PlatformName for formatting and/or just to NokiaSrl.String() when
// passing to.
func (p PlatformName) String() string {
	return string(p)
}

const (
	AethraAtosnt         PlatformName = "aethra_atosnt"
	AlcatelAos           PlatformName = "alcatel_aos"
	AristaEos            PlatformName = "arista_eos"
	ArubaAoscx           PlatformName = "aruba_aoscx"
	ArubaWlc             PlatformName = "aruba_wlc"
	CiscoAireos          PlatformName = "cisco_aireos"
	CiscoAsa             PlatformName = "cisco_asa"
	CiscoCbs             PlatformName = "cisco_cbs"
	CiscoFtd             PlatformName = "cisco_ftd"
	CiscoIosxe           PlatformName = "cisco_iosxe"
	CiscoIosxr           PlatformName = "cisco_iosxr"
	CiscoNxos            PlatformName = "cisco_nxos"
	CumulusLinux         PlatformName = "cumulus_linux"
	CumulusVtysh         PlatformName = "cumulus_vtysh"
	DatacomDmos          PlatformName = "datacom_dmos"
	DatacomDmswitch      PlatformName = "datacom_dmswitch"
	Default              PlatformName = "default"
	DellEmc              PlatformName = "dell_emc"
	DellEnterprisesonic  PlatformName = "dell_enterprisesonic"
	DlinkOs              PlatformName = "dlink_os"
	EdgecoreEcs          PlatformName = "edgecore_ecs"
	EltexEsr             PlatformName = "eltex_esr"
	FortinetFortios      PlatformName = "fortinet_fortios"
	FortinetWlc          PlatformName = "fortinet_wlc"
	HpComware            PlatformName = "hp_comware"
	HuaweiSmartax        PlatformName = "huawei_smartax"
	HuaweiVrp            PlatformName = "huawei_vrp"
	IpinfusionOcnos      PlatformName = "ipinfusion_ocnos"
	JuniperJunos         PlatformName = "juniper_junos"
	MikrotikRouteros     PlatformName = "mikrotik_routeros"
	NokiaSrlinux         PlatformName = "nokia_srlinux"
	NokiaSros            PlatformName = "nokia_sros"
	NokiaSrosClassic     PlatformName = "nokia_sros_classic"
	NokiaSrosClassicAram PlatformName = "nokia_sros_classic_aram"
	PaloaltoPanos        PlatformName = "paloalto_panos"
	RaisecomRos          PlatformName = "raisecom_ros"
	RuckusFastiron       PlatformName = "ruckus_fastiron"
	RuckusUnleashed      PlatformName = "ruckus_unleashed"
	RuijieRgos           PlatformName = "ruijie_rgos"
	SiemensRoxii         PlatformName = "siemens_roxii"
	VersaFlexvnf         PlatformName = "versa_flexvnf"
	VyosVyos             PlatformName = "vyos_vyos"
	ZyxelDslam           PlatformName = "zyxel_dslam"
)

// GetPlatformNames is used to get the "core" (as in embedded in assets and used in testing)
// platform names. If your platform isn't listed here, or you need to have a tweaked platform
// definition you can always pass a definition filename instead of a platform name on Cli
// creation.
func GetPlatformNames() []string {
	return []string{
		string(AethraAtosnt),
		string(AlcatelAos),
		string(AristaEos),
		string(ArubaAoscx),
		string(ArubaWlc),
		string(CiscoAireos),
		string(CiscoAsa),
		string(CiscoCbs),
		string(CiscoFtd),
		string(CiscoIosxe),
		string(CiscoIosxr),
		string(CiscoNxos),
		string(CumulusLinux),
		string(CumulusVtysh),
		string(DatacomDmos),
		string(DatacomDmswitch),
		string(Default),
		string(DellEmc),
		string(DellEnterprisesonic),
		string(DlinkOs),
		string(EdgecoreEcs),
		string(EltexEsr),
		string(FortinetFortios),
		string(FortinetWlc),
		string(HpComware),
		string(HuaweiSmartax),
		string(HuaweiVrp),
		string(IpinfusionOcnos),
		string(JuniperJunos),
		string(MikrotikRouteros),
		string(NokiaSrlinux),
		string(NokiaSros),
		string(NokiaSrosClassic),
		string(NokiaSrosClassicAram),
		string(PaloaltoPanos),
		string(RaisecomRos),
		string(RuckusFastiron),
		string(RuckusUnleashed),
		string(RuijieRgos),
		string(SiemensRoxii),
		string(VersaFlexvnf),
		string(VyosVyos),
		string(ZyxelDslam),
	}
}
