Sending Configs
===============

Sending configs is easy! Just use the `SendConfig`, `SendConfigs`, or `SendConfigsFromFile` 
methods. The first accepts a single string of configs and will split them on new lines, sending 
each line individually. The `SendConfigs` method accepts a slice of configs, and the 
`SendConfigsFromFile` of course accepts a file path which it loads and sends line by line to the 
device.

Note that there are some "types" of configurations that will cause some issues for 
scrapli/scrapligo -- these are "vi-like" configurations. This is most commonly seen when 
configuring banners (like motd banner). These types of config modes "break" scrapligo because 
there is no prompt painted after each line is entered, and scrapligo never sends subsequent 
lines of config until the prompt is "re-found". In order to not break scrapligo in these types 
of config sections you can send them with the `WithEager` option. This "eager" mode causes 
scrapligo to no longer wait until it sees the prompt pattern after sending each config line. 
Generally you probably don't want to use eager unless you need to as it can cause scrapligo to 
fire the configs to the device much too quickly which sometimes causes devices to... lose their 
mind for lack of a better explanation!