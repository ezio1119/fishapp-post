package repo

import (
	"github.com/ezio1119/fishapp-post/interfaces/controllers/event"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/golang/protobuf/ptypes"
)

func convPostCreatedProto(p *models.Post) (*event.PostCreated, error) {
	cAt, err := ptypes.TimestampProto(p.CreatedAt)
	if err != nil {
		return nil, err
	}
	uAt, err := ptypes.TimestampProto(p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mAt, err := ptypes.TimestampProto(p.MeetingAt)
	if err != nil {
		return nil, err
	}
	return &event.PostCreated{
		Id:                p.ID,
		Title:             p.Title,
		Content:           p.Content,
		FishingSpotTypeId: p.FishingSpotTypeID,
		FishTypeIds:       p.FishTypeIDs,
		PrefectureId:      p.PrefectureID,
		MeetingPlaceId:    p.MeetingPlaceID,
		MeetingAt:         mAt,
		MaxApply:          p.MaxApply,
		UserId:            p.UserID,
		CreatedAt:         cAt,
		UpdatedAt:         uAt,
	}, nil
}
