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
