package interactor

import (
	"context"
	"fmt"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostInteractor interface {
	GetPost(ctx context.Context, id int64) (*models.Post, error)
	ListPosts(ctx context.Context, p *models.Post, pageSize int64, pageToken string, filter *models.PostFilter) ([]*models.Post, string, error)
	CreatePost(ctx context.Context, p *models.Post) (string, error)
	UpdatePost(ctx context.Context, p *models.Post) error
	DeletePost(ctx context.Context, id int64) error

	GetApplyPost(ctx context.Context, id int64) (*models.ApplyPost, error)
	ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error)
	BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, applyPost *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, id int64) error
}

type postInteractor struct {
	postRepo      repo.PostRepo
	applyPostRepo repo.ApplyPostRepo
	outboxRepo    repo.OutboxRepo
	ctxTimeout    time.Duration
}

func NewPostInteractor(
	pr repo.PostRepo,
	ar repo.ApplyPostRepo,
	or repo.OutboxRepo,
	timeout time.Duration,
) PostInteractor {
	return &postInteractor{pr, ar, or, timeout}
}

func (i *postInteractor) GetPost(ctx context.Context, id int64) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	p, err := i.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
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
	nextToken := ""
	if len(list) == int(pageSize) {
		list = list[:pageSize-1]
		nextToken = genPageTokenFromID(list[len(list)-1].ID)
	}

	return list, nextToken, nil
}

func (i *postInteractor) CreatePost(ctx context.Context, p *models.Post) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := i.postRepo.CreatePost(ctx, p); err != nil {
		return "", err
	}

	sagaID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	s := newCreatePostSagaState(p, i.outboxRepo, sagaID.String())
	// fmt.Println(s.FSM.Current())
	// if err := s.FSM.Event("UploadImage"); err != nil {
	// 	fmt.Println(err)
	// }
	fmt.Println(s.FSM.Current())
	if err := s.FSM.Event("CreateRoom"); err != nil {
		return "", err
	}
	fmt.Println(s.FSM.Current())
	return sagaID.String(), nil
}

func (i *postInteractor) UpdatePost(ctx context.Context, p *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	// postに紐づいているapply_postをカウント
	cnt, err := i.applyPostRepo.CountApplyPostsByPostID(ctx, p.ID)
	if err != nil {
		return err
	}
	if cnt > p.MaxApply {
		return status.Errorf(codes.FailedPrecondition, "got max_apply is %d but already have %d apply", p.MaxApply, cnt)
	}

	now := time.Now()
	p.UpdatedAt = now

	oldP, err := i.postRepo.GetPostByID(ctx, p.ID)
	if err != nil {
		return err
	}
	if err := i.postRepo.UpdatePost(ctx, p); err != nil {
		return err
	}
	// 結果整合性
	cnt, err = i.applyPostRepo.CountApplyPostsByPostID(ctx, p.ID)
	if err != nil {
		if err := i.postRepo.UpdatePost(ctx, oldP); err != nil {
			return err
		}
		return err
	}
	if cnt > p.MaxApply {
		if err := i.postRepo.UpdatePost(ctx, oldP); err != nil {
			return err
		}
		return status.Errorf(codes.FailedPrecondition, "got max_apply is %d but already have %d apply", p.MaxApply, cnt)
	}
	// 完全なデータで返すため
	p.CreatedAt = oldP.CreatedAt
	p.UserID = oldP.UserID
	return nil
}

func (i *postInteractor) DeletePost(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if err := i.postRepo.DeletePost(ctx, id); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) GetApplyPost(ctx context.Context, id int64) (*models.ApplyPost, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	return i.applyPostRepo.GetApplyPostByID(ctx, id)
}

func (i *postInteractor) ListApplyPosts(ctx context.Context, a *models.ApplyPost) ([]*models.ApplyPost, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if a.UserID != 0 {
		return i.applyPostRepo.ListApplyPostsByUserID(ctx, a.UserID)
	}
	if a.PostID != 0 {
		return i.applyPostRepo.ListApplyPostsByPostID(ctx, a.PostID)
	}
	return nil, nil
}

func (i *postInteractor) BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	return i.applyPostRepo.BatchGetApplyPostsByPostIDs(ctx, postIDs)
}

func (i *postInteractor) CreateApplyPost(ctx context.Context, a *models.ApplyPost) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	cnt, err := i.applyPostRepo.CountApplyPostsByPostID(ctx, a.PostID)
	if err != nil {
		return err
	}
	p, err := i.postRepo.GetPostByID(ctx, a.PostID)
	if err != nil {
		return err
	}
	if p.MaxApply <= cnt {
		return status.Error(codes.FailedPrecondition, "already reached max_apply limit")
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	if err := i.applyPostRepo.CreateApplyPost(ctx, a); err != nil {
		return err
	}
	// 結果整合性
	cnt, err = i.applyPostRepo.CountApplyPostsByPostID(ctx, a.PostID)
	if err != nil {
		// 補償トランザクション
		if err := i.applyPostRepo.DeleteApplyPost(ctx, a.ID); err != nil {
			return err
		}
		return err
	}
	p, err = i.postRepo.GetPostByID(ctx, a.PostID)
	if err != nil {
		if err := i.applyPostRepo.DeleteApplyPost(ctx, a.ID); err != nil {
			return err
		}
		return err
	}
	if p.MaxApply < cnt {
		if err := i.applyPostRepo.DeleteApplyPost(ctx, a.ID); err != nil {
			return err
		}
		return status.Error(codes.FailedPrecondition, "already reached max_apply limit")
	}
	return nil
}

func (i *postInteractor) DeleteApplyPost(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if err := i.applyPostRepo.DeleteApplyPost(ctx, id); err != nil {
		return err
	}
	return nil
}
