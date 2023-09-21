package main

import (
	"context"
	"crypto/tls"
	"log"
	"strings"

	"github.com/anilsenay/go-opentelemetry-example/internal/handlers"
	"github.com/anilsenay/go-opentelemetry-example/internal/repositories"
	"github.com/anilsenay/go-opentelemetry-example/internal/services"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

const OTEL_EXPORTER_OTLP_ENDPOINT = "https://bcab868a8913431690102bb3d3a20696.apm.europe-west3.gcp.cloud.es.io:443"
const OTEL_EXPORTER_OTLP_HEADERS = "SECRET_TOKEN"

func initTracer() *sdktrace.TracerProvider {
	collectorURL := strings.Replace(OTEL_EXPORTER_OTLP_ENDPOINT, "https://", "", 1)
	secretToken := OTEL_EXPORTER_OTLP_HEADERS

	secureOption := otlptracegrpc.WithInsecure()
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
			otlptracegrpc.WithHeaders(map[string]string{
				"Authorization": "Bearer " + secretToken,
			}),
			otlptracegrpc.WithTLSCredentials(credentials.NewTLS(&tls.Config{})),
		),
	)
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

	// dependency injection
	database := initDb()
	todoRepository := repositories.NewTodoRepository(database)
	todoService := services.NewTodoService(todoRepository)
	todoHandler := handlers.NewTodoHandler(todoService)

	// fiber app
	app := fiber.New()
	app.Use(otelfiber.Middleware())

	todoHandler.SetRoutes(app)

	app.Listen(":3000")
}
