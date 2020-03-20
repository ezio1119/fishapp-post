package interactor

import (
	"context"
	"fmt"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostInteractor interface {
	GetPost(ctx context.Context, id int64) (*models.Post, error)
	ListPosts(ctx context.Context, p *models.Post, pageSize int64, pageToken string, filter *models.PostFilter) ([]*models.Post, string, error)
	CreatePost(ctx context.Context, p *models.Post) error
	UpdatePost(ctx context.Context, p *models.Post) error
	DeletePost(ctx context.Context, id int64, userID int64) error

	ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, applyPost *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, id int64, userID int64) error
}

type postInteractor struct {
	postRepo   repo.PostRepo
	ctxTimeout time.Duration
}

func NewPostInteractor(pr repo.PostRepo, timeout time.Duration) *postInteractor {
	return &postInteractor{pr, timeout}
}

func (i *postInteractor) ListPosts(ctx context.Context, p *models.Post, pageSize int64, pageToken string, f *models.PostFilter) ([]*models.Post, string, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	if pageSize == 0 {
		pageSize = conf.C.Sv.DefaultPageSize
	}
	pageSize++
	var cursor int64
	if pageToken != "" {
		var err error
		cursor, err = extractIDFromPageToken(pageToken)
		fmt.Println(cursor)
		if err != nil {
			return nil, "", err
		}
	}

	list, err := i.postRepo.ListPosts(ctx, p, pageSize, cursor, f)
	if err != nil {
		return nil, "", err
	}
	if len(list) == int(pageSize) {
		list = list[:pageSize-1]
		pageToken = genPageTokenFromID(list[len(list)-1].ID)
	}

	return list, pageToken, nil
}

func (i *postInteractor) GetPost(ctx context.Context, id int64) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	return i.postRepo.GetPostWithPostsFishType(ctx, id)
}

func (i *postInteractor) CreatePost(ctx context.Context, p *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	if err := i.postRepo.CreatePost(ctx, p); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) UpdatePost(ctx context.Context, p *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	res, err := i.postRepo.GetPostWithPostsFishType(ctx, p.ID)
	if err != nil {
		return err
	}
	if res.UserID != p.UserID {
		return status.Errorf(codes.PermissionDenied, "user_id=%d does not have permission to update post_id=%d", p.UserID, res.ID)
	}
	listA, err := i.postRepo.ListApplyPosts(ctx, &models.ApplyPost{PostID: p.ID})
	if err != nil {
		return err
	}
	if len(listA) > int(p.MaxApply) {
		return status.Error(codes.InvalidArgument, "there are more apply than max_apply")
	}
	postsFishTypeIDs := make([]int64, len(res.PostsFishTypes))
	for i, f := range res.PostsFishTypes {
		postsFishTypeIDs[i] = f.ID
	}
	if err := i.postRepo.BatchDeletePostsFishType(ctx, postsFishTypeIDs); err != nil {
		return err
	}
	if err := i.postRepo.UpdatePost(ctx, p); err != nil {
		return err
	}
	p.CreatedAt = res.CreatedAt
	return nil
}

func (i *postInteractor) DeletePost(ctx context.Context, id int64, uID int64) error {

	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	res, err := i.postRepo.GetPost(ctx, id)
	if err != nil {
		return err
	}
	if res.UserID != uID {
		return status.Errorf(codes.PermissionDenied, "user_id=%d does not have permission to Delete post_id=%d", uID, id)
	}

	if err := i.postRepo.DeletePost(ctx, id); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) ListApplyPosts(ctx context.Context, a *models.ApplyPost) ([]*models.ApplyPost, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if a.UserID == 0 && a.PostID == 0 {
		return nil, status.Error(codes.InvalidArgument, "enter a value for either user_id or post_id")
	}

	return i.postRepo.ListApplyPosts(ctx, a)
}

func (i *postInteractor) CreateApplyPost(ctx context.Context, a *models.ApplyPost) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	p, err := i.postRepo.GetPost(ctx, a.PostID)
	if err != nil {
		return err
	}
	if p.UserID == a.UserID {
		return status.Error(codes.PermissionDenied, "cannot apply your own post")
	}
	res, err := i.postRepo.ListApplyPosts(ctx, &models.ApplyPost{PostID: a.PostID})
	if err != nil {
		return err
	}
	if len(res) >= int(p.MaxApply) {
		return status.Error(codes.InvalidArgument, "cannot apply because upper limit")
	}
	if err := i.postRepo.CreateApplyPost(ctx, a); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) DeleteApplyPost(ctx context.Context, id int64, uID int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	a, err := i.postRepo.GetApplyPost(ctx, id)
	if err != nil {
		return err
	}
	if a.UserID != uID {
		return status.Errorf(codes.PermissionDenied, "user_id=%d does not have permission to Delete apply_post_id=%d", uID, id)
	}
	if err := i.postRepo.DeleteApplyPost(ctx, a.UserID); err != nil {
		return err
	}
	return nil
}
