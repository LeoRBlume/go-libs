package logger

import (
	"net/http"

	"github.com/google/uuid"
)

// TraceMiddleware injects a trace ID into the context of every incoming HTTP request.
// It reads the X-Correlation-ID request header; if absent, a UUID v4 is generated.
// The resolved trace ID is always echoed back in the X-Correlation-ID response header.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Correlation-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		w.Header().Set("X-Correlation-ID", traceID)
		ctx := WithTraceID(r.Context(), traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
