package repo

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
)

type PostsFishTypeRepo interface {
	ListPostsFishTypesByPostID(ctx context.Context, postID int64) ([]*models.PostsFishType, error)
}
