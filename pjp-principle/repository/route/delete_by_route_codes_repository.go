package route

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeRepository) DeleteByRouteCodes(ctx context.Context, tx *gorm.DB, routeCodes []int, custId string) {
	if len(routeCodes) == 0 {
		return // tidak perlu hapus apapun
	}

	db := tx.WithContext(ctx)
	if err := db.Where("route_code IN ? AND cust_id = ?", routeCodes, custId).Delete(&model.Route{}).Error; err != nil {
		helper.ErrorPanic(err)
	}
}
