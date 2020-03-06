package controllers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/post/mocks"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	now := ptypes.TimestampNow()
	ctx := context.TODO()
	tests := []struct {
		name    string
		in      *post_grpc.CreateReq
		out     *post_grpc.Post
		setMock func(m *mocks.Usecase)
		wantErr bool
		err     error
	}{
		{
			name: "正常に動作すること",
			in: &post_grpc.CreateReq{
				Title:   "title",
				Content: "content",
				UserId:  1,
			},
			out: &post_grpc.Post{
				Id:        1,
				Title:     "title",
				Content:   "content",
				UserId:    1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			setMock: func(m *mocks.Usecase) {
				mPost := &models.Post{
					Title:   "title",
					Content: "content",
					UserID:  1,
				}
				mPostProto := &post_grpc.Post{
					Id:        1,
					Title:     "title",
					Content:   "content",
					UserId:    1,
					CreatedAt: now,
					UpdatedAt: now,
				}
				m.On("Create", ctx, mPost).Return(mPostProto, nil).Once()
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "interactorエラーをgRPCエラーでラップすること",
			in: &post_grpc.CreateReq{
				Title:   "title",
				Content: "content",
				UserId:  1,
			},
			out: nil,
			setMock: func(m *mocks.Usecase) {
				mPost := &models.Post{
					Title:   "title",
					Content: "content",
					UserID:  1,
				}
				m.On("Create", ctx, mPost).Return(nil, fmt.Errorf("interactor error")).Once()
			},
			wantErr: true,
			err:     models.WrapOnGrpcErr(fmt.Errorf("interactor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mPostInter := new(mocks.Usecase)
			tt.setMock(mPostInter)
			c := controllers.NewPostController(mPostInter)
			p, err := c.Create(ctx, tt.in)
			if tt.wantErr {
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, tt.out, p)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, p)

			mPostInter.AssertExpectations(t)
		})
	}
}

func TestGetByID(t *testing.T) {
	now := ptypes.TimestampNow()
	ctx := context.TODO()
	tests := []struct {
		name    string
		in      *post_grpc.ID
		out     *post_grpc.Post
		setMock func(m *mocks.Usecase)
		wantErr bool
		err     error
	}{
		{
			name: "正常に動作すること",
			in: &post_grpc.ID{
				Id: 1,
			},
			out: &post_grpc.Post{
				Id:        1,
				Title:     "title",
				Content:   "content",
				UserId:    1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			setMock: func(m *mocks.Usecase) {
				mPostProto := &post_grpc.Post{
					Id:        1,
					Title:     "title",
					Content:   "content",
					UserId:    1,
					CreatedAt: now,
					UpdatedAt: now,
				}
				m.On("GetByID", ctx, int64(1)).Return(mPostProto, nil).Once()
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "interactorエラーをgRPCエラーでラップすること",
			in: &post_grpc.ID{
				Id: 1,
			},
			out: nil,
			setMock: func(m *mocks.Usecase) {
				m.On("GetByID", ctx, int64(1)).Return(nil, fmt.Errorf("interactor error")).Once()
			},
			wantErr: true,
			err:     models.WrapOnGrpcErr(fmt.Errorf("interactor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mPostInter := new(mocks.Usecase)
			tt.setMock(mPostInter)
			c := controllers.NewPostController(mPostInter)
			p, err := c.GetByID(ctx, tt.in)
			if tt.wantErr {
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, tt.out, p)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, p)

			mPostInter.AssertExpectations(t)
		})
	}
}

func TestGetList(t *testing.T) {
	nowProto := ptypes.TimestampNow()
	now, err := ptypes.Timestamp(nowProto)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		in      *post_grpc.ListReq
		out     *post_grpc.ListPost
		setMock func(m *mocks.Usecase)
		wantErr bool
		err     error
	}{
		{
			name: "正常に動作すること",
			in: &post_grpc.ListReq{
				Datetime: nowProto,
				Num:      1,
			},
			out: &post_grpc.ListPost{
				Posts: []*post_grpc.Post{
					{
						Id:        1,
						Title:     "title",
						Content:   "content",
						UserId:    1,
						CreatedAt: nowProto,
						UpdatedAt: nowProto,
					},
				},
			},
			setMock: func(m *mocks.Usecase) {
				mListPost := &post_grpc.ListPost{
					Posts: []*post_grpc.Post{
						{
							Id:        1,
							Title:     "title",
							Content:   "content",
							UserId:    1,
							CreatedAt: nowProto,
							UpdatedAt: nowProto,
						},
					},
				}
				m.On("GetList", ctx, now, int64(1)).Return(mListPost, nil).Once()
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "interactorエラーをgRPCエラーでラップすること",
			in: &post_grpc.ListReq{
				Datetime: nowProto,
				Num:      1,
			},
			out: nil,
			setMock: func(m *mocks.Usecase) {
				m.On("GetList", ctx, now, int64(1)).Return(nil, fmt.Errorf("interactor error")).Once()
			},
			wantErr: true,
			err:     models.WrapOnGrpcErr(fmt.Errorf("interactor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mPostInter := new(mocks.Usecase)
			tt.setMock(mPostInter)
			c := controllers.NewPostController(mPostInter)
			p, err := c.GetList(ctx, tt.in)
			if tt.wantErr {
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, tt.out, p)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, p)

			mPostInter.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	nowProto := ptypes.TimestampNow()
	ctx := context.TODO()
	tests := []struct {
		name    string
		in      *post_grpc.UpdateReq
		out     *post_grpc.Post
		setMock func(*mocks.Usecase)
		wantErr bool
		err     error
	}{
		{
			name: "正常に動作すること",
			in: &post_grpc.UpdateReq{
				Id:      1,
				Title:   "title",
				Content: "content",
				UserId:  1,
			},
			out: &post_grpc.Post{
				Id:        1,
				Title:     "title",
				Content:   "content",
				UserId:    1,
				CreatedAt: nowProto,
				UpdatedAt: nowProto,
			},
			setMock: func(m *mocks.Usecase) {
				mPost := &models.Post{
					ID:      1,
					Title:   "title",
					Content: "content",
				}
				mPostProto := &post_grpc.Post{
					Id:        1,
					Title:     "title",
					Content:   "content",
					UserId:    1,
					CreatedAt: nowProto,
					UpdatedAt: nowProto,
				}
				m.On("Update", ctx, mPost, int64(1)).Return(mPostProto, nil).Once()
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "interactorエラーをgRPCエラーでラップすること",
			in: &post_grpc.UpdateReq{
				Id:      1,
				Title:   "title",
				Content: "content",
				UserId:  1,
			},
			out: nil,
			setMock: func(m *mocks.Usecase) {
				mPost := &models.Post{
					ID:      1,
					Title:   "title",
					Content: "content",
				}
				userID := int64(1)
				m.On("Update", ctx, mPost, userID).Return(nil, fmt.Errorf("interactor error")).Once()
			},
			wantErr: true,
			err:     models.WrapOnGrpcErr(fmt.Errorf("interactor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mPostInter := new(mocks.Usecase)
			tt.setMock(mPostInter)
			c := controllers.NewPostController(mPostInter)
			p, err := c.Update(ctx, tt.in)
			if tt.wantErr {
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, tt.out, p)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, p)

			mPostInter.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()
	tests := []struct {
		name    string
		in      *post_grpc.DeleteReq
		out     *wrappers.BoolValue
		setMock func(*mocks.Usecase)
		wantErr bool
		err     error
	}{
		{
			name: "正常に動作すること",
			in: &post_grpc.DeleteReq{
				Id:     1,
				UserId: 1,
			},
			out: &wrappers.BoolValue{Value: true},
			setMock: func(m *mocks.Usecase) {
				m.On("Delete", ctx, int64(1), int64(1)).Return(nil).Once()
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "interactorエラーをgRPCエラーでラップすること",
			in: &post_grpc.DeleteReq{
				Id:     1,
				UserId: 1,
			},
			out: nil,
			setMock: func(m *mocks.Usecase) {
				m.On("Delete", ctx, int64(1), int64(1)).Return(fmt.Errorf("interactor error")).Once()
			},
			wantErr: true,
			err:     models.WrapOnGrpcErr(fmt.Errorf("interactor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mPostInter := new(mocks.Usecase)
			tt.setMock(mPostInter)
			c := controllers.NewPostController(mPostInter)
			out, err := c.Delete(ctx, tt.in)
			if tt.wantErr {
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, tt.out, out)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.out, out)

			mPostInter.AssertExpectations(t)
		})
	}
}
