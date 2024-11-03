package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// var (
// 	// Counter for total gRPC requests
// 	grpcRequestsTotal = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "grpc_requests_total",
// 			Help: "Total number of gRPC requests received",
// 		},
// 		[]string{"method", "status"},
// 	)

// 	// Histogram for gRPC request latency
// 	grpcRequestLatency = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "grpc_request_latency_seconds",
// 			Help:    "Histogram of latencies for gRPC requests",
// 			Buckets: prometheus.DefBuckets, // Default buckets for latency
// 		},
// 		[]string{"method"},
// 	)

// 	// Histogram for request size
// 	grpcRequestSize = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "grpc_request_size_bytes",
// 			Help:    "Histogram of gRPC request sizes in bytes",
// 			Buckets: prometheus.ExponentialBuckets(100, 10, 6), // Custom buckets for request size
// 		},
// 		[]string{"method"},
// 	)

// 	// Histogram for response size
// 	grpcResponseSize = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "grpc_response_size_bytes",
// 			Help:    "Histogram of gRPC response sizes in bytes",
// 			Buckets: prometheus.ExponentialBuckets(100, 10, 6), // Custom buckets for response size
// 		},
// 		[]string{"method"},
// 	)
// )

// func init() {
// 	// Register the Prometheus metrics
// 	// prometheus.MustRegister(grpcRequestsTotal)
// 	// prometheus.MustRegister(grpcRequestLatency)
// 	// prometheus.MustRegister(grpcRequestSize)
// 	// prometheus.MustRegister(grpcResponseSize)
// }

// Define your MetricsInterceptor with the correct signature
func MetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("Intercepting request for method: %s", info.FullMethod)
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Measure latency
	latency := time.Since(start).Seconds()
	log.Printf("Request processed in %f seconds", latency)

	// Capture status code
	statusCode := status.Code(err).String()
	log.Printf("Status code: %s", statusCode)

	return resp, err
}
