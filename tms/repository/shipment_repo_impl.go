package repository

import (
	"context"
	"errors"
	"fmt"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"strings"

	"gorm.io/gorm"
)

type ShipmentRepoImpl struct {
	Db *gorm.DB
}

func NewShipmentRepoImpl(db *gorm.DB) ShipmentRepo {
	return &ShipmentRepoImpl{Db: db}
}

func (repo *ShipmentRepoImpl) Insert(ctx context.Context, data model.Shipment) error {
	result := repo.Db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return fmt.Errorf("failed to insert shipment: %w", result.Error)
	}
	return nil
}

func (repo *ShipmentRepoImpl) GetLastShipmentNo() (string, error) {
	var data model.Shipment

	if err := repo.Db.Last(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}

	return data.ShipmentNo, nil
}

func (repo *ShipmentRepoImpl) FindByShipmentNo(ctx context.Context, shipmentNo string) (model.Shipment, error) {
	var data model.Shipment

	result := repo.Db.WithContext(ctx).
		Select(
			"tms.shipments.id",
			"tms.shipments.shipment_no",
			"tms.shipments.vehicle_no",
			"tms.shipments.vehicle_type",
			"tms.shipments.vehicle_name",
			"tms.shipments.driver_id",
			"tms.shipments.driver_name",
			"tms.shipments.helper_id",
			"tms.shipments.helper_name",
			"tms.shipments.driver_name",
			"tms.shipments.vehicle_type",
			"tms.shipments.delivery_date",
			"tms.shipments.created_at",
			"tms.shipments.status",
			"tms.shipments.start",
		).
		Where("shipment_no = ?", shipmentNo).
		Preload("ShipmentInvoices").
		First(&data)

	if result.Error != nil {
		return data, result.Error
	}

	if result.RowsAffected == 0 {
		return data, errors.New("shipment not found")
	}

	return data, nil
}

func (repo *ShipmentRepoImpl) FindAll(ctx context.Context, dataFilter entity.ShipmentQueryFilter) []model.Shipment {
	var data []model.Shipment

	query := repo.Db.Model(&data).
		Select(
			"tms.shipments.id",
			"tms.shipments.shipment_no",
			"tms.shipments.vehicle_id",
			"tms.shipments.vehicle_no",
			"tms.shipments.vehicle_type",
			"tms.shipments.vehicle_name",
			"tms.shipments.driver_id",
			"tms.shipments.driver_name",
			"tms.shipments.helper_id",
			"tms.shipments.helper_name",
			"tms.shipments.driver_name",
			"tms.shipments.delivery_date",
			"tms.shipments.created_at",
			"tms.shipments.status",
			"tms.shipments.shipment_type",
		).
		Preload("ShipmentInvoices")

	if dataFilter.OutletName != "" {
		query = query.Joins("LEFT JOIN tms.shipment_invoices ON tms.shipment_invoices.shipment_no = tms.shipments.shipment_no").
			Where("tms.shipment_invoices.outlet_name ILIKE ?", "%"+dataFilter.OutletName+"%")
	}

	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		query = query.Where("tms.shipments.delivery_date BETWEEN ? AND ?", dataFilter.StartDate, dataFilter.EndDate)
	}

	if dataFilter.DriverID != 0 {
		query = query.Where("tms.shipments.driver_id = ?", dataFilter.DriverID)
	}

	if dataFilter.VehicleID != 0 {
		query = query.Where("tms.shipments.vehicle_id = ?", dataFilter.VehicleID)
	}

	if dataFilter.DriverName != "" {
		query = query.Where("tms.shipments.driver_name = ?", dataFilter.DriverName)
	}

	if dataFilter.ShipmentNo != "" {
		query = query.Where("tms.shipments.shipment_no = ?", dataFilter.ShipmentNo)
	}

	if dataFilter.DeliveryDate != "" {
		query = query.Where("tms.shipments.delivery_date = ?", dataFilter.DeliveryDate)
	}

	if dataFilter.CustID != "" {
		query = query.Where("tms.shipments.cust_id = ?", dataFilter.CustID)
	}

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
		query = query.Order("tms.shipments.id DESC")
	}

	result := query.WithContext(ctx).Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo *ShipmentRepoImpl) UpdateByQuery(ctx context.Context, column string, query any, data model.Shipment) error {
	result := repo.Db.WithContext(ctx).Model(&data).Where(column+" = ?", query).Where("shipment_no = ?", data.ShipmentNo).Debug().Updates(&data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentRepoImpl) CountSummary(ctx context.Context, query int) (shipment, trip, finished, inProgress int, err error) {
	var data struct {
		Trip       int
		Shipment   int
		Finished   int
		InProgress int
	}

	queryRaw := `WITH shipment_data AS (
		SELECT
			s.delivery_date,
            s.status,
            s.driver_id
		FROM tms.shipments s
		WHERE s.driver_id = ? AND s.delivery_date = CURRENT_DATE
	)
	SELECT
		COUNT(*) AS shipment,
		COUNT(CASE WHEN status = 'Finished' THEN 1 END) AS finished,
		COUNT(CASE WHEN status = 'In Progress' THEN 1 END) AS in_progress
	FROM shipment_data;`

	err = repo.Db.WithContext(ctx).Raw(queryRaw, query).Scan(&data).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	queryTrip := `WITH shipment_data AS (
		SELECT
			si.outlet_id,
			s.delivery_date
		FROM tms.shipment_invoices si
        JOIN tms.shipments s ON si.shipment_no = s.shipment_no
		WHERE s.delivery_date = CURRENT_DATE
	)
	SELECT
		COUNT(DISTINCT outlet_id) AS trip
	FROM shipment_data;`

	err = repo.Db.WithContext(ctx).Raw(queryTrip).Scan(&data.Trip).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return data.Shipment, data.Trip, data.Finished, data.InProgress, nil
}

func (repo *ShipmentRepoImpl) FindByDriverID(ctx context.Context, driverId int) []model.Shipment {
	var data []model.Shipment
	result := repo.Db.WithContext(ctx).
		Where("driver_id = ?", driverId).
		Where("delivery_date = CURRENT_DATE").
		Order("id ASC").
		Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo *ShipmentRepoImpl) DeleteByQuery(ctx context.Context, column string, value any) error {
	var data model.Shipment
	result := repo.Db.WithContext(ctx).Where(column+" = ?", value).Delete(&data)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("shipment is not found")
	}

	return nil
}

func (repo *ShipmentRepoImpl) DeleteBulk(ctx context.Context, shipmentNo []string) error {
	var data model.Shipment
	result := repo.Db.WithContext(ctx).Where("shipment_no IN ?", shipmentNo).Delete(&data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ShipmentRepoImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := repo.Db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (repo *ShipmentRepoImpl) InsertWithTx(tx *gorm.DB, data model.Shipment) error {
	return tx.Create(&data).Error
}

func (repo *ShipmentRepoImpl) FindByColumns(ctx context.Context, selects []string, columns []string, queries []any) (model.Shipment, error) {
	if len(columns) != len(queries) {
		return model.Shipment{}, errors.New("columns and queries length mismatch")
	}

	var data model.Shipment
	db := repo.Db.WithContext(ctx)
	if len(selects) == 0 {
		selects = []string{"tms.shipments.*"}
	}

	db = db.Table("tms.shipments").Select(selects)
	for i, column := range columns {
		query := queries[i]
		if column == "delivery_date" && query == "CURRENT_DATE" {
			db = db.Where(column + " = CURRENT_DATE")
		} else {
			db = db.Where(column+" = ?", query)
		}
	}
	result := db.First(&data)

	if result.Error != nil {
		return data, errors.New("shipments not found")
	}

	return data, nil
}

func (repo *ShipmentRepoImpl) FindByEndTimes(ctx context.Context, columns []string, queries []any) (data []model.Shipment, err error) {
	if len(columns) != len(queries) {
		return nil, errors.New("columns and queries length mismatch")
	}

	db := repo.Db.WithContext(ctx)
	for i, column := range columns {
		query := queries[i]
		if column == "delivery_date" && query == "CURRENT_DATE" {
			db = db.Where(column + " = CURRENT_DATE")
		} else {
			db = db.Where(column+" = ?", query)
		}
	}

	result := db.Where("status = ?", "Finished").Order("start ASC").Find(&data)

	if result.Error != nil {
		return data, errors.New("shipments not found")
	}

	return data, nil
}

func (repo *ShipmentRepoImpl) FindByStartTimes(ctx context.Context, columns []string, queries []any) (data []model.Shipment, err error) {
	if len(columns) != len(queries) {
		return nil, errors.New("columns and queries length mismatch")
	}

	db := repo.Db.WithContext(ctx)
	for i, column := range columns {
		query := queries[i]
		if column == "delivery_date" && query == "CURRENT_DATE" {
			db = db.Where(column + " = CURRENT_DATE")
		} else {
			db = db.Where(column+" = ?", query)
		}
	}

	result := db.Order("finish DESC").Find(&data)

	if result.Error != nil {
		return data, errors.New("shipments not found")
	}

	return data, nil
}
