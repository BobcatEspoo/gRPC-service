package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"log"
	"sync"
)

var (
	EndpointMetrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_endpoint_requests_total",
			Help: "Total number of requests to specific gRPC endpoints",
		},
		[]string{"endpoint"},
	)
	once sync.Once
)

func InitMetrics() {
	once.Do(func() {
		prometheus.MustRegister(EndpointMetrics)
	})
}

func InitOpenTelemetry() {
	_, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("MyGRPCService"),
			semconv.ServiceVersionKey.String("v1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
}
