package repo

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/usecase/repo"
)

type imageUploaderRepo struct {
	client *storage.Client
}

func NewImageUploaderRepo(c *storage.Client) repo.ImageUploaderRepo {
	return &imageUploaderRepo{c}
}

func (r *imageUploaderRepo) UploadImage(ctx context.Context, image io.Reader, objName string) (string, error) {

	wc := r.client.Bucket(conf.C.Gcs.BucketName).Object(objName).NewWriter(ctx)

	if _, err := io.Copy(wc, image); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	return conf.C.Gcs.BucketURL + objName, nil
}
