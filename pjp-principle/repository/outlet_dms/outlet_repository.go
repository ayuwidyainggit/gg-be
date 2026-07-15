package outletdms

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type OutletDmsRepository interface {
	GetOutletDms(ctx context.Context, tx *gorm.DB, dataFilter model.OutletQueryFilter, custId string) []model.OutletDms
	CountOutletDms(ctx context.Context, tx *gorm.DB, filter model.OutletQueryFilter, custId string) int64
}

type outletDmsRepository struct {
}

func NewOutletDmsRepository() OutletDmsRepository {
	return &outletDmsRepository{}
}

func (repo *outletDmsRepository) buildQuery(ctx context.Context, tx *gorm.DB, filter model.OutletQueryFilter, custId string) *gorm.DB {
	query := tx.WithContext(ctx).Model(&model.OutletDms{}).Where("cust_id ILIKE ? AND is_active = ? AND verification_status = ?", custId+"%", true, 1)

	if filter.OutletCode != "" {
		query = query.Where("outlet_code = ?", filter.OutletCode)
	}

	if filter.OutletID != 0 {
		query = query.Where("outlet_id = ?", filter.OutletID)
	}

	if filter.SalesTeamID != "" {
		query = query.Where("sales_team_id = ?", filter.SalesTeamID)
	}

	if filter.Query != "" {
		query = query.Where("outlet_name ILIKE ? OR outlet_code ILIKE ?", "%"+filter.Query+"%", "%"+filter.Query+"%")
	}

	return query
}
