package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/ezio1119/fishapp-post/interfaces/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/golang/protobuf/ptypes"
)

func convPostProto(p *models.Post) (*post_grpc.Post, error) {
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
	fishTypeIds := make([]int64, len(p.PostsFishTypes))
	for i, t := range p.PostsFishTypes {
		fishTypeIds[i] = t.FishTypeID
	}
	aProto, err := convListApplyPostsProto(p.ApplyPosts)
	if err != nil {
		return nil, err
	}
	return &post_grpc.Post{
		Id:                p.ID,
		Title:             p.Title,
		Content:           p.Content,
		FishingSpotTypeId: p.FishingSpotTypeID,
		FishTypeIds:       fishTypeIds,
		PrefectureId:      p.PrefectureID,
		MeetingPlaceId:    p.MeetingPlaceID,
		MeetingAt:         mAt,
		MaxApply:          p.MaxApply,
		ApplyPosts:        aProto,
		UserId:            p.UserID,
		CreatedAt:         cAt,
		UpdatedAt:         uAt,
	}, nil
}

func convListPostsProto(list []*models.Post) ([]*post_grpc.Post, error) {
	listP := make([]*post_grpc.Post, len(list))
	for i, p := range list {
		pProto, err := convPostProto(p)
		if err != nil {
			return nil, err
		}
		listP[i] = pProto
	}
	return listP, nil
}

func convApplyPostProto(a *models.ApplyPost) (*post_grpc.ApplyPost, error) {
	cAt, err := ptypes.TimestampProto(a.CreatedAt)
	if err != nil {
		return nil, err
	}
	uAt, err := ptypes.TimestampProto(a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &post_grpc.ApplyPost{
		Id:        a.ID,
		PostId:    a.PostID,
		UserId:    a.UserID,
		CreatedAt: cAt,
		UpdatedAt: uAt,
	}, nil
}

func convListApplyPostsProto(list []*models.ApplyPost) ([]*post_grpc.ApplyPost, error) {
	fmt.Printf("%#v\n", list)
	listA := make([]*post_grpc.ApplyPost, len(list))
	for i, a := range list {
		aP, err := convApplyPostProto(a)
		if err != nil {
			return nil, err
		}
		listA[i] = aP
	}
	return listA, nil
}

func convPostFilter(f *post_grpc.ListPostsReq_Filter) (*models.PostFilter, error) {
	postF := &models.PostFilter{CanApply: f.CanApply, FishTypeIDs: f.FishTypeIds}

	if f.MeetingAtFrom != nil {
		mAtFrom, err := ptypes.Timestamp(f.MeetingAtFrom)
		if err != nil {
			return nil, err
		}
		postF.MeetingAtFrom = mAtFrom.In(time.Local)
	}

	if f.MeetingAtTo != nil {
		mAtTo, err := ptypes.Timestamp(f.MeetingAtTo)
		if err != nil {
			return nil, err
		}
		postF.MeetingAtTo = mAtTo.In(time.Local)
	}

	switch f.OrderBy {
	case post_grpc.ListPostsReq_Filter_ASC:
		postF.OrderBy = models.OrderByAsc
	case post_grpc.ListPostsReq_Filter_DESC:
		postF.OrderBy = models.OrderByDesc
	default:
		return nil, errors.New("ascac")
	}
	switch f.SortBy {
	case post_grpc.ListPostsReq_Filter_CREATED_AT:
		postF.SortBy = models.SortByID
	case post_grpc.ListPostsReq_Filter_MEETING_AT:
		postF.SortBy = models.SortByMeetingAt
	default:
		return nil, errors.New("ascac")
	}
	return postF, nil
}
