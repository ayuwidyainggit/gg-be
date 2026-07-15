package destinationhistory

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationHistoryRepository) CreateBulk(ctx context.Context, tx *gorm.DB, outlets []model.DestinationHistory) {
	if len(outlets) == 0 {
		return
	}
	tx.WithContext(ctx).Create(&outlets)
}
