package interactor

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/entry"
	"github.com/ezio1119/fishapp-post/entry/controllers/entry_post_grpc"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type entryInteractor struct {
	entryRepository entry.Repository
	entryPresenter  entry.Presenter
	postRepository  post.Repository
	ctxTimeout      time.Duration
}

func NewEntryInteractor(er entry.Repository, ep entry.Presenter, pp post.Repository, t time.Duration) entry.Usecase {
	return &entryInteractor{er, ep, pp, t}
}

func (i *entryInteractor) Create(ctx context.Context, e *models.Entry) (*entry_post_grpc.Entry, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	post, err := i.postRepository.GetByID(ctx, e.PostID)
	if err != nil {
		return nil, err
	}
	if post.UserID == e.UserID {
		return nil, status.Error(codes.InvalidArgument, "cannot enter your own post")
	}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	if err := i.entryRepository.Create(ctx, e); err != nil {
		return nil, err
	}
	return i.entryPresenter.TransformEntryProto(e)
}

func (i *entryInteractor) Delete(ctx context.Context, id int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	e, err := i.entryRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if e.UserID != userID {
		return status.Error(codes.PermissionDenied, "do not have permission to delete this entry")
	}
	return i.entryRepository.Delete(ctx, id)
}
