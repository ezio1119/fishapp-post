package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
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
		if err := rows.Scan(&f.PostID, &f.FishTypeID); err != nil {
			return nil, err
		}
		result = append(result, f)
	}

	return result, nil
}

func (r *postRepo) fillPostWithFishTypeIDs(ctx context.Context, p *models.Post) error {
	query := `SELECT fish_type_id
           	FROM posts_fish_types
						 WHERE post_id = ?`
	rows, err := r.db.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
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
			return err
		}
		ids = append(ids, id)
	}
	p.FishTypeIDs = ids

	return nil
}

func (r *postRepo) fillListPostsWithFishTypeIDs(ctx context.Context, posts []*models.Post) error {
	query := `SELECT post_id, fish_type_id
                        FROM posts_fish_types
                        WHERE post_id IN(?` + strings.Repeat(",?", len(posts)-1) + ")"

	args := make([]interface{}, len(posts))
	for i, p := range posts {
		args[i] = p.ID
	}
	fishes, err := r.fetchPostsFishTypes(ctx, query, args...)
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

func (r *postRepo) batchCreatePostsFishTypesTX(ctx context.Context, tx *sql.Tx, p *models.Post) error {
	query := `INSERT INTO posts_fish_types(post_id, fish_type_id, created_at, updated_at)
                     VALUES (?, ?, ?, ?)` + strings.Repeat(", (?, ?, ?, ?)", len(p.FishTypeIDs)-1)

	args := make([]interface{}, len(p.FishTypeIDs)*4)
	for i, fID := range p.FishTypeIDs {
		args[i*4] = p.ID
		args[i*4+1] = fID
		args[i*4+2] = p.CreatedAt
		args[i*4+3] = p.UpdatedAt
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if int(rowCnt) != len(p.FishTypeIDs) {
		return fmt.Errorf("expected %d row affected, got %d rows affected", len(p.FishTypeIDs), rowCnt)
	}

	return nil
}

func (r *postRepo) deletePostsFishTypesByPostIDTX(ctx context.Context, tx *sql.Tx, pID int64) error {
	query := "DELETE FROM posts_fish_types WHERE post_id = ?"
	res, err := tx.ExecContext(ctx, query, pID)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *postRepo) CreatePost(ctx context.Context, p *models.Post) error {
	query := `INSERT posts SET title=?, content=?, fishing_spot_type_id=?, prefecture_id=?, meeting_place_id=?, meeting_at=?, max_apply=?, user_id=?, updated_at=?, created_at=?`
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err != nil {
		log.Fatal(err)
	}
	res, err := tx.ExecContext(ctx, query, p.Title, p.Content, p.FishingSpotTypeID, p.PrefectureID, p.MeetingPlaceID, p.MeetingAt, p.MaxApply, p.UserID, p.UpdatedAt, p.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if int(rowCnt) != 1 {
		tx.Rollback()
		return fmt.Errorf("expected %d row affected, got %d rows affected", 1, rowCnt)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	p.ID = lastID

	// postsFishTypesをbatch insert
	if err := r.batchCreatePostsFishTypesTX(ctx, tx, p); err != nil {
		tx.Rollback()
		return err
	}

	pProto, err := convPostCreatedProto(p)
	if err != nil {
		tx.Rollback()
		return err
	}
	eventData, err := protojson.Marshal(pProto)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := createOutboxTX(ctx, tx, &models.Outbox{
		EventType:     "post.created",
		EventData:     eventData,
		AggregateID:   strconv.FormatInt(p.ID, 10),
		AggregateType: "post",
	}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
	if err := r.fillPostWithFishTypeIDs(ctx, list[0]); err != nil {
		return nil, err
	}
	return list[0], nil
}

func (r *postRepo) ListPosts(ctx context.Context, p *models.Post, num int64, cursor int64, f *models.PostFilter) ([]*models.Post, error) {
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
	posts, err := r.fetchPosts(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if err := r.fillListPostsWithFishTypeIDs(ctx, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) UpdatePost(ctx context.Context, p *models.Post) error {
	query := `UPDATE posts SET title=?, content=?, fishing_spot_type_id=?, prefecture_id=?, meeting_place_id=?, meeting_at=?, max_apply=?, updated_at=?
                        WHERE id = ?`
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}
	res, err := tx.ExecContext(ctx, query, p.Title, p.Content, p.FishingSpotTypeID, p.PrefectureID, p.MeetingPlaceID, p.MeetingAt, p.MaxApply, p.UpdatedAt, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowCnt != 1 {
		tx.Rollback()
		return fmt.Errorf("expected %d row affected, got %d rows affected", 1, rowCnt)
	}

	if err := r.deletePostsFishTypesByPostIDTX(ctx, tx, p.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.batchCreatePostsFishTypesTX(ctx, tx, p); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *postRepo) DeletePost(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowCnt != 1 {
		tx.Rollback()
		return fmt.Errorf("expected %d row affected, got %d rows affected", 1, rowCnt)
	}

	query = "DELETE FROM posts_fish_types WHERE post_id = ?"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
