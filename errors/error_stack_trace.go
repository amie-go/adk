package exterr

import (
	"context"
	"errors"
	"runtime"

	"github.com/amie-go/adk/common"
	"github.com/amie-go/adk/options"
)

// WithStackTrace returns a decorator that adds a stack trace to the error.
func WithStackTrace(opts ...options.With[stackTraceConfig]) func(error) error {
	var config = options.NewWithDefaults(context.Background(), setDefaults, opts...)
	return func(err error) error {
		return &stackTraceError{
			inner:       err,
			stackFrames: config.generateFn(2+config.skip, config.maxDepth),
		}
	}
}

// GetStackTrace returns the stack trace associated to the error, if exists.
// TOTHINK: StackTraceFromErr, StackTraceFrom, StackFramesFrom
func GetStackTrace(err error) []runtime.Frame {
	for v := err; err != nil; {
		if e, ok := v.(*stackTraceError); ok {
			return e.stackFrames
		}
		v = errors.Unwrap(v)
	}
	return nil
}

// ---------------------------------------------------------
// Implementation

type stackTraceError struct {
	inner       error
	stackFrames []runtime.Frame
}

func (e stackTraceError) Error() string {
	return e.inner.Error()
}

func (e stackTraceError) Unwrap() error {
	return e.inner
}

// ---------------------------------------------------------
// WithOptions

type stackTraceConfig struct {
	skip       int
	maxDepth   int
	generateFn func(int, int) []runtime.Frame
}

func setDefaults(dst *stackTraceConfig) {
	dst.skip = 0
	dst.maxDepth = common.MaxStackDepth
	dst.generateFn = common.NewStackTrace
}

func WithSkip(value int) options.WithFn[stackTraceConfig] {
	return func(dst *stackTraceConfig) {
		if value >= 0 {
			dst.skip = value
		}
	}
}

func WithMaxDepth(value int) options.WithFn[stackTraceConfig] {
	return func(dst *stackTraceConfig) {
		if value >= 0 {
			dst.maxDepth = value
		}
	}
}

func WithFramesGenerator(fn func(int, int) []runtime.Frame) options.WithFn[stackTraceConfig] {
	return func(dst *stackTraceConfig) {
		if fn != nil {
			dst.generateFn = fn
		}
	}
}
