<p align=center><a href=""><img src=scrapligo.svg?sanitize=true/></a></p>

[![Go Report](https://img.shields.io/badge/go%20report-A%2B-blue?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://goreportcard.com/report/github.com/scrapli/scrapligo)
[![License: MIT](https://img.shields.io/badge/License-MIT-blueviolet.svg?style=flat-square)](https://opensource.org/licenses/MIT)


---

**Source Code**: <a href="https://github.com/scrapli/scrapligo" target="_blank">https://github.com/scrapli/scrapligo</a>

**Examples**: <a href="https://github.com/scrapli/scrapligo/tree/main/examples" target="_blank">https://github.com/scrapli/scrapligo/tree/main/examples</a>

**Go Docs**: <a href="https://pkg.go.dev/github.com/scrapli/scrapligo" target="_blank">https://pkg.go.dev/github.com/scrapli/scrapligo</a>

---

scrapligo -- scrap(e c)li (but in go!) --  is a Go library focused on connecting to devices, specifically network devices
(routers/switches/firewalls/etc.) via SSH and NETCONF.


#### Key Features:

- __Easy__: It's easy to get going with scrapligo -- if you are familiar with go and/or scrapli you are already most of 
  the way there! Check out the examples linked above to get started! 
- __Fast__: Do you like to go fast? Of course you do! All of scrapli is built with speed in mind, but this port of 
  scrapli to go is of course even faster than its python sibling! 
- __But wait, there's more!__: Have NETCONF devices in your environment, but love the speed and simplicity of
  scrapli? You're in luck! NETCONF support is built right into scrapligo!
- __Sounds great, but I'm a Pythonista__: No worries! scrapligo is inspired by (and is sort of a 
  port of) [scrapli](https://github.com/carlmontanari/scrapli), so check that out!


## Running the Examples

You need [Go 1.16+](https://golang.org/doc/install) installed. Clone the repo and `go run` any 
of the examples in the [examples](/examples) folder. Below are a few example outputs.

### Executing a number of commands (from a file)

```bash
$  go run examples/base_driver/main.go
found prompt: 
csr1000v-1#


sent command 'show version', output received:
 Cisco IOS XE Software, Version 16.09.03
Cisco IOS Software [Fuji], Virtual XE Software (X86_64_LINUX_IOSD-UNIVERSALK9-M), Version 16.9.3, RELEASE SOFTWARE (fc2)
Technical Support: http://www.cisco.com/techsupport
Copyright (c) 1986-2019 by Cisco Systems, Inc.
Compiled Wed 20-Mar-19 07:56 by mcpre
...
```

### Parsing a command output

For more details, check out [Network automation options in Go with scrapligo](https://netdevops.me/2021/network-automation-options-in-go-with-scrapligo/).

```yaml
$  go run examples/network_driver/textfsm/main.go
Hostname: csr1000v-1
SW Version: 16.9.3
Uptime: 18 minutes
```


## Code Example

```go
package main

import (
  "fmt"
  "github.com/scrapli/scrapligo/driver/options"
  "github.com/scrapli/scrapligo/platform"
)

func main() {
    p, err := platform.NewPlatform(
        // cisco_iosxe refers to the included cisco iosxe platform definition
        "cisco_iosxe",
        "sandbox-iosxe-latest-1.cisco.com",
        options.WithAuthNoStrictKey(),
        options.WithAuthUsername("developer"),
        options.WithAuthPassword("C1sco12345"),
    )
    if err != nil {
        fmt.Printf("failed to create platform; error: %+v\n", err)

		return
    }

    d, err := p.GetNetworkDriver()
    if err != nil {
        fmt.Printf("failed to fetch network driver from the platform; error: %+v\n", err)

        return
    }

    err = d.Open()
    if err != nil {
        fmt.Printf("failed to open driver; error: %+v\n", err)

        return
    }

    defer d.Close()

    r, err := d.SendCommand("show version")
    if err != nil {
        fmt.Printf("failed to send command; error: %+v\n", err)
        return
    }
  
    fmt.Printf(
        "sent command '%s', output received (SendCommand):\n %s\n\n\n",
        r.Input,
        r.Result,
    )
}
```

## Migrating From Pre 1.0.0

Scrapligo has had a very significant overhaul from v0.x.x versions to the v1.0.0 version, while 
the user facing API stayed *similar* (with the very notable exception to actual import paths), 
here are some (maybe not fully inclusive) list of changes to take note of:

- "cfg" behavior has moved to its own [repo]("https://github.com/scrapli/scrapligocfg)
  - "cfg" also changed a fair bit -- but mostly just from an organizational perspective, so 
    imports are not where you left them, but everything else should be more or less the same!  
- There is no more "base" driver -- this has all been rolled into the "generic" driver.
- All driver options now live in "driver/options" -- these are options that you can pass during 
  driver creation to modify the drivers behavior.
- All "operation" options now live in "driver/opoptions" -- these are options like 
  "WithPrivilegeLevel" and similar.
- Channel level operations no longer require explicit arguments and instead accept options (you 
  *probably* shouldn't be using the Channel directly anyway, so this shouldn't matter much!).
- All bool options now accept no arguments, when these options are provided they simply negate 
  the default behavior rather than accept true/false and set that.
- Logging is no longer "global" for scrapligo, and is instead set per driver -- much more like 
  scrapli Python.
- On "X" (OnOpen/OnClose) functions now have two flavors -- Generic and Network -- which 
  correspond to the driver type they are attached to. For the most part you can get away with 
  Generic OnX functions as you are usually just disabling paging or setting terminal flags, but 
  the Network flavor exists if you need to care about privilege levels and such.
- On "X" functions are now set via options (as opposed to manually setting them like pre 1.0.0).
- There are no more "core" platforms/drivers -- these have all moved to simple YAML definitions. 
  The core drivers you know and love still exist as an embedded YAML asset.
- Platforms now support "variants" (as in scrapli community (Python)) which allow for multiple 
  "variations" of a driver/platform type.




<small>* gopher artwork by [@egonelbre](https://github.com/egonelbre/gophers)</small>
