package server

import (
	"RatesMicroservice/internal/metrics"
	"RatesMicroservice/internal/service"
	"context"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	promExporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Failed to create Prometheus exporter: %v", err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(promExporter))
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shut down MeterProvider: %v", err)
		}
	}()
	metrics.InitMetrics()
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error in listing port: %v", err)
	}
	Db, err := AccessToDB()
	if err != nil {
		log.Fatalf("Error in accessing to DB: %v", err)
	}
	MainServer := service.NewRateService(Db)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	service.RegisterRatesServiceServer(grpcServer, MainServer)

	grpc_prometheus.Register(grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram()

	healthServer.SetServingStatus("RateService", grpc_health_v1.HealthCheckResponse_SERVING)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("server on port: 50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Error in running on port: %v", err)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics_telemetry", promExporter)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Error in running on port: %v", err)
	}
	// Ожидание системного сигнала
	<-stop
	log.Println("Receiving signal for server stop...")

	// Graceful Shutdown
	go func() {
		time.Sleep(5 * time.Second) // Можно задать тайм-аут для завершения активных операций
		log.Println("Shutting down...")
		os.Exit(1)
	}()
	grpcServer.GracefulStop()

	log.Println("server gracefully shutdown")
}
func AccessToDB() (*sql.DB, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file ", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL not install in env")
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error in access to db: %v", err)
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database is anavailable: %v", err)
	}

	log.Println("Fine access to db")
	return db, nil
}

//Проверка Health Check
//Вы можете использовать grpc-health-probe для проверки состояния сервера. Установите утилиту:
//
//bash
//Копировать код
//go install github.com/grpc-ecosystem/grpc-health-probe@latest
//Затем проверьте статус Health Check:
//
//bash
//Копировать код
//grpc-health-probe -addr=localhost:50051
