package entry

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
)

type Repository interface {
	Create(ctx context.Context, e *models.Entry) error
	GetByID(ctx context.Context, id int64) (*models.Entry, error)
	Delete(ctx context.Context, id int64) error
}
