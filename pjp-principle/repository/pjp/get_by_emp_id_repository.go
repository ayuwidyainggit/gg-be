package pjp

import (
	"context"
	"errors"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) GetPjpIdByEmpId(ctx context.Context, tx *gorm.DB, empId int, currentCustomerId string) model.Pjp {
	var pjp model.Pjp

	err := tx.WithContext(ctx).
		Where("cust_id = ? AND salesman_id = ?", currentCustomerId, empId).
		First(&pjp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.ErrorPanic(errors.New("pjp not found"))
		}
		helper.ErrorPanic(errors.New("error getting pjp"))
	}

	return pjp
}
