Default Logging
===============

The `options.WithDefaultLogger` function applies a simple logging instance with a single log 
emitter of `log.Print` at the "info" level. If you wanted to do fancier things with logging you 
can create your own `logging.Instance` at whatever log level and with whatever log emitter 
functions you want rather than using this default logger.