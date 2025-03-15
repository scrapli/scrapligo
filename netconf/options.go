package netconf

// Option defines a functional option for a netconf rpc.
type Option func(o any)

// WithSourceType set the source datastore type for the rpc.
func WithSourceType(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOption:
			to.source = t
		}
	}
}

// WithTargetType set the target datastore type for the rpc.
func WithTargetType(t DatastoreType) Option {
	return func(o any) {
		switch to := o.(type) {
		case editConfigOption:
			to.target = t
		}
	}
}

// WithFilter apply a filter for the rpc.
func WithFilter(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOption:
			to.filter = s
		}
	}
}

// WithFilterType apply a filter type for the rpc.
func WithFilterType(t FilterType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOption:
			to.filterType = t
		}
	}
}

// WithFilterNamespacePrefix apply namespace prefix for a filter namespace.
func WithFilterNamespacePrefix(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOption:
			to.filterNamespacePrefix = s
		}
	}
}

// WithFilterNamespace apply a namespace for a filter.
func WithFilterNamespace(s string) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOption:
			to.filterNamespace = s
		}
	}
}

// WithDefaultsType apply a defaults type for the rpc.
func WithDefaultsType(t DefaultsType) Option {
	return func(o any) {
		switch to := o.(type) {
		case getConfigOption:
			to.defaultsType = t
		}
	}
}
