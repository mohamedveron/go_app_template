package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	tracer "go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/mohamedveron/go_app_template/"

// TelemetryConfig holds configuration for OpenTelemetry
type TelemetryConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	EnableStdout   bool
}

// DefaultTelemetryConfig returns a default telemetry configuration
func DefaultTelemetryConfig() *TelemetryConfig {
	return &TelemetryConfig{
		ServiceName:    "mohamedveron/go_app_template",
		ServiceVersion: "1.0.0",
		Environment:    getEnv("ENVIRONMENT", "development"),
		OTLPEndpoint:   getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		EnableStdout:   getEnv("OTEL_ENABLE_STDOUT", "false") == "true",
	}
}

// InitTelemetry initializes OpenTelemetry tracing
func InitTelemetry(ctx context.Context, config *TelemetryConfig) (*trace.TracerProvider, error) {
	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create span exporters
	var exporters []trace.SpanExporter

	// OTLP HTTP exporter if endpoint is provided
	if config.OTLPEndpoint != "" {
		otlpExporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(config.OTLPEndpoint),
		)
		if err != nil {
			slog.Warn("Failed to create OTLP exporter", "error", err, "endpoint", config.OTLPEndpoint)
		} else {
			exporters = append(exporters, otlpExporter)
			slog.Info("OTLP exporter configured", "endpoint", config.OTLPEndpoint)
		}
	}

	// Stdout exporter if enabled
	if config.EnableStdout {
		stdoutExporter, err := stdouttrace.New()
		if err != nil {
			slog.Warn("Failed to create stdout exporter", "error", err)
		} else {
			exporters = append(exporters, stdoutExporter)
			slog.Info("Stdout trace exporter enabled")
		}
	}

	// If no exporters are configured, use a no-op exporter
	if len(exporters) == 0 {
		slog.Info("No trace exporters configured, traces will not be exported")
	}

	// Create batch span processor for each exporter
	var spanProcessors []trace.SpanProcessor
	for _, exporter := range exporters {
		spanProcessors = append(spanProcessors, trace.NewBatchSpanProcessor(exporter))
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
	)

	// Add span processors
	for _, sp := range spanProcessors {
		tp.RegisterSpanProcessor(sp)
	}

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.TraceContext{})

	slog.Info("OpenTelemetry initialized",
		"service", config.ServiceName,
		"version", config.ServiceVersion,
		"environment", config.Environment,
		"exporters", len(exporters))

	return tp, nil
}

// Shutdown gracefully shuts down the telemetry
func Shutdown(ctx context.Context, tp *trace.TracerProvider) error {
	if tp == nil {
		return nil
	}

	slog.Info("Shutting down OpenTelemetry...")
	return tp.Shutdown(ctx)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Tracer() tracer.Tracer {
	return otel.Tracer(tracerName)
}
