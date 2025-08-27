package calendargrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/logger"
	"google.golang.org/grpc"
)

// LoggingUnaryInterceptor returns a unary interceptor that logs details about the request.
func LoggingUnaryInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		logger.Info(fmt.Sprintf("gRPC call start: %s", info.FullMethod))

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			logger.Info(fmt.Sprintf("gRPC call error: %s | duration: %s | error: %v", info.FullMethod, duration, err))
		} else {
			logger.Info(fmt.Sprintf("gRPC call end: %s | duration: %s", info.FullMethod, duration))
		}

		return resp, err
	}
}
