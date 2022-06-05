Basics
======

This example shows the basics of creating a `network.Driver` object and sending some inputs and
commands.

Note that by default scrapligo does strict SSH key checking, in most lab scenarios you probably
want to disable this by passing the `options.WithAuthNoStrictKey()` option to the `New` driver
function.

Unlike the "generic" driver, the "network" driver has an additional requirement that you *must* 
pass a default desired privilege level as well as a map of privilege levels. In *most* cases you 
will not be instantiating a network driver directly, but instead doing so via the platforms 
factory (see the platforms example directory).
