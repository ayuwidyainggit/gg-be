package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type UserLocationRepository interface {
	FindFirst() (*model.MSendLocation, error)
	CreateLocation(ctx context.Context, user *model.UserLocation) error
}

type userLocationRepository struct {
	db *gorm.DB
}

func NewUserLocationRepository(db *gorm.DB) UserLocationRepository {
	return &userLocationRepository{
		db: db,
	}
}

func (r *userLocationRepository) FindFirst() (*model.MSendLocation, error) {
	var records model.MSendLocation

	err := r.db.First(&records).Error
	if err != nil {
		return nil, err
	}

	return &records, nil
}

func (r *userLocationRepository) CreateLocation(ctx context.Context, user *model.UserLocation) error {
	return r.db.WithContext(ctx).Create(&user).Error
}
