package log

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

func InitSlog(level slog.Level, addSource, dev bool) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource:   addSource,
		Level:       level,
		ReplaceAttr: ReplaceAttr,
	}

	// var logger *slog.Logger
	var handler slog.Handler
	if dev {
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
		handler = slog.NewJSONHandler(os.Stderr, &opts)
	}
	logger := slog.New(NewHandler(handler))
	slog.SetDefault(logger)
	return logger
}
