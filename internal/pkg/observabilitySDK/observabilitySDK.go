package observabilitySDK

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const componentName = "location-history"

const accessTokenValue = "access_token_from_otel_provider"

var tracer trace.Tracer

var shutDownFunctions []func() error

// Resource is the part of SDK that defines the object that generated the telemetry signal
func createResource(ctx context.Context) (*resource.Resource, error) {
	// build attribute list based on available environment vars from the platform (e.g. kubernetes)
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(componentName)))
}

// Exporter is the part of SDK that exports telemetry out of application to a backend (e.g. jaeger etc)
// Here we are using the OTLP (https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/otlp/otlptrace)
// exporter
func otlpGRPCSpanExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	return otlptracegrpc.New(ctx,
		secureOption,
		//otlptracegrpc.WithInsecure(),
		//otlptracegrpc.WithEndpoint("https://localhost:8360"),
		otlptracegrpc.WithEndpoint("ingest.lightstep.com:443"),
		otlptracegrpc.WithHeaders(exporterHeaders()),
		otlptracegrpc.WithCompressor(gzip.Name))
}

func exporterHeaders() map[string]string {
	return map[string]string{
		"lightstep-access-token": accessTokenValue,
	}
}

// SpanProcessors are pipelines that will receive spans from Tracer and pass it along to Exporters
// BatchSpanProcessor is responsible for batching spans before flushing it out to exporter to be sent on wire
func batchSpanProcessor(exporter *otlptrace.Exporter) sdktrace.SpanProcessor {
	return sdktrace.NewBatchSpanProcessor(exporter)
}

// TraceProvider provides a `Tracer` to the instrumentation module used in target code
func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, func() error, error) {
	traceExporter, err := otlpGRPCSpanExporter(ctx)
	if err != nil {
		return nil, nil, err
	}
	spanProcessor := batchSpanProcessor(traceExporter)

	res, err := createResource(ctx)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()))

	return tp, func() error {
		_ = spanProcessor.Shutdown(ctx)
		return traceExporter.Shutdown(ctx)
	}, nil
}

// Propagator is used to extract and inject context data (traceContext, baggage) from messages
// crossing service/application boundaries
func textMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{})
}

func InitOtel(ctx context.Context) error {
	logger := &DefaultLogger{}
	errorHandler := &defaultHandler{logger: logger}
	otel.SetErrorHandler(errorHandler)

	//launcher.ConfigureOpentelemetry()
	tp, shutdown, err := newTracerProvider(ctx)
	if err != nil {
		tracer = trace.NewNoopTracerProvider().Tracer("")
		return err
	}
	shutDownFunctions = append(shutDownFunctions, shutdown)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(textMapPropagator())

	// Create the tracer once
	tracer = otel.Tracer(componentName)
	return nil
}

func Tracer() trace.Tracer {
	return tracer
}

func Shutdown() {
	for _, shutdown := range shutDownFunctions {
		if err := shutdown(); err != nil {
			log.Warnf("failed to shutdown opentelemetry exporters")
		}
	}
}
