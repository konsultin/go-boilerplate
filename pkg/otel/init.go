package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
)

// InitTracerProvider initializes an OTLP trace exporter and registers the trace provider globally.
// It returns a shutdown function that should be called when the service is stopping.
func InitTracerProvider(ctx context.Context, serviceName, serviceVersion, collectorEndpoint string) (func(context.Context) error, error) {
	// If no endpoint is provided, we can either skip OTEL or return error.
	// For boilerplate safety, if empty, we just return no-op.
	if collectorEndpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	headers := map[string]string{
		"User-Agent": fmt.Sprintf("%s/%s", serviceName, serviceVersion),
	}

	// Insecure for Jaeger local. For prod, might need TLS config.
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(collectorEndpoint),
			otlptracegrpc.WithInsecure(),                   // Used for local Jaeger/Collector without TLS
			otlptracegrpc.WithDialOption(grpc.WithBlock()), // Wait for connection
			otlptracegrpc.WithHeaders(headers),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otlp trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			semconv.TelemetrySDKLanguageGo,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample 100% for dev. In prod, adjust this.
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set global propagator to W3C Trace Context (Standard)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}
