package pjp

import (
	"context"
	"scyllax-pjp/constant"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"errors"

	"gorm.io/gorm"
)

func (repo *pjpRepository) Update(ctx context.Context, tx *gorm.DB, pjp model.Pjp) {
	db := tx.WithContext(ctx)

	// Cek duplikat pjp_code untuk cust_id yang sama selain ID ini
	var existing model.Pjp
	dupCheckErr := db.
		Where("pjp_code = ? AND cust_id = ? AND id != ?", pjp.PjpCode, pjp.CustID, pjp.ID).
		First(&existing).Error

	if dupCheckErr == nil {
		// Duplikat ditemukan
		helper.ErrorPanic(errors.New(constant.ErrPjpCodeExists))
	} else if !errors.Is(dupCheckErr, gorm.ErrRecordNotFound) {
		// Error lain (bukan not found)
		helper.ErrorPanic(dupCheckErr)
	}

	// Update data
	if err := db.Model(&pjp).Updates(pjp).Error; err != nil {
		helper.ErrorPanic(err)
	}
}
