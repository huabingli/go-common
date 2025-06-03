package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-cz/devslog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitSlog 初始化日志器，支持开发模式和生产模式，日志级别和源信息选择
func InitSlog(level slog.Level, addSource, dev bool, requestIDKey string) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource:   addSource,
		Level:       level,
		ReplaceAttr: ReplaceAttr(dev),
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

type LoggerConfig struct {
	Level        slog.Level
	AddSource    bool
	OutputType   string // "json", "text", "dev"
	LogPath      string // 如果非空，则写入文件（使用 lumberjack 管理）
	MaxSizeMB    int    // 日志文件最大大小 (MB)
	MaxBackups   int    // 最多保留的旧文件数量
	MaxAgeDays   int    // 文件最大保留天数
	Compress     bool   // 是否压缩旧日志
	Console      bool   // 是否输出到控制台
	RequestIDKey string
	ErrorStack   bool // 是否输出错误的stacktrace
}

func NewLogger(cfg LoggerConfig) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource:   cfg.AddSource,
		Level:       cfg.Level,
		ReplaceAttr: ReplaceAttr(cfg.ErrorStack),
	}
	var writers []io.Writer
	// 控制台输出
	if cfg.Console {
		writers = append(writers, os.Stderr)
	}
	// 默认输出是 stderr
	// var output io.Writer = os.Stderr
	// 如果指定了路径，使用 lumberjack 管理输出文件
	if cfg.LogPath != "" {
		dir := filepath.Dir(cfg.LogPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Sprintf("无法创建日志目录 %s: %v", dir, err))
		}
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.LogPath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
		writers = append(writers, fileWriter)
	}

	if len(writers) == 0 {
		// 如果没有指定输出，默认使用 stderr
		writers = append(writers, os.Stderr)
	}
	combinedWriter := io.MultiWriter(writers...)

	var handler slog.Handler

	switch cfg.OutputType {
	case "json":
		handler = slog.NewJSONHandler(combinedWriter, &opts)
	case "text":
		handler = slog.NewTextHandler(combinedWriter, &opts)
	case "dev":
		handler = devslog.NewHandler(
			combinedWriter, &devslog.Options{
				HandlerOptions:     &opts,
				MaxSlicePrintSize:  50,
				SortKeys:           true,
				NewLineAfterLog:    true,
				StringIndentation:  false,
				MaxErrorStackTrace: 0,
				StringerFormatter:  true,
			},
		)
	default:
		handler = slog.NewJSONHandler(combinedWriter, &opts)
	}
	logger := slog.New(NewHandler(handler, cfg.RequestIDKey))
	slog.SetDefault(logger)
	return logger
}
