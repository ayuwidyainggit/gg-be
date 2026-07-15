package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) FindDataByID(ctx context.Context, tx *gorm.DB, outletVisitId int) (data model.OutletVisitListPrinciple, err error) {
	result := tx.WithContext(ctx).
		Table("pjp_principles.outlet_visit_list").
		Where("id = ?", outletVisitId).
		Scan(&data)

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}
