package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) UpdateOutletVisitListByID(ctx context.Context, tx *gorm.DB, id int64, data model.OutletVisitListPrinciple) error {
	err := tx.WithContext(ctx).
		Where("id = ?", id).
		Updates(&data).Error

	if err != nil {
		return err
	}

	return nil
}
