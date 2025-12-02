package exterr_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"testing"

	exterr "github.com/amie-go/adk/errors"
	"github.com/stretchr/testify/assert"
)

// Article:
// - https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
// - https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
// - https://www.dolthub.com/blog/2023-11-10-stack-traces-in-go/

// packages:
// - Multi Error https://github.com/hashicorp/go-multierror?tab=readme-ov-file
// - https://github.com/pkg/errors/blob/master/errors.go#L145
// - https://github.com/go-errors/errors
// - https://github.com/ztrue/tracerr
// - https://github.com/src-d/go-errors

var (
	// Some Cause/Sentinel errors
	ErrNotFound         = errors.New("not found")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNotAuthorized    = errors.New("not authorized")
)

// ---------------------------------------------------------

// var WithStackTraceCustom = WithStackTrace(WithSlogKey("toto"))
var WithStackTraceCustom = exterr.WithStackTrace()

func TestDecorate(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	var ctx = context.Background()

	var fn = func() {
		err := exterr.New(ErrNotFound, WithStackTraceCustom, exterr.WithMetadata(exterr.KV{"foo", "bar"}))
		err = exterr.New(err, exterr.WithMetadata(exterr.KV{"foo2", "bar2"}))
		fmt.Println(err)
		slog.Info("slog info msg", "err", err)
		exterr.Log(err)

		exterr.LogWithAttrCtx(ctx, err, exterr.Attrs)

		info, ok := debug.ReadBuildInfo()
		fmt.Println(ok)
		fmt.Println(info)
	}

	fn()

	var fn2 = func() {
		err := exterr.WithMetadata(exterr.KV{"foo", "bar"})(WithStackTraceCustom(ErrNotAuthenticated))
		slog.LogAttrs(ctx, slog.LevelError, err.Error(),
			exterr.AttrsAppender(err,
				exterr.AttrsFromMetadata,
				exterr.AttrsFromStackTrace("stacktrace", nil),
			)...,
		)
	}

	fn2()
}

func TestErrNotModified(t *testing.T) {
	// Create an error
	var errNotFound = errors.New("not found")
	// Decorate that error
	var err = exterr.New(errNotFound, exterr.WithStackTrace())
	if assert.NotNil(t, err) {
		// Check that the error is decorated
		assert.ErrorIs(t, err, errNotFound)
		// Check that input pointer error is not modified
		assert.NotEqual(t, err, errNotFound)
	}

	// IsDecorated(err)
	// GetDecorations(err)
}

func TestWithStackTrace(t *testing.T) {
	t.Run("WithStackTrace with New", func(t *testing.T) {
		var err = exterr.New(ErrNotFound, exterr.WithStackTrace())
		var frames = exterr.GetStackTrace(err)
		if assert.NotEmpty(t, frames) {
			assert.Equal(t, "github.com/amie-go/adk/errors_test.TestWithStackTrace.func1", frames[0].Function)
		}
	})

	t.Run("WithStackTrace with wrap error", func(t *testing.T) {
		var err = exterr.WithStackTrace()(ErrNotFound)
		var frames = exterr.GetStackTrace(err)
		exterr.GetStackTrace(err)
		if assert.NotEmpty(t, frames) {
			assert.Equal(t, "github.com/amie-go/adk/errors_test.TestWithStackTrace.func2", frames[0].Function)
		}
	})

	t.Run("WithStackTrace custom with wrap error", func(t *testing.T) {
		var WithStackTraceCustom = exterr.WithStackTrace()

		var err = WithStackTraceCustom(ErrNotFound)
		var frames = exterr.GetStackTrace(err)
		if assert.NotEmpty(t, frames) {
			assert.Equal(t, "github.com/amie-go/adk/errors_test.TestWithStackTrace.func3", frames[0].Function)
		}
	})
}
