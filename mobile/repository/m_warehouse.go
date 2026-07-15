package repository

import (
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryMWarehouseImpl struct {
		*gorm.DB
	}
)

type MWarehouseRepository interface {
	FindOneByID(whID int64) (model.MWarehouse, error)
}

func NewMWarehouseRepository(db *gorm.DB) MWarehouseRepository {
	return &RepositoryMWarehouseImpl{db}
}

func (repository *RepositoryMWarehouseImpl) FindOneByID(whID int64) (warehouse model.MWarehouse, err error) {
	err = repository.Select("wh_id, wh_code, wh_name, is_active").
		Table("mst.m_warehouse").
		Where("mst.m_warehouse.wh_id = ?", whID).
		Take(&warehouse).Error
	return warehouse, err
}
