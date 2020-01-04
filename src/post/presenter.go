package post

import (
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

type Presenter interface {
	TransformPostProto(po *models.Post) (*post_grpc.Post, error)
	TransformListPostProto(listPost []*models.Post) (*post_grpc.ListPost, error)
}
