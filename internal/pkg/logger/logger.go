package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var defaultLogger *slog.Logger

// Init initializes the logger with JSON format
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Error logs an error message with context
func Error(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	defaultLogger.ErrorContext(ctx, msg, args...)
}

// Info logs an info message with context
func Info(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	defaultLogger.InfoContext(ctx, msg, args...)
}

// Warn logs a warning message with context
func Warn(ctx context.Context, msg string, args ...any) {
	args = appendRequestID(ctx, args)
	defaultLogger.WarnContext(ctx, msg, args...)
}

// appendRequestID adds request_id to log args if present in context
func appendRequestID(ctx context.Context, args []any) []any {
	if ctx == nil {
		return args
	}
	if reqID := ctx.Value(RequestIDKey); reqID != nil {
		args = append(args, "request_id", reqID)
	}
	return args
}
