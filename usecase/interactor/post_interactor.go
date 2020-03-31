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

func (i *postInteractor) fillPostWithFishTypeIDs(ctx context.Context, p *models.Post) error {
	ids, err := i.postRepo.ListFishTypeIDsByPostID(ctx, p.ID)
	if err != nil {
		return err
	}
	p.FishTypeIDs = ids
	return nil
}

func (i *postInteractor) fillListPostsWithFishTypeIDs(ctx context.Context, posts []*models.Post) error {
	pIDs := make([]int64, len(posts))
	for i, p := range posts {
		pIDs[i] = p.ID
	}
	fishes, err := i.postRepo.BatchListPostsFishTypesByPostIDs(ctx, pIDs)
	if err != nil {
		return err
	}
	for _, p := range posts {
		for _, f := range fishes {
			if p.ID == f.PostID {
				p.FishTypeIDs = append(p.FishTypeIDs, f.FishTypeID)
			}
		}
	}
	return nil
}

func (i *postInteractor) GetPost(ctx context.Context, id int64) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	p, err := i.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.fillPostWithFishTypeIDs(ctx, p); err != nil {
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

	list, err := i.postRepo.List(ctx, p, pageSize, cursor, f)
	if err != nil {
		return nil, "", err
	}
	nextToken := ""
	if len(list) == int(pageSize) {
		list = list[:pageSize-1]
		nextToken = genPageTokenFromID(list[len(list)-1].ID)
	}

	if len(list) != 0 {
		if err = i.fillListPostsWithFishTypeIDs(ctx, list); err != nil {
			return nil, "", err
		}
	}
	return list, nextToken, nil
}

func (i *postInteractor) CreatePost(ctx context.Context, p *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := i.postRepo.CreatePost(ctx, p); err != nil {
		return err
	}
	return nil
}

func (i *postInteractor) UpdatePost(ctx context.Context, p *models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()

	// res, err := i.postRepo.GetByID(ctx, p.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// postに紐づいているapply_postをカウント
	cnt, err := i.applyPostRepo.CountByPostID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	fmt.Println("キトや！！", cnt)

	if cnt > p.MaxApply {
		return nil, status.Errorf(codes.InvalidArgument, "got max_apply is %d but already have %d apply", p.MaxApply, cnt)
	}
	panic("ccas")
	// postに紐づいているPostsFishTypeをすべて消す

	now := time.Now()
	p.UpdatedAt = now
	// それから新しいPostsFishTypeをアソシエーションでinsertする。batch insertはされない。
	if err := i.postRepo.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	if err := i.postsFishTypeRepo.DeleteByPostID(ctx, p.ID); err != nil {
		return nil, err
	}
	panic("ccas")
	// apply_post,
	// return i.postRepo.GetPostWithChildlen(ctx, p.ID)
}

func (i *postInteractor) DeletePost(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, i.ctxTimeout)
	defer cancel()
	if err := i.postRepo.Delete(ctx, id); err != nil {
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
