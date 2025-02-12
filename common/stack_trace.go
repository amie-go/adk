package common

import (
	"runtime"
)

const MaxStackDepth = 50

func NewStackTrace(skip, maxDepth int) []runtime.Frame {
	// compute the stack offset and depth
	if maxDepth <= 0 {
		return nil
	}
	var startOffset = 2
	if skip > 0 {
		startOffset += skip
	}

	// prepare the stack trace
	var pcs = make([]uintptr, startOffset+maxDepth)
	var recorded = runtime.Callers(startOffset, pcs)
	var frames = runtime.CallersFrames(pcs)

	// generate the stack trace
	var result = make([]runtime.Frame, 0, recorded)
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		result = append(result, frame)
	}
	return result
}
