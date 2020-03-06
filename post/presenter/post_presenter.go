package presenter

import (
	"github.com/ezio1119/fishapp-post/post"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"github.com/golang/protobuf/ptypes"
)

type postPresenter struct{}

func NewPostPresenter() post.Presenter {
	return &postPresenter{}
}

func (*postPresenter) transformPostProto(p *models.Post) (*post_grpc.Post, error) {
	updatedAt, err := ptypes.TimestampProto(p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	createdAt, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post_grpc.Post{
		Id:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		UserId:    p.UserID,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (*postPresenter) transformApplyPostProto(a *models.ApplyPost) (*post_grpc.ApplyPost, error) {
	updatedAt, err := ptypes.TimestampProto(a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	createdAt, err := ptypes.TimestampProto(a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post_grpc.ApplyPost{
		Id:        a.ID,
		PostId:    a.PostID,
		UserId:    a.UserID,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (pp *postPresenter) ConvertGetPostByIDRes(p *models.Post) (*post_grpc.GetPostByIDRes, error) {
	pProto, err := pp.transformPostProto(p)
	if err != nil {
		return nil, err
	}
	return &post_grpc.GetPostByIDRes{Post: pProto}, nil
}
func (pp *postPresenter) ConvertGetListPostsRes(posts []*models.Post) (*post_grpc.GetListPostsRes, error) {
	postsProto := make([]*post_grpc.Post, len(posts))
	for i, p := range posts {
		pProto, err := pp.transformPostProto(p)
		if err != nil {
			return nil, err
		}
		postsProto[i] = pProto
	}
	return &post_grpc.GetListPostsRes{Posts: postsProto}, nil
}
func (pp *postPresenter) ConvertCreatePostRes(p *models.Post) (*post_grpc.CreatePostRes, error) {
	pProto, err := pp.transformPostProto(p)
	if err != nil {
		return nil, err
	}
	return &post_grpc.CreatePostRes{Post: pProto}, nil
}
func (pp *postPresenter) ConvertUpdatePostRes(p *models.Post) (*post_grpc.UpdatePostRes, error) {
	pProto, err := pp.transformPostProto(p)
	if err != nil {
		return nil, err
	}
	return &post_grpc.UpdatePostRes{Post: pProto}, nil
}
func (pp *postPresenter) ConvertDeletePostRes(b bool) (*post_grpc.DeletePostRes, error) {
	return &post_grpc.DeletePostRes{Success: b}, nil
}

func (pp *postPresenter) ConvertCreateApplyPostRes(a *models.ApplyPost) (*post_grpc.CreateApplyPostRes, error) {
	aProto, err := pp.transformApplyPostProto(a)
	if err != nil {
		return nil, err
	}
	return &post_grpc.CreateApplyPostRes{ApplyPost: aProto}, nil
}

func (pp *postPresenter) ConvertDeleteApplyPostRes(b bool) (*post_grpc.DeleteApplyPostRes, error) {
	return &post_grpc.DeleteApplyPostRes{Success: b}, nil
}
