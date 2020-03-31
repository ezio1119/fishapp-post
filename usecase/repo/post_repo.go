package repo

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
)

type PostRepo interface {
	List(ctx context.Context, p *models.Post, num int64, cursor int64, filter *models.PostFilter) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id int64) (*models.Post, error)
	GetPostCanApply(ctx context.Context, id int64) (*models.Post, error)
	// GetPostWithChildlen(ctx context.Context, id int64) (*models.Post, error)
	UpdatePost(ctx context.Context, p *models.Post) error
	CreatePost(ctx context.Context, p *models.Post) error
	Delete(ctx context.Context, id int64) error
	DeletePostsFishTypesByPostID(ctx context.Context, postID int64) error
	ListFishTypeIDsByPostID(ctx context.Context, postID int64) ([]int64, error)
	BatchListPostsFishTypesByPostIDs(ctx context.Context, postIDs []int64) ([]*models.PostsFishType, error)

	GetApplyPost(ctx context.Context, pID int64, uID int64) (*models.ApplyPost, error)
	BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error)
	ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, p *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, pID int64, uID int64) error
}
