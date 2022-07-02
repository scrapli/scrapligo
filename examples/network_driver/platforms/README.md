Platforms
=========

scrapli and scrapligo have always had support for a few popular network operating systems. 
"Support" in this context means that there have been some sane default settings relevant to 
these platforms -- things like sane privilege levels, proper disabling of paging, and a default 
understanding of the "normal" privilege level for sending "show" commands. These platforms are:

- Cisco IOSXE
- Cisco IOSXR
- Cisco NXOS
- Arista EOS
- Juniper JunOS

In earlier versions of scrapligo there was a driver layer called "core". At this layer there 
were structs representing each of these network operating systems. This is/was perfectly 
workable, however it was not very flexible. Scrapligo now supports platform definition via YAML 
or JSON, and embeds YAML platform definitions for each of the above platforms plus Nokia SRLinux 
by default. These platform definitions live in scrapligo/assets.

To instantiate a "platform" you can use the `platform.NewPlatform` function. You must pass a 
platform name or file in addition to your "normal" driver creation options. The value returned 
from this NewPlatform function is a *Platform instance. Depending on the type of platform you 
are instantiating you then need to call the `GetNetworkDriver` or `GetGenericDriver` method of 
the Platform, this will return, as you may expect, the corresponding driver type. After this, 
all operations proceed as "normal".

As mentioned, first argument to the `platform.NewPlatform` function is either a platform name or 
a file path. The valid options for platform name are:

- cisco_iosxe
- cisco_iosxr
- cisco_nxos
- arista_eos
- juniper_junos
- nokia_srl

When passing any of the above platform names, scrapligo will refer to the platform definitions 
included in the binary as assets. If you preferred to provide your own platform definition, 
you can do so by simply passing a file path.