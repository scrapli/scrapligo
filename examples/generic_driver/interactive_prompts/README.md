Interactive Prompts
===================

Sometimes you may need to "interact" with devices -- or put another way, you may need to respond 
to a prompt from the device. The `SendInteractive` method is used to handle these things. This 
method accepts a slice of `channel.SendInteractiveEvent` which define the input, the expected 
output, and whether the device will "hide" the input (as is the case with password prompts). 
This is a fairly simple/dumb method, but works well enough. If you need a more 
elaborate/advanced way to handle a multitude of prompts/outputs you may want to check out the 
`SendWithCallbacks` method instead.