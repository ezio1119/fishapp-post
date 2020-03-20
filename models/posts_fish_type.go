package models

import "time"

type PostsFishType struct {
	ID         int64
	PostID     int64
	FishTypeID int64
	UpdatedAt  time.Time
	CreatedAt  time.Time
}
