package controllers

import (
	"context"

	"github.com/ezio1119/fishapp-post/entry"
	"github.com/ezio1119/fishapp-post/entry/controllers/entry_post_grpc"
	"github.com/ezio1119/fishapp-post/models"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

type entryController struct {
	entryInteractor entry.Usecase
}

func NewEntryController(eu entry.Usecase) entry_post_grpc.EntryServiceServer {
	return &entryController{eu}
}

func (c *entryController) Create(ctx context.Context, in *entry_post_grpc.CreateReq) (*entry_post_grpc.Entry, error) {
	return c.entryInteractor.Create(ctx, &models.Entry{
		PostID: in.PostId,
		UserID: in.UserId,
	})
}

func (c *entryController) Delete(ctx context.Context, in *entry_post_grpc.DeleteReq) (*wrappers.BoolValue, error) {
	if err := c.entryInteractor.Delete(ctx, in.Id, in.UserId); err != nil {
		return nil, err
	}
	return &wrappers.BoolValue{
		Value: true,
	}, nil
}
