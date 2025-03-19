package netconf

// Option defines a functional option for a netconf rpc.
type Option func(o any)

// WithDatastore set the datastore type for the rpc.
func WithDatastore(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getDataOptions:
			to.datastore = t
		}
	}
}

// WithSourceType set the source datastore type for the rpc.
func WithSourceType(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOptions:
			to.source = t
		case copyConfigOptions:
			to.target = t
		}
	}
}

// WithTargetType set the target datastore type for the rpc.
func WithTargetType(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case editConfigOptions:
			to.target = t
		case copyConfigOptions:
			to.target = t
		case deleteConfigOptions:
			to.target = t
		case lockOptions:
			to.target = t
		case unlockOptions:
			to.target = t
		}
	}
}

// WithFilter apply a filter for the rpc.
func WithFilter(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOptions:
			to.filter = s
		case getOptions:
			to.filter = s
		}
	}
}

// WithFilterType apply a filter type for the rpc.
func WithFilterType(t FilterType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOptions:
			to.filterType = t
		case getOptions:
			to.filterType = t
		}
	}
}

// WithFilterNamespacePrefix apply namespace prefix for a filter namespace.
func WithFilterNamespacePrefix(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOptions:
			to.filterNamespacePrefix = s
		case getOptions:
			to.filterNamespacePrefix = s
		}
	}
}

// WithFilterNamespace apply a namespace for a filter.
func WithFilterNamespace(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOptions:
			to.filterNamespace = s
		case getOptions:
			to.filterNamespace = s
		}
	}
}

// WithDefaultsType apply a defaults type for the rpc.
func WithDefaultsType(t DefaultsType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOptions:
			to.defaultsType = t
		case getOptions:
			to.defaultsType = t
		}
	}
}

// WithSchemaFormat apply a schema format for the rpc.
func WithSchemaFormat(t SchemaFormat) Option {
	return func(o any) {
		switch to := o.(type) {
		case getSchemaOptions:
			to.format = t
		}
	}
}

// WithVersion apply a version argument for the rpc.
func WithVersion(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getSchemaOptions:
			to.version = s
		}
	}
}

// WithStartTime apply a start time argument for the rpc.
func WithStartTime(i uint64) Option {
	return func(o any) {
		switch to := o.(type) {
		case createSubscriptionOptions:
			to.startTime = i
		}
	}
}

// WithStopTime apply a stop time argument for the rpc.
func WithStopTime(i uint64) Option {
	return func(o any) {
		switch to := o.(type) {
		case createSubscriptionOptions:
			to.stopTime = i
		}
	}
}

// WithPeriod apply a period argument for the rpc.
func WithPeriod(i uint64) Option {
	return func(o any) {
		switch to := o.(type) {
		case establishSubscriptionOptions:
			to.period = i
		}
	}
}

// WithDSCP apply a dscp argument for the rpc.
func WithDSCP(i uint8) Option {
	return func(o any) {
		switch to := o.(type) {
		case establishSubscriptionOptions:
			to.dscp = i
		}
	}
}

// WithWeighting apply a weighting argument for the rpc.
func WithWeighting(i uint8) Option {
	return func(o any) {
		switch to := o.(type) {
		case establishSubscriptionOptions:
			to.dscp = i
		}
	}
}

// WithDependency apply a dependency argument for the rpc.
func WithDependency(i uint32) Option {
	return func(o any) {
		switch to := o.(type) {
		case establishSubscriptionOptions:
			to.dependency = i
		}
	}
}

// WithEncoding apply an encoding argument for the rpc.
func WithEncoding(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case establishSubscriptionOptions:
			to.encoding = s
		}
	}
}

// WithConfigFilter apply the config filter option for the rpc.
func WithConfigFilter() Option {
	return func(o any) {
		switch to := o.(type) {
		case getDataOptions:
			to.configFilter = true
		}
	}
}

// WithMaxDepth apply the max depth option for the rpc.
func WithMaxDepth(i uint32) Option {
	return func(o any) {
		switch to := o.(type) {
		case getDataOptions:
			to.maxDepth = i
		}
	}
}

// WithOrigin apply the with origin option for the rpc.
func WithOrigin() Option {
	return func(o any) {
		switch to := o.(type) {
		case getDataOptions:
			to.withOrigin = true
		}
	}
}
