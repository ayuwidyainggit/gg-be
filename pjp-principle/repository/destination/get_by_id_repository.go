package destination

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationRepository) GetById(ctx context.Context, tx *gorm.DB, id int, custId string) model.Destination {
	var route model.Destination

	result := tx.WithContext(ctx).First(&route, "id = ? AND cust_id = ?", id, custId)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return route
}
