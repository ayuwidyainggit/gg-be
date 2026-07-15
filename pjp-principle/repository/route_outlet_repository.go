package repository

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RouteOutletRepository interface {
	CreateBulk(ctx context.Context, outlets []model.Destination) error
	FindByPjpId(ctx context.Context, pjpId int) ([]model.Destination, error)

	FindAll(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) []model.Destination
	FindAllEnhance(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) []model.Pjp
	FindByRouteCode(ctx context.Context, code int) (model.Destination, error)
	FindById(ctx context.Context, id int) (model.Destination, error)
	Insert(ctx context.Context, route model.Destination)
	FindByRouteCodeAndPjpCode(ctx context.Context, routeCode int, pjpCode int) ([]model.Destination, error)
	UpdatePivot(ctx context.Context, route model.Destination)
	UpdatePjp(ctx context.Context, route model.Destination) error
	Save(ctx context.Context, route model.Destination) error
	DeleteByDestinationCode(ctx context.Context, route model.Destination) error
	DeleteByDestinationCodeAdditional(ctx context.Context, route model.DestinationAdditional) error
	FindByPjpCode(ctx context.Context, pjpCode int) (model.Destination, error)
	FindByPjpCodeEnhance(ctx context.Context, pjpCode int) ([]model.Destination, error)
	FindByRouteCodes(ctx context.Context, routeCode, pjpCode int) []model.Destination
	Update(ctx context.Context, code int, name string)
	UpdateOrCreate(ctx context.Context, route model.Destination)
	Create(ctx context.Context, route model.Destination)
	CreateAdditionalRoute(ctx context.Context, route model.DestinationAdditional)
	Count(ctx context.Context, currentCustomerId string) (int64, error)
	CountAllEnhance(ctx context.Context, currentCustomerId string) (int64, error)
	UpdateNewRoute(ctx context.Context, route model.Destination)
	GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (data []model.Destination)
	FindByRouteCodeAndDestinationIDAndPjpNull(ctx context.Context, routeCode int, DestinationID int) (*model.Destination, error)
	UpdatePjpRouteOutlet(ctx context.Context, route model.Destination) error
	FindAllDestinationIDByPjpId(ctx context.Context, pjpIds []int) []int
	FindAllDestinationIDByPjpIdToday(ctx context.Context, pjpIds []int) []int
	MobileCancelAddOutletToRoute(ctx context.Context, route model.DestinationAdditional, tx ...*gorm.DB) error
	BeginTx(ctx context.Context) (*gorm.DB, error)
	SearchDestinationIDByPjpId(ctx context.Context, pjpIds []int, search string) []int
}
