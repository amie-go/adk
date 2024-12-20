package options

import "context"

// NewWithDefaults creates a new instance of the configuration, applies the default values to it,
// and applies the options to it.
func NewWithDefaults[T any](ctx context.Context, setDefaultFn func(*T), opts ...With[T]) *T {
	var config = new(T)
	Apply(ctx, config, SetDefaults(setDefaultFn))
	Apply(ctx, config, opts...)
	return config
}

// SetDefaults calls the SetDefaults method of the configuration if it exists.
func SetDefaults[T any](setDefaultFn func(*T)) WithFn[T] {
	return func(target *T) {
		// If SetDefault function provided not nil, call it on the target and return
		if setDefaultFn != nil {
			setDefaultFn(target)
			return
		}

		// Apply the default SetDefault method of the target if it exists
		if casted, ok := any(target).(interface{ SetDefaults() }); ok && casted != nil {
			casted.SetDefaults()
		}
	}
}
