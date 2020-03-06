package infrastructure

import (
	"github.com/ezio1119/fishapp-post/middleware"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGrpcServer(middL middleware.Middleware, postController post_grpc.PostServiceServer) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middL.LoggerInterceptor(),
			middL.ValidatorInterceptor(),
			middL.RecoveryInterceptor(),
		)),
	)
	post_grpc.RegisterPostServiceServer(server, postController)
	reflection.Register(server)
	return server
}
