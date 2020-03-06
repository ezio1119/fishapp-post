package interactor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/post/interactor"
	"github.com/ezio1119/fishapp-post/post/mocks"
)

func TestGetList(t *testing.T) {
	mPostRepo := new(mocks.Repository)
	mPostPre := new(mocks.Presenter)
	mListPost := []*models.Post{}
	mListPostProto := &post_grpc.ListPost{}
	tim := time.Time{}

	t.Run("正常に動作すること", func(t *testing.T) {
		mPostRepo.On("GetList", mock.Anything, tim, mock.AnythingOfType("int64")).Return(mListPost, nil).Once()
		mPostPre.On("TransformListPostProto", mListPost).Return(mListPostProto, nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.GetList(context.TODO(), tim, 0)

		assert.NoError(t, err)
		assert.NotNil(t, p)

		mPostRepo.AssertExpectations(t)
		mPostPre.AssertExpectations(t)
	})

	t.Run("リポジトリエラーを伝搬すること", func(t *testing.T) {
		mPostRepo.On("GetList", mock.Anything, tim, mock.AnythingOfType("int64")).Return(nil, fmt.Errorf("Unexpected")).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.GetList(context.TODO(), tim, 0)

		assert.EqualError(t, err, "Unexpected")
		assert.Nil(t, p)

		mPostRepo.AssertExpectations(t)
	})
}
func TestGetByID(t *testing.T) {
	mPostRepo := new(mocks.Repository)
	mPostPre := new(mocks.Presenter)
	mPost := &models.Post{}
	mPostProto := &post_grpc.Post{}

	t.Run("正常に動作すること", func(t *testing.T) {
		mPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mPost, nil).Once()
		mPostPre.On("TransformPostProto", mPost).Return(mPostProto, nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.GetByID(context.TODO(), 0)

		assert.NoError(t, err)
		assert.NotNil(t, p)

		mPostRepo.AssertExpectations(t)
		mPostPre.AssertExpectations(t)
	})
	t.Run("リポジトリエラーを伝搬すること", func(t *testing.T) {
		mPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, fmt.Errorf("unexpected error")).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.GetByID(context.TODO(), 0)

		assert.EqualError(t, err, "unexpected error")
		assert.Nil(t, p)

		mPostRepo.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	mPostRepo := new(mocks.Repository)
	mPostPre := new(mocks.Presenter)
	mPost := &models.Post{}
	mPostProto := &post_grpc.Post{}

	t.Run("正常に動作すること", func(t *testing.T) {
		mPostRepo.On("Create", mock.Anything, mPost).Return(nil).Once()
		mPostPre.On("TransformPostProto", mPost).Return(mPostProto, nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.Create(context.TODO(), mPost)

		assert.NoError(t, err)
		assert.NotNil(t, p)

		mPostRepo.AssertExpectations(t)
		mPostPre.AssertExpectations(t)
	})

	t.Run("リポジトリエラーを伝搬すること", func(t *testing.T) {
		mPostRepo.On("Create", mock.Anything, mPost).Return(fmt.Errorf("Unexpected")).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.Create(context.TODO(), mPost)

		assert.EqualError(t, err, "Unexpected")
		assert.Nil(t, p)

		mPostRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mPostRepo := new(mocks.Repository)
	mPostPre := new(mocks.Presenter)
	mPost := &models.Post{}
	mPostProto := &post_grpc.Post{}

	t.Run("正常に動作すること", func(t *testing.T) {
		mPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mPost, nil).Once().
			On("Update", mock.Anything, mPost).Return(nil).Once()
		mPostPre.On("TransformPostProto", mPost).Return(mPostProto, nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.Update(context.TODO(), mPost, 0)

		assert.NoError(t, err)
		assert.NotNil(t, p)

		mPostRepo.AssertExpectations(t)
		mPostPre.AssertExpectations(t)
	})

	t.Run("権限がないときにエラーが発生すること", func(t *testing.T) {
		mPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mPost, nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		p, err := i.Update(context.TODO(), mPost, 19)

		e := models.WrapOnPostInterErr(&models.UpdatePostPermissionDenied{UserID: 19})
		assert.EqualError(t, err, e.Error())
		assert.Nil(t, p)

		mPostRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mPostRepo := new(mocks.Repository)
	mPostPre := new(mocks.Presenter)
	mPost := &models.Post{}

	t.Run("正常に動作すること", func(t *testing.T) {
		mPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mPost, nil).Once().
			On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		err := i.Delete(context.TODO(), 0, 0)

		assert.NoError(t, err)

		mPostRepo.AssertExpectations(t)
	})
	t.Run("権限がないときにエラーが発生すること", func(t *testing.T) {
		mPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mPost, nil).Once()
		i := interactor.NewPostInteractor(mPostRepo, mPostPre, time.Second*2)
		e := models.WrapOnPostInterErr(&models.DeletePostPermissionDenied{UserID: 19})
		err := i.Delete(context.TODO(), 0, 19)

		assert.EqualError(t, err, e.Error())

		mPostRepo.AssertExpectations(t)
	})
}
