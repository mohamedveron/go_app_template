package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// LoggingMiddleware returns a gin middleware that logs requests with trace context
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Add trace information to the context if available
		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("path", path),
			slog.String("remote_addr", c.ClientIP()),
		}

		// Add trace information if available
		if span := trace.SpanFromContext(c.Request.Context()); span.SpanContext().IsValid() {
			spanCtx := span.SpanContext()
			attrs = append(attrs,
				slog.String("trace_id", spanCtx.TraceID().String()),
				slog.String("span_id", spanCtx.SpanID().String()),
			)
		}

		// Process request
		c.Next()

		// Log the request
		duration := time.Since(start)
		status := c.Writer.Status()

		attrs = append(attrs,
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.Int("response_size", c.Writer.Size()),
		)

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.LogAttrs(c.Request.Context(), level, "HTTP request processed", attrs...)
	}
}
