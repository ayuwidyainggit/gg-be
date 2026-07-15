package distributordms

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *distributorDmsRepository) CountDistributorDms(ctx context.Context, tx *gorm.DB, filter model.DistributorQueryFilter, custId string) int64 {
	var count int64
	query := repo.buildQuery(ctx, tx, filter, custId)

	result := query.Count(&count)
	helper.ErrorPanic(result.Error)

	return count
}
