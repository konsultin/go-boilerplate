package otel

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Middleware wraps a fasthttp request handler to create spans for each request.
func Middleware(serviceName string) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			// Extract context from headers
			reqHeader := make(map[string]string)
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				reqHeader[string(key)] = string(value)
			})

			// Inject into Go context
			goCtx := propagator.Extract(context.Background(), propagation.MapCarrier(reqHeader))

			// Start span with manual attributes
			opts := []trace.SpanStartOption{
				trace.WithAttributes(
					attribute.String("http.method", string(ctx.Method())),
					attribute.String("http.url", string(ctx.RequestURI())),
					attribute.String("http.host", string(ctx.Host())),
					attribute.String("http.route", string(ctx.Path())),
					attribute.String("http.user_agent", string(ctx.UserAgent())),
				),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			spanName := string(ctx.Method()) + " " + string(ctx.Path())
			goCtx, span := tracer.Start(goCtx, spanName, opts...)
			defer span.End()

			// Pass traced context to UserValue so handlers can use it
			ctx.SetUserValue("otelCtx", goCtx)

			next(ctx)

			// Record status
			status := ctx.Response.StatusCode()
			span.SetAttributes(attribute.Int("http.status_code", status))
			if status >= 500 {
				span.SetStatus(codes.Error, "Internal Server Error")
			}
		}
	}
}
