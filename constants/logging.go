package constants

const (
	// ScrapligoDebug is the env var when if set to anything enables debug logging -- this is
	// almost entirely for ffi related bits as all other logging would be handled by providing
	// a logger callback to libscrapli.
	ScrapligoDebug = "SCRAPLIGO_DEBUG"
)
