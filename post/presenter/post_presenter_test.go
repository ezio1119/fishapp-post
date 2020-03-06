package presenter_test

import (
	"testing"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/post/presenter"
	"github.com/stretchr/testify/assert"
)

func TestTransformPostProto(t *testing.T) {
	p := presenter.NewPostPresenter()
	tests := []struct {
		in  *models.Post
		out *post_grpc.Post
	}{
		{in: &models.Post{}, out: &post_grpc.Post{}},
	}
	for _, tt := range tests {
		p, err := p.TransformPostProto(tt.in)
		assert.NoError(t, err)
		assert.NotNil(t, p)
	}
}
