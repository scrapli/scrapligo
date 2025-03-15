package netconf

// Option defines a functional option for a netconf rpc.
type Option func(o any)

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
		case lockUnlockOptions:
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
