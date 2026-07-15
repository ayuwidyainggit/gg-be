package repository

import (
	"context"
	"scyllax-tms/entity"
	"scyllax-tms/model"

	"gorm.io/gorm"
)

type ShipmentRepo interface {
	Insert(ctx context.Context, data model.Shipment) error
	FindByShipmentNo(ctx context.Context, shipmentNo string) (model.Shipment, error)
	// FindByOrderNo(ctx context.Context, orderNo string) (model.Shipment, error)
	FindAll(ctx context.Context, dataFilter entity.ShipmentQueryFilter) []model.Shipment
	DeleteByQuery(ctx context.Context, column string, value any) error
	DeleteBulk(ctx context.Context, shipmentNo []string) error
	GetLastShipmentNo() (string, error)
	UpdateByQuery(ctx context.Context, column string, query any, data model.Shipment) error
	FindByColumns(ctx context.Context, selects []string, columns []string, queries []any) (model.Shipment, error)
	CountSummary(ctx context.Context, query int) (shipment, trip, finished, inProgress int, err error)
	FindByDriverID(ctx context.Context, driverId int) []model.Shipment
	BeginTx(ctx context.Context) (*gorm.DB, error)
	InsertWithTx(tx *gorm.DB, data model.Shipment) error
	FindByStartTimes(ctx context.Context, columns []string, queries []any) (data []model.Shipment, err error)
	FindByEndTimes(ctx context.Context, columns []string, queries []any) (data []model.Shipment, err error)
}
