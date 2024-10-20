package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	result, err := handler(ctx, req)
	startTime := time.Now()
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protocol", "gRPC").
		Str("method", info.FullMethod).
		Dur("duration", duration).
		Int("status_code", int(statusCode)).
		Str("status", statusCode.String()).
		Msg("received from gRPC client")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		duration := time.Since(startTime)

		logger := log.Info()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			statusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		if rec.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)

		}

		logger.
			Str("protocol", "HTTP").
			Str("method", req.Method).
			Dur("duration", duration).
			Int("status_code", int(rec.statusCode)).
			Str("status", http.StatusText(rec.statusCode)).
			Msg("received from HTTP client")
	})
}
