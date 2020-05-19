package controllers

import (
	"context"

	"github.com/ezio1119/fishapp-post/usecase/interactor"
)

type sagaReplyController struct {
	sagaReplyInteractor interactor.SagaReplyInteractor
}

type SagaReplyController interface {
	RoomCreated(ctx context.Context, sagaID string) error
	CreateRoomFailed()
}

func NewSagaReplyController(i interactor.SagaReplyInteractor) SagaReplyController {
	return &sagaReplyController{i}
}

func (c *sagaReplyController) RoomCreated(ctx context.Context, sagaID string) error {
	return c.sagaReplyInteractor.RoomCreated(ctx, sagaID)
}

func (c *sagaReplyController) CreateRoomFailed() {}
