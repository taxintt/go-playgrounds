package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var counter metric.Int64Counter

func main() {
	// create echo instance
	e := echo.New()

	ctx := context.Background()
	shutdown := setupCounter(ctx)
	defer shutdown(ctx)

	// create middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		counter.Add(context.Background(), 100)
		return c.String(http.StatusOK, "Hello, World!")
	})

	// start server
	e.Logger.Fatal(e.Start(":8000"))
}

func setupCounter(ctx context.Context) func(context.Context) error {
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "sample-cloud-run-app"
	}
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("Error creating resource: %s", err)
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Error creating exporter: %s", err)
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(r),
	)

	meter := provider.Meter("example.com/metrics")
	counter, err = meter.Int64Counter("sidecar-sample-counter")
	if err != nil {
		log.Fatalf("Error creating counter: %s", err)
	}
	return provider.Shutdown
}
