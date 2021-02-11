package log

import (
	"context"
)

// AdapterInterface is implemented by all logging adapters.
type AdapterInterface interface {

	// ContextWithTraceID attaches a trace id to context that can be later read by the logger.
	ContextWithTraceID(ctx context.Context, id string) context.Context

	// ContextWithTracePoint attaches an appendable trace path to context that can be later read by the logger.
	ContextWithTracePoint(ctx context.Context, point string) context.Context

	// Error logs a message as of error type.
	Error(ctx context.Context, message string, options ...interface{})

	// Debug logs a message as of debug type.
	Debug(ctx context.Context, message string, options ...interface{})

	// Info logs a message as of information type.
	Info(ctx context.Context, message string, options ...interface{})

	// Warn logs a message as of warning type.
	Warn(ctx context.Context, message string, options ...interface{})

	// Destruct will close the logger gracefully releasing all resources.
	Destruct()
}
