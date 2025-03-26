package log

import (
	"context"
	"log/slog"
)

type Handler struct {
	handler slog.Handler
}

func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{
		handler: handler,
	}
}
func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.handler.WithAttrs(attrs)}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	if attrs, ok := ctx.Value("requestID").(string); ok {
		record.AddAttrs(slog.String("request_id", attrs))
	}
	return h.handler.Handle(ctx, record)
}
