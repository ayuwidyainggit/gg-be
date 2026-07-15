package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type RegionRepository interface {
	List(ctx context.Context, custID string) ([]model.Region, error)
}

type regionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepository{
		db: db,
	}
}

func (r *regionRepository) List(ctx context.Context, custID string) ([]model.Region, error) {
	var records []model.Region

	query := `
		SELECT cust_id, region_id, region_code, region_name
		FROM mst.m_region
		WHERE cust_id LIKE ?
			AND is_active = true
			AND is_del = false`

	err := r.db.WithContext(ctx).Raw(query, custID+"%").Scan(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
