package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
)

type postsFishTypeRepo struct {
	db *sql.DB
}

func NewPostsFishTypeRepo(db *sql.DB) repo.PostsFishTypeRepo {
	return &postsFishTypeRepo{db}
}

func (r *postsFishTypeRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.PostsFishType, error) {
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

func (r *postsFishTypeRepo) ListByPostID(ctx context.Context, pID int64) ([]*models.PostsFishType, error) {
	query := `SELECT id, post_id, fish_type_id, created_at, updated_at
						FROM posts_fish_types
						WHERE post_id = ?`
	return r.fetch(ctx, query, pID)
}

func (r *postsFishTypeRepo) BatchCreate(ctx context.Context, pID int64, fIDs []int64, now time.Time) error {
	query := `INSERT INTO posts_fish_types(post_id, fish_type_id, created_at, updated_at)
						VALUES (?, ?, ?, ?)` + strings.Repeat(", (?, ?, ?, ?)", len(fIDs)-1)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	args := make([]interface{}, len(fIDs)*4)
	for i, fID := range fIDs {
		args[i*4] = pID
		args[i*4+1] = fID
		args[i*4+2] = now
		args[i*4+3] = now
	}

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if int(rowCnt) != len(fIDs) {
		return fmt.Errorf("expected %d row affected, got %d rows affected", len(fIDs), rowCnt)
	}
	return nil
}

func (r *postsFishTypeRepo) DeleteByPostID(ctx context.Context, postID int64) error {
	query := "DELETE FROM posts_fish_types WHERE post_id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, postID)
	if err != nil {
		return err
	}
	return nil
}
