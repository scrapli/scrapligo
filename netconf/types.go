package netconf

// DatastoreType is an enum(ish) representing the name of a datastore on a device.
type DatastoreType uint8

const (
	DatastoreTypeConventional DatastoreType = iota
	DatastoreTypeRunning
	DatastoreTypeCandidate
	DatastoreTypeStartup
	DatastoreTypeIntended
	DatastoreTypeDynamic
	DatastoreTypeOperational
)

// FilterType is an enum(ish) representing the name of a filter type.
type FilterType uint8

const (
	FilterTypeSubtree FilterType = iota
	FilterTypeXpath
)

// DefaultsType is an enum(ish) representing the name of a filter type.
type DefaultsType uint8

const (
	DefaultsTypeReportAll DefaultsType = iota
	DefaultsTypeReportAllTagged
	DefaultsTypeTrim
	DefaultsTypeExplicit
)

// SchemaFormat is an enum(ish) representing the name of a schema format.
type SchemaFormat uint8

const (
	SchemaFormatXsd SchemaFormat = iota
	SchemaFormatYang
	SchemaFormatYin
	SchemaFormatRng
	SchemaFormatRnc
)

// ConfigFilter is an enum(ish) representing the valid config-filter options.
type ConfigFilter uint8

const (
	ConfigFilterTrue ConfigFilter = iota
	ConfigFilterFalse
)

// DefaultOperation is an enum(ish) representing the name of a default operation field.
type DefaultOperation uint8

const (
	DefaultOperationMerge DefaultOperation = iota
	DefaultOperationReplace
	DefaultOperationNone
)

// TestOption is an enum(ish) representing the name of a default operation field.
type TestOption uint8

const (
	TestOptionTestThenSet TestOption = iota
	TestOptionSet
)

// ErrorOption is an enum(ish) representing the name of a default operation field.
type ErrorOption uint8

const (
	ErrorOptionStopOnError ErrorOption = iota
	ErrorOptionContinueOnError
	ErrorOptionRollbackOnError
)
