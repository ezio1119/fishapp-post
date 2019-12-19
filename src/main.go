package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/middleware"
	_postGrpcDeliver "github.com/ezio1119/fishapp-post/post/delivery/grpc"
	_postRepo "github.com/ezio1119/fishapp-post/post/repository"
	_postUcase "github.com/ezio1119/fishapp-post/post/usecase"
	_ "github.com/go-sql-driver/mysql"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func main() {
	conf.Readconf()
	CONNECT := conf.C.Db.User + ":" + conf.C.Db.Pass + "@(" + conf.C.Db.Host + ":" + conf.C.Db.Port + ")/" + conf.C.Db.Name + "?" + conf.C.Db.ConnOpt
	dbConn, err := sql.Open(conf.C.Db.Dbms, CONNECT)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	postRepo := _postRepo.NewMysqlPostRepository(dbConn)
	timeoutContext := time.Duration(conf.C.Sv.Timeout) * time.Second
	postUcase := _postUcase.NewPostUsecase(postRepo, timeoutContext)

	middL := middleware.InitMiddleware()

	gserver := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middL.LoggerInterceptor(),
			middL.AuthInterceptor(),
			middL.ValidatorInterceptor(),
			middL.RecoveryInterceptor(),
		)),
	)
	_postGrpcDeliver.NewPostServerGrpc(gserver, postUcase)

	list, err := net.Listen("tcp", ":"+conf.C.Sv.Port)
	if err != nil {
		log.Fatal(err)
	}

	err = gserver.Serve(list)
	if err != nil {
		log.Fatal(err)
	}
}
