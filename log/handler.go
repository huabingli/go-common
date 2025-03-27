package log

import (
	"context"
	"log/slog"
)

type Handler struct {
	handler      slog.Handler
	requestIDKey string // 这里可以是 string 或其他类型
}

func NewHandler(handler slog.Handler, requestIDKey string) slog.Handler {
	return Handler{
		handler:      handler,
		requestIDKey: requestIDKey,
	}
}
func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{
		handler:      h.handler.WithAttrs(attrs),
		requestIDKey: h.requestIDKey,
	}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return Handler{
		handler:      h.handler.WithGroup(name),
		requestIDKey: h.requestIDKey,
	}
}
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if requestID, ok := ctx.Value(h.requestIDKey).(string); ok {
		record.AddAttrs(slog.String("request_id", requestID))
	}
	return h.handler.Handle(ctx, record)
}
