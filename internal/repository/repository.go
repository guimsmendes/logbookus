package repository

import (
	"context"
	"guimsmendes/personal/logbookus/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCity(ctx context.Context) (model.City, error) {
	var city model.City
	err := r.db.WithContext(ctx).First(&city).Error
	return city, err
}
