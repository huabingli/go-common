package log

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

// InitSlog 初始化日志器，支持开发模式和生产模式，日志级别和源信息选择
func InitSlog(level slog.Level, addSource, dev bool, requestIDKey string) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource:   addSource,
		Level:       level,
		ReplaceAttr: ReplaceAttr,
	}

	// var logger *slog.Logger
	var handler slog.Handler
	if dev {
		// 开发环境下使用漂亮的控制台输出
		handler = devslog.NewHandler(
			os.Stderr, &devslog.Options{
				HandlerOptions:     &opts,
				MaxSlicePrintSize:  50,
				SortKeys:           true,
				NewLineAfterLog:    true,
				StringIndentation:  false,
				MaxErrorStackTrace: 0,
				StringerFormatter:  true,
			},
		)
	} else {
		// 生产环境下使用 JSON 格式输出
		handler = slog.NewJSONHandler(os.Stderr, &opts)
	}
	logger := slog.New(NewHandler(handler, requestIDKey))
	slog.SetDefault(logger)
	return logger
}
