package repo

import (
	"context"
	"database/sql"

	"github.com/ezio1119/fishapp-post/interfaces/controllers/event"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"github.com/golang/protobuf/ptypes"
)

type outboxRepo struct {
	db *sql.DB
}

func NewOutboxRepo(db *sql.DB) repo.OutboxRepo {
	return &outboxRepo{db}
}

func convPostCreatedProto(p *models.Post) (*event.PostCreated, error) {
	cAt, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
	}
	uAt, err := ptypes.TimestampProto(p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mAt, err := ptypes.TimestampProto(p.MeetingAt)
	if err != nil {
		return nil, err
	}
	return &event.PostCreated{
		Id:                p.ID,
		Title:             p.Title,
		Content:           p.Content,
		FishingSpotTypeId: p.FishingSpotTypeID,
		FishTypeIds:       p.FishTypeIDs,
		PrefectureId:      p.PrefectureID,
		MeetingPlaceId:    p.MeetingPlaceID,
		MeetingAt:         mAt,
		MaxApply:          p.MaxApply,
		UserId:            p.UserID,
		CreatedAt:         cAt,
		UpdatedAt:         uAt,
	}, nil
}

// func publishPostCreatedEvent(ctx context.Context, tx *sql.Tx, p *models.Post) error {
// 	pProto, err := convPostCreatedProto(p)
// 	if err != nil {
// 		return err
// 	}
// 	eventData, err := protojson.Marshal(pProto)
// 	if err != nil {
// 		return err
// 	}

// 	if err := createOutboxTX(ctx, tx, &models.Outbox{
// 		EventType:     "create.room",
// 		EventData:     eventData,
// 		AggregateID:   strconv.FormatInt(p.ID, 10),
// 		AggregateType: "chat",
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *outboxRepo) CreateOutbox(ctx context.Context, o *models.Outbox, i *models.SagaInstance) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	query := `INSERT outbox SET id=?, event_type=?, event_data=?, aggregate_id=?, aggregate_type=?`
	res, err := tx.ExecContext(ctx, query, o.ID, o.EventType, o.EventData, o.AggregateID, o.AggregateType)
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

	if err := createSagaInstanceTX(ctx, tx, i); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
