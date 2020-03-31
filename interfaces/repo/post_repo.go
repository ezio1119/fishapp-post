package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) repo.PostRepo {
	return &postRepo{db}
}

func (r *postRepo) fetchPosts(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	result := make([]*models.Post, 0)
	for rows.Next() {
		p := new(models.Post)
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.FishingSpotTypeID,
			&p.PrefectureID,
			&p.MeetingPlaceID,
			&p.MeetingAt,
			&p.MaxApply,
			&p.UserID,
			&p.UpdatedAt,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *postRepo) fetchPostsFishTypes(ctx context.Context, query string, args ...interface{}) ([]*models.PostsFishType, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	result := make([]*models.PostsFishType, 0)
	for rows.Next() {
		f := new(models.PostsFishType)
		err = rows.Scan(
			&f.ID,
			&f.PostID,
			&f.FishTypeID,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}

	return result, nil
}

func (r *postRepo) CreatePost(ctx context.Context, p *models.Post) error {
	query := `INSERT INTO posts(title, content, fishing_spot_type_id, prefecture_id, meeting_place_id, meeting_at, max_apply, user_id, updated_at, created_at)
						VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, p.Title, p.Content, p.FishingSpotTypeID, p.PrefectureID, p.MeetingPlaceID, p.MeetingAt, p.MaxApply, p.UserID, p.UpdatedAt, p.CreatedAt)
	if err != nil {
		e, ok := err.(*mysql.MySQLError)
		if ok {
			if e.Number == 1062 {
				err = status.Error(codes.AlreadyExists, err.Error())
			}
		}
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = lastID
	return nil
}

func (r *postRepo) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `SELECT id, title, content, fishing_spot_type_id, prefecture_id, meeting_place_id, meeting_at, max_apply, user_id, updated_at, created_at
						FROM posts
						WHERE id = ?`

	list, err := r.fetchPosts(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "post with id='%d' is not found", id)
	}
	return list[0], nil
}

func (r *postRepo) GetPostCanApply(ctx context.Context, id int64) (*models.Post, error) {
	// p := &models.Post{}
	// if err := r.db.Where("id = ? AND max_apply > ?", id, r.db.Table("apply_posts").Select("COUNT(*)").Where("post_id = ?", id).QueryExpr()).Error; err != nil {
	// 	if gorm.IsRecordNotFoundError(err) {
	// 		err = status.Errorf(codes.InvalidArgument, "cannot apply post_id ='%d' because upper limit", id)
	// 	}
	// 	return nil, err
	// }
	// return p, nil
	panic("not implemented")
}

func (r *postRepo) List(ctx context.Context, p *models.Post, num int64, cursor int64, f *models.PostFilter) ([]*models.Post, error) {
	sq := sq.Select("id, title, content, fishing_spot_type_id, prefecture_id, meeting_place_id, meeting_at, max_apply, user_id, updated_at, created_at").
		From("posts").
		GroupBy("posts.id").
		Limit(uint64(num))

	if p.FishingSpotTypeID != 0 {
		sq = sq.Where("fishing_spot_type_id = ?", p.FishingSpotTypeID)
	}
	if p.PrefectureID != 0 {
		sq = sq.Where("prefecture_id = ?", p.PrefectureID)
	}
	if p.UserID != 0 {
		sq = sq.Where("user_id = ?", p.UserID)
	}
	if f.CanApply {
		sq = sq.LeftJoin("apply_posts ON posts.id = apply_posts.post_id").
			Having("count(apply_posts.id) < posts.max_apply")
	}
	if f.FishTypeIDs != nil {
		sq = sq.Join("posts_fish_types ON posts.id = posts_fish_types.post_id").
			Where("posts_fish_types.fish_type_id IN(?)", f.FishTypeIDs).
			Having("count(posts_fish_types.fish_type_id) = ?", len(f.FishTypeIDs))
	}
	if !f.MeetingAtFrom.IsZero() && !f.MeetingAtTo.IsZero() {
		sq = sq.Where("meeting_at BETWEEN ? AND ?", f.MeetingAtFrom, f.MeetingAtTo)
	}

	if cursor != 0 {
		switch f.SortBy {
		case models.SortByID:
			if f.OrderBy == models.OrderByAsc {
				sq = sq.Where("posts.id > ?", cursor).
					OrderBy("id asc")
			}
			if f.OrderBy == models.OrderByDesc {
				sq = sq.Where("posts.id < ?", cursor).
					OrderBy("id desc")
			}
		// meeting_atはユニークではないため、同じ値の場合を考えidでも絞り込む
		case models.SortByMeetingAt:
			p, err := r.GetPostByID(ctx, cursor)
			if err != nil {
				return nil, err
			}
			switch f.OrderBy {
			case models.OrderByAsc:
				sq = sq.Where("meeting_at >= ?", p.MeetingAt).
					Where("meeting_at > ? or posts.id > ?", p.MeetingAt, cursor).
					OrderBy("meeting_at asc, id asc")
			case models.OrderByDesc:
				sq = sq.Where("meeting_at <= ?", p.MeetingAt).
					Where("meeting_at < ? or posts.id < ?", p.MeetingAt, cursor).
					OrderBy("meeting_at desc, id desc")
			}
		}
	}
	if cursor == 0 {
		sq = sq.OrderBy(fmt.Sprintf("%s %s", f.SortBy, f.OrderBy))
	}
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}
	fmt.Println(query)
	return r.fetchPosts(ctx, query, args...)
}

func (r *postRepo) UpdatePost(ctx context.Context, p *models.Post) error {
	query := `UPDATE posts SET title=?, content=?, fishing_spot_type_id=?, prefecture_id=?, meeting_place_id=?, meeting_at=?, max_apply=?, updated_at=?
						WHERE id = ?`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}
	res, err := stmt.ExecContext(ctx, p.Title, p.Content, p.FishingSpotTypeID, p.PrefectureID, p.MeetingPlaceID, p.MeetingAt, p.MaxApply, p.UpdatedAt)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return fmt.Errorf("expected %d row affected, got %d rows affected", 1, rowCnt)
	}
	return nil
}

func (r *postRepo) DeletePostsFishTypesByPostID(ctx context.Context, pID int64) error {
	// return r.db.Where("post_id = ?", pID).Delete(&models.PostsFishType{}).Error
	panic("not implemented")
}

func (r *postRepo) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {

		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowCnt != 1 {
		err = fmt.Errorf("expected %d row affected, got %d rows affected", 1, rowCnt)
		return err
	}

	return nil
}

func (r *postRepo) GetApplyPost(ctx context.Context, pID int64, uID int64) (*models.ApplyPost, error) {
	// a := &models.ApplyPost{}
	// if err := r.db.Where("user_id = ? AND post_id = ?", uID, pID).Take(a).Error; err != nil {
	// 	if gorm.IsRecordNotFoundError(err) {
	// 		err = status.Errorf(codes.NotFound, "apply_post with post_id='%d' user_id='%d' is not found", pID, uID)
	// 	}
	// 	return nil, err
	// }
	// return a, nil
	panic("not implemented")
}

func (r *postRepo) ListApplyPostsByPostID(ctx context.Context, pID int64) ([]*models.ApplyPost, error) {
	// applyPosts := []*models.ApplyPost{}
	// if err := r.db.Model(&models.Post{ID: pID}).Related(&applyPosts).Error; err != nil {
	// 	return nil, err
	// }
	// return applyPosts, nil
	panic("not implemented")
}

func (r *postRepo) ListApplyPostsByUserID(ctx context.Context, uID int64) ([]*models.ApplyPost, error) {
	// applyPosts := []*models.ApplyPost{}
	// if err := r.db.Where("user_id = ?", uID).Find(&applyPosts).Error; err != nil {
	// 	return nil, err
	// }
	// return applyPosts, nil
	panic("not implemented")
}

func (r *postRepo) ListApplyPosts(ctx context.Context, a *models.ApplyPost) ([]*models.ApplyPost, error) {
	// applyPosts := []*models.ApplyPost{}
	// if err := r.db.Where(a).Find(&applyPosts).Error; err != nil {
	// 	return nil, err
	// }

	// return applyPosts, nil
	panic("not implemented")
}

func (r *postRepo) BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error) {
	// applyPosts := []*models.ApplyPost{}
	// if err := r.db.Where("post_id IN (?)", postIDs).Order("created_at DESC").Find(&applyPosts).Error; err != nil {
	// 	return nil, err
	// }
	// return applyPosts, nil
	panic("not implemented")
}

func (r *postRepo) CreateApplyPost(ctx context.Context, a *models.ApplyPost) error {
	// if err := r.db.Create(a).Error; err != nil {
	// 	e, ok := err.(*mysql.MySQLError)
	// 	if ok {
	// 		if e.Number == 1062 {
	// 			err = status.Error(codes.AlreadyExists, err.Error())
	// 		}
	// 	}
	// 	return err
	// }
	// return nil
	panic("not implemented")
}

func (r *postRepo) DeleteApplyPost(ctx context.Context, pID int64, uID int64) error {
	// return r.db.Where("user_id = ? AND post_id = ?", uID, pID).Delete(&models.ApplyPost{}).Error
	panic("not implemented")
}

func (r *postRepo) BatchListPostsFishTypesByPostIDs(ctx context.Context, pIDs []int64) ([]*models.PostsFishType, error) {
	query := `SELECT id, post_id, fish_type_id, created_at, updated_at
						FROM posts_fish_types
						WHERE post_id IN(?` + strings.Repeat(",?", len(pIDs)-1) + ")"
	args := make([]interface{}, len(pIDs))
	for i, pID := range pIDs {
		args[i] = pID
	}
	return r.fetchPostsFishTypes(ctx, query, args...)
}

func (r *postRepo) ListFishTypeIDsByPostID(ctx context.Context, postID int64) ([]int64, error) {
	query := `SELECT fish_type_id
						FROM posts_fish_types
						WHERE post_id = ?`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
