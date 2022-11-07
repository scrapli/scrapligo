NETCONF Periodic Subscription
======

Using NETCONF, we can subscribe to continuous updates from a YANG datastore. This is accomplished
easily using the `EstablishPeriodicSubscription` method, which requires two parameters. 
An `xpath` filter is provided to tell the target device which YANG datastore to send updates for. 
The `period` value specifies the interval in which the device should push updates.

In the example code we subscribe to `/interfaces/interface[name=\"GigabitEthernet1\"]/statistics`
which contains operational state data on interface `GigabitEthernet1` (in/out packets, errors, etc). 
The period is specified as `1000` centiseconds, or 10 seconds. Any subscription messages received
from the device can be retrieved using `GetSubscriptionMessages` and passing the subscription ID.