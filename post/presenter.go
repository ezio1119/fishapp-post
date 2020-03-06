package post

import (
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

type Presenter interface {
	ConvertGetPostByIDRes(p *models.Post) (*post_grpc.GetPostByIDRes, error)
	ConvertGetListPostsRes(p []*models.Post) (*post_grpc.GetListPostsRes, error)
	ConvertCreatePostRes(p *models.Post) (*post_grpc.CreatePostRes, error)
	ConvertUpdatePostRes(p *models.Post) (*post_grpc.UpdatePostRes, error)
	ConvertDeletePostRes(bool) (*post_grpc.DeletePostRes, error)
	ConvertCreateApplyPostRes(applyPost *models.ApplyPost) (*post_grpc.CreateApplyPostRes, error)
	ConvertDeleteApplyPostRes(bool) (*post_grpc.DeleteApplyPostRes, error)
}
