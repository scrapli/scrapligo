Custom Logging
===============

If you want to use a custom logging setup with one or more outputs you can create a `logging.
Instance` and pass it to the driver via the `WithLogger` option. The logging instance can accept 
as many "loggers" as you want to give it -- these "loggers" can be any function that accepts a 
variadic of interface. Check out the example to see this in action.