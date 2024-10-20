package options

// With is the interface contract to apply modification on options
type With[T any] interface {
	Apply(*T)
}

// WithFn type is an adapter to allow the use of ordinary functions as options modifiers
// If fn is a function with the appropriate signature, WithFn(fn) is a [With] that calls fn.
type WithFn[T any] func(*T)

func (fn WithFn[T]) Apply(options *T) {
	if fn != nil {
		fn(options)
	}
}
