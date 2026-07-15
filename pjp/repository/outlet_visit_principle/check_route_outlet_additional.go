package outlet_visit_principle

import (
	"context"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) CheckRouteOutletAdditional(ctx context.Context, tx *gorm.DB, destinationCode string, destinationId int) (bool, error) {
	var count int64
	err := tx.Table("pjp_principles.destinations_additional").
		Where("destination_code = ? AND destination_id = ?", destinationCode, destinationId).
		Count(&count).Error

	return count > 0, err
}
