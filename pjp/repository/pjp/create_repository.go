package pjp

import (
	"context"
	"fmt"
	"scyllax-pjp/constant"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) Create(ctx context.Context, tx *gorm.DB, pjp model.Pjp) model.Pjp {
	db := tx.WithContext(ctx)

	// Cek apakah kombinasi pjp_code + cust_id sudah ada
	var existing model.Pjp
	err := db.Where("pjp_code = ? AND cust_id = ?", pjp.PjpCode, pjp.CustID).First(&existing).Error

	if err == nil {
		helper.ErrorPanic(fmt.Errorf(constant.ErrPjpCodeExists))
	}
	if err != gorm.ErrRecordNotFound {
		helper.ErrorPanic(err)
	}

	if err := db.Create(&pjp).Error; err != nil {
		helper.ErrorPanic(err)
	}

	return pjp
}
