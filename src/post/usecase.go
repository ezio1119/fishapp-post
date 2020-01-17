package post

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

// Usecase represent the post's usecases
type Usecase interface {
	GetList(ctx context.Context, datetime time.Time, num int64) (*post_grpc.ListPost, error)
	GetByID(ctx context.Context, id int64) (*post_grpc.Post, error)
	Update(ctx context.Context, p *models.Post, userID int64) (*post_grpc.Post, error)
	Create(ctx context.Context, p *models.Post) (*post_grpc.Post, error)
	Delete(ctx context.Context, id int64, userID int64) error
	// CreateEntry(ctx context.Context, e *models.Entry) error
	// DeleteEntry(ctx context.Context, id int64, userID int64) error
}
