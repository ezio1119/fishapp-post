package post

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
)

// Usecase represent the post's usecases
type Usecase interface {
	GetList(ctx context.Context, datetime time.Time, num int64) ([]*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	Update(ctx context.Context, p *models.Post) error
	Create(ctx context.Context, p *models.Post) error
	Delete(ctx context.Context, id int64, userID int64) error
}
