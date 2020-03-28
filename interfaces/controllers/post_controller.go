package controllers

import (
	"context"

	"github.com/ezio1119/fishapp-post/interfaces/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/interactor"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
)

type postController struct {
	postInteractor interactor.PostInteractor
}

func NewPostController(pu interactor.PostInteractor) *postController {
	return &postController{pu}
}

func (c *postController) GetPost(ctx context.Context, in *post_grpc.GetPostReq) (*post_grpc.Post, error) {
	p, err := c.postInteractor.GetPost(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return convPostProto(p)
}

func (c *postController) ListPosts(ctx context.Context, in *post_grpc.ListPostsReq) (*post_grpc.ListPostsRes, error) {
	f, err := convPostFilter(in.Filter)
	if err != nil {
		return nil, err
	}
	list, nextToken, err := c.postInteractor.ListPosts(ctx, &models.Post{
		FishingSpotTypeID: in.Filter.FishingSpotTypeId,
		PrefectureID:      in.Filter.PrefectureId,
		UserID:            in.Filter.UserId,
	}, in.PageSize, in.PageToken, f)
	if err != nil {
		return nil, err
	}
	listProto, err := convListPostsProto(list)
	if err != nil {
		return nil, err
	}

	return &post_grpc.ListPostsRes{Posts: listProto, NextPageToken: nextToken}, nil
}

func (c *postController) CreatePost(ctx context.Context, in *post_grpc.CreatePostReq) (*post_grpc.Post, error) {
	t, err := ptypes.Timestamp(in.MeetingAt)
	if err != nil {
		return nil, err
	}
	pfishTypes := make([]*models.PostsFishType, len(in.FishTypeIds))
	for i, f := range in.FishTypeIds {
		pfishTypes[i] = &models.PostsFishType{FishTypeID: f}
	}
	p := &models.Post{
		Title:             in.Title,
		Content:           in.Content,
		FishingSpotTypeID: in.FishingSpotTypeId,
		PostsFishTypes:    pfishTypes,
		PrefectureID:      in.PrefectureId,
		MeetingPlaceID:    in.MeetingPlaceId,
		MeetingAt:         t,
		MaxApply:          in.MaxApply,
		UserID:            in.UserId,
	}
	if err := c.postInteractor.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	return convPostProto(p)
}

func (c *postController) UpdatePost(ctx context.Context, in *post_grpc.UpdatePostReq) (*post_grpc.Post, error) {
	t, err := ptypes.Timestamp(in.MeetingAt)
	if err != nil {
		return nil, err
	}
	pfishTypes := make([]*models.PostsFishType, len(in.FishTypeIds))
	for i, f := range in.FishTypeIds {
		pfishTypes[i] = &models.PostsFishType{FishTypeID: f}
	}
	p := &models.Post{
		ID:                in.Id,
		Title:             in.Title,
		Content:           in.Content,
		FishingSpotTypeID: in.FishingSpotTypeId,
		PostsFishTypes:    pfishTypes,
		PrefectureID:      in.PrefectureId,
		MeetingPlaceID:    in.MeetingPlaceId,
		MeetingAt:         t,
		MaxApply:          in.MaxApply,
	}
	post, err := c.postInteractor.UpdatePost(ctx, p)
	if err != nil {
		return nil, err
	}
	return convPostProto(post)
}

func (c *postController) DeletePost(ctx context.Context, in *post_grpc.DeletePostReq) (*empty.Empty, error) {
	if err := c.postInteractor.DeletePost(ctx, in.Id); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (c *postController) GetApplyPost(ctx context.Context, in *post_grpc.GetApplyPostReq) (*post_grpc.ApplyPost, error) {
	a, err := c.postInteractor.GetApplyPost(ctx, in.PostId, in.UserId)
	if err != nil {
		return nil, err
	}
	return convApplyPostProto(a)
}

func (c *postController) ListApplyPosts(ctx context.Context, in *post_grpc.ListApplyPostsReq) (*post_grpc.ListApplyPostsRes, error) {
	list, err := c.postInteractor.ListApplyPosts(ctx, &models.ApplyPost{
		UserID: in.Filter.UserId,
		PostID: in.Filter.PostId,
	})
	if err != nil {
		return nil, err
	}
	listProto, err := convListApplyPostsProto(list)
	if err != nil {
		return nil, err
	}
	return &post_grpc.ListApplyPostsRes{ApplyPosts: listProto}, nil
}

func (c *postController) CreateApplyPost(ctx context.Context, in *post_grpc.CreateApplyPostReq) (*post_grpc.ApplyPost, error) {
	a := &models.ApplyPost{
		PostID: in.PostId,
		UserID: in.UserId,
	}
	err := c.postInteractor.CreateApplyPost(ctx, a)
	if err != nil {
		return nil, err
	}
	return convApplyPostProto(a)
}

func (c *postController) DeleteApplyPost(ctx context.Context, in *post_grpc.DeleteApplyPostReq) (*empty.Empty, error) {
	if err := c.postInteractor.DeleteApplyPost(ctx, in.PostId, in.UserId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
