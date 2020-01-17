package models

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WrapOnPostRepoErr(err error) error {
	return fmt.Errorf("post repository error: %w", err)
}

func WrapOnPostInterErr(err error) error {
	return fmt.Errorf("post interactor error: %w", err)
}

func WrapOnGrpcErr(err error) error {
	var postNotFound *PostNotFound
	var postsNotFound *PostsNotFound
	var rowsAffectedErr *RowsAffected
	if errors.As(err, &postNotFound) || errors.As(err, &postsNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.As(err, &rowsAffectedErr) {
		return status.Error(codes.Internal, err.Error())
	}
	return status.Error(codes.Unknown, err.Error())
}

type UpdatePostPermissionDenied struct {
	UserID int64
	PostID int64
}

func NewUpdatePostPermissionDeniedErr(uID int64, pID int64) error {
	return &UpdatePostPermissionDenied{uID, pID}
}

func (e *UpdatePostPermissionDenied) Error() string {
	return fmt.Sprintf("user_id=%d does not have permission to Update post_id=%d", e.UserID, e.PostID)
}

type DeletePostPermissionDenied struct {
	UserID int64
	PostID int64
}

func NewDeletePostPermissionDeniedErr(uID int64, pID int64) error {
	return &DeletePostPermissionDenied{uID, pID}
}

func (e *DeletePostPermissionDenied) Error() string {
	return fmt.Sprintf("user_id=%d does not have permission to Delete post_id=%d", e.UserID, e.PostID)
}

type RowsAffected struct {
	Rows int64
}

func NewRowsAffectedErr(rows int64) error {
	return &RowsAffected{rows}
}

func (e *RowsAffected) Error() string {
	return fmt.Sprintf("weird Behaviour. total affected: %d", e.Rows)
}

type PostNotFound struct {
	ID int64
}

func NewPostNotFoundErr(id int64) error {
	return &PostNotFound{id}
}

func (e *PostNotFound) Error() string {
	return fmt.Sprintf("post_id=%d is not found", e.ID)
}

type PostsNotFound struct{}

func NewPostsNotFoundErr() error {
	return &PostsNotFound{}
}

func (e *PostsNotFound) Error() string {
	return fmt.Sprintf("posts are not found")
}
