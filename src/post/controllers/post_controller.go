package controllers

import (
	"context"

	"github.com/golang/protobuf/ptypes"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type postController struct {
	postInteractor post.Usecase
}

func NewPostController(pu post.Usecase) post_grpc.PostServiceServer {
	return &postController{pu}
}

func (c *postController) Create(ctx context.Context, in *post_grpc.CreateReq) (*post_grpc.Post, error) {
	post, err := c.postInteractor.Create(ctx, &models.Post{
		Title:   in.Title,
		Content: in.Content,
		UserID:  in.UserId,
	})
	if err != nil {
		return nil, models.WrapOnGrpcErr(err)
	}
	return post, nil
}

func (c *postController) GetByID(ctx context.Context, in *post_grpc.ID) (*post_grpc.Post, error) {
	post, err := c.postInteractor.GetByID(ctx, in.Id)
	if err != nil {
		return nil, models.WrapOnGrpcErr(err)
	}
	return post, nil
}

func (c *postController) GetList(ctx context.Context, in *post_grpc.ListReq) (*post_grpc.ListPost, error) {
	datetime, err := ptypes.Timestamp(in.Datetime)
	if err != nil {
		return nil, models.WrapOnGrpcErr(err)
	}
	post, err := c.postInteractor.GetList(ctx, datetime, in.Num)
	if err != nil {
		return nil, models.WrapOnGrpcErr(err)
	}
	return post, nil
}

func (c *postController) Update(ctx context.Context, in *post_grpc.UpdateReq) (*post_grpc.Post, error) {
	post, err := c.postInteractor.Update(ctx, &models.Post{
		ID:      in.Id,
		Title:   in.Title,
		Content: in.Content,
	}, in.UserId)
	if err != nil {
		return nil, models.WrapOnGrpcErr(err)
	}
	return post, nil
}

func (c *postController) Delete(ctx context.Context, in *post_grpc.DeleteReq) (*wrappers.BoolValue, error) {
	if err := c.postInteractor.Delete(ctx, in.Id, in.UserId); err != nil {
		return nil, models.WrapOnGrpcErr(err)
	}
	return &wrappers.BoolValue{
		Value: true,
	}, nil
}
