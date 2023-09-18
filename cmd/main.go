package main

import (
	"context"
	"log"

	"github.com/anilsenay/go-opentelemetry-example/internal/handlers"
	"github.com/anilsenay/go-opentelemetry-example/internal/repositories"
	"github.com/anilsenay/go-opentelemetry-example/internal/services"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func initTracer() *sdktrace.TracerProvider {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("todo-app"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initMetrics() (*metric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	return provider, nil
}

func initDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"))
	if err != nil {
		panic("failed to connect database")
	}
	err = db.Use(tracing.NewPlugin())
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	_, err := initMetrics()
	if err != nil {
		log.Fatal(err)
	}

	// dependency injection
	database := initDb()
	todoRepository := repositories.NewTodoRepository(database)
	todoService := services.NewTodoService(todoRepository)
	todoHandler := handlers.NewTodoHandler(todoService)

	// fiber app
	app := fiber.New()
	app.Use(otelfiber.Middleware())

	todoHandler.SetRoutes(app)

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Listen(":3000")
}
