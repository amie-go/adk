package options

import "context"

// New creates a new instance of the configuration with a context, and applies the options to it.
func New[T any](ctx context.Context, opts ...With[T]) *T {
	var config = new(T)
	Apply(ctx, config, opts...)
	return config
}

// Apply applies the options to the configuration with a context.
func Apply[T any](ctx context.Context, config *T, opts ...With[T]) {
	if config == nil {
		return
	}
	for _, v := range opts {
		if v != nil {
			v.Apply(ctx, config)
		}
	}
}
