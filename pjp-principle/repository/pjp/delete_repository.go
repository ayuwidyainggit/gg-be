package pjp

import (
	"context"
	"fmt"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custID string) {
	db := tx.WithContext(ctx)

	// Cek apakah data dengan pjp_id dan cust_id ada
	var existing model.Pjp
	err := db.Where("id = ? AND cust_id = ?", pjpId, custID).First(&existing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.ErrorPanic(fmt.Errorf("data dengan PJP id %d dan customer ID %s tidak ditemukan", pjpId, custID))
		}
		helper.ErrorPanic(err)
	}

	// Hapus data
	if err := db.Delete(&existing).Error; err != nil {
		helper.ErrorPanic(err)
	}
}
