package constants

// Version is the version of scrapligo. Set with build flags, so leave at 0.0.0.
var Version = "0.0.0" //nolint: gochecknoglobals

// LibScrapliVersion is the version of libscrapli scrapligo was built with. Set with build flags,
// so leave at 0.0.0.
var LibScrapliVersion = "0.0.1" //nolint: gochecknoglobals

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
)

const (
	// PermissionsOwnerReadWrite is the permissions for owner read/write nobody else anything.
	PermissionsOwnerReadWrite = 0o600

	// PermissionsOwnerReadWriteEveryoneRead is the permissions for owner read/write, everyone
	// else read.
	PermissionsOwnerReadWriteEveryoneRead = 0o644

	// PermissionsOwnerReadWriteExecute is the permissions for owner read/write nobody else
	// anything.
	PermissionsOwnerReadWriteExecute = 0o700
)
