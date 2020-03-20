package middleware

import (
	"google.golang.org/grpc"
)

type Middleware interface {
	LoggerInterceptor() grpc.UnaryServerInterceptor
	RecoveryInterceptor() grpc.UnaryServerInterceptor
	ValidatorInterceptor() grpc.UnaryServerInterceptor
}

type middleware struct{}

func InitMiddleware() Middleware {
	return &middleware{}
}
