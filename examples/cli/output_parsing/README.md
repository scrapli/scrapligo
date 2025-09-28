# Output Parsing

This example shows using [textfsm](https://github.com/google/textfsm) (optional with
[ntc-templates](https://github.com/networktocode/ntc-templates)) to turn unstructured text output
from a device into a structured object. Scrapli in this case is not really doing any parsing for you
or anything, but rather providing some convenience functions to more simply parse the data.

*Note:* unlike scrappli (py) there is no auto template lookup with ntc-templates, you must pass the
template path/url.
