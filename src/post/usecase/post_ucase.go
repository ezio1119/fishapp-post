package usecase

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postUsecase struct {
	postRepo       post.Repository
	contextTimeout time.Duration
}

// NewPostUsecase will create new an postUsecase object representation of post.Usecase interface
func NewPostUsecase(p post.Repository, timeout time.Duration) post.Usecase {
	return &postUsecase{
		postRepo:       p,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (p *postUsecase) GetList(ctx context.Context, datetime time.Time, num int64) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	listPost, err := p.postRepo.GetList(ctx, datetime, num)
	if err != nil {
		return nil, err
	}

	return listPost, nil
}

func (p *postUsecase) GetByID(ctx context.Context, id int64) (*models.Post, error) {

	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	res, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *postUsecase) Create(ctx context.Context, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	if err := p.postRepo.Create(ctx, post); err != nil {
		return err
	}
	return nil
}

func (p *postUsecase) Update(ctx context.Context, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	res, err := p.postRepo.GetByID(ctx, post.Id)
	if err != nil {
		return err
	}
	if res.UserId != post.UserId {
		return status.Error(codes.Unauthenticated, "do not have permission to update this post")
	}
	now := time.Now()
	post.UpdatedAt = now
	if err := p.postRepo.Update(ctx, post); err != nil {
		return err
	}
	post.CreatedAt = res.CreatedAt
	return nil
}

func (p *postUsecase) Delete(ctx context.Context, id int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	res, err := p.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if res.UserId != userID {
		return status.Error(codes.Unauthenticated, "do not have permission to delete this post")
	}
	return p.postRepo.Delete(ctx, id)
}
