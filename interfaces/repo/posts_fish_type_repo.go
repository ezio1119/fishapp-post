package repo

import "github.com/jinzhu/gorm"

type postsFishTypeRepo struct {
	db *gorm.DB
}

func NewPostsFishTypeRepo(db *gorm.DB) *postsFishTypeRepo {
	return &postsFishTypeRepo{db}
}