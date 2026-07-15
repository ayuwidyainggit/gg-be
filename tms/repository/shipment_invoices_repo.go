package repository

import (
	"context"
	"scyllax-tms/entity"
	"scyllax-tms/model"

	"gorm.io/gorm"
)

type ShipmentInvoicesRepo interface {
	Insert(ctx context.Context, data model.ShipmentInvoices) error
	Update(ctx context.Context, data model.ShipmentInvoices) error
	FindByShipmentNoAndOutletName(ctx context.Context, shipmentNo string, outletName string) ([]model.ShipmentInvoices, error)
	FindByTwoColumns(ctx context.Context, columnFirst, columnSecond string, queryFirst, querySecond any) ([]model.ShipmentInvoices, error)
	FindByOneColumn(ctx context.Context, column string, query any) (model.ShipmentInvoices, error)
	FindByOutletId(ctx context.Context, params []interface{}) []entity.OutletResponse
	FindAll(ctx context.Context, dataFilter entity.ShipmentInvoicesQueryFilter) []model.ShipmentInvoices
	FindByDriverID(ctx context.Context, driverId int) []model.ShipmentInvoices
	FindByColumns(ctx context.Context, selects []string, columns []string, queries []any) (model.ShipmentInvoices, error)
	UpdateByColumnAt(ctx context.Context, outletId int, data model.ShipmentInvoices) error
	UpdateSkip(ctx context.Context, outletId int, data model.ShipmentInvoices) error
	UpdateByProduct(ctx context.Context, productId int, data model.ShipmentInvoices) error
	UpdateReject(ctx context.Context, Ids []int, data model.ShipmentInvoices) error
	UpdateRejectPartial(ctx context.Context, Id int, data model.ShipmentInvoices, tx ...*gorm.DB) error
	UpdatePickUp(ctx context.Context, Ids []int, data model.ShipmentInvoices) error
	UpdatePickUpPartial(ctx context.Context, Id int, data model.ShipmentInvoices) error
	GetAllReject(ctx context.Context, dataFilter entity.RejectQueryFilter) []model.ShipmentInvoices
	InsertWithTx(tx *gorm.DB, data model.ShipmentInvoices) error
	RejectCancel(ctx context.Context, shipmentInvId []int, data model.ShipmentInvoices) error
	CountReport(ctx context.Context, dataFilter entity.DriverReportQueryFilter) (shipment, finished, skipped, trip int, progress float64, err error)
	GetReport(ctx context.Context, dataFilter entity.DriverReportQueryFilter) (data []entity.SkippedReason, err error)
	FindAllByInvoiceNo(ctx context.Context) []string

	GetListShipmentNo(ctx context.Context) []entity.ShipmentNoDropdown
	GetListReasons(ctx context.Context) []entity.ReasonDropdown
	GetListOutlet(ctx context.Context) []entity.OutletDropdown
	GetListDriver(ctx context.Context) []entity.DriverNameDropdown
	GetListProductCode(ctx context.Context) []entity.ProductCodeDropdown
	GetShipmentReportSummary(ctx context.Context, dataFilter entity.ShipmentReportQueryFilter) (data []model.Shipment, err error)
	GetShipmentReportDetail(ctx context.Context, dataFilter entity.ShipmentReportDetailQueryFilter) (data []model.Shipment, err error)
	GetShipmentReportReject(ctx context.Context, dataFilter entity.ShipmentReportRejectlQueryFilter) (data []model.Shipment, err error)

	FindAllOrderNoByShipmentNo(ctx context.Context, shipmentNo string, outletId int) (result []string)
	FindAllOrderNoById(ctx context.Context, shipmentInvId []int) (result []string)
	GetAllOrderNoByShipmentNo(ctx context.Context, shipmentNo string) (result []string)
	FindByShipmentNo(ctx context.Context, shipmentNo string) ([]model.ShipmentInvoices, error)
	GetAllOrderNo(ctx context.Context) []model.ShipmentInvoices

	UpdateColumnAt(ctx context.Context, column string, currentTime *int, outletId int, data model.ShipmentInvoices)
	UpdateByColumnAtUnload(ctx context.Context, outletId int, data model.ShipmentInvoices) error

	FindTodoList(ctx context.Context, outletId int, shipmentNo string) (entity.TravelListResponse, error)
	BeginTx(ctx context.Context) (*gorm.DB, error)
}
