package controllers

import (
	"context"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

type postController struct {
	postInteractor post.Usecase
}

func NewPostController(pu post.Usecase) post_grpc.PostServiceServer {
	return &postController{pu}
}

func (c *postController) CreatePost(ctx context.Context, in *post_grpc.CreatePostReq) (*post_grpc.CreatePostRes, error) {
	res, err := c.postInteractor.CreatePost(ctx, &models.Post{
		Title:   in.Title,
		Content: in.Content,
		UserID:  in.UserId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *postController) GetPostByID(ctx context.Context, in *post_grpc.GetPostByIDReq) (*post_grpc.GetPostByIDRes, error) {
	res, err := c.postInteractor.GetPostByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *postController) GetListPosts(ctx context.Context, in *post_grpc.GetListPostsReq) (*post_grpc.GetListPostsRes, error) {
	res, err := c.postInteractor.GetListPosts(ctx, in.Num)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *postController) UpdatePost(ctx context.Context, in *post_grpc.UpdatePostReq) (*post_grpc.UpdatePostRes, error) {
	res, err := c.postInteractor.UpdatePost(ctx, &models.Post{
		ID:      in.Id,
		Title:   in.Title,
		Content: in.Content,
	}, in.UserId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *postController) DeletePost(ctx context.Context, in *post_grpc.DeletePostReq) (*post_grpc.DeletePostRes, error) {
	res, err := c.postInteractor.DeletePost(ctx, in.Id, in.UserId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *postController) CreateApplyPost(ctx context.Context, in *post_grpc.CreateApplyPostReq) (*post_grpc.CreateApplyPostRes, error) {
	res, err := c.postInteractor.CreateApplyPost(ctx, &models.ApplyPost{
		PostID: in.PostId,
		UserID: in.UserId,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (c *postController) DeleteApplyPost(ctx context.Context, in *post_grpc.DeleteApplyPostReq) (*post_grpc.DeleteApplyPostRes, error) {
	res, err := c.postInteractor.DeleteApplyPost(ctx, in.Id, in.UserId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
