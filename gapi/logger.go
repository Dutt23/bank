package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	log.Info().Str("protocol", "Msg").Str("method", info.FullMethod).Msg("Recieved grpc request")

	result, err := handler(ctx, req)

	duration := time.Since(startTime)
	statusCode := status.Code(err)

	logger := log.Info()

	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "Msg").Str("method", info.FullMethod).Str("status_text", statusCode.String()).Int("status_code", int(statusCode)).Dur("duration", duration).Msg("Recieved grpc request")
	return result, err
}
