package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("demo-go-app"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = tp.Tracer("demo-go-app")

	return tp, nil
}

func processOrder(ctx context.Context, orderID string) error {
	ctx, span := tracer.Start(ctx, "processOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.String("order.type", "online"),
	)

	if err := validateOrder(ctx, orderID); err != nil {
		span.RecordError(err)
		return err
	}

	if err := checkInventory(ctx, orderID); err != nil {
		span.RecordError(err)
		return err
	}

	if err := processPayment(ctx, orderID); err != nil {
		span.RecordError(err)
		return err
	}

	if err := shipOrder(ctx, orderID); err != nil {
		span.RecordError(err)
		return err
	}

	span.SetAttributes(attribute.String("order.status", "completed"))
	return nil
}

func validateOrder(ctx context.Context, orderID string) error {
	_, span := tracer.Start(ctx, "validateOrder")
	defer span.End()

	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	span.SetAttributes(attribute.Bool("order.valid", true))
	return nil
}

func checkInventory(ctx context.Context, orderID string) error {
	_, span := tracer.Start(ctx, "checkInventory")
	defer span.End()

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	available := rand.Intn(10)
	span.SetAttributes(
		attribute.Int("inventory.available", available),
		attribute.String("warehouse.location", "US-EAST"),
	)

	if available < 3 {
		return fmt.Errorf("insufficient inventory")
	}

	return nil
}

func processPayment(ctx context.Context, orderID string) error {
	_, span := tracer.Start(ctx, "processPayment")
	defer span.End()

	time.Sleep(time.Duration(rand.Intn(150)) * time.Millisecond)

	amount := rand.Float64() * 1000
	span.SetAttributes(
		attribute.Float64("payment.amount", amount),
		attribute.String("payment.method", "credit_card"),
		attribute.String("payment.status", "approved"),
	)

	return nil
}

func shipOrder(ctx context.Context, orderID string) error {
	_, span := tracer.Start(ctx, "shipOrder")
	defer span.End()

	time.Sleep(time.Duration(rand.Intn(80)) * time.Millisecond)

	span.SetAttributes(
		attribute.String("shipping.carrier", "USPS"),
		attribute.String("tracking.number", fmt.Sprintf("TRACK-%s", orderID)),
	)

	return nil
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "HTTP POST /order")
	defer span.End()

	orderID := fmt.Sprintf("ORD-%d", rand.Intn(10000))

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.route", "/order"),
	)

	err := processOrder(ctx, orderID)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order %s processed successfully\n", orderID)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
