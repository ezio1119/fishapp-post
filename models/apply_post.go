package models

import "time"

type ApplyPost struct {
	ID        int64
	UserID    int64
	PostID    int64
	UpdatedAt time.Time
	CreatedAt time.Time
}