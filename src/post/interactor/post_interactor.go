package interactor

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

type postInteractor struct {
	postRepository post.Repository
	postPresenter  post.Presenter
	ctxTimeout     time.Duration
}

// NewPostInteractor will create new an postInteractor object representation of post.Usecase interface
func NewPostInteractor(pr post.Repository, pp post.Presenter, timeout time.Duration) post.Usecase {
	return &postInteractor{pr, pp, timeout}
}

func (p *postInteractor) GetList(ctx context.Context, datetime time.Time, num int64) (*post_grpc.ListPost, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	listPost, err := p.postRepository.GetList(ctx, datetime, num)
	if err != nil {
		return nil, err
	}
	return p.postPresenter.TransformListPostProto(listPost)
}

func (p *postInteractor) GetByID(ctx context.Context, id int64) (*post_grpc.Post, error) {

	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	post, err := p.postRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return p.postPresenter.TransformPostProto(post)
}

func (p *postInteractor) Create(ctx context.Context, post *models.Post) (*post_grpc.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	if err := p.postRepository.Create(ctx, post); err != nil {
		return nil, err
	}
	return p.postPresenter.TransformPostProto(post)
}

func (p *postInteractor) Update(ctx context.Context, post *models.Post, userID int64) (*post_grpc.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	res, err := p.postRepository.GetByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	if res.UserID != userID {
		return nil, models.WrapOnPostInterErr(&models.UpdatePostPermissionDenied{PostID: post.ID, UserID: userID})
	}
	now := time.Now()
	post.UpdatedAt = now
	if err := p.postRepository.Update(ctx, post); err != nil {
		return nil, err
	}
	post.CreatedAt = res.CreatedAt
	return p.postPresenter.TransformPostProto(post)
}

func (p *postInteractor) Delete(ctx context.Context, id int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	res, err := p.postRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if res.UserID != userID {
		return models.WrapOnPostInterErr(&models.DeletePostPermissionDenied{PostID: id, UserID: userID})
	}
	return p.postRepository.Delete(ctx, id)
}
