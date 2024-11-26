package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
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
