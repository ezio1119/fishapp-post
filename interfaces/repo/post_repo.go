package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *postRepo {
	return &postRepo{db}
}

func (r *postRepo) CreatePost(ctx context.Context, p *models.Post) error {
	if err := r.db.Create(p).Error; err != nil {
		e, ok := err.(*mysql.MySQLError)
		if ok {
			if e.Number == 1062 {
				err = status.Error(codes.AlreadyExists, err.Error())
			}
		}
		return err
	}
	return nil
}

func (r *postRepo) GetPostWithChildlen(ctx context.Context, id int64) (*models.Post, error) {
	p := &models.Post{}
	if err := r.db.Take(p, id).Related(&p.PostsFishTypes).Related(&p.ApplyPosts).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = status.Errorf(codes.NotFound, "post with id='%d' is not found", id)
		}
		return nil, err
	}
	return p, nil
}

func (r *postRepo) GetPost(ctx context.Context, id int64) (*models.Post, error) {
	p := &models.Post{}
	if err := r.db.Take(p, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = status.Errorf(codes.NotFound, "post with id='%d' is not found", id)
		}
		return nil, err
	}
	return p, nil
}

func (r *postRepo) ListPosts(ctx context.Context, p *models.Post, num int64, cursor int64, f *models.PostFilter) ([]*models.Post, error) {
	list := []*models.Post{}
	tx := r.db.Table("posts").
		Select("posts.*").
		Group("posts.id").
		Where(p).
		Limit(num)

	if !f.MeetingAtFrom.IsZero() && !f.MeetingAtTo.IsZero() {
		tx = tx.Where("meeting_at BETWEEN ? AND ?", f.MeetingAtFrom, f.MeetingAtTo)
	}

	if f.CanApply {
		tx = tx.Joins("left join apply_posts on posts.id = apply_posts.post_id").
			Having("count(apply_posts.id) < posts.max_apply")
	}

	if f.FishTypeIDs != nil {
		tx = tx.Joins("inner join posts_fish_types on posts.id = posts_fish_types.post_id").
			Where("posts_fish_types.fish_type_id IN(?)", f.FishTypeIDs).
			Having("count(posts_fish_types.fish_type_id) = ?", len(f.FishTypeIDs))
	}

	if cursor != 0 {
		switch f.SortBy {
		case models.SortByID:
			if f.OrderBy == models.OrderByAsc {
				tx = tx.Where("posts.id > ?", cursor).
					Order("id asc")
			}
			if f.OrderBy == models.OrderByDesc {
				tx = tx.Where("posts.id < ?", cursor).
					Order("id desc")
			}
		// meeting_atはユニークではないため、同じ値の場合を考えidでも絞り込む
		case models.SortByMeetingAt:
			mAt := []time.Time{}
			// ２箇所で "meeting_at" を使っているからサブクエリじゃなくてpluckつかった
			if err := r.db.New().Model(&models.Post{ID: cursor}).Pluck("meeting_at", &mAt).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {
					err = status.Errorf(codes.NotFound, "post with id='%d' is not found", cursor)
				}
				return nil, err
			}
			if f.OrderBy == models.OrderByAsc {
				tx = tx.Where("meeting_at >= ?", mAt[0]).
					Where("meeting_at > ? or posts.id > ?", mAt[0], cursor).
					Order("meeting_at asc, id asc")
			}
			if f.OrderBy == models.OrderByDesc {
				tx = tx.Where("meeting_at <= ?", mAt[0]).
					Where("meeting_at < ? or posts.id < ?", mAt[0], cursor).
					Order("meeting_at desc, id desc")
			}
		}
	}
	if cursor == 0 {
		tx = tx.Order(fmt.Sprintf("%s %s", f.SortBy, f.OrderBy))
	}

	if err := tx.Preload("ApplyPosts").Preload("PostsFishTypes").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *postRepo) UpdatePost(ctx context.Context, p *models.Post) error {
	if err := r.db.Model(p).Updates(p).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = status.Errorf(codes.NotFound, "post with id='%d' is not found", p.ID)
		}
		return err
	}
	return nil
}

func (r *postRepo) BatchDeletePostsFishType(ctx context.Context, ids []int64) error {
	if err := r.db.Where(ids).Delete(&models.PostsFishType{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepo) DeletePost(ctx context.Context, id int64) error {
	if err := r.db.Delete(&models.Post{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepo) GetApplyPost(ctx context.Context, id int64) (*models.ApplyPost, error) {
	a := &models.ApplyPost{}
	if err := r.db.Take(a, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = status.Errorf(codes.NotFound, "apply_post with id='%d' is not found", id)
		}
		return nil, err
	}
	return a, nil
}

func (r *postRepo) ListApplyPostsByPostID(ctx context.Context, pID int64) ([]*models.ApplyPost, error) {
	applyPosts := []*models.ApplyPost{}
	if err := r.db.Model(&models.Post{ID: pID}).Related(&applyPosts).Error; err != nil {
		return nil, err
	}
	return applyPosts, nil
}

func (r *postRepo) ListApplyPostsByUserID(ctx context.Context, uID int64) ([]*models.ApplyPost, error) {
	applyPosts := []*models.ApplyPost{}
	if err := r.db.Where("user_id = ?", uID).Find(&applyPosts).Error; err != nil {
		return nil, err
	}
	return applyPosts, nil
}

func (r *postRepo) ListApplyPosts(ctx context.Context, a *models.ApplyPost) ([]*models.ApplyPost, error) {
	applyPosts := []*models.ApplyPost{}
	if err := r.db.Where(a).Find(&applyPosts).Error; err != nil {
		return nil, err
	}

	return applyPosts, nil
}

func (r *postRepo) BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error) {
	applyPosts := []*models.ApplyPost{}
	if err := r.db.Where("post_id IN (?)", postIDs).Order("created_at DESC").Find(&applyPosts).Error; err != nil {
		return nil, err
	}
	return applyPosts, nil
}

func (r *postRepo) CreateApplyPost(ctx context.Context, a *models.ApplyPost) error {
	if err := r.db.Create(a).Error; err != nil {
		e, ok := err.(*mysql.MySQLError)
		if ok {
			if e.Number == 1062 {
				err = status.Error(codes.AlreadyExists, err.Error())
			}
		}
		return err
	}
	return nil
}

func (r *postRepo) DeleteApplyPost(ctx context.Context, id int64) error {
	if err := r.db.Delete(&models.ApplyPost{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
