package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
)

type sagaInstanceRepo struct {
	db *sql.DB
}

func NewSagaInstanceRepo(db *sql.DB) repo.SagaInstanceRepo {
	return &sagaInstanceRepo{db}
}

func createSagaInstanceTX(ctx context.Context, tx *sql.Tx, i *models.SagaInstance) error {
	query := `INSERT saga_instance SET id=?, saga_type=?, saga_data=?, current_state=?`
	res, err := tx.ExecContext(ctx, query, i.ID, i.SagaType, i.SagaData, i.CurrentState)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if int(rowCnt) != 1 {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *sagaInstanceRepo) UpdateSagaInstance(ctx context.Context, s *models.SagaInstance) error {
	return nil
}

func (r *sagaInstanceRepo) GetSagaInstance(ctx context.Context, sagaID string) (*models.SagaInstance, error) {
	query := `SELECT id, saga_type, saga_data, current_state FROM saga_instance WHERE id=?`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	i := &models.SagaInstance{}

	err = stmt.QueryRowContext(ctx, sagaID).Scan(&i.ID, &i.SagaType, &i.SagaData, &i.CurrentState)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("no saga_instance with id %s", sagaID)
	case err != nil:
		return nil, err
	}

	return i, nil
}
