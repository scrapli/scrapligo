package netconf

const (
	// DefaultStreamValue is the default value for "stream" field on create/establish/modify
	// subscription rpcs.
	DefaultStreamValue = "NETCONF"
)

// Option defines a functional option for a netconf rpc.
type Option func(o any)

// WithExtraNamespaces sets additional namespaces for the rpc. The value is a list of pairs of
// strings representing the prefix and URI.
func WithExtraNamespaces(e [][2]string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *rawRPCOptions:
			to.extraNamespaces = e
		}
	}
}

// WithDatastore set the datastore type for the rpc.
func WithDatastore(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getDataOptions:
			to.datastore = t
		case *lockOptions:
			to.target = t
		case *unlockOptions:
			to.target = t
		}
	}
}

// WithSourceType set the source datastore type for the rpc.
func WithSourceType(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getConfigOptions:
			to.source = t
		case *copyConfigOptions:
			to.target = t
		}
	}
}

// WithTargetType set the target datastore type for the rpc.
func WithTargetType(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case *editConfigOptions:
			to.target = t
		case *copyConfigOptions:
			to.target = t
		case *deleteConfigOptions:
			to.target = t
		case *lockOptions:
			to.target = t
		case *unlockOptions:
			to.target = t
		}
	}
}

// WithFilter apply a filter for the rpc.
func WithFilter(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getConfigOptions:
			to.filter = s
		case *getOptions:
			to.filter = s
		case *getDataOptions:
			to.filter = s
		}
	}
}

// WithFilterType apply a filter type for the rpc.
func WithFilterType(t FilterType) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getConfigOptions:
			to.filterType = t
		case *getOptions:
			to.filterType = t
		}
	}
}

// WithFilterNamespacePrefix apply namespace prefix for a filter namespace.
func WithFilterNamespacePrefix(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getConfigOptions:
			to.filterNamespacePrefix = s
		case *getOptions:
			to.filterNamespacePrefix = s
		}
	}
}

// WithFilterNamespace apply a namespace for a filter.
func WithFilterNamespace(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getConfigOptions:
			to.filterNamespace = s
		case *getOptions:
			to.filterNamespace = s
		}
	}
}

// WithDefaultsType apply a defaults type for the rpc.
func WithDefaultsType(t DefaultsType) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getConfigOptions:
			to.defaultsType = t
		case *getOptions:
			to.defaultsType = t
		}
	}
}

// WithSchemaFormat apply a schema format for the rpc.
func WithSchemaFormat(t SchemaFormat) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getSchemaOptions:
			to.format = t
		}
	}
}

// WithVersion apply a version argument for the rpc.
func WithVersion(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getSchemaOptions:
			to.version = s
		}
	}
}

// WithConfigFilter apply the config filter option for the rpc.
func WithConfigFilter(t ConfigFilter) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getDataOptions:
			to.configFilter = t
		}
	}
}

// WithMaxDepth apply the max depth option for the rpc.
func WithMaxDepth(i uint32) Option {
	return func(o any) {
		switch to := o.(type) {
		case *getDataOptions:
			to.maxDepth = i
		}
	}
}

// WithOrigin apply the with origin option for the rpc.
func WithOrigin() Option {
	return func(o any) {
		switch to := o.(type) {
		case *getDataOptions:
			to.withOrigin = true
		}
	}
}
