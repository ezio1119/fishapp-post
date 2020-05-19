package interactor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
)

type sagaReplyInteractor struct {
	sagaRepo   repo.SagaInstanceRepo
	outboxRepo repo.OutboxRepo
}

func NewSagaReplyInteractor(sr repo.SagaInstanceRepo, or repo.OutboxRepo) SagaReplyInteractor {
	return &sagaReplyInteractor{sr, or}
}

type SagaReplyInteractor interface {
	RoomCreated(ctx context.Context, sagaID string) error
}

func (i *sagaReplyInteractor) RoomCreated(ctx context.Context, sagaID string) error {
	sagaIn, err := i.sagaRepo.GetSagaInstance(ctx, sagaID)
	if err != nil {
		return err
	}

	p := &models.Post{}
	if err := json.Unmarshal(sagaIn.SagaData, p); err != nil {
		return err
	}
	saga := newCreatePostSagaState(p, i.outboxRepo, sagaID)
	fmt.Println(saga.FSM.Current())
	saga.FSM.SetState(sagaIn.CurrentState)
	fmt.Println(saga.FSM.Current())
	return nil
	// CreatePostSagaStateをゲットする
	//
}
