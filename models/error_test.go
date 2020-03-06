package models_test

import (
	"errors"
	"testing"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/stretchr/testify/assert"
)

func TestWrapOnPostRepoErr(t *testing.T) {
	tests := []struct {
		in  error
		out error
	}{
		{errors.New("some error"), errors.New("post repository error: some error")},
		{errors.New("some error 2"), errors.New("post repository error: some error 2")},
	}
	for _, tt := range tests {
		err := models.WrapOnPostRepoErr(tt.in)
		assert.EqualError(t, err, tt.out.Error())
	}
}

func TestWrapOnPostInterErr(t *testing.T) {
	tests := []struct {
		in  error
		out error
	}{
		{errors.New("some error"), errors.New("post interactor error: some error")},
		{errors.New("some error 2"), errors.New("post interactor error: some error 2")},
	}
	for _, tt := range tests {
		err := models.WrapOnPostInterErr(tt.in)
		assert.EqualError(t, err, tt.out.Error())
	}
}

func TestNewUpdatePostPermissionDeniedErr(t *testing.T) {
	tests := []struct {
		in  []int64
		out error
	}{
		{[]int64{1, 1}, errors.New("user_id=1 does not have permission to Update post_id=1")},
		{[]int64{100, 300}, errors.New("user_id=100 does not have permission to Update post_id=300")},
	}
	for _, tt := range tests {
		err := models.NewUpdatePostPermissionDeniedErr(tt.in[0], tt.in[1])
		assert.EqualError(t, err, tt.out.Error())
	}
}
