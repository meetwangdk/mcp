package dataframe

// Option is an interface that represents a data frame convert to table option
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

type options struct {
	formatTime    bool
	includeFields []string
	excludeFields []string
	renameFields  map[string]string
	sortByField   string
	sortDesc      bool
}

// WithFormatTime is an option that format timestamp to time string
func WithFormatTime() Option {
	return optionFunc(func(o *options) {
		o.formatTime = true
	})
}

// WithIncludeFields is an option that include fields
func WithIncludeFields(fields ...string) Option {
	return optionFunc(func(o *options) {
		o.includeFields = fields
	})
}

// WithExcludeFields is an option that exclude fields
func WithExcludeFields(fields ...string) Option {
	return optionFunc(func(o *options) {
		o.excludeFields = fields
	})
}

// WithRenameFields is an option that rename fields
func WithRenameFields(fields map[string]string) Option {
	return optionFunc(func(o *options) {
		o.renameFields = fields
	})
}

// WithSortByField is an option that sort by field
func WithSortByField(field string, desc bool) Option {
	return optionFunc(func(o *options) {
		o.sortByField = field
		o.sortDesc = desc
	})
}
