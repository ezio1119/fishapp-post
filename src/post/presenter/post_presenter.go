package presenter

import (
	"github.com/ezio1119/fishapp-post/post"
	"github.com/golang/protobuf/ptypes"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
)

type postPresenter struct{}

func NewPostPresenter() post.Presenter {
	return &postPresenter{}
}

func (*postPresenter) TransformPostProto(po *models.Post) (*post_grpc.Post, error) {
	updatedAt, err := ptypes.TimestampProto(po.UpdatedAt)
	if err != nil {
		return nil, err
	}
	createdAt, err := ptypes.TimestampProto(po.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post_grpc.Post{
		Id:        po.ID,
		Title:     po.Title,
		Content:   po.Content,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
		UserId:    po.UserID,
	}, nil
}

func (p *postPresenter) TransformListPostProto(listPost []*models.Post) (*post_grpc.ListPost, error) {
	listProto := make([]*post_grpc.Post, len(listPost))
	for i, post := range listPost {
		postProto, err := p.TransformPostProto(post)
		if err != nil {
			return nil, err
		}
		listProto[i] = postProto
	}
	return &post_grpc.ListPost{
		Posts: listProto,
	}, nil
}
