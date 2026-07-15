package pjp

import (
	"context"
	"errors"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) GetPjpsByEmpCodes(ctx context.Context, tx *gorm.DB, empCode []string, currentCustomerId string) []model.Pjp {
	var pjps []model.Pjp

	err := tx.WithContext(ctx).
		Where("cust_id = ? AND salesman_code IN ?", currentCustomerId, empCode).
		Find(&pjps).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.ErrorPanic(errors.New("pjp not found"))
		}
		helper.ErrorPanic(errors.New("error getting pjp"))
	}

	return pjps
}
