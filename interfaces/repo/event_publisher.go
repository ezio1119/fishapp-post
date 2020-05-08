package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ezio1119/fishapp-post/models"
)

func createOutboxTX(ctx context.Context, tx *sql.Tx, o *models.Outbox) error {
	query := `INSERT outbox SET event_type=?, event_data=?, aggregate_id=?, aggregate_type=?`
	res, err := tx.ExecContext(ctx, query, o.EventType, o.EventData, o.AggregateID, o.AggregateType)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if int(rowCnt) != 1 {
		return fmt.Errorf("expected %d row affected, got %d rows affected", 1, rowCnt)
	}
	return nil
}
