package entry

import (
	"context"

	"github.com/ezio1119/fishapp-post/entry/controllers/entry_post_grpc"
	"github.com/ezio1119/fishapp-post/models"
)

type Usecase interface {
	Create(ctx context.Context, e *models.Entry) (*entry_post_grpc.Entry, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
