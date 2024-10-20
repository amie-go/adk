package options

// New creates a new instance of the options
// It sets default values if the options implement the SetDefault method,
// applies the options and validates them if the options implement the Validate method
// It returns the options or an error if the options are not valid
func New[T any](opts ...With[T]) (*T, error) {
	// Create a new instance of the options
	var options = new(T)

	// Set default values
	if casted, ok := any(options).(interface{ SetDefault() }); ok && casted != nil {
		casted.SetDefault()
	}

	// Apply options
	Apply(options, opts...)

	// Validate options
	if casted, ok := any(options).(interface{ Validate() error }); ok && casted != nil {
		if err := casted.Validate(); err != nil {
			return nil, err
		}
	}

	return options, nil
}

// Apply applies the options to the instance
func Apply[T any](options *T, opts ...With[T]) {
	if options == nil {
		return
	}
	for _, v := range opts {
		if v != nil {
			v.Apply(options)
		}
	}
}
