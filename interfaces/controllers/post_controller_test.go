package controllers

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	mock "github.com/ezio1119/fishapp-post/interfaces/controllers/mock"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
)

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewMockPostServiceClient(ctrl)
	c.EXPECT().CreatePost(gomock.Any(), gomock.Any())
	stream, err := c.CreatePost(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	req := &pb.CreatePostReq{
		Request: &pb.CreatePostReq_Details{
			Details: &pb.CreatePostReqDetails{
				Title:             "ccsd",
				Content:           "cdcdsds",
				FishingSpotTypeId: 1,
				FishTypeIds:       []int64{1, 3},
				PrefectureId:      1,
				MeetingPlaceId:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				MeetingAt:         ptypes.TimestampNow(),
				MaxApply:          3,
				UserId:            1,
			}},
	}

	if err := stream.Send(req); err != nil {
		t.Fatal(err)
	}

	image, err := os.Open("/app/images/724-2.jpg")
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)

	defer image.Close()

	for {
		n, err := image.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.CreatePostReq{
			Request: &pb.CreatePostReq_ImageChunk{
				ImageChunk: buf[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

	post, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(post)
}
