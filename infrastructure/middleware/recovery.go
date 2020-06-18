package middleware

import (
	"github.com/ezio1119/fishapp-post/conf"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (*middleware) UnaryRecoveryInterceptor() grpc.UnaryServerInterceptor {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{}

	if !conf.C.Sv.Debug {
		opts = append(opts, grpc_recovery.WithRecoveryHandler(customFunc))
	}
	return grpc_recovery.UnaryServerInterceptor(opts...)
}

func (*middleware) StreamRecoveryInterceptor() grpc.StreamServerInterceptor {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{}

	if !conf.C.Sv.Debug {
		opts = append(opts, grpc_recovery.WithRecoveryHandler(customFunc))
	}
	return grpc_recovery.StreamServerInterceptor(opts...)
}
