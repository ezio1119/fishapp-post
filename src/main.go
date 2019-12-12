package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	_postGrpcDeliver "github.com/ezio1119/fishapp-post/post/delivery/grpc"
	_postRepo "github.com/ezio1119/fishapp-post/post/repository"
	_postUcase "github.com/ezio1119/fishapp-post/post/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type env struct {
	DbPass     string `required:"true" split_words:"true"`
	DbDbms     string `required:"true" split_words:"true"`
	DbUser     string `required:"true" split_words:"true"`
	DbName     string `required:"true" split_words:"true"`
	DbPort     string `required:"true" split_words:"true"`
	DbHost     string `required:"true" split_words:"true"`
	DbConnOpt  string `required:"true" split_words:"true"`
	Timeout    int64  `required:"true"`
	ListenPort string `required:"true" split_words:"true"`
}

func main() {
	var env env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err)
	}
	CONNECT := env.DbUser + ":" + env.DbPass + "@(" + env.DbHost + ":" + env.DbPort + ")/" + env.DbName + "?" + env.DbConnOpt
	dbConn, err := sql.Open(env.DbDbms, CONNECT)
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
	timeoutContext := time.Duration(env.Timeout) * time.Second
	postUcase := _postUcase.NewPostUsecase(postRepo, timeoutContext)

	gserver := grpc.NewServer()
	_postGrpcDeliver.NewPostServerGrpc(gserver, postUcase)

	list, err := net.Listen("tcp", ":"+env.ListenPort)
	if err != nil {
		log.Fatal(err)
	}

	err = gserver.Serve(list)
	if err != nil {
		log.Fatal(err)
	}
}
