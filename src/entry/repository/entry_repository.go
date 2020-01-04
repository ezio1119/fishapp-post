package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ezio1119/fishapp-post/entry"
	"github.com/ezio1119/fishapp-post/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type entryRepository struct {
	conn *sql.DB
}

func NewEntryRepository(conn *sql.DB) entry.Repository {
	return &entryRepository{conn}
}

func (r *entryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Entry, error) {
	rows, err := r.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*models.Entry, 0)
	for rows.Next() {
		e := new(models.Entry)
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

func (r *entryRepository) Create(ctx context.Context, e *models.Entry) error {
	query := `INSERT entries SET post_id=?, user_id=?, created_at=?, updated_at=?`
	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(ctx, e.PostID, e.UserID, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = lastID
	return nil
}

func (r *entryRepository) GetByID(ctx context.Context, id int64) (*models.Entry, error) {
	query := `SELECT id, post_id, user_id, created_at, updated_at
							FROM entries WHERE id = ?`
	list, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Entry with id='%d' is not found", id))
	}
	res := list[0]
	return res, nil
}

func (r *entryRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM entries WHERE id = ?"
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
