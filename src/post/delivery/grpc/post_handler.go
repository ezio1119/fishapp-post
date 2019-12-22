package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/ezio1119/fishapp-post/post/delivery/grpc/post_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	PUsecase post.Usecase
}

func NewPostServerGrpc(gserver *grpc.Server, us post.Usecase) {
	postServer := &server{
		PUsecase: us,
	}
	post_grpc.RegisterPostServiceServer(gserver, postServer)
	reflection.Register(gserver)
}

func (s *server) transformPostRPC(po *models.Post) (*post_grpc.Post, error) {
	updatedAt, err := ptypes.TimestampProto(po.UpdatedAt)
	if err != nil {
		return nil, err
	}
	createdAt, err := ptypes.TimestampProto(po.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post_grpc.Post{
		Id:        po.Id,
		Title:     po.Title,
		Content:   po.Content,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
		UserId:    po.UserId,
	}, nil
}

func (s *server) Create(ctx context.Context, in *post_grpc.CreateReq) (*post_grpc.Post, error) {
	post := &models.Post{
		Title:   in.Title,
		Content: in.Content,
		UserId:  in.UserId,
	}
	if err := s.PUsecase.Create(ctx, post); err != nil {
		return nil, err
	}
	postRPC, err := s.transformPostRPC(post)
	if err != nil {
		return nil, err
	}
	return postRPC, nil
}

func (s *server) GetByID(ctx context.Context, in *post_grpc.ID) (*post_grpc.Post, error) {
	post, err := s.PUsecase.GetByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	postRPC, err := s.transformPostRPC(post)
	if err != nil {
		return nil, err
	}
	return postRPC, nil
}

func (s *server) GetList(ctx context.Context, in *post_grpc.ListReq) (*post_grpc.ListPost, error) {
	datetime, err := ptypes.Timestamp(in.Datetime)
	if err != nil {
		return nil, err
	}
	list, err := s.PUsecase.GetList(ctx, datetime, in.Num)
	if err != nil {
		return nil, err
	}
	listRPC := make([]*post_grpc.Post, len(list))
	for i, post := range list {
		postRPC, err := s.transformPostRPC(post)
		if err != nil {
			return nil, err
		}
		listRPC[i] = postRPC
	}
	return &post_grpc.ListPost{
		Posts: listRPC,
	}, nil
}

func (s *server) Update(ctx context.Context, in *post_grpc.UpdateReq) (*post_grpc.Post, error) {
	post := &models.Post{
		Id:      in.Id,
		Title:   in.Title,
		Content: in.Content,
		UserId:  in.UserId,
	}
	if err := s.PUsecase.Update(ctx, post); err != nil {
		return nil, err
	}
	postRPC, err := s.transformPostRPC(post)
	if err != nil {
		return nil, err
	}
	return postRPC, nil
}

func (s *server) Delete(ctx context.Context, in *post_grpc.DeleteReq) (*post_grpc.DeleteRes, error) {
	if err := s.PUsecase.Delete(ctx, in.Id, in.UserId); err != nil {
		return nil, err
	}
	return &post_grpc.DeleteRes{
		Deleted: true,
	}, nil
}
