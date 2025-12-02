package xruntime

import (
	"runtime"
)

const MaxStackDepth = 50

func NewStackTrace(skip, maxDepth uint32) []runtime.Frame {
	var startOffset = 2 + skip

	// prepare the stack trace
	var pcs = make([]uintptr, startOffset+maxDepth)
	var recorded = runtime.Callers(int(startOffset), pcs)
	var frames = runtime.CallersFrames(pcs)

	// generate the stack trace
	var result = make([]runtime.Frame, 0, recorded)
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		result = append(result, frame)
	}
	return result
}

/*
func NewAutoStackTrace(skip, maxDepth int) []runtime.Frame {
	// compute the stack offset and depth
	if maxDepth <= 0 {
		return nil
	}
	var startOffset = 2

	// prepare the stack trace
	var pcs = make([]uintptr, startOffset+maxDepth)
	var recorded = runtime.Callers(startOffset, pcs)
	var frames = runtime.CallersFrames(pcs)
	// generate the stack trace
	frame, more := frames.Next()

	var searchStartFn = true //skip < 0
	if searchStartFn {
		for ; more; frame, more = frames.Next() {
			// search for func github.com/amie-go/adk/errors.WithStackTrace.func1 to start after that
			if frame.Function != "github.com/amie-go/adk/errors.WithStackTrace.func1" {
				continue
			}
			frame, more = frames.Next()
			// skip function github.com/amie-go/adk/errors.New
			if frame.Function == "github.com/amie-go/adk/errors.New" {
				frame, more = frames.Next()
			}
			break
		}
	}

	// generate the stack trace
	var result = make([]runtime.Frame, 0, recorded)
	for ; more; frame, more = frames.Next() {
		result = append(result, frame)
	}

	return result
}

func NewAutoStackTrace2(skip, maxDepth int) []runtime.Frame {
	var frames = NewStackTrace(skip, maxDepth)

	// filter by removing frames with function github.com/amie-go/adk/errors.WithStackTrace.func1
	// and github.com/amie-go/adk/errors.New
	var posStart = slices.IndexFunc(frames, func(frame runtime.Frame) bool {
		return frame.Function == "github.com/amie-go/adk/errors.WithStackTrace.func1"
	})
	if posStart >= 0 && (posStart+1) < len(frames) {
		// check if next frame is github.com/amie-go/adk/errors.New
		if frames[posStart+1].Function == "github.com/amie-go/adk/errors.New" {
			posStart++
		}
	}

	posStart++
	if posStart < 0 {
		posStart = 0
	}
	posStart += skip
	if posStart >= len(frames) {
		return nil
	}
	return frames[posStart:]
}
*/
