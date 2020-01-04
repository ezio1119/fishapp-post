package presenter

import (
	"github.com/ezio1119/fishapp-post/entry"
	"github.com/ezio1119/fishapp-post/entry/controllers/entry_post_grpc"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/golang/protobuf/ptypes"
)

type entryPresenter struct{}

func NewEntryPresenter() entry.Presenter {
	return &entryPresenter{}
}

func (*entryPresenter) TransformEntryProto(e *models.Entry) (*entry_post_grpc.Entry, error) {
	updatedAt, err := ptypes.TimestampProto(e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	createdAt, err := ptypes.TimestampProto(e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &entry_post_grpc.Entry{
		Id:        e.ID,
		PostId:    e.PostID,
		UserId:    e.UserID,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}
