package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	log.Info().Str("protocol", "grpc").Str("method", info.FullMethod).Msg("Recieved grpc request")

	result, err := handler(ctx, req)

	duration := time.Since(startTime)
	statusCode := status.Code(err)

	logger := log.Info()

	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").Str("method", info.FullMethod).Str("status_text", statusCode.String()).Int("status_code", int(statusCode)).Dur("duration", duration).Msg("Completed grpc request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.Header()
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
			Body:           nil,
		}
		handler.ServeHTTP(rec, req)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("error", rec.Body)
		}
		duration := time.Since(startTime)
		logger.Str("protocol", "http").Str("path", req.RequestURI).Str("status_text", http.StatusText(rec.StatusCode)).Int("status_code", rec.StatusCode).Dur("duration", duration).Msg("Completed a HTTP request")
	})
}
