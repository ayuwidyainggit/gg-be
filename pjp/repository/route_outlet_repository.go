package repository

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RouteOutletRepository interface {
	CreateBulk(ctx context.Context, outlets []model.RouteOutlet) error
	FindByPjpId(ctx context.Context, pjpId int) ([]model.RouteOutlet, error)

	FindAll(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) []model.RouteOutlet
	FindAllEnhance(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) []model.Pjp
	FindByRouteCode(ctx context.Context, code int) (model.RouteOutlet, error)
	FindById(ctx context.Context, id int) (model.RouteOutlet, error)
	Insert(ctx context.Context, route model.RouteOutlet)
	FindByRouteCodeAndPjpCode(ctx context.Context, routeCode int, pjpCode int) ([]model.RouteOutlet, error)
	UpdatePivot(ctx context.Context, route model.RouteOutlet)
	UpdatePjp(ctx context.Context, route model.RouteOutlet) error
	Save(ctx context.Context, route model.RouteOutlet) error
	DeleteByOutletCode(ctx context.Context, route model.RouteOutlet) error
	DeleteByOutletCodeAdditional(ctx context.Context, route model.RouteOutletAdditional) error
	FindByPjpCode(ctx context.Context, pjpCode int) (model.RouteOutlet, error)
	FindByPjpCodeEnhance(ctx context.Context, pjpCode int) ([]model.RouteOutlet, error)
	FindByRouteCodes(ctx context.Context, routeCode, pjpCode int) []model.RouteOutlet
	Update(ctx context.Context, code int, name string)
	UpdateOrCreate(ctx context.Context, route model.RouteOutlet)
	Create(ctx context.Context, route model.RouteOutlet)
	CreateAdditionalRoute(ctx context.Context, route model.RouteOutletAdditional)
	Count(ctx context.Context, currentCustomerId string) (int64, error)
	CountAllEnhance(ctx context.Context, currentCustomerId string) (int64, error)
	UpdateNewRoute(ctx context.Context, route model.RouteOutlet)
	GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (data []model.RouteOutlet)
	FindByRouteCodeAndOutletIDAndPjpNull(ctx context.Context, routeCode int, outletID int) (*model.RouteOutlet, error)
	UpdatePjpRouteOutlet(ctx context.Context, route model.RouteOutlet) error
	FindAllOutletIdByPjpId(ctx context.Context, pjpIds []int) []int
	FindAllOutletIdByPjpIdToday(ctx context.Context, pjpIds []int) []int
	MobileCancelAddOutletToRoute(ctx context.Context, route model.RouteOutletAdditional, tx ...*gorm.DB) error
	BeginTx(ctx context.Context) (*gorm.DB, error)
	SearchOutletIdByPjpId(ctx context.Context, pjpIds []int, search string) []int
}
