package repo

import (
	"context"
	"database/sql"
	"log"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
)

type applyPostRepo struct {
	db *sql.DB
}

func NewApplyPostRepo(db *sql.DB) repo.ApplyPostRepo {
	return &applyPostRepo{db}
}

func (r *applyPostRepo) GetApplyPost(ctx context.Context, pID int64, uID int64) (*models.ApplyPost, error) {
	panic("not implemented")
}
func (r *applyPostRepo) BatchGetApplyPostsByPostIDs(ctx context.Context, postIDs []int64) ([]*models.ApplyPost, error) {
	panic("not implemented")
}
func (r *applyPostRepo) ListApplyPosts(ctx context.Context, applyPost *models.ApplyPost) ([]*models.ApplyPost, error) {
	panic("not implemented")
}
func (r *applyPostRepo) CreateApplyPost(ctx context.Context, p *models.ApplyPost) error {
	panic("not implemented")
}
func (r *applyPostRepo) DeleteApplyPost(ctx context.Context, pID int64, uID int64) error {
	panic("not implemented")
}

func (r *applyPostRepo) CountByPostID(ctx context.Context, postID int64) (int64, error) {
	query := `SELECT COUNT(*)
					 FROM apply_posts
					 WHERE post_id = ?`
	var cnt int64
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return 0, err
		}
	}

	return cnt, nil
}
