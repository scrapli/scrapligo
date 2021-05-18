package netconf

// Default values for netconf filter/default options.
const (
	DefaultNetconfOptionsFilter      = ""
	DefaultNetconfOptionsFilterType  = "subtree"
	DefaultNetconfOptionsDefaultType = ""
)

// Options struct representing options for netconf operations.
type Options struct {
	Filter      string
	FilterType  string
	DefaultType string
}

// Option netconf operation options.
type Option func(*Options)

// WithNetconfFilter add filter to a netconf operation.
func WithNetconfFilter(filter string) Option {
	return func(o *Options) {
		o.Filter = filter
	}
}

// WithNetconfFilterType add filter-type to a netconf operation.
func WithNetconfFilterType(filterType string) Option {
	return func(o *Options) {
		o.FilterType = filterType
	}
}

// WithNetconfDefaultType add default type to a netconf operation.
func WithNetconfDefaultType(defaultType string) Option {
	return func(o *Options) {
		o.DefaultType = defaultType
	}
}
