package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postRepository struct {
	conn *sql.DB
}

func NewPostRepository(conn *sql.DB) post.Repository {
	return &postRepository{conn}
}

func (r *postRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
	rows, err := r.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*models.Post, 0)
	for rows.Next() {
		p := new(models.Post)
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.UserID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *postRepository) GetList(ctx context.Context, datetime time.Time, num int64) ([]*models.Post, error) {
	query := `SELECT id, title, content, user_id, created_at, updated_at
							FROM posts WHERE created_at > ? ORDER BY created_at DESC LIMIT ? `
	res, err := r.fetch(ctx, query, datetime, num)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *postRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `SELECT id, title, content, user_id, created_at, updated_at
  						FROM posts WHERE id = ?`
	list, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Post with id='%d' is not found", id))
	}
	res := list[0]
	return res, nil
}
func (r *postRepository) Create(ctx context.Context, p *models.Post) error {
	query := `INSERT posts SET title=?, content=?, user_id=?, created_at=?, updated_at=?`
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, p.Title, p.Content, p.UserID, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = lastID
	return nil
}

func (r *postRepository) Update(ctx context.Context, p *models.Post) error {
	query := `UPDATE posts SET title=?, content=?, updated_at=? WHERE id = ?`
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	result, err := stmt.ExecContext(ctx, p.Title, p.Content, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return status.Error(codes.Unknown, fmt.Sprintf("Weird  Behaviour. Total Affected: %d", rows))
	}
	return nil
}

func (r *postRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return status.Error(codes.Unknown, fmt.Sprintf("Weird  Behaviour. Total Affected: %d", rows))
	}
	return nil
}
