package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
)

var Meter = otel.Meter("my-service-meter")

func Record(ctx context.Context, name string) func(time.Time) {
	return func(start time.Time) {
		h, _ := Meter.Float64Histogram(
			name,
		)
		h.Record(ctx, time.Since(start).Seconds())
	}
}
