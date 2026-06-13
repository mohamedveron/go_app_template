package observability

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/mohamedveron/go_app_template/logging"
	"go.opentelemetry.io/otel/trace"
)

// ContextHandler is a slog handler that adds trace information to log records
type ContextHandler struct {
	handler slog.Handler
}

// NewContextHandler creates a new context-aware slog handler
func NewContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{handler: handler}
}

// Enabled implements slog.Handler
func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler and adds trace information to the log record
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add trace information if available
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		spanCtx := span.SpanContext()
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs implements slog.Handler
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{handler: h.handler.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{handler: h.handler.WithGroup(name)}
}

// InitLogger initializes the structured logger with OpenTelemetry integration
func InitLogger(format string) {
	// Create handler options
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Check if debug logging is enabled
	if os.Getenv("DEBUG") == "true" {
		opts.Level = slog.LevelDebug
	}

	// Choose handler based on format
	var handler slog.Handler
	switch format {
	case "pretty":
		// Use custom pretty handler
		handler = NewPrettyHandler(os.Stdout, opts)
	case "json":
		fallthrough
	default:
		// Use JSON handler for structured logging (default)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	contextHandler := NewContextHandler(handler)

	logger := slog.New(contextHandler)
	slog.SetDefault(logger)
}

// LoggerWithContext returns a logger that will include trace information in log entries
func LoggerWithContext(ctx context.Context) *slog.Logger {
	return slog.With()
}

// PrettyHandler is a simple pretty handler for logs
type PrettyHandler struct {
	opts      slog.HandlerOptions
	writer    io.Writer
	attrs     []slog.Attr
	formatter *logging.PrettyFormatter
}

// NewPrettyHandler creates a new pretty handler
func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		opts:      *opts,
		writer:    w,
		formatter: logging.NewPrettyFormatter(),
	}
}

// Enabled implements slog.Handler
func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle implements slog.Handler
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// Collect all attributes
	attrs := make(map[string]interface{})

	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value
		return true
	})

	// Add handler-level attributes
	for _, attr := range h.attrs {
		attrs[attr.Key] = attr.Value
	}

	// Create log entry
	entry := logging.LogEntry{
		Timestamp: r.Time,
		Level:     logging.SlogLevelToLogLevel(r.Level),
		Message:   r.Message,
		Attrs:     attrs,
	}

	// Format and write
	formatted := h.formatter.Format(entry)
	_, err := h.writer.Write([]byte(formatted))
	return err
}

// WithAttrs implements slog.Handler
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &PrettyHandler{
		opts:      h.opts,
		writer:    h.writer,
		attrs:     newAttrs,
		formatter: h.formatter,
	}
}

// WithGroup implements slog.Handler
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// For simplicity, just return the same handler
	return h
}
