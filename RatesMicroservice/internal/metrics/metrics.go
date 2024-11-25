package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
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

func InitializeTracerProvider(serviceName string) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")),
	)
	if err != nil {
		return nil, err
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			serviceName,
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithBatcher(exp),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}
