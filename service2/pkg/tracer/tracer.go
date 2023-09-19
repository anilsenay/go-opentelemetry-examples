package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer = otel.Tracer("fiber-server-2")

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer.Start(ctx, spanName, opts...)
}

func Trace(ctx *context.Context, spanName string, opts ...trace.SpanStartOption) func() {
	c, span := Tracer.Start(*ctx, spanName, opts...)
	*ctx = c
	return func() {
		span.End()
	}
}
