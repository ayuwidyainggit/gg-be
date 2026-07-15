package pjp

import (
	"context"
	"errors"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) GetPjpById(ctx context.Context, tx *gorm.DB, pjpId int, currentCustomerId string) model.Pjp {
	var pjp model.Pjp
	db := tx.WithContext(ctx)

	err := db.Where("id = ? AND cust_id = ?", pjpId, currentCustomerId).First(&pjp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.ErrorPanic(errors.New("pjp not found"))
		}
		helper.ErrorPanic(err)
	}

	return pjp
}
