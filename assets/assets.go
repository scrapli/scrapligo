package assets

import "embed"

// Assets is the embedded assets objects for the included platform yaml data.
//
//go:embed definitions/*
var Assets embed.FS

// GetPlatformNames is used to get the "core" (as in embedded in assets and used in testing)
// platform names. If your platform isn't listed here, or you need to have a tweaked platform
// definition you can always pass a definition filename instead of a platform name on Driver
// creation.
func GetPlatformNames() []string {
	return []string{
		"arista_eos",
		//"aruba_wlc",
		"cisco_iosxe",
		//"cisco_iosxr",
		//"cisco_nxos",
		//"cumulus_linux",
		//"cumulus_vtysh",
		//"hp_comware",
		//"huawei_vrp",
		//"ipinfusion_ocnos",
		//"juniper_junos",
		"nokia_srl",
		//"nokia_sros",
		//"nokia_sros_classic",
		//"paloalto_panos",
		//"ruijie_rgos",
		//"vyatta_vyos",
	}
}
