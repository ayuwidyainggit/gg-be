package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"strings"

	"gorm.io/gorm"
)

type ShipmentInvoicesRepoImpl struct {
	Db *gorm.DB
}

func NewShipmentInvoicesRepoImpl(db *gorm.DB) ShipmentInvoicesRepo {
	return &ShipmentInvoicesRepoImpl{Db: db}
}

func (repo *ShipmentInvoicesRepoImpl) Insert(ctx context.Context, data model.ShipmentInvoices) error {
	result := repo.Db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return fmt.Errorf("failed to insert shipment invoices: %w", result.Error)
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) Update(ctx context.Context, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).WithContext(ctx).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) FindByShipmentNoAndOutletName(ctx context.Context, shipmentNo string, outletName string) ([]model.ShipmentInvoices, error) {
	var data []model.ShipmentInvoices

	result := repo.Db.WithContext(ctx).Where("shipment_no = ?", shipmentNo).Where("outlet_name = ?", outletName).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(data) == 0 {
		return nil, errors.New("record not found")
	}

	return data, nil
}

func (repo *ShipmentInvoicesRepoImpl) FindByOutletId(ctx context.Context, params []interface{}) []entity.OutletResponse {
	var data []entity.OutletResponse

	result := repo.Db.Raw(`
				SELECT 
				si.outlet_id,
				si.outlet_code,
				si.outlet_name,
				si.outlet_address,
				si.outlet_status,
				s.shipment_no,
				COUNT(si.product_id) AS total_product,
				COUNT(CASE WHEN si.status = 'Delivery' THEN 1 END) AS total_product_delivery,
				COUNT(CASE WHEN si.status = 'Pick Up' THEN 1 END) AS total_product_pickup
			FROM
				tms.shipment_invoices si
			JOIN
				tms.shipments s ON s.shipment_no = si.shipment_no
			WHERE
				s.driver_id = ?
			AND
				si.outlet_id = ?
			AND
				s.shipment_no = ?
			GROUP BY
				si.outlet_id,
				si.outlet_code,
				si.outlet_name,
				si.outlet_address,
				si.outlet_status,
				s.shipment_no;
    `, params...).WithContext(ctx).Scan(&data)
	helper.ErrorPanic(result.Error)

	return data
}

func (repo *ShipmentInvoicesRepoImpl) FindAll(ctx context.Context, dataFilter entity.ShipmentInvoicesQueryFilter) []model.ShipmentInvoices {
	var data []model.ShipmentInvoices

	query := repo.Db.Model(&data).Preload("Shipment")

	if dataFilter.SalesmanName != "" {
		query = query.Where("tms.shipment_invoices.salesman_name ILIKE ?", "%"+dataFilter.SalesmanName+"%")
	}

	if dataFilter.OutletName != "" {
		query = query.Where("tms.shipment_invoices.outlet_name ILIKE ?", "%"+dataFilter.OutletName+"%")
	}

	if dataFilter.OutletCode != "" {
		query = query.Where("tms.shipment_invoices.outlet_code ILIKE ?", "%"+dataFilter.OutletCode+"%")
	}

	if dataFilter.OutletID != 0 {
		query = query.Where("tms.shipment_invoices.outlet_id = ?", +dataFilter.OutletID)
	}

	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		query = query.Where("tms.shipment_invoices.delivery_date BETWEEN ? AND ?", dataFilter.StartDate, dataFilter.EndDate)
	}

	if dataFilter.DriverID != 0 {
		query = query.Joins("LEFT JOIN tms.shipments ON tms.shipments.shipment_no = tms.shipment_invoices.shipment_no").
			Where("tms.shipments.driver_id = ?", dataFilter.DriverID)
	}

	if dataFilter.ShipmentNo != "" {
		query = query.Where("tms.shipment_invoices.shipment_no ILIKE ? ", "%"+dataFilter.ShipmentNo+"%")
	}

	if dataFilter.CustID != "" {
		query = query.Where("tms.shipment_invoices.cust_id ILIKE ? ", "%"+dataFilter.CustID+"%")
	}

	if dataFilter.ProductID != 0 {
		query = query.Where("tms.shipment_invoices.product_id = ?", +dataFilter.ProductID)
	}

	if dataFilter.ProductName != "" {
		query = query.Where("tms.shipment_invoices.product_name ILIKE ?", "%"+dataFilter.ProductName+"%")
	}

	if dataFilter.Status != "" {
		query = query.Where("REPLACE(tms.shipment_invoices.status, ' ', '') ILIKE ?", "%"+strings.ReplaceAll(dataFilter.Status, " ", "")+"%")
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
		query = query.Order("tms.shipment_invoices.id DESC")
	}

	result := query.WithContext(ctx).Find(&data)
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}
	return data
}

func (repo *ShipmentInvoicesRepoImpl) FindByDriverID(ctx context.Context, driverId int) []model.ShipmentInvoices {
	var data []model.ShipmentInvoices
	result := repo.Db.WithContext(ctx).Where("driver_id = ?", driverId).Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo *ShipmentInvoicesRepoImpl) FindByColumns(ctx context.Context, selects []string, columns []string, queries []any) (model.ShipmentInvoices, error) {
	if len(columns) != len(queries) {
		return model.ShipmentInvoices{}, errors.New("columns and queries length mismatch")
	}

	var data model.ShipmentInvoices
	db := repo.Db.WithContext(ctx)
	if len(selects) == 0 {
		selects = []string{"tms.shipment_invoices.*"}
	}

	db = db.Table("tms.shipment_invoices").Select(selects)
	for i, column := range columns {
		query := queries[i]
		if column == "delivery_date" && query == "CURRENT_DATE" {
			db = db.Where(column + " = CURRENT_DATE")
		} else {
			db = db.Where(column+" = ?", query)
		}
	}

	// Prioritas data dengan unload_at atau pickup_at yang ada
	db = db.Order("unload_at IS NULL, pickup_at IS NULL") // NULL values sorted last.

	result := db.First(&data)

	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *ShipmentInvoicesRepoImpl) UpdateByColumnAt(ctx context.Context, outletId int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).
		Where("outlet_id = ?", outletId).
		Where("shipment_no = ?", data.ShipmentNo).
		Updates(data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (repo *ShipmentInvoicesRepoImpl) UpdateByColumnAtUnload(ctx context.Context, outletId int, data model.ShipmentInvoices) error {
	var existingRecord model.ShipmentInvoices
	result := repo.Db.WithContext(ctx).
		Where("outlet_id = ?", outletId).
		Where("shipment_no = ?", data.ShipmentNo).
		First(&existingRecord)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	// Update only if product_status is NULL or "-"
	result = repo.Db.Model(&data).
		WithContext(ctx).
		Where("outlet_id = ?", outletId).
		Where("shipment_no = ?", data.ShipmentNo).
		Where("product_status IS NULL OR product_status = '-'").
		Updates(data)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ShipmentInvoicesRepoImpl) UpdateSkip(ctx context.Context, outletId int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).
		Where("outlet_id = ?", outletId).
		Where("shipment_no = ?", data.ShipmentNo).
		Select("shipment_no", "outlet_status", "skip_at", "skip_reason", "in_outlet", "updated_at").
		Updates(data)

	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (repo *ShipmentInvoicesRepoImpl) FindByTwoColumns(ctx context.Context, columnFirst, columnSecond string, queryFirst, querySecond any) ([]model.ShipmentInvoices, error) {
	var data []model.ShipmentInvoices
	result := repo.Db.WithContext(ctx).Where(columnFirst+" = ?", queryFirst).Where(columnSecond+" = ?", querySecond).Find(&data)

	if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if result.Error != nil {
		return nil, errors.New("shipment invoices not found")
	}

	return data, nil
}

func (repo *ShipmentInvoicesRepoImpl) FindByOneColumn(ctx context.Context, column string, query any) (model.ShipmentInvoices, error) {
	var data model.ShipmentInvoices
	result := repo.Db.WithContext(ctx).Where(column+" = ?", query).Debug().Find(&data)

	if result.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	if result.Error != nil {
		return data, errors.New("shipment invoices not found")
	}

	return data, nil
}

func (repo *ShipmentInvoicesRepoImpl) GetAllReject(ctx context.Context, dataFilter entity.RejectQueryFilter) []model.ShipmentInvoices {
	var data []model.ShipmentInvoices

	query := repo.Db.Table("tms.shipment_invoices").
		Select("tms.shipment_invoices.id, tms.shipment_invoices.product_id, tms.shipment_invoices.product_name, tms.shipment_invoices.product_status, tms.shipment_invoices.sku, tms.shipment_invoices.qty1, tms.shipment_invoices.qty2, tms.shipment_invoices.qty3, tms.shipment_invoices.conv_unit1, tms.shipment_invoices.conv_unit2, tms.shipment_invoices.conv_unit3, tms.shipment_invoices.unit_id1, tms.shipment_invoices.unit_id2, tms.shipment_invoices.unit_id3, tms.shipment_invoices.reason_id, tms.shipment_invoices.reason_name, tms.shipment_invoices.outlet_id, tms.shipments.driver_id, tms.shipment_invoices.qty_reject_1, tms.shipment_invoices.qty_reject_2, tms.shipment_invoices.qty_reject_3").
		Joins("LEFT JOIN tms.shipments ON tms.shipment_invoices.shipment_no = tms.shipments.shipment_no")

	if dataFilter.DriverID != 0 {
		query = query.Where("tms.shipments.driver_id = ?", dataFilter.DriverID)
	}

	if dataFilter.ReasonID != 0 {
		query = query.Where("tms.shipment_invoices.reason_id = ?", dataFilter.ReasonID)
	}

	if dataFilter.OutletID != 0 {
		query = query.Where("tms.shipment_invoices.outlet_id = ?", dataFilter.OutletID)
	}

	if dataFilter.ProductName != "" {
		query = query.Where("tms.shipment_invoices.product_name ILIKE ?", "%"+dataFilter.ProductName+"%")
	}

	if dataFilter.ShipmentNo != "" {
		query = query.Where("tms.shipment_invoices.shipment_no ILIKE ?", "%"+dataFilter.ShipmentNo+"%")
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
		query = query.Order("tms.shipment_invoices.id DESC")
	}

	result := query.WithContext(ctx).Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo *ShipmentInvoicesRepoImpl) UpdateByProduct(ctx context.Context, productId int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).Where("product_id = ?", productId).Select("qty").Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) InsertWithTx(tx *gorm.DB, data model.ShipmentInvoices) error {
	return tx.Create(&data).Error
}

func (repo *ShipmentInvoicesRepoImpl) UpdateReject(ctx context.Context, Ids []int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).Where("id IN ?", Ids).Select("product_status", "reason_id", "reason_name", "unload_at").Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) UpdateRejectPartial(ctx context.Context, Id int, data model.ShipmentInvoices, tx ...*gorm.DB) error {
	db := repo.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	result := db.Model(&data).
		WithContext(ctx).Where("id = ?", Id).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) UpdatePickUp(ctx context.Context, Ids []int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).Where("id IN ?", Ids).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) UpdatePickUpPartial(ctx context.Context, Id int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).Where("id = ?", Id).Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) RejectCancel(ctx context.Context, shipmentInvId []int, data model.ShipmentInvoices) error {
	result := repo.Db.Model(&data).
		WithContext(ctx).Where("id IN ?", shipmentInvId).Select("reason_id", "reason_name").Updates(data)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ShipmentInvoicesRepoImpl) CountReport(ctx context.Context, dataFilter entity.DriverReportQueryFilter) (shipment, finished, skipped, trip int, progress float64, err error) {
	var query string

	reportToday := `s.delivery_date = CURRENT_DATE`
	reportMonth := `EXTRACT(MONTH FROM si.delivery_date) = EXTRACT(MONTH FROM CURRENT_DATE)`

	rawQuery := `
		WITH countDriverReport AS (
			SELECT
				si.shipment_no,
				s.delivery_date,
				s.status,
				s.driver_id,
				si.outlet_status,
				si.outlet_id,
				si.leave_at,
				COUNT(DISTINCT si.outlet_id) AS total_trip
			FROM tms.shipments s
			JOIN tms.shipment_invoices si ON s.shipment_no = si.shipment_no
			WHERE
				s.driver_id = ?
				AND %s
			GROUP BY
				si.shipment_no, s.delivery_date, s.status, s.driver_id, si.outlet_status, si.outlet_id, si.leave_at
		)
		SELECT
			COUNT(DISTINCT shipment_no) AS shipment,
			COUNT(DISTINCT CASE WHEN outlet_status = 'Finished' THEN shipment_no END) AS finished,
			COUNT(DISTINCT CASE WHEN outlet_status = 'Skipped' THEN shipment_no END) AS skipped,
			COUNT(total_trip) AS trip,
			CASE
				WHEN COUNT(total_trip) = 0 THEN 0
				ELSE ROUND(COUNT(DISTINCT CASE WHEN outlet_status = 'Finished' THEN shipment_no END) * 100.0 / COUNT(total_trip), 2)
			END AS progress
		FROM countDriverReport;
	`

	switch dataFilter.Period {
	case "today":
		query = fmt.Sprintf(rawQuery, reportToday)
	case "month":
		query = fmt.Sprintf(rawQuery, reportMonth)
	default:
		return 0, 0, 0, 0, 0, nil
	}

	rows, err := repo.Db.WithContext(ctx).Raw(query, dataFilter.DriverID).Rows()
	if err != nil {
		return 0, 0, 0, 0, 0, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&shipment, &finished, &skipped, &trip, &progress)
		if err != nil {
			return 0, 0, 0, 0, 0, err
		}
	}

	return shipment, finished, skipped, trip, progress, nil
}

func (repo *ShipmentInvoicesRepoImpl) GetReport(ctx context.Context, dataFilter entity.DriverReportQueryFilter) (data []entity.SkippedReason, err error) {
	var query string

	reportToday := `si.delivery_date = CURRENT_DATE`
	reportMonth := `EXTRACT(MONTH FROM si.delivery_date) = EXTRACT(MONTH FROM CURRENT_DATE)`
	reasonMonth := `AND si.skip_reason IS NOT NULL`

	rawQuery := `
		SELECT
			si.skip_reason,
			si.delivery_date,
			COUNT(DISTINCT CASE WHEN si.skip_reason IS NOT NULL THEN si.outlet_id END) AS count,
            ROUND(COUNT(DISTINCT CASE WHEN si.skip_reason IS NOT NULL THEN si.outlet_id END) * 100.0 / COUNT(DISTINCT CASE WHEN si.outlet_status = 'Skipped' THEN si.shipment_no END), 2) AS percentage
		FROM
			tms.shipment_invoices si
		JOIN
			tms.shipments s ON si.shipment_no = s.shipment_no
		WHERE
			s.driver_id = ?
			AND si.outlet_status = 'Skipped'
		    %s
			AND %s
		GROUP BY
			si.skip_reason, si.delivery_date
		ORDER BY
			count DESC;
	`

	switch dataFilter.Period {
	case "today":
		query = fmt.Sprintf(rawQuery, "", reportToday)
	case "month":
		query = fmt.Sprintf(rawQuery, reasonMonth, reportMonth)
	}

	rows, err := repo.Db.WithContext(ctx).Raw(query, dataFilter.DriverID).Rows()
	if err != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var reason entity.SkippedReason
		err := rows.Scan(&reason.SkipReason, &reason.DeliveryDate, &reason.Count, &reason.Percentage)
		//fmt.Printf("SkipReason: %d, DeliverDate: %d, OutletStatus: %d, Count: %d, Percentage: %d\n", *reason.ReasonID, *reason.ReasonName, reason.DeliveryDate, reason.OutletStatus, reason.Count, reason.Percentage)
		if err != nil {
			return nil, err
		}
		data = append(data, reason)
	}

	return data, nil
}

func (repo ShipmentInvoicesRepoImpl) FindAllOrderNoByShipmentNo(ctx context.Context, shipmentNo string, outletId int) (result []string) {
	var orderNos []string

	queryResult := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("DISTINCT order_no").
		Where("shipment_no = ? AND outlet_id = ?", shipmentNo, outletId).
		Debug().
		Pluck("order_no", &orderNos)

	helper.ErrorPanic(queryResult.Error)

	return orderNos
}

func (repo ShipmentInvoicesRepoImpl) FindAllOrderNoById(ctx context.Context, shipmentInvId []int) (result []string) {
	var orderNos []string

	queryResult := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("DISTINCT order_no").
		Where("id IN ?", shipmentInvId).
		Debug().
		Pluck("order_no", &orderNos)

	helper.ErrorPanic(queryResult.Error)

	return orderNos
}

func (repo ShipmentInvoicesRepoImpl) FindAllByInvoiceNo(ctx context.Context) []string {
	var invoice_no []string

	result := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("invoice_no").
		Debug().
		Find(&invoice_no)

	helper.ErrorPanic(result.Error)
	return invoice_no

}

func (repo ShipmentInvoicesRepoImpl) GetListShipmentNo(ctx context.Context) []entity.ShipmentNoDropdown {
	var shipment_no []entity.ShipmentNoDropdown

	result := repo.Db.WithContext(ctx).
		Table("tms.shipments").
		Select("DISTINCT(shipment_no)").
		Debug().
		Find(&shipment_no)

	helper.ErrorPanic(result.Error)
	return shipment_no
}

func (repo ShipmentInvoicesRepoImpl) GetListReasons(ctx context.Context) []entity.ReasonDropdown {
	var reasonName []entity.ReasonDropdown

	result := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("DISTINCT(reason_name)").
		Debug().
		Find(&reasonName)

	helper.ErrorPanic(result.Error)
	return reasonName
}

func (repo ShipmentInvoicesRepoImpl) GetListOutlet(ctx context.Context) []entity.OutletDropdown {
	var outlet []entity.OutletDropdown

	result := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("DISTINCT outlet_code, outlet_name").
		Debug().
		Find(&outlet)

	helper.ErrorPanic(result.Error)
	return outlet
}

func (repo ShipmentInvoicesRepoImpl) GetListDriver(ctx context.Context) []entity.DriverNameDropdown {
	var driver []entity.DriverNameDropdown

	result := repo.Db.WithContext(ctx).
		Table("tms.shipments").
		Select("DISTINCT(driver_name)").
		Debug().
		Find(&driver)

	helper.ErrorPanic(result.Error)
	return driver
}

func (repo ShipmentInvoicesRepoImpl) GetListProductCode(ctx context.Context) []entity.ProductCodeDropdown {
	var product []entity.ProductCodeDropdown

	result := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("DISTINCT(product_code)").
		Debug().
		Find(&product)

	helper.ErrorPanic(result.Error)
	return product
}

func (repo ShipmentInvoicesRepoImpl) GetAllOrderNoByShipmentNo(ctx context.Context, shipmentNo string) []string {
	var orderNos []string

	queryResult := repo.Db.WithContext(ctx).
		Table("tms.shipment_invoices").
		Select("DISTINCT order_no").
		Where("shipment_no = ?", shipmentNo).
		Debug().
		Find(&orderNos)

	helper.ErrorPanic(queryResult.Error)

	return orderNos
}

func (repo ShipmentInvoicesRepoImpl) FindByShipmentNo(ctx context.Context, shipmentNo string) ([]model.ShipmentInvoices, error) {
	var data []model.ShipmentInvoices

	result := repo.Db.WithContext(ctx).
		Table("(SELECT DISTINCT ON (order_no) * FROM tms.shipment_invoices WHERE shipment_no = ?) as distinct_invoices", shipmentNo).
		Find(&data)

	if result.Error != nil {
		return data, result.Error
	}

	if result.RowsAffected == 0 {
		return data, errors.New("shipment not found")
	}

	return data, nil
}

func (repo ShipmentInvoicesRepoImpl) GetAllOrderNo(ctx context.Context) []model.ShipmentInvoices {
	var data []model.ShipmentInvoices

	querySql := `
		SELECT order_no, product_status
		FROM tms.shipment_invoices
		WHERE order_no IN (
			SELECT order_no
			FROM tms.shipment_invoices
			GROUP BY order_no
			HAVING (
				-- Case 1: Contains '-' in product_status
				(COUNT(DISTINCT product_status) = 1 AND MAX(product_status) = '-')
				OR
				-- Case 2: All product_status are 'Receive'
				(COUNT(DISTINCT product_status) = 1 AND MAX(product_status) = 'Receive')
				OR
				-- Case 3: All product_status are 'Reject Partial'
				(COUNT(DISTINCT product_status) = 1 AND MAX(product_status) = 'Reject Partial')
				OR
				-- Case 4: All product_status are 'Reject All'
				(COUNT(DISTINCT product_status) = 1 AND MAX(product_status) = 'Reject All')
			)
		);`

	result := repo.Db.WithContext(ctx).Raw(querySql).Scan(&data)

	helper.ErrorPanic(result.Error)

	return data
}

// TODO Shipment Report Repo Impl

func (repo ShipmentInvoicesRepoImpl) GetShipmentReportSummary(ctx context.Context, dataFilter entity.ShipmentReportQueryFilter) (data []model.Shipment, err error) {

	var shipments []model.Shipment

	query := repo.Db.Model(&shipments).
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
			"tms.shipments.delivery_date",
			"tms.shipments.created_at",
			"tms.shipments.status",
			"tms.shipments.shipment_type",
			"tms.shipments.start",
			"tms.shipments.finish",
		).
		Preload("ShipmentInvoices")

	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		query = query.Where("tms.shipments.delivery_date BETWEEN ? AND ?", dataFilter.StartDate, dataFilter.EndDate)
	}

	if dataFilter.ShipmentNo != "" {
		query = query.Where("tms.shipments.shipment_no = ?", dataFilter.ShipmentNo)
	}

	if dataFilter.DriverName != "" {
		query = query.Where("tms.shipments.driver_name = ?", dataFilter.DriverName)
	}

	result := query.WithContext(ctx).Find(&shipments)
	if err := result.Error; err != nil {
		log.Println("Error fetching shipment report summary:", err)
		return nil, err
	}

	return shipments, nil

}

func (repo ShipmentInvoicesRepoImpl) GetShipmentReportDetail(ctx context.Context, dataFilter entity.ShipmentReportDetailQueryFilter) (data []model.Shipment, err error) {
	var shipments []model.Shipment

	query := repo.Db.Model(&shipments).
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
			"tms.shipments.delivery_date",
			"tms.shipments.created_at",
			"tms.shipments.status",
			"tms.shipments.shipment_type",
			"tms.shipments.start",
			"tms.shipments.finish",
		).
		Preload("ShipmentInvoices")

	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		query = query.Where("tms.shipments.delivery_date BETWEEN ? AND ?", dataFilter.StartDate, dataFilter.EndDate)
	}

	if dataFilter.ShipmentNo != "" {
		query = query.Where("tms.shipments.shipment_no = ?", dataFilter.ShipmentNo)
	}

	if dataFilter.DriverName != "" {
		query = query.Where("tms.shipments.driver_name = ?", dataFilter.DriverName)
	}

	if dataFilter.VisitedStatus != "" {
		query = query.Where("tms.shipments.status = ?", dataFilter.VisitedStatus)
	}

	if dataFilter.ReceivedStatus != "" {
		query = query.Where("tms.shipments_invoices.status = ?", dataFilter.ReceivedStatus)
	}

	if dataFilter.OutletName != "" {
		query = query.Where("tms.shipments_invoices.outlet_name = ?", dataFilter.OutletName)
	}

	result := query.WithContext(ctx).Find(&shipments)
	if err := result.Error; err != nil {
		log.Println("Error fetching shipment report summary:", err)
		return nil, err
	}

	return shipments, nil
}

func (repo ShipmentInvoicesRepoImpl) GetShipmentReportReject(ctx context.Context, dataFilter entity.ShipmentReportRejectlQueryFilter) (data []model.Shipment, err error) {
	var shipments []model.Shipment

	query := repo.Db.Model(&shipments).
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
			"tms.shipments.delivery_date",
			"tms.shipments.created_at",
			"tms.shipments.status",
			"tms.shipments.shipment_type",
			"tms.shipments.start",
			"tms.shipments.finish",
			"tms.shipment_invoices.qty_reject_1",
			"tms.shipment_invoices.qty_reject_2",
			"tms.shipment_invoices.qty_reject_3",
		).
		Joins("JOIN tms.shipment_invoices ON tms.shipments.shipment_no = tms.shipment_invoices.shipment_no").
		Where("tms.shipment_invoices.product_status LIKE ?", "%Reject%").
		Preload("ShipmentInvoices")

	if dataFilter.StartDate != "" && dataFilter.EndDate != "" {
		query = query.Where("tms.shipments.delivery_date BETWEEN ? AND ?", dataFilter.StartDate, dataFilter.EndDate)
	}

	if dataFilter.ShipmentNo != "" {
		query = query.Where("tms.shipments.shipment_no = ?", dataFilter.ShipmentNo)
	}

	if dataFilter.DriverName != "" {
		query = query.Where("tms.shipments.driver_name = ?", dataFilter.DriverName)
	}

	if dataFilter.OutletName != "" {
		query = query.Where("tms.shipments_invoices.outlet_name = ?", dataFilter.OutletName)
	}

	if dataFilter.ProductCode != "" {
		query = query.Where("tms.shipments_invoices.product_code = ?", dataFilter.ProductCode)
	}

	if dataFilter.Reason != "" {
		query = query.Where("tms.shipments_invoices.reason_name = ?", dataFilter.Reason)
	}

	result := query.WithContext(ctx).Find(&shipments)
	if err := result.Error; err != nil {
		log.Println("Error fetching shipment report summary:", err)
		return nil, err
	}

	return shipments, nil
}

func (repo *ShipmentInvoicesRepoImpl) UpdateColumnAt(ctx context.Context, column string, currentTime *int, outletId int, data model.ShipmentInvoices) {

	result := repo.Db.WithContext(ctx).Exec(`
		UPDATE tms.shipment_invoices
		SET `+column+` = ?
		WHERE outlet_id = ? AND shipment_no = ?
	`, currentTime, outletId, data.ShipmentNo)

	helper.ErrorPanic(result.Error)
}

func (repo *ShipmentInvoicesRepoImpl) FindTodoList(ctx context.Context, outletId int, shipmentNo string) (entity.TravelListResponse, error) {

	var data entity.TravelListResponse

	result := repo.Db.Raw(`
				SELECT
				order_no,
				COUNT(*) AS product_count,
				CASE
					WHEN MAX(unload_at) IS NULL THEN 'Didnt unload yet'
					WHEN MAX(product_status) = 'Reject Partial' THEN 'Some product rejected by outlet'
					ELSE '10 ' || COUNT(*) || ' Product rejected by outlet'
				END AS unload_desc,
				MAX(arrive_at) AS arrive_at,
				BOOL_OR(in_outlet) AS in_outlet,
				MAX(leave_at) AS leave_at,
				MAX(on_hold) AS on_hold,
				MAX(pickup_at) AS pickup_at,
				MAX(resume_at) AS resume_at,
				MAX(skip_at) AS skip_at,
				MAX(skip_reason) AS skip_reason,
				MAX(unload_at) AS unload_at
			FROM
				tms.shipment_invoices
			WHERE
				outlet_id = ?
				AND shipment_no = ?
			GROUP BY
				order_no
    `, outletId, shipmentNo).WithContext(ctx).Scan(&data)
	helper.ErrorPanic(result.Error)

	return data, nil
}


func (repo *ShipmentInvoicesRepoImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := repo.Db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}