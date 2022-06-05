On Open/Close
=============

Sometimes you may have some tasks that should be done upon opening a connection or prior to 
closing a connection. The obvious example of this is disabling paging (ex: "term len 0"), but 
there could be other tasks such as responding to some EULA/prompt, entering a special config 
mode, disabling console logging, or gracefully tearing down a user session (this last one being 
an "on close" type of thing of course). The `OnOpen` and `OnClose` options give you the ability 
to pass a function that accepts a single argument of the driver you are creating, thereby 
giving you to access the "send" methods allowing you to send "term len 0" or whatever other 
inputs you need.

You could of course always accomplish these types of "on open" and "on close" activities 
"normally" by simply sending commands/inputs after driver creation, but the goal of these 
arguments are to give a place to handle "boilerplate" type tasks without having to muddle up 
your programs with these boring details.