package post

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
)

type Repository interface {
	GetListPosts(ctx context.Context, num int64) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id int64) (*models.Post, error)
	UpdatePost(ctx context.Context, p *models.Post) error
	CreatePost(ctx context.Context, p *models.Post) error
	DeletePost(ctx context.Context, id int64) error

	// GetApplyPostByID(ctx context.Context, id int64) (*models.ApplyPost, error)
	GetApplyPostByID(ctx context.Context, id int64) (*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, p *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, id int64) error
}
