package repo

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
)

type ApplyPostRepo interface {
	GetApplyPost(ctx context.Context, pID int64, uID int64) (*models.ApplyPost, error)
	BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error)
	ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, p *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, pID int64, uID int64) error
}
