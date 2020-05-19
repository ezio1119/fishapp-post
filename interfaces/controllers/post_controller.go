package controllers

import (
	"context"

	"github.com/ezio1119/fishapp-post/interfaces/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/interactor"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (c *postController) CreatePost(ctx context.Context, in *post_grpc.CreatePostReq) (*post_grpc.CreatePostRes, error) {
	mAt, err := ptypes.Timestamp(in.MeetingAt)
	if err != nil {
		return nil, err
	}
	p := &models.Post{
		Title:             in.Title,
		Content:           in.Content,
		FishingSpotTypeID: in.FishingSpotTypeId,
		FishTypeIDs:       in.FishTypeIds,
		PrefectureID:      in.PrefectureId,
		MeetingPlaceID:    in.MeetingPlaceId,
		MeetingAt:         mAt,
		MaxApply:          in.MaxApply,
		UserID:            in.UserId,
	}
	sagaID, err := c.postInteractor.CreatePost(ctx, p)
	if err != nil {
		return nil, err
	}
	pProto, err := convPostProto(p)
	if err != nil {
		return nil, err
	}
	return &post_grpc.CreatePostRes{
		Post:   pProto,
		SagaId: sagaID,
	}, nil
}

func (c *postController) UpdatePost(ctx context.Context, in *post_grpc.UpdatePostReq) (*post_grpc.Post, error) {
	mAt, err := ptypes.Timestamp(in.MeetingAt)
	if err != nil {
		return nil, err
	}
	p := &models.Post{
		ID:                in.Id,
		Title:             in.Title,
		Content:           in.Content,
		FishingSpotTypeID: in.FishingSpotTypeId,
		FishTypeIDs:       in.FishTypeIds,
		PrefectureID:      in.PrefectureId,
		MeetingPlaceID:    in.MeetingPlaceId,
		MeetingAt:         mAt,
		MaxApply:          in.MaxApply,
	}
	if err := c.postInteractor.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	return convPostProto(p)
}

func (c *postController) DeletePost(ctx context.Context, in *post_grpc.DeletePostReq) (*empty.Empty, error) {
	if err := c.postInteractor.DeletePost(ctx, in.Id); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (c *postController) GetApplyPost(ctx context.Context, in *post_grpc.GetApplyPostReq) (*post_grpc.ApplyPost, error) {
	a, err := c.postInteractor.GetApplyPost(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return convApplyPostProto(a)
}

func (c *postController) ListApplyPosts(ctx context.Context, in *post_grpc.ListApplyPostsReq) (*post_grpc.ListApplyPostsRes, error) {
	if (in.Filter.UserId == 0 && in.Filter.PostId == 0) || (in.Filter.UserId != 0 && in.Filter.PostId != 0) {
		return nil, status.Error(codes.InvalidArgument, "invalid ListApplyPostsReq.Filter.PostId, ListApplyPostsReq.Filter.UserId: value must be set either user_id or post_id")
	}
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

func (c *postController) BatchGetApplyPostsByPostIDs(ctx context.Context, in *post_grpc.BatchGetApplyPostsByPostIDsReq) (*post_grpc.BatchGetApplyPostsByPostIDsRes, error) {
	list, err := c.postInteractor.BatchGetApplyPostsByPostIDs(ctx, in.PostIds)
	if err != nil {
		return nil, err
	}
	listProto, err := convListApplyPostsProto(list)
	if err != nil {
		return nil, err
	}
	return &post_grpc.BatchGetApplyPostsByPostIDsRes{ApplyPosts: listProto}, nil
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
	if err := c.postInteractor.DeleteApplyPost(ctx, in.Id); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
