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
	UpdatePost(ctx context.Context, p *models.Post) (*models.Post, error)
	DeletePost(ctx context.Context, id int64) error

	GetApplyPost(ctx context.Context, postID int64, userID int64) (*models.ApplyPost, error)
	ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, applyPost *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, postID int64, userID int64) error
}

type postInteractor struct {
	postRepo          repo.PostRepo
	postsFishTypeRepo repo.PostsFishTypeRepo
	applyPostRepo     repo.ApplyPostRepo
	ctxTimeout        time.Duration
}

func NewPostInteractor(
	pr repo.PostRepo,
	pfr repo.PostsFishTypeRepo,
	ar repo.ApplyPostRepo,
	timeout time.Duration,
) PostInteractor {
	return &postInteractor{pr, pfr, ar, timeout}
}

func (i *postInteractor) GetPost(ctx context.Context, id int64) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	p, err := i.postRepo.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	f, err := i.postsFishTypeRepo.ListPostsFishTypesByPostID(id)
	if err != nil {
		return nil, err
	}
	p.PostsFishTypes = f
	return p, nil
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
	if list, err = i.applyPostRepo.FillPostWithApplyPost(); err != nil {
		return nil, nil, err
	}
	return list, pageToken, nil
}

func (i *postInteractor) CreatePost(ctx context.Context, p *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	if err := i.postRepo.CreatePost(ctx, p); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) UpdatePost(ctx context.Context, p *models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	res, err := i.postRepo.GetPost(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	if res.MaxApply > p.MaxApply {
		return nil, status.Error(codes.InvalidArgument, "there are more apply than max_apply")
	}
	// postに紐づいているPostsFishTypeをすべて消す
	if err := i.postRepo.DeletePostsFishTypesByPostID(ctx, p.ID); err != nil {
		return nil, err
	}
	// それから新しいPostsFishTypeをアソシエーションでinsertする。batch insertはされない。
	if err := i.postRepo.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	// apply_post,
	return i.postRepo.GetPostWithChildlen(ctx, p.ID)
}

func (i *postInteractor) DeletePost(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if err := i.postRepo.DeletePost(ctx, id); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) GetApplyPost(ctx context.Context, pID int64, uID int64) (*models.ApplyPost, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	p, err := i.postRepo.GetApplyPost(ctx, pID, uID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (i *postInteractor) ListApplyPosts(ctx context.Context, a *models.ApplyPost) ([]*models.ApplyPost, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	// バリデーション 二つのフィールドどちらか入ってないとエラー proto-gen-validateでできなかった
	if a.UserID == 0 && a.PostID == 0 {
		return nil, status.Error(codes.InvalidArgument, "enter a value for either user_id or post_id")
	}

	return i.postRepo.ListApplyPosts(ctx, a)
}

func (i *postInteractor) CreateApplyPost(ctx context.Context, a *models.ApplyPost) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if _, err := i.postRepo.GetPostCanApply(ctx, a.PostID); err != nil {
		return err
	}
	if err := i.postRepo.CreateApplyPost(ctx, a); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) DeleteApplyPost(ctx context.Context, pID int64, uID int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if err := i.postRepo.DeleteApplyPost(ctx, pID, uID); err != nil {
		return err
	}
	return nil
}
