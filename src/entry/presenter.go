package entry

import (
	"github.com/ezio1119/fishapp-post/entry/controllers/entry_post_grpc"
	"github.com/ezio1119/fishapp-post/models"
)

type Presenter interface {
	TransformEntryProto(po *models.Entry) (*entry_post_grpc.Entry, error)
}
