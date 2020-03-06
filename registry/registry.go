package registry

import (
	"database/sql"
	"time"

	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

type registry struct {
	db         *sql.DB
	ctxTimeout time.Duration
}

type Registry interface {
	NewPostController() post_grpc.PostServiceServer
}

func NewRegistry(conn *sql.DB, t time.Duration) Registry {
	return &registry{conn, t}
}
