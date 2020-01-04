package main

import (
	"log"
	"net"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/infrastructure"
	"github.com/ezio1119/fishapp-post/middleware"
	"github.com/ezio1119/fishapp-post/registry"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConn := infrastructure.NewMysqlConn()
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	ctxTimeout := time.Duration(conf.C.Sv.Timeout) * time.Second
	r := registry.NewRegistry(dbConn, ctxTimeout)
	server := infrastructure.NewGrpcServer(
		middleware.InitMiddleware(),
		r.NewPostController(),
		r.NewEntryController(),
	)
	list, err := net.Listen("tcp", ":"+conf.C.Sv.Port)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(list); err != nil {
		log.Fatal(err)
	}
}
