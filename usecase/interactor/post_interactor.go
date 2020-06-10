package interactor

import (
	"bytes"
	"context"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/interactor/saga"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostInteractor interface {
	GetPost(ctx context.Context, id int64) (*models.Post, error)
	ListPosts(ctx context.Context, p *models.Post, pageSize int64, pageToken string, filter *models.PostFilter) ([]*models.Post, string, error)
	CreatePost(ctx context.Context, p *models.Post, imageBufs map[int64]*bytes.Buffer) (string, error)
	UpdatePost(ctx context.Context, p *models.Post) error
	DeletePost(ctx context.Context, id int64) error

	GetApplyPost(ctx context.Context, id int64) (*models.ApplyPost, error)
	ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error)
	BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error)
	CreateApplyPost(ctx context.Context, applyPost *models.ApplyPost) error
	DeleteApplyPost(ctx context.Context, id int64) error
}

type postInteractor struct {
	postRepo              repo.PostRepo
	imageUploaderRepo     repo.ImageUploaderRepo
	applyPostRepo         repo.ApplyPostRepo
	transactionRepo       repo.TransactionRepo
	createPostSagaManager *saga.CreatePostSagaManager
	ctxTimeout            time.Duration
}

func NewPostInteractor(
	pr repo.PostRepo,
	ir repo.ImageUploaderRepo,
	ar repo.ApplyPostRepo,
	tr repo.TransactionRepo,
	sm *saga.CreatePostSagaManager,
	timeout time.Duration,
) PostInteractor {
	return &postInteractor{pr, ir, ar, tr, sm, timeout}
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

func (i *postInteractor) CreatePost(ctx context.Context, p *models.Post, imageBufs map[int64]*bytes.Buffer) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	for _, buf := range imageBufs {
		url, err := i.imageUploaderRepo.UploadImage(ctx, buf, uuid.New().String())
		if err != nil {
			return "", nil
		}
		p.Images = append(p.Images, &models.Image{
			URL:       url,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	for _, f := range p.PostsFishTypes {
		f.CreatedAt = now
		f.UpdatedAt = now
	}

	ctx, err := i.transactionRepo.BeginTx(ctx)
	if err != nil {
		return "", nil
	}

	defer func() {
		if recover() != nil {
			i.transactionRepo.Roolback(ctx)
		}
	}()

	if err := i.postRepo.CreatePost(ctx, p); err != nil {
		i.transactionRepo.Roolback(ctx)
		return "", err
	}

	if err := i.postRepo.BatchCreatePostsFishTypes(ctx, p); err != nil {
		i.transactionRepo.Roolback(ctx)
		return "", err
	}

	if err := i.postRepo.BatchCreateImage(ctx, p); err != nil {
		i.transactionRepo.Roolback(ctx)
		return "", err
	}

	ctx, err = i.transactionRepo.Commit(ctx)
	if err != nil {
		return "", err
	}

	sagaID := uuid.New().String()
	state := saga.NewCreatePostSagaState(p, "init", sagaID)

	s := i.createPostSagaManager.NewCreatePostSagaManager(state)
	// 同期

	// if err := s.FSM.Event("UploadImage", ctx, image); err != nil {
	// 	return "", err
	// }
	// 非同期
	if err := s.FSM.Event("CreateRoom", ctx); err != nil {
		return "", err
	}

	return sagaID, nil
}

func (i *postInteractor) UpdatePost(ctx context.Context, p *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	now := time.Now()
	p.UpdatedAt = now

	ctx, err := i.transactionRepo.BeginTx(ctx)
	if err != nil {
		return nil
	}

	defer func() {
		if recover() != nil {
			i.transactionRepo.Roolback(ctx)
		}
	}()

	cnt, err := i.applyPostRepo.CountApplyPostsByPostID(ctx, p.ID)
	if err != nil {
		return err
	}

	if cnt > p.MaxApply {
		return status.Errorf(codes.FailedPrecondition, "got max_apply is %d but already have %d apply", p.MaxApply, cnt)
	}

	if err := i.postRepo.UpdatePost(ctx, p); err != nil {
		i.transactionRepo.Roolback(ctx)
		return err
	}

	if err := i.postRepo.DeletePostsFishTypesByPostID(ctx, p.ID); err != nil {
		i.transactionRepo.Roolback(ctx)
		return err
	}

	if err := i.postRepo.BatchCreatePostsFishTypes(ctx, p); err != nil {
		i.transactionRepo.Roolback(ctx)
		return err
	}

	ctx, err = i.transactionRepo.Commit(ctx)
	if err != nil {
		return err
	}

	// 完全なデータで返すため
	p, err = i.postRepo.GetPostByID(ctx, p.ID)
	if err != nil {
		return err
	}
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

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	ctx, err := i.transactionRepo.BeginTx(ctx)
	if err != nil {
		return nil
	}

	defer func() {
		if recover() != nil {
			i.transactionRepo.Roolback(ctx)
		}
	}()

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

	if err := i.applyPostRepo.CreateApplyPost(ctx, a); err != nil {
		i.transactionRepo.Roolback(ctx)
		return err
	}

	ctx, err = i.transactionRepo.Commit(ctx)
	if err != nil {
		return err
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
