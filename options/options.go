package options

import "context"

// New creates a new instance of the configuration, and applies the options to it.
func New[T any](opts ...With[T]) *T {
	return NewCtx(context.Background(), opts...)
}

// NewCtx creates a new instance of the configuration with a context, and applies the options to it.
func NewCtx[T any](ctx context.Context, opts ...With[T]) *T {
	var config = new(T)
	ApplyCtx(ctx, config, opts...)
	return config
}

// Apply applies the options to the configuration.
func Apply[T any](config *T, opts ...With[T]) {
	ApplyCtx(context.Background(), config, opts...)
}

// ApplyCtx applies the options to the configuration with a context.
func ApplyCtx[T any](ctx context.Context, config *T, opts ...With[T]) {
	if config == nil {
		return
	}
	for _, v := range opts {
		if v != nil {
			v.Apply(ctx, config)
		}
	}
}
