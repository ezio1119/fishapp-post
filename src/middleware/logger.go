package middleware

import (
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func (m *GoMiddleware) LoggerInterceptor(debug bool) grpc.UnaryServerInterceptor {
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}
	var zapLogger *zap.Logger
	if debug {
		zapLogger, _ = zap.NewDevelopment()
	} else {
		zapLogger, _ = zap.NewProduction()
	}

	grpc_zap.ReplaceGrpcLogger(zapLogger)
	return grpc_zap.UnaryServerInterceptor(zapLogger, opts...)
}
