package repo

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
)

type PostsFishTypeRepo interface {
	ListByPostID(ctx context.Context, postID int64) ([]*models.PostsFishType, error)
	BatchCreate(ctx context.Context, pID int64, fIDs []int64, now time.Time) error
	
	DeleteByPostID(ctx context.Context, postID int64) error
}
