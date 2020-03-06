package post

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

// Usecase represent the post's usecases
type Usecase interface {
	GetPostByID(ctx context.Context, id int64) (*post_grpc.GetPostByIDRes, error)
	GetListPosts(ctx context.Context, num int64) (*post_grpc.GetListPostsRes, error)
	CreatePost(ctx context.Context, p *models.Post) (*post_grpc.CreatePostRes, error)
	UpdatePost(ctx context.Context, p *models.Post, userID int64) (*post_grpc.UpdatePostRes, error)
	DeletePost(ctx context.Context, id int64, userID int64) (*post_grpc.DeletePostRes, error)
	CreateApplyPost(ctx context.Context, applyPost *models.ApplyPost) (*post_grpc.CreateApplyPostRes, error)
	DeleteApplyPost(ctx context.Context, id int64, userID int64) (*post_grpc.DeleteApplyPostRes, error)
}
