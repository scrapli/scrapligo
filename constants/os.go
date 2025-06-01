package constants

const (
	// Darwin for the goos platform darwin. For test reasons mostly.
	Darwin = "darwin"

	// PermissionsOwnerReadWrite is the permissions for owner read/write nobody else anything.
	PermissionsOwnerReadWrite = 0o600

	// PermissionsOwnerReadWriteEveryoneRead is the permissions for owner read/write, everyone
	// else read.
	PermissionsOwnerReadWriteEveryoneRead = 0o644

	// PermissionsOwnerReadWriteExecute is the permissions for owner read/write nobody else
	// anything.
	PermissionsOwnerReadWriteExecute = 0o700
)
