package interactor

import (
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (i *postInteractor) GetListPosts(ctx context.Context, num int64) (*post_grpc.GetListPostsRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	posts, err := i.postRepository.GetListPosts(ctx, num)
	if err != nil {
		return nil, err
	}
	return i.postPresenter.ConvertGetListPostsRes(posts)
}

func (i *postInteractor) GetPostByID(ctx context.Context, id int64) (*post_grpc.GetPostByIDRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	p, err := i.postRepository.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return i.postPresenter.ConvertGetPostByIDRes(p)
}

func (i *postInteractor) CreatePost(ctx context.Context, p *models.Post) (*post_grpc.CreatePostRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := i.postRepository.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	return i.postPresenter.ConvertCreatePostRes(p)
}

func (i *postInteractor) UpdatePost(ctx context.Context, p *models.Post, uID int64) (*post_grpc.UpdatePostRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	res, err := i.postRepository.GetPostByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	if res.UserID != uID {
		return nil, status.Errorf(codes.PermissionDenied, "user_id=%d does not have permission to Update post_id=%d", uID, p.ID)
	}
	p.UpdatedAt = time.Now()
	if err := i.postRepository.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	p.CreatedAt = res.CreatedAt // pの中にCreatedAtが含まれないため、GetPostByIDの値を挿入
	p.UserID = res.UserID
	return i.postPresenter.ConvertUpdatePostRes(p)
}

func (i *postInteractor) DeletePost(ctx context.Context, id int64, uID int64) (*post_grpc.DeletePostRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	res, err := i.postRepository.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if res.UserID != uID {
		return nil, status.Errorf(codes.PermissionDenied, "user_id=%d does not have permission to Delete post_id=%d", uID, id)
	}

	if err := i.postRepository.DeletePost(ctx, id); err != nil {
		return nil, err
	}
	return i.postPresenter.ConvertDeletePostRes(true)
}

func (i *postInteractor) CreateApplyPost(ctx context.Context, a *models.ApplyPost) (*post_grpc.CreateApplyPostRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	p, err := i.postRepository.GetPostByID(ctx, a.PostID)
	if err != nil {
		return nil, err
	}
	if p.UserID == a.UserID {
		return nil, status.Error(codes.PermissionDenied, "cannot apply your own post")
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	if err := i.postRepository.CreateApplyPost(ctx, a); err != nil {
		return nil, err
	}
	return i.postPresenter.ConvertCreateApplyPostRes(a)
}

func (i *postInteractor) DeleteApplyPost(ctx context.Context, aID int64, uID int64) (*post_grpc.DeleteApplyPostRes, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	a, err := i.postRepository.GetApplyPostByID(ctx, aID)
	if err != nil {
		return nil, err
	}
	if a.UserID != uID {
		return nil, status.Errorf(codes.PermissionDenied, "user_id=%d does not have permission to Delete apply_post_id=%d", uID, aID)
	}
	if err := i.postRepository.DeleteApplyPost(ctx, aID); err != nil {
		return nil, err
	}
	return i.postPresenter.ConvertDeleteApplyPostRes(true)
}
