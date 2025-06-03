package log

import (
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func ReplaceAttr(errorStack bool) func([]string, slog.Attr) slog.Attr {
	return func(_ []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))
			}
			return a

		case "err", "error":
			if !errorStack {
				return a
			}
			if src, ok := a.Value.Any().(error); ok && src != nil {
				a.Value = slog.GroupValue(
					slog.String("msg", src.Error()),
					slog.String("type", reflect.TypeOf(src).String()),
					slog.String("stack", stacktrace()), // 建议统一叫 "stack"
				)
			}
			return a

		default:
			return a
		}
	}

}

func stacktrace() string {
	buf := make([]byte, 64*1024)
	n := runtime.Stack(buf, false)
	trace := string(buf[:n])

	// 去掉 "goroutine xxx" 那一行
	if idx := strings.Index(trace, "\n"); idx != -1 {
		return trace[idx+1:]
	}
	return trace
}
