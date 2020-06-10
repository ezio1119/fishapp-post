package main

import (
	"context"
	"net"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/infrastructure"
	"github.com/ezio1119/fishapp-post/infrastructure/middleware"
	"github.com/ezio1119/fishapp-post/infrastructure/sqlhandler"
	"github.com/ezio1119/fishapp-post/interfaces/controllers"
	"github.com/ezio1119/fishapp-post/interfaces/repo"
	"github.com/ezio1119/fishapp-post/usecase/interactor"
	"github.com/ezio1119/fishapp-post/usecase/interactor/saga"
)

func main() {
	ctx := context.Background()
	dbConn, err := infrastructure.NewMySQLDB()
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	natsConn, err := infrastructure.NewNatsStreamingConn()
	if err != nil {
		panic(err)
	}
	defer natsConn.Close()

	gcsClient, err := infrastructure.NewGCSClient(ctx)
	if err != nil {
		panic(err)
	}
	defer gcsClient.Close()

	ctxTimeout := time.Duration(conf.C.Sv.Timeout) * time.Second
	sqlHandler := sqlhandler.NewSqlHandler(dbConn)

	createPostSagaManager := saga.InitCreatePostSagaManager(
		repo.NewOutboxRepo(sqlHandler),
		repo.NewPostRepo(sqlHandler),
		repo.NewSagaInstanceRepo(sqlHandler),
		repo.NewTransactionRepo(sqlHandler),
		repo.NewImageUploaderRepo(gcsClient),
	)

	pController := controllers.NewPostController(
		interactor.NewPostInteractor(
			repo.NewPostRepo(sqlHandler),
			repo.NewImageUploaderRepo(gcsClient),
			repo.NewApplyPostRepo(sqlHandler),
			repo.NewTransactionRepo(sqlHandler),
			createPostSagaManager,
			ctxTimeout,
		))
	server := infrastructure.NewGrpcServer(
		middleware.InitMiddleware(),
		pController,
	)

	rController := controllers.NewSagaReplyController(
		interactor.NewSagaReplyInteractor(
			createPostSagaManager,
			repo.NewSagaInstanceRepo(sqlHandler),
		),
	)

	if err := infrastructure.StartSubscribeCreatePostSagaReply(natsConn, rController); err != nil {
		panic(err)
	}

	list, err := net.Listen("tcp", ":"+conf.C.Sv.Port)
	if err != nil {
		panic(err)
	}

	if err := server.Serve(list); err != nil {
		panic(err)
	}
}
