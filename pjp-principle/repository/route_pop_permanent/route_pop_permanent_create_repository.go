package routepoppermanent

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routePopPermanentRepository) CreateBulk(ctx context.Context, tx *gorm.DB, routePopPermanents []model.RoutePopPermanent) {
	if len(routePopPermanents) == 0 {
		return
	}
	tx.WithContext(ctx).Create(&routePopPermanents)
}
