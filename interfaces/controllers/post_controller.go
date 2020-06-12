package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/pb"
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

func (c *postController) GetPost(ctx context.Context, in *pb.GetPostReq) (*pb.Post, error) {
	p, err := c.postInteractor.GetPost(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return convPostProto(p)
}

func (c *postController) ListPosts(ctx context.Context, in *pb.ListPostsReq) (*pb.ListPostsRes, error) {
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

	return &pb.ListPostsRes{Posts: listProto, NextPageToken: nextToken}, nil
}

func (c *postController) CreatePost(stream pb.PostService_CreatePostServer) error {

	ctx := stream.Context()
	in := &pb.CreatePostReqInfo{}

	imageBufs := map[int64]*bytes.Buffer{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			goto END
		}
		if err != nil {
			return err
		}

		switch x := req.Data.(type) {
		case *pb.CreatePostReq_Info:
			*in = *x.Info
		case *pb.CreatePostReq_ImageChunk:

			chunkData := x.ImageChunk.ChunkData
			chunkNum := x.ImageChunk.ChunkNum

			if imageBufs[chunkNum] == nil {
				imageBufs[chunkNum] = &bytes.Buffer{}
			}

			if _, err := imageBufs[chunkNum].Write(chunkData); err != nil {
				return err
			}

		default:
			return fmt.Errorf("CreatePostReq.Request has unexpected type %T", x)
		}
	}

END:

	mAt, err := ptypes.Timestamp(in.MeetingAt)
	if err != nil {
		return err
	}

	p := &models.Post{
		Title:             in.Title,
		Content:           in.Content,
		FishingSpotTypeID: in.FishingSpotTypeId,
		PostsFishTypes:    models.ConvPostsFishTypes(in.FishTypeIds),
		PrefectureID:      in.PrefectureId,
		MeetingPlaceID:    in.MeetingPlaceId,
		MeetingAt:         mAt,
		MaxApply:          in.MaxApply,
		UserID:            in.UserId,
	}

	sagaID, err := c.postInteractor.CreatePost(ctx, p, imageBufs)
	if err != nil {
		return err
	}
	pProto, err := convPostProto(p)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.CreatePostRes{
		Post:   pProto,
		SagaId: sagaID,
	})
}

func (c *postController) UpdatePost(stream pb.PostService_UpdatePostServer) error {
	ctx := stream.Context()
	in := &pb.UpdatePostReqInfo{}

	imageBufs := map[int64]*bytes.Buffer{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			goto END
		}
		if err != nil {
			return err
		}

		switch x := req.Data.(type) {
		case *pb.UpdatePostReq_Info:
			*in = *x.Info
		case *pb.UpdatePostReq_ImageChunk:

			chunkData := x.ImageChunk.ChunkData
			chunkNum := x.ImageChunk.ChunkNum

			if imageBufs[chunkNum] == nil {
				imageBufs[chunkNum] = &bytes.Buffer{}
			}

			if _, err := imageBufs[chunkNum].Write(chunkData); err != nil {
				return err
			}

		default:
			return fmt.Errorf("CreatePostReq.Request has unexpected type %T", x)
		}
	}

END:

	mAt, err := ptypes.Timestamp(in.MeetingAt)
	if err != nil {
		return err
	}

	p := &models.Post{
		ID:                in.Id,
		Title:             in.Title,
		Content:           in.Content,
		FishingSpotTypeID: in.FishingSpotTypeId,
		PostsFishTypes:    models.ConvPostsFishTypes(in.FishTypeIds),
		PrefectureID:      in.PrefectureId,
		MeetingPlaceID:    in.MeetingPlaceId,
		MeetingAt:         mAt,
		MaxApply:          in.MaxApply,
	}
	if err := c.postInteractor.UpdatePost(ctx, p, imageBufs, in.DeleteImageIds); err != nil {
		return err
	}

	pProto, err := convPostProto(p)
	if err != nil {
		return err
	}

	return stream.SendAndClose(pProto)
}

func (c *postController) DeletePost(ctx context.Context, in *pb.DeletePostReq) (*empty.Empty, error) {
	if err := c.postInteractor.DeletePost(ctx, in.Id); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (c *postController) GetApplyPost(ctx context.Context, in *pb.GetApplyPostReq) (*pb.ApplyPost, error) {
	a, err := c.postInteractor.GetApplyPost(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return convApplyPostProto(a)
}

func (c *postController) ListApplyPosts(ctx context.Context, in *pb.ListApplyPostsReq) (*pb.ListApplyPostsRes, error) {
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
	return &pb.ListApplyPostsRes{ApplyPosts: listProto}, nil
}

func (c *postController) BatchGetApplyPostsByPostIDs(ctx context.Context, in *pb.BatchGetApplyPostsByPostIDsReq) (*pb.BatchGetApplyPostsByPostIDsRes, error) {
	list, err := c.postInteractor.BatchGetApplyPostsByPostIDs(ctx, in.PostIds)
	if err != nil {
		return nil, err
	}
	listProto, err := convListApplyPostsProto(list)
	if err != nil {
		return nil, err
	}
	return &pb.BatchGetApplyPostsByPostIDsRes{ApplyPosts: listProto}, nil
}

func (c *postController) CreateApplyPost(ctx context.Context, in *pb.CreateApplyPostReq) (*pb.ApplyPost, error) {
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

func (c *postController) DeleteApplyPost(ctx context.Context, in *pb.DeleteApplyPostReq) (*empty.Empty, error) {
	if err := c.postInteractor.DeleteApplyPost(ctx, in.Id); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
