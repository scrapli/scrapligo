Privilege Levels
================

When using a "network" driver, scrapligo has the concept of privilege levels. Each privilege 
level corresponds to a "mode" or a prompt flavor on a device. Using Cisco IOSXE as an example, 
there are at least three privilege levels -- "exec" (> prompt), "privilege-exec" (# prompt), and 
"configuration" (config# prompt). scrapligo keeps track of the current privilege level, and it 
understands how to traverse the privilege levels. When you send a "command" scrapligo will 
auto-acquire the "default desired privilege level". This default privilege leve is more or less 
the "normal" privilege level you would operate in when connecting to a device. When sending a 
"config" scrapligo will auto acquire the "configuration" privilege level.

While you probably won't need to do so often, you can explicitly acquire a privilege level with 
the "AcquirePriv" method. If you need to send configs in a non-standard configuration level, 
such as "configuration-exclusive" you can pass the `WithPrivilegeLevel` option to any config 
methods. Note that the privilege levels must be defined/exist in the platform definition you are 
building your connection from!