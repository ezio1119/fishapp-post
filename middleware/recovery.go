package middleware

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*middleware) RecoveryInterceptor() grpc.UnaryServerInterceptor {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	return grpc_recovery.UnaryServerInterceptor(opts...)
}
