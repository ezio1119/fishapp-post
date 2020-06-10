package infrastructure

import (
	"github.com/ezio1119/fishapp-post/infrastructure/middleware"
	"github.com/ezio1119/fishapp-post/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGrpcServer(middL middleware.Middleware, postController pb.PostServiceServer) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middL.UnaryLogingInterceptor(),
			middL.UnaryRecoveryInterceptor(),
			middL.UnaryValidationInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			middL.StreamLogingInterceptor(),
			middL.StreamRecoveryInterceptor(),
			middL.StreamValidationInterceptor(),
		)),
	)

	pb.RegisterPostServiceServer(server, postController)
	reflection.Register(server)
	return server
}
