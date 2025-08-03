package netconf

// DatastoreType is an enum(ish) representing the name of a datastore on a device.
type DatastoreType string

// String (stringer) method for DatastoreType.
func (t DatastoreType) String() string {
	return string(t)
}

const (
	DatastoreTypeConventional DatastoreType = "conventional"
	DatastoreTypeRunning      DatastoreType = "running"
	DatastoreTypeCandidate    DatastoreType = "candidate"
	DatastoreTypeStartup      DatastoreType = "startup"
	DatastoreTypeIntended     DatastoreType = "intended"
	DatastoreTypeDynamic      DatastoreType = "dynamic"
	DatastoreTypeOperational  DatastoreType = "operational"
)

// FilterType is an enum(ish) representing the name of a filter type.
type FilterType string

// String (stringer) method for FilterType.
func (t FilterType) String() string {
	return string(t)
}

const (
	FilterTypeSubtree FilterType = "subtree"
	FilterTypeXpath   FilterType = "xpath"
)

// DefaultsType is an enum(ish) representing the name of a filter type.
type DefaultsType string

// String (stringer) method for DefaultsType.
func (t DefaultsType) String() string {
	return string(t)
}

const (
	// DefaultsTypeUnset is the option to *not* set a defaults type -- we need to have a string
	// to pass to zig abi so we just specify unset if it would be null in zig native.
	DefaultsTypeUnset           DefaultsType = "unset"
	DefaultsTypeReportAll       DefaultsType = "report-all"
	DefaultsTypeReportAllTagged DefaultsType = "report-all-tagged"
	DefaultsTypeTrim            DefaultsType = "trim"
	DefaultsTypeExplicit        DefaultsType = "explicit"
)

// SchemaFormat is an enum(ish) representing the name of a schema format.
type SchemaFormat string

// String (stringer) method for SchemaFormat.
func (t SchemaFormat) String() string {
	return string(t)
}

const (
	SchemaFormatXsd  SchemaFormat = "xsd"
	SchemaFormatYang SchemaFormat = "yang"
	SchemaFormatYin  SchemaFormat = "yin"
	SchemaFormatRng  SchemaFormat = "rng"
	SchemaFormatRnc  SchemaFormat = "rnc"
)

// ConfigFilter is an enum(ish) representing the valid config-filter options.
type ConfigFilter string

// String (stringer) method for ConfigFilter.
func (t ConfigFilter) String() string {
	return string(t)
}

const (
	ConfigFilterTrue  = "true"
	ConfigFilterFalse = "false"
	ConfigFilterUnset = "unset"
)

// DefaultOperation is an enum(ish) representing the name of a default operation field.
type DefaultOperation string

// String (stringer) method for DefaultOperation.
func (t DefaultOperation) String() string {
	return string(t)
}

const (
	DefaultOperationMerge   DefaultOperation = "merge"
	DefaultOperationReplace DefaultOperation = "replace"
	DefaultOperationNone    DefaultOperation = "none"
)

// TestOption is an enum(ish) representing the name of a default operation field.
type TestOption string

// String (stringer) method for TestOption.
func (t TestOption) String() string {
	return string(t)
}

const (
	TestOptionTestThenSet TestOption = "test-then-set"
	TestOptionSet         TestOption = "set"
)

// ErrorOption is an enum(ish) representing the name of a default operation field.
type ErrorOption string

// String (stringer) method for DefaultOperation.
func (t ErrorOption) String() string {
	return string(t)
}

const (
	ErrorOptionStopOnError     ErrorOption = "stop-on-error"
	ErrorOptionContinueOnError ErrorOption = "continue-on-error"
	ErrorOptionRollbackOnError ErrorOption = "rollback-on-error"
)
