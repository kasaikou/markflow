package docstak

import (
	"context"
	"log/slog"
)

type ctxLogger struct{}

var ctxLoggerKey = ctxLogger{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	return ctx.Value(ctxLoggerKey).(*slog.Logger)
}
