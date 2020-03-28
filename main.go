package main

import (
	"log"
	"net"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/infrastructure"
	"github.com/ezio1119/fishapp-post/infrastructure/middleware"
	"github.com/ezio1119/fishapp-post/interfaces/controllers"
	"github.com/ezio1119/fishapp-post/interfaces/repo"
	"github.com/ezio1119/fishapp-post/usecase/interactor"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConn := infrastructure.NewGormDB()
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	ctxTimeout := time.Duration(conf.C.Sv.Timeout) * time.Second
	pController := controllers.NewPostController(
		interactor.NewPostInteractor(
			repo.NewPostRepo(dbConn),
			repo.NewPostsFishTypeRepo(dbConn),
			repo.NewApplyPostRepo(dbConn),
			ctxTimeout,
		))
	server := infrastructure.NewGrpcServer(
		middleware.InitMiddleware(),
		pController,
	)
	list, err := net.Listen("tcp", ":"+conf.C.Sv.Port)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(list); err != nil {
		log.Fatal(err)
	}
}
