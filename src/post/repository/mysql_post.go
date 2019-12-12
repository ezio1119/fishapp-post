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

type mysqlPostRepository struct {
	Conn *sql.DB
}

// NewMysqlPostRepository will create an object that represent the post.Repository interface
func NewMysqlPostRepository(Conn *sql.DB) post.Repository {
	return &mysqlPostRepository{Conn}
}

func (m *mysqlPostRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*models.Post, 0)
	for rows.Next() {
		p := new(models.Post)
		err = rows.Scan(
			&p.Id,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.UserId,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (m *mysqlPostRepository) GetList(ctx context.Context, datetime time.Time, num int64) ([]*models.Post, error) {
	query := `SELECT id, title, content, created_at, updated_at, user_id
							FROM posts WHERE created_at > ? ORDER BY created_at DESC LIMIT ? `
	res, err := m.fetch(ctx, query, datetime, num)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *mysqlPostRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `SELECT id, title, content, created_at, updated_at, user_id
  						FROM posts WHERE id = ?`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Post with id='%d' is not found", id))
	}
	res := list[0]
	return res, nil
}
func (m *mysqlPostRepository) Create(ctx context.Context, p *models.Post) error {
	query := `INSERT posts SET title=?, content=?, user_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, p.Title, p.Content, p.UserId)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = lastID
	return nil
}

func (m *mysqlPostRepository) Update(ctx context.Context, p *models.Post) error {
	query := `UPDATE posts SET title=?, content=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	result, err := stmt.ExecContext(ctx, p.Title, p.Content, p.Id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return status.Error(codes.Unknown, fmt.Sprintf("No affected rows"))
	}
	if rows > 1 {
		return status.Error(codes.Unknown, "More affected rows than expected")
	}
	return nil
}

func (m *mysqlPostRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM posts WHERE id = ?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
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
	if rows == 0 {
		return status.Error(codes.Unknown, fmt.Sprintf("No affected rows"))
	}
	if rows > 1 {
		return status.Error(codes.Unknown, "More affected rows than expected")
	}
	return nil
}
