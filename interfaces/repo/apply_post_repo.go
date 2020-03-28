package repo

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/jinzhu/gorm"
)

type applyPostRepo struct {
	db *gorm.DB
}

func NewApplyPostRepo(db *gorm.DB) *applyPostRepo {
	return &applyPostRepo{db}
}

func (r *applyPostRepo) GetApplyPost(ctx context.Context, pID int64, uID int64) (*models.ApplyPost, error) {
	panic("not implemented")
}
func (r *applyPostRepo) BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error) {
	panic("not implemented")
}
func (r *applyPostRepo) ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error) {
	panic("not implemented")
}
func (r *applyPostRepo) CreateApplyPost(ctx context.Context, p *models.ApplyPost) error {
	panic("not implemented")
}
func (r *applyPostRepo) DeleteApplyPost(ctx context.Context, pID int64, uID int64) error {
	panic("not implemented")
}
