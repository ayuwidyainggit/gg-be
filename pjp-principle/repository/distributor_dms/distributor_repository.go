package distributordms

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type DistributorDmsRepository interface {
	GetDistributorDms(ctx context.Context, tx *gorm.DB, dataFilter model.DistributorQueryFilter, custId string) []model.DistributorDms
	CountDistributorDms(ctx context.Context, tx *gorm.DB, filter model.DistributorQueryFilter, custId string) int64
}

type distributorDmsRepository struct {
}

func NewDistributorDmsRepository() DistributorDmsRepository {
	return &distributorDmsRepository{}
}

func (repo *distributorDmsRepository) buildQuery(ctx context.Context, tx *gorm.DB, filter model.DistributorQueryFilter, custId string) *gorm.DB {
	query := tx.WithContext(ctx).Model(&model.DistributorDms{}).Where("cust_id ILIKE ? AND is_active = ?", custId+"%", true)

	if filter.DistributorCode != "" {
		query = query.Where("distributor_code = ?", filter.DistributorCode)
	}

	if filter.DistributorID != 0 {
		query = query.Where("distributor_id = ?", filter.DistributorID)
	}

	if filter.SalesTeamID != "" {
		query = query.Where("sales_team_id = ?", filter.SalesTeamID)
	}

	if filter.Query != "" {
		query = query.Where("distributor_name ILIKE ? OR distributor_code ILIKE ?", "%"+filter.Query+"%", "%"+filter.Query+"%")
	}

	return query
}
