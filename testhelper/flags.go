package testhelper

import (
	"flag"
)

// Record is the flag indicating if we should record the fixture for the unit test data from the
// target (eos) test device.
var Record = flag.Bool("record", false, "record") //nolint: gochecknoglobals

// Update is the flag indicating if golden files should be updated.
var Update = flag.Bool("update", false, "update the golden files") //nolint: gochecknoglobals

// Flags handles parsing test flags.
func Flags() {
	flag.Parse()
}
