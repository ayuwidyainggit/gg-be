package pjp

import (
	"context"
	"errors"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) GetPjpIdByEmpCode(ctx context.Context, tx *gorm.DB, empCode string, currentCustomerId string) model.Pjp {
	var pjp model.Pjp

	err := tx.WithContext(ctx).
		Where("cust_id = ? AND salesman_code = ?", currentCustomerId, empCode).
		First(&pjp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.ErrorPanic(errors.New("pjp not found"))
		}
		helper.ErrorPanic(errors.New("error getting pjp"))
	}

	return pjp
}
