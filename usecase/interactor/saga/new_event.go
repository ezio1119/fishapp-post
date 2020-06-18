package saga

import (
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

func newCreateRoomEvent(c *pb.CreateRoom) (*models.Outbox, error) {
	eventData, err := protojson.Marshal(c)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return &models.Outbox{
		ID:        uuid.New().String(),
		EventType: "create.room",
		EventData: eventData,
		Channel:   "create.room",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func newPostApprovedEvent(p *pb.Post) (*models.Outbox, error) {
	jsonPost, err := protojson.Marshal(p)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return &models.Outbox{
		ID:        uuid.New().String(),
		EventType: "post.approved",
		EventData: jsonPost,
		Channel:   "post.approved",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
