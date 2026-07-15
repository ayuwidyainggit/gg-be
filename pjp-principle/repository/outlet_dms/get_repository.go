package outletdms

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strconv"

	"gorm.io/gorm"
)

func (repo *outletDmsRepository) GetOutletDms(ctx context.Context, tx *gorm.DB, filter model.OutletQueryFilter, custId string) []model.OutletDms {
	var outletDms []model.OutletDms

	query := repo.buildQuery(ctx, tx, filter, custId)

	// Sorting
	if filter.Sort != "" {
		query = query.Order(filter.Sort)
	} else {
		query = query.Order("outlet_id DESC")
	}

	// Pagination
	if filter.Page != "" && filter.Limit != "" {
		page, err1 := strconv.Atoi(filter.Page)
		limit, err2 := strconv.Atoi(filter.Limit)
		if err1 == nil && err2 == nil && page > 0 && limit > 0 {
			offset := (page - 1) * limit
			query = query.Offset(offset).Limit(limit)
		}
	}

	result := query.Find(&outletDms)
	helper.ErrorPanic(result.Error)

	return outletDms
}
