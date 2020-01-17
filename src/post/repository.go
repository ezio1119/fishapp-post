package post

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
)

type Repository interface {
	GetList(ctx context.Context, datetime time.Time, num int64) ([]*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	Update(ctx context.Context, p *models.Post) error
	Create(ctx context.Context, p *models.Post) error
	Delete(ctx context.Context, id int64) error
}
