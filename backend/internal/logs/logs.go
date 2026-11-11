package logs

import (
	"context"
	"log/slog"
	"os"
)

type attrsHandler struct {
	slog.Handler
}

func AttrsHandler(handler slog.Handler) *attrsHandler {
	return &attrsHandler{
		Handler: handler,
	}
}

func (h *attrsHandler) Handle(ctx context.Context, record slog.Record) error {
	if requestID := RequestIDFromContext(ctx); requestID != "" {
		record.AddAttrs(slog.String("request_id", requestID))
	}

	return h.Handler.Handle(ctx, record)
}

func (h *attrsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &attrsHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *attrsHandler) WithGroup(name string) slog.Handler {
	return &attrsHandler{Handler: h.Handler.WithGroup(name)}
}

func NewDefaultLogger() *slog.Logger {
	return slog.New(
		AttrsHandler(
			slog.NewJSONHandler(os.Stdout, nil),
		),
	)
}
