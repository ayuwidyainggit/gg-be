package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"scyllax-tms/model"
)

type ShipmentOrderStatusRepoImpl struct {
	Db *gorm.DB
}

func NewShipmentOrderStatusRepoImpl(db *gorm.DB) ShipmentOrderStatusRepo {
	return &ShipmentOrderStatusRepoImpl{Db: db}
}

func (repo *ShipmentOrderStatusRepoImpl) CreateOrUpdate(ctx context.Context, data model.ShipmentOrderStatus) error {
	var existingRecord model.ShipmentOrderStatus
	err := repo.Db.WithContext(ctx).Where("order_no = ?", data.OrderNo).First(&existingRecord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result := repo.Db.WithContext(ctx).Create(&data)
			return result.Error
		}
		return err
	}

	existingRecord.StatusOrder = data.StatusOrder

	result := repo.Db.WithContext(ctx).Save(&existingRecord)
	return result.Error
}
