package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type postRepository struct {
	conn *sql.DB
}

func NewPostRepository(conn *sql.DB) post.Repository {
	return &postRepository{conn}
}

func (r *postRepository) fetchPost(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
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

func (r *postRepository) GetListPosts(ctx context.Context, num int64) ([]*models.Post, error) {
	query := `SELECT id, title, content, user_id, created_at, updated_at
							FROM posts ORDER BY created_at DESC LIMIT ?`
	list, err := r.fetchPost(ctx, query, num)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Error(codes.NotFound, "no posts found")
	}
	return list, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `SELECT id, title, content, user_id, created_at, updated_at
							FROM posts WHERE id = ?`
	list, err := r.fetchPost(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "post_id=%d is not found", id)
	}
	res := list[0]
	return res, nil
}

func (r *postRepository) CreatePost(ctx context.Context, p *models.Post) error {
	query := `INSERT posts SET title=?, content=?, user_id=?, created_at=?, updated_at=?`
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, p.Title, p.Content, p.UserID, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("%d rows affected", rows)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = lastID
	return nil
}

func (r *postRepository) UpdatePost(ctx context.Context, p *models.Post) error {
	query := `UPDATE posts SET title=?, content=?, updated_at=? WHERE id = ?`
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
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
		return fmt.Errorf("%d rows affected", rows)
	}
	return nil
}

func (r *postRepository) DeletePost(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = ?`
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
		return fmt.Errorf("%d rows affected", rows)
	}
	return nil
}

func (r *postRepository) fetchApplyPost(ctx context.Context, query string, args ...interface{}) ([]*models.ApplyPost, error) {
	rows, err := r.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.ApplyPost, 0)
	for rows.Next() {
		e := new(models.ApplyPost)
		err = rows.Scan(
			&e.ID,
			&e.PostID,
			&e.UserID,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return result, nil
}

func (r *postRepository) CreateApplyPost(ctx context.Context, e *models.ApplyPost) error {
	query := `INSERT apply_post SET post_id=?, user_id=?, created_at=?, updated_at=?`
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, e.PostID, e.UserID, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		e, ok := err.(*mysql.MySQLError)
		if ok {
			if e.Number == 1062 {
				err = status.Error(codes.AlreadyExists, err.Error())
			}
		}
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = lastID
	return nil
}

func (r *postRepository) GetApplyPostByID(ctx context.Context, id int64) (*models.ApplyPost, error) {
	query := `SELECT id, post_id, user_id, created_at, updated_at
							FROM apply_post WHERE id = ?`
	list, err := r.fetchApplyPost(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "apply_post with id='%d' is not found", id)
	}
	res := list[0]
	return res, nil
}

func (r *postRepository) DeleteApplyPost(ctx context.Context, id int64) error {
	query := `DELETE FROM apply_post WHERE id = ?`
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
		return fmt.Errorf("%d rows affected", rows)
	}
	return nil
}

// func (r *entryRepository) GetListByPostID(ctx context.Context, postID int64) ([]*models.Entry, error) {
// 	query := `SELECT id, post_id, user_id, created_at, updated_at
// 	FROM entries WHERE post_id = ? ORDER BY created_at DESC`
// 	res, err := r.fetch(ctx, query, postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }
