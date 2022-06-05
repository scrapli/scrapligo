Basics
======

The `netconf.Driver` is, as you may expect, a driver for doing NETCONF things. This driver 
implements many of the standard NETCONF RPCs, including GetConfig, Get, Commit, Discard, etc..

Note that the default port is always 22 -- so if your NETCONF server is listening on 830 (or 
anything else for that matter), make sure you pass the `WithPort` option during driver creation.