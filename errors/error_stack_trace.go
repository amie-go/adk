package exterr

import (
	"context"
	"errors"
	"runtime"

	"github.com/amie-go/adk/options"
	"github.com/amie-go/adk/xruntime"
)

// WithStackTrace returns a decorator that adds a stack trace to the error.
func WithStackTraceOld(opts ...options.With[stackTraceConfig]) func(error) error {
	var config = options.NewWithDefaults(context.Background(), setDefaults, opts...)
	return func(err error) error {
		return &stackTraceError{
			inner:       err,
			stackFrames: config.generateFn(config.skip, config.maxDepth),
		}
	}
}

func WithStackTrace(opts ...options.With[stackTraceConfig]) func(error) error {
	var config = options.NewWithDefaults(context.Background(), setDefaults, opts...)
	return func(err error) error {
		// add 1 to skip the current function
		var frames = config.generateFn(1+config.skip, config.maxDepth)
		// check if first frame is the current function
		if len(frames) > 0 && frames[0].Function == "github.com/amie-go/adk/errors.WithStackTrace.func1" {
			frames = frames[1:]
		}
		// remove decorator new function it is next function
		if len(frames) > 0 && frames[0].Function == "github.com/amie-go/adk/errors.New" {
			frames = frames[1:]
		}

		return &stackTraceError{inner: err, stackFrames: frames}
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

type StackGenerator func(uint32, uint32) []runtime.Frame

type stackTraceConfig struct {
	skip       uint32
	maxDepth   uint32
	generateFn StackGenerator
}

func setDefaults(dst *stackTraceConfig) {
	dst.skip = 0
	dst.maxDepth = xruntime.MaxStackDepth
	dst.generateFn = xruntime.NewStackTrace
}

func WithSkip(value uint32) options.WithFn[stackTraceConfig] {
	return func(dst *stackTraceConfig) { dst.skip = value }
}

func WithMaxDepth(value uint32) options.WithFn[stackTraceConfig] {
	return func(dst *stackTraceConfig) { dst.maxDepth = value }
}

func WithFramesGenerator(fn StackGenerator) options.WithFn[stackTraceConfig] {
	return func(dst *stackTraceConfig) {
		if fn != nil {
			dst.generateFn = fn
		}
	}
}
