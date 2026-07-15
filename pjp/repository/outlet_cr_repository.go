package repository

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type OutletCrRepository interface {
	CreateOutletCr(ctx context.Context, tx *gorm.DB, outletCr model.OutletCr) (int64, error)
	CreateOutletCrDet(ctx context.Context, tx *gorm.DB, outletCrDet model.OutletCrDet) error
	GetOutletLocation(ctx context.Context, db *gorm.DB, custId string, outletId int64) (latitude, longitude string, err error)
}

type OutletCrRepositoryImpl struct {
	Db *gorm.DB
}

func NewOutletCrRepository(db *gorm.DB) OutletCrRepository {
	return &OutletCrRepositoryImpl{
		Db: db,
	}
}

// CreateOutletCr inserts a new outlet change request record
func (repo *OutletCrRepositoryImpl) CreateOutletCr(ctx context.Context, tx *gorm.DB, outletCr model.OutletCr) (int64, error) {
	result := tx.WithContext(ctx).Create(&outletCr)
	if result.Error != nil {
		return 0, result.Error
	}
	return outletCr.OutletCrId, nil
}

// CreateOutletCrDet inserts a new outlet change request detail record
func (repo *OutletCrRepositoryImpl) CreateOutletCrDet(ctx context.Context, tx *gorm.DB, outletCrDet model.OutletCrDet) error {
	result := tx.WithContext(ctx).Create(&outletCrDet)
	return result.Error
}

// GetOutletLocation retrieves latitude and longitude from mst.m_outlet table
func (repo *OutletCrRepositoryImpl) GetOutletLocation(ctx context.Context, db *gorm.DB, custId string, outletId int64) (latitude, longitude string, err error) {
	var outlet struct {
		Latitude  string
		Longitude string
	}

	result := db.WithContext(ctx).
		Table("mst.m_outlet").
		Select("latitude", "longitude").
		Where("cust_id LIKE ? AND outlet_id = ? AND is_del = false", custId+"%", outletId).
		First(&outlet)

	if result.Error != nil {
		return "", "", result.Error
	}

	return outlet.Latitude, outlet.Longitude, nil
}
