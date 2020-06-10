package models

import "time"

type Image struct {
	ID        int64
	URL       string
	PostID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
