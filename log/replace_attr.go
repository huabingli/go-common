package log

import (
	"log/slog"
	"reflect"
	"runtime"
	"time"
)

func ReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey { // 格式化 key 为 "time" 的属性值
		if t, ok := a.Value.Any().(time.Time); ok {
			a.Value = slog.StringValue(t.Format(`2006-01-02 15:04:05.000`))
		}
	}
	if a.Key == "err" || a.Key == "error" {
		if src, ok := a.Value.Any().(error); ok {
			values := []slog.Attr{
				slog.String("msg", src.Error()),
				slog.String("type", reflect.TypeOf(src).String()),
				slog.String("stacktrace", stacktrace()),
			}
			a.Value = slog.GroupValue(values...)
		}
	}
	return a
}

func stacktrace() string {
	stackInfo := make([]byte, 1024*1024)

	if stackSize := runtime.Stack(stackInfo, false); stackSize > 0 {
		return string(stackInfo[:stackSize])
	}

	return ""
}
