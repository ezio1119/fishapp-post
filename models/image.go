package models

import "time"

type Image struct {
	ID        int64
	Name      string
	PostID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
