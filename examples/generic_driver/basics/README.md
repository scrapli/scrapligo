Basics
======

This example shows the basics of creating a `generic.Driver` object and sending some inputs and 
commands. 

Note that by default scrapligo does strict SSH key checking, in most lab scenarios you probably 
want to disable this by passing the `options.WithAuthNoStrictKey()` option to the `New` driver 
function.
