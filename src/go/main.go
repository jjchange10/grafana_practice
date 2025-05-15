package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer(ctx context.Context) func() {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("alloy.alloy.svc.cluster.local:4317"), // AlloyやCollector
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("hello-api"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	// 手動でspanを取る（中の処理計測用）
	_, span := otel.Tracer("hello-handler").Start(req.Context(), "handle-hello")
	defer span.End()

	// 属性付けてみる
	span.SetAttributes(attribute.String("custom.attribute", "hello-world"))

	// 実際の処理
	fmt.Fprintln(w, "Hello, World")
}

func main() {
	ctx := context.Background()
	shutdown := initTracer(ctx)
	defer shutdown()

	// 標準http.HandleにOpenTelemetry Middlewareを追加
	handler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "hello-endpoint")

	http.Handle("/hello", handler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
