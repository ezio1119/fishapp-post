package interactor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/interactor/saga"
	"github.com/ezio1119/fishapp-post/usecase/repo"
)

type sagaReplyInteractor struct {
	createPostSagaManager *saga.CreatePostSagaManager
	sagaInstanceRepo      repo.SagaInstanceRepo
}

func NewSagaReplyInteractor(m *saga.CreatePostSagaManager, sr repo.SagaInstanceRepo) SagaReplyInteractor {
	return &sagaReplyInteractor{m, sr}
}

type SagaReplyInteractor interface {
	RoomCreated(ctx context.Context, sagaID string) error
	CreateRoomFailed(ctx context.Context, sagaID string, errMsg string) error
}

func (i *sagaReplyInteractor) RoomCreated(ctx context.Context, sagaID string) error {
	sagaIn, err := i.sagaInstanceRepo.GetSagaInstance(ctx, sagaID)
	if err != nil {
		return err
	}
	p := &models.Post{}
	if err := json.Unmarshal(sagaIn.SagaData, p); err != nil {
		return err
	}

	state := saga.NewCreatePostSagaState(p, sagaIn.CurrentState, sagaID)
	s := i.createPostSagaManager.NewCreatePostSagaManager(state)

	if err := s.FSM.Event("ApprovePost", ctx); err != nil {
		return err
	}
	return nil

}

func (i *sagaReplyInteractor) CreateRoomFailed(ctx context.Context, sagaID string, errMsg string) error {
	log.Printf("error: %s", errMsg)
	sagaIn, err := i.sagaInstanceRepo.GetSagaInstance(ctx, sagaID)
	if err != nil {
		return err
	}
	p := &models.Post{}
	if err := json.Unmarshal(sagaIn.SagaData, p); err != nil {
		return err
	}
	state := saga.NewCreatePostSagaState(p, sagaIn.CurrentState, sagaID)
	s := i.createPostSagaManager.NewCreatePostSagaManager(state)
	if err := s.FSM.Event("RejectPost", ctx); err != nil {
		return err
	}

	return nil
}
