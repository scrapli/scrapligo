package constants

const (
	// LibScrapliPathOverrideEnv holds the key of the environment variable, that when set will force
	// the ffi loader to load libscrapli from the provided path.
	LibScrapliPathOverrideEnv = "LIBSCRAPLI_PATH"

	// LibScrapliCacheOverrideEnv holds the key for the environment variable that can be used to
	// override where we store/look for the libscrapli dynamic library file on disk.
	LibScrapliCacheOverrideEnv = "LIBSCRAPLI_CACHE_PATH"

	// XdgCacheHomeEnv is the key for env var for XDG_CAHCE_HOME -- we use this to try to see where
	// a user would want us to cache the libscrapli dynamic library file.
	XdgCacheHomeEnv = "XDG_CACHE_HOME"

	// HomeEnv is the key for the HOME env var, we use this and/or the XdgCacheHomeEnv to figure
	// out where to put the libscrapli dynamic library file.
	HomeEnv = "HOME"

	// LibScrapliDelimiter is the delimiter used to delim string lists going to/from the ffi -- this
	// is used so we dont have to pass lengths of slices back and forth over the ffi. Its a bit
	// hacky but works.
	LibScrapliDelimiter = "__libscrapli__"
)
