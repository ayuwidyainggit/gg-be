package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type AreaRepository interface {
	List(ctx context.Context, custID string) ([]model.Area, error)
}

type areaRepository struct {
	db *gorm.DB
}

func NewAreaRepository(db *gorm.DB) AreaRepository {
	return &areaRepository{
		db: db,
	}
}

func (r *areaRepository) List(ctx context.Context, custID string) ([]model.Area, error) {
	var records []model.Area

	query := `
		SELECT cust_id, area_id, area_code, area_name, region_id, official_id
		FROM mst.m_area
		WHERE cust_id LIKE ?
			AND is_active = true
			AND is_del = false`

	err := r.db.WithContext(ctx).Raw(query, custID+"%").Scan(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
