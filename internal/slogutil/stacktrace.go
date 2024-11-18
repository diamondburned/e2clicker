package slogutil

import (
	"fmt"
	"log/slog"
	"runtime"
)

// StackTrace returns a slog attribute that contains the stack trace of the
// current goroutine.
func StackTrace(skip int) slog.Attr {
	pc := make([]uintptr, 20)
	pc = pc[:runtime.Callers(skip+1, pc)]

	frames := runtime.CallersFrames(pc)
	attr := make([]slog.Attr, 0, len(pc))

	for i := 0; ; i-- {
		f, more := frames.Next()

		group := make([]slog.Attr, 1, 2)
		group[0] = slog.String("location", fmt.Sprintf("%s:%d", f.File, f.Line))
		if f.Function != "" {
			group = append(group, slog.String("function", f.Function))
		}

		attr = append(attr, slog.Any(
			fmt.Sprintf("[%d]", i),
			slog.GroupValue(group...),
		))

		if !more {
			break
		}
	}

	return slog.Any("stack", slog.GroupValue(attr...))
}
