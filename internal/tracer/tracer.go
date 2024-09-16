package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewExporter creates a new OTLP trace exporter
func NewExporter(endpoint string) (tracesdk.SpanExporter, error) {
	// headers := map[string]string{
	// 	"content-type": "application/json",
	// }

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(endpoint),
			// otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}
	return exporter, nil
}

// NewTraceProvider creates a new TracerProvider
func NewTraceProvider(exp tracesdk.SpanExporter, serviceName string) (*tracesdk.TracerProvider, error) {
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(
			exp,
			tracesdk.WithMaxExportBatchSize(tracesdk.DefaultMaxExportBatchSize),
			tracesdk.WithBatchTimeout(tracesdk.DefaultScheduleDelay*time.Millisecond),
			tracesdk.WithMaxExportBatchSize(tracesdk.DefaultMaxExportBatchSize),
		),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)
	return tracerProvider, nil
}

// InitTracer initializes the OpenTelemetry tracer
func InitTracer(endpoint, serviceName string) (*tracesdk.TracerProvider, error) {
	exporter, err := NewExporter(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	tracerProvider, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace provider: %w", err)
	}

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}
