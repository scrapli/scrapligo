Channel Logging
===============

It is possible to log all reads/writes to the underlying transport via the ChannelLog parameter. 
Pass an `io.Writer` object to your driver creation, and you can read the bytes back out of that 
object when you are done with the connection. 