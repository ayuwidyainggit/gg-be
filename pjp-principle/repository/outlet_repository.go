package repository

import (
	"context"
	"fmt"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strings"

	"github.com/WinterYukky/gorm-extra-clause-plugin/exclause"
	"gorm.io/gorm"
)

type OutletRepositoryImpl struct {
	Db *gorm.DB
}

type OutletRepository interface {
	GetOutlet(ctx context.Context, dataFilter model.DmsQueryFilter) []model.Outlet
}

func NewOutletRepository(Db *gorm.DB) OutletRepository {
	return &OutletRepositoryImpl{Db: Db}
}

func (repo OutletRepositoryImpl) GetOutlet(ctx context.Context, dataFilter model.DmsQueryFilter) []model.Outlet {
	var data []model.Outlet

	query := repo.Db.Clauses(exclause.NewWith("cte", repo.Db.Table("mst.m_outlet"))).
		Select("outlet_id", "outlet_name", "outlet_code", "outlet_status", "latitude", "longitude", "address1", "avg_sales_week").
		Table("cte").
		Where("NOT EXISTS (SELECT 1 FROM pjp_principles.destinations WHERE cte.outlet_id = pjp_principles.destinations.outlet_id)")

	if dataFilter.Sort != "" {
		sortBy := ""
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf("%s %s, ", colSort[0], colSort[1])
			}
		}

		sortBy = strings.TrimSuffix(sortBy, ", ")
		query = query.Order(sortBy)
	} else {
		query = query.Order("cte.outlet_id DESC")
	}

	result := query.WithContext(ctx).Scan(&data)
	helper.ErrorPanic(result.Error)

	return data

}
