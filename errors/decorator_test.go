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

	var fn = func() {
		err := exterr.New(ErrNotFound, WithStackTraceCustom, exterr.WithMetadata(exterr.KV{"foo", "bar"}))
		err = exterr.New(err, exterr.WithMetadata(exterr.KV{"foo2", "bar2"}))
		fmt.Println(err)
		slog.Info("slog info msg", "err", err)
		exterr.Log(err)

		exterr.LogWithAttrCtx(context.Background(), err, exterr.GetLogAttrs)

		info, ok := debug.ReadBuildInfo()
		fmt.Println(ok)
		fmt.Println(info)
	}

	fn()
}
