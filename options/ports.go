package options

import "context"

// With is the interface contract to apply modification on options
type With[T any] interface {
	Apply(context.Context, *T)
}

// WithFn type is an adapter to allow the use of ordinary functions as options modifiers
// If fn is a function with the appropriate signature, WithFn(fn) is a [With] that calls fn.
type WithFn[T any] func(*T)

func (fn WithFn[T]) Apply(ctx context.Context, options *T) {
	if fn != nil {
		fn(options)
	}
}

// WithCtxFn type is an adapter to allow the use of ordinary functions as options modifiers with context
// If fn is a function with the appropriate signature, WithCtxFn(fn) is a [With] that calls fn.
type WithCtxFn[T any] func(context.Context, *T)

func (fn WithCtxFn[T]) Apply(ctx context.Context, options *T) {
	if fn != nil {
		fn(ctx, options)
	}
}
