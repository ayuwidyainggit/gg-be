package routepopdaily

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *routePopDailyRepository) Create(ctx context.Context, tx *gorm.DB, routePopDaily model.RoutePopDaily) {
	result := tx.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "year"},
			{Name: "week"},
			{Name: "date"},
			{Name: "day"},
			{Name: "route_code"},
			{Name: "pjp_id"},
			{Name: "pjp_code"},
			{Name: "cust_id"},
			{Name: "status"},
		},
		UpdateAll: true,
	}).Create(&routePopDaily)

	helper.ErrorPanic(result.Error)
}
