package constants

// Version is the version of scrapligo. Set during release via ci.
var Version = "2.0.0-alpha.1"

// LibScrapliVersion is the version of libscrapli scrapligo was built with. Should be set prior to
// a release via build/update_all.sh to the version of libscrapli bundled in assets.
var LibScrapliVersion = "0.0.1-beta.8"

// ScrapliDefinitionsVersion is the version of scrapli definitions embedded in assets in this build.
// This should be set prior to a release via build/update_all.sh.
var ScrapliDefinitionsVersion = "1a3d771"
