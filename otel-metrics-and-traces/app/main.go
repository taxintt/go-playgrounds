package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/contrib/detectors/gcp"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var counter metric.Int64Counter

func main() {
	// create echo instance
	e := echo.New()
	ctx := context.Background()

	// create counter
	meterProvider := newMeterProvider(ctx)
	meter := meterProvider.Meter("example.com/metrics")
	counter, err := meter.Int64Counter("sidecar-sample-counter")
	if err != nil {
		log.Fatalf("Error creating counter: %s", err)
	}

	// create tracer
	traceProvider := newTraceProvider(ctx)
	tracer := traceProvider.Tracer("example.com/trace")

	// create middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		// create span
		_, span := tracer.Start(ctx, "op1")
		defer span.End()

		time.Sleep(1000 * time.Millisecond)

		// increment counter
		counter.Add(context.Background(), 100)
		return c.String(http.StatusOK, "Hello, World!")
	})

	// start server
	e.Logger.Fatal(e.Start(":8000"))
}

func newTraceProvider(ctx context.Context) *sdktrace.TracerProvider {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	var exporter sdktrace.SpanExporter
	var err error

	if isCloudRun := os.Getenv("K_SERVICE") != ""; isCloudRun {
		exporter, err = texporter.New(texporter.WithProjectID(projectID))
		if err != nil {
			log.Fatalf("texporter.New: %v", err)
		}
	} else {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			log.Fatalf("stdouttrace.New: %v", err)
		}
	}

	res, err := resource.New(ctx,
		// Use the GCP resource detector to detect information about the GCP platform
		resource.WithDetectors(gcp.NewDetector()),
		// Keep the default detectors
		resource.WithTelemetrySDK(),
		// Add your own custom attributes to identify your application
		resource.WithAttributes(
			semconv.ServiceNameKey.String("my-application"),
		),
	)
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer provider.Shutdown(ctx) // flushes any pending spans, and closes connections.

	return provider
}

func newMeterProvider(ctx context.Context) *sdkmetric.MeterProvider {
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "sample-cloud-run-app"
	}
	res, err := resource.Merge(
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
		sdkmetric.WithResource(res),
	)
	defer provider.Shutdown(ctx)

	return provider
}
