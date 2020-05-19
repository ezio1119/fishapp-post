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
	dbConn := infrastructure.NewMySQLDB()
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
			repo.NewApplyPostRepo(dbConn),
			repo.NewOutboxRepo(dbConn),
			ctxTimeout,
		))
	server := infrastructure.NewGrpcServer(
		middleware.InitMiddleware(),
		pController,
	)
	rController := controllers.NewSagaReplyController(
		interactor.NewSagaReplyInteractor(
			repo.NewSagaInstanceRepo(dbConn),
			repo.NewOutboxRepo(dbConn)),
	)

	natsConn, err := infrastructure.NewNatsStreamingConn()
	if err != nil {
		log.Fatal(err)
	}
	infrastructure.StartSubscribeSagaReply(natsConn, rController)
	// if err := replyC.RoomCreated(context.Background(), "f08ce297-98de-11ea-9ddf-0242ac120005"); err != nil {
	// 	log.Fatal(err)
	// }
	list, err := net.Listen("tcp", ":"+conf.C.Sv.Port)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(list); err != nil {
		log.Fatal(err)
	}
}
