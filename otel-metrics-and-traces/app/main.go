package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	meter  = otel.Meter("github.com/taxintt/otel-metrics-demo")
	tracer = otel.Tracer("github.com/taxintt/otel-traces-demo")

	shutdownFuncs []func(context.Context) error
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// create echo instance
	e := echo.New()
	ctx := context.Background()

	// error handler
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	defer func() {
		err = errors.Join(err, shutdown(context.Background()))
	}()

	// create counter
	err = newMetricProvider(ctx, handleErr)
	counter, err := meter.Int64Counter("demo-app-counter")
	if err != nil {
		handleErr(err)
		return
	}

	// create tracer
	err = newTraceProvider(ctx, handleErr)
	if err != nil {
		handleErr(err)
		return
	}

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
	e.Logger.Fatal(e.Start(":" + os.Getenv("ENV_PORT")))

	// graceful shutdown
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()

	return
}

func newMetricProvider(ctx context.Context, handleErr func(err error)) error {
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
		handleErr(err)
		return err
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
	)
	if err != nil {
		handleErr(err)
		return err
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Second))),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)
	shutdownFuncs = append(shutdownFuncs, provider.Shutdown)

	return nil
}

func newTraceProvider(ctx context.Context, handleErr func(err error)) error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	var exporter sdktrace.SpanExporter
	var err error

	if isCloudRun := os.Getenv("K_SERVICE") != ""; isCloudRun {
		exporter, err = texporter.New(texporter.WithProjectID(projectID))
		if err != nil {
			handleErr(err)
			return err
		}
	} else {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			handleErr(err)
			return err
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
		handleErr(err)
		return err
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	shutdownFuncs = append(shutdownFuncs, provider.Shutdown)

	return nil
}
