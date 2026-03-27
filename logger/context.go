package logger

import "context"

type ctxKey int

const (
	ctxKeyTraceID ctxKey = iota
	ctxKeyUserID
)

// WithTraceID returns a new context with the given trace ID injected.
// The value is automatically included in every log entry produced from this context.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyTraceID, id)
}

// WithUserID returns a new context with the given user ID injected.
// The value is automatically included in every log entry produced from this context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, id)
}

func extractFromContext(ctx context.Context) (traceID, userID string) {
	if v, ok := ctx.Value(ctxKeyTraceID).(string); ok {
		traceID = v
	}
	if v, ok := ctx.Value(ctxKeyUserID).(string); ok {
		userID = v
	}
	return
}
