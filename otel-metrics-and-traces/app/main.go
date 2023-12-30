package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/contrib/detectors/gcp"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	// create echo instance
	e := echo.New()
	ctx := context.Background()

	// create counter
	// meterProvider := newMeterProvider(ctx)
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "sample-local-app"
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
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		log.Fatalf("Error creating exporter: %s", err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	defer provider.Shutdown(ctx)
	otel.SetMeterProvider(provider)

	meter := otel.Meter("github.com/taxintt/otel-metrics-demo")
	counter, err := meter.Int64Counter("demo-app/counter")
	if err != nil {
		log.Fatalf("Error creating counter: %s", err)
	}

	// create tracer
	traceProvider := newTraceProvider(ctx)
	tracer := traceProvider.Tracer("github.com/taxintt/otel-traces-demo")
	otel.SetTracerProvider(traceProvider)

	// create middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		// create span
		_, span := tracer.Start(ctx, "op1")
		defer span.End()

		time.Sleep(1000 * time.Millisecond)

		// increment counter
		counter.Add(ctx, 100)
		return c.String(http.StatusOK, "Hello, World!")
	})

	// start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("ENV_PORT")))

	// graceful shutdown
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
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
			semconv.ServiceNameKey.String("sample-local-app"),
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
		serviceName = "sample-local-app"
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
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		log.Fatalf("Error creating exporter: %s", err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	defer provider.Shutdown(ctx)

	return provider
}
