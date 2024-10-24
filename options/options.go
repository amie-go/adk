package options

// New creates a new instance of the configuration, and applies the options to it.
func New[T any](opts ...With[T]) *T {
	var config = new(T)
	Apply(config, opts...)
	return config
}

// Apply applies the options to the configuration.
func Apply[T any](config *T, opts ...With[T]) {
	if config == nil {
		return
	}
	for _, v := range opts {
		if v != nil {
			v.Apply(config)
		}
	}
}
