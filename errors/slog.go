package exterr

import (
	"context"
	"errors"
	"log/slog"
	"runtime"
)

func Log(err error) {
	LogCtx(context.Background(), err)
}

func LogCtx(ctx context.Context, err error) {
	LogWithAttrCtx(ctx, err, AttrsFromMetadata, AttrsFromStackTrace("stacktrace", nil))
}

// -> Can filter on some type and not other...
// -> LogCtx(ctx, err, MetadataAttrs, StackTraceAttrs("stacktrace"))
// -> LogCtx(ctx, err, GetLogAttrs)
func LogWithAttrCtx(ctx context.Context, err error, fns ...AttrsFn) {
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, err.Error(), AttrsAppender(err, fns...)...)
	}
}

type AttrsFn func(error) []slog.Attr

func AttrsAppender(err error, filterFns ...AttrsFn) (attrs []slog.Attr) {
	for _, v := range filterFns {
		if v != nil {
			attrs = append(attrs, v(err)...)
		}
	}
	return
}

// ---------------------------------------------------------

// AttrsFromMetadata returns the metadata of the error as slog attributes.
func AttrsFromMetadata(err error) (attrs []slog.Attr) {
	for _, v := range GetMetadata(err) {
		attrs = append(attrs, slog.Any(v.Key, v.Value))
	}
	return
}

// AttrsFromStackTrace returns the stack trace of the error as slog attributes.
func AttrsFromStackTrace(key string, format func([]runtime.Frame) any) AttrsFn {
	var fn = format
	if fn == nil {
		fn = stackFrameToLogSource
	}
	return func(err error) []slog.Attr {
		var stack = GetStackTrace(err)
		if len(stack) == 0 {
			return nil
		}
		return []slog.Attr{slog.Any(key, fn(stack))}
	}
}

func stackFrameToLogSource(stack []runtime.Frame) any {
	var result = make([]slog.Source, 0, len(stack))

	for _, v := range stack {
		result = append(result, slog.Source{
			Function: v.Function,
			File:     v.File,
			Line:     v.Line,
		})
	}
	return result
}

// Attrs returns the common attributes of the error as slog attributes.
func Attrs(errToBrowse error) (attrs []slog.Attr) {
	for err := errToBrowse; err != nil; err = errors.Unwrap(err) {
		if castedErr, ok := err.(interface{ GetLogAttr() slog.Attr }); ok {
			attrs = append(attrs, castedErr.GetLogAttr())
		} else if castedErr, ok := err.(interface{ GetLogAttrs() []slog.Attr }); ok {
			attrs = append(attrs, castedErr.GetLogAttrs()...)
		}
	}
	return
}

/*
func GetLogAttrs(err error, filterFns ...FilterFn) (attrs []slog.Attr) {
	for e := err; e != nil; e = errors.Unwrap(e) {
		for _, v := range filterFns {
			if v != nil {
				attrs = append(attrs, v(e)...)
			}
		}
	}
	return
}

type FilterLogsAttr interface{ GetLogAttr() slog.Attr }
type FilterLogsAttrs interface{ GetLogAttrs() []slog.Attr }

func GetLogAttrs2(errToBrowse error) (attrs []slog.Attr) {
	for err := errToBrowse; err != nil; err = errors.Unwrap(err) {
		if castedErr, ok := err.(FilterLogsAttr); ok {
			attrs = append(attrs, castedErr.GetLogAttr())
		} else if castedErr, ok := err.(FilterLogsAttrs); ok {
			attrs = append(attrs, castedErr.GetLogAttrs()...)
		}
	}
	return
}
*/

/*
func Map(errToBrowse error) (attrs []slog.Attr) {
	for err := errToBrowse; err != nil; err = errors.Unwrap(err) {
		if castedErr, ok := err.(FilterLogsAttrs); ok {
			attrs = append(attrs, castedErr.GetLogAttrs()...)
		}
	}
	return
}
*/
