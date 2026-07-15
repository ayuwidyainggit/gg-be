package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/str"
	"mobile/pkg/times"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryOrderImpl struct {
		*gorm.DB
	}
)
type OrderRepository interface {
	Store(c context.Context, data *model.Order) error
	StoreNoOrder(c context.Context, data *model.NoOrder) error
	StoreDetail(c context.Context, data *model.OrderDetail) error
	StoreStock(c context.Context, data *model.Stock) error
	FindByNo(RoNo string, custId string) (realOrder model.OrderList, err error)
	FindDetail(RoNo string, custId string) (details []model.OrderDetailRead, err error)
	FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error)
	FindAllNoOrderByCustId(dataFilter entity.NoOrderQueryFilter) ([]model.NoOrderList, int64, int, error)
	FindAllOutletBySalesmanId(dataFilter entity.OrderQueryFilter) ([]model.OutletBySalesman, int64, int, error)

	Update(c context.Context, RoNo string, data model.Order) error
	DeleteDetailNotInIDs(c context.Context, RoNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.OrderDetail) error
	Delete(c context.Context, custId string, RoNo string, deletedBy int64) error

	FindOneProductByProductIdAndCustId(productId int64, custId string, parentCustId string) (model.ProductConversion, error)
	FindProductByListID(ctx context.Context, productIDs []int64) (products []model.ProductOrder, err error)
	CountAllRoByCustId(custId string, roDate string) (int, error)

	SummaryTotalBySalesmanAndDate(salesmanID int, startDate string, endDate string, dataStatus []int, parentCustId string) (summaryOrder model.SummaryOrder, err error)
	GetSalesman(custID string, empID int64) (salesman model.Salesman, err error)
	GetNextInvoiceNumber(c context.Context, date, custId string) (int, error)
	GetNextRoNumber(c context.Context, custId string, date string) (int, error)
	GetStockByWarehouseProductCustomer(ctx context.Context, warehouseID int64, productID int64, customerID string) (float64, error)
	UpdateWarehouseStock(c context.Context, whId int64, proId int64, custId string, qty float64) error
	GetInvoiceByNumbers(custID, parentCustID string, empID int64, isCollection bool) ([]model.OrderListInvoice, error)
	GetSummarySales(custID string, salesmanID int, startDate string, endDate string) (float64, error)
}

func NewOrderRepo(db *gorm.DB) *RepositoryOrderImpl {
	return &RepositoryOrderImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderImpl) Store(c context.Context, data *model.Order) error {
	fmt.Println("Ro Number Repo: ", data.RoNo)
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) StoreNoOrder(c context.Context, data *model.NoOrder) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) StoreDetail(c context.Context, data *model.OrderDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) StoreStock(c context.Context, data *model.Stock) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderImpl) FindByNo(roNo string, custId string) (realOrder model.OrderList, err error) {
	err = repository.
		Select(`sls.order.*, 
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name,
			sls.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.order.ro_no = ? AND sls.order.cust_id=?", roNo, custId).
		Take(&realOrder).Error
	return realOrder, err
}

func (repository *RepositoryOrderImpl) FindDetail(roNo string, custId string) (details []model.OrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("ro_no = ? AND sls.order_detail.cust_id=?", roNo, custId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderImpl) FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderList, int64, int, error) {
	var ro []model.OrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ro_no")
	query := repository.Select(
		`sls.order.*, 
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, ot.outlet_name, 
			sls.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.order.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id=?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.RoFrom != nil && dataFilter.RoTo != nil {
		query.Where("sls.order.ro_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.RoFrom), str.UnixTimestampToUtcTime(*dataFilter.RoTo))
		queryCount.Where("sls.order.ro_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.RoFrom), str.UnixTimestampToUtcTime(*dataFilter.RoTo))
	}

	if dataFilter.InvoiceFrom != nil && dataFilter.InvoiceTo != nil {
		query.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
		queryCount.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.order.ro_no=?", dataFilter.Query)
		query.Where("sls.order.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.order.data_status in ?", dataFilter.Status)
		query.Where("sls.order.data_status in ?", dataFilter.Status)
	}

	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ro_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryOrderImpl) FindAllNoOrderByCustId(dataFilter entity.NoOrderQueryFilter) ([]model.NoOrderList, int64, int, error) {
	var ro []model.NoOrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("no_order_id")
	query := repository.Select(
		`sls.no_order.*, 
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, ot.outlet_name, tor.taking_order_name, tor.image_url,
			sls.sales_name`).
		Joins("left join sys.m_user us on us.user_id = sls.no_order.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.no_order.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.no_order.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_taking_order tor on tor.taking_order_id = sls.no_order.taking_order_id AND tor.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.no_order.cust_id=?", dataFilter.CustId)
	query.Where("sls.no_order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.no_order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.no_order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.no_order.ro_no=?", dataFilter.Query)
		query.Where("sls.no_order.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.no_order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.no_order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.no_order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.no_order.outlet_id in ?", dataFilter.OutletID)
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("no_order_id DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryOrderImpl) Update(c context.Context, RoNo string, data model.Order) error {
	result := repository.model(c).Model(&data).Where("ro_no=?", RoNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryOrderImpl) DeleteDetailNotInIDs(c context.Context, RoNo string, IDs []int64) error {
	var Details model.OrderDetail
	err := repository.model(c).Where("ro_no=? AND order_detail_id not in (?) ", RoNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryOrderImpl) UpdateDetail(c context.Context, Details *model.OrderDetail) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryOrderImpl) Delete(c context.Context, custId string, RoNo string, deletedBy int64) error {
	var data model.Order
	result := repository.model(c).Model(&data).Where("ro_no=? AND cust_id = ? AND is_del= ? ", RoNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderImpl) FindOneProductByProductIdAndCustId(productId int64, custId string, parentCustId string) (productConversion model.ProductConversion, err error) {
	err = repository.
		Select(`mst.m_product.cust_id,
			mst.m_product.pro_id,
			mst.m_product.conv_unit2, 
			mst.m_product.conv_unit3, 
			mst.m_product.conv_unit4, 
			mst.m_product.conv_unit5`).
		Where("mst.m_product.pro_id = ? AND mst.m_product.cust_id=?", productId, parentCustId).
		Take(&productConversion).Error
	return productConversion, err
}

func (repository *RepositoryOrderImpl) CountAllRoByCustId(custId string, roDate string) (int, error) {
	var total int64

	err := repository.
		Model(&model.OrderList{}).
		Where("cust_id = ?", custId).
		Where("ro_date = ?", roDate).
		Count(&total).Error

	return int(total), err
}

func (repository *RepositoryOrderImpl) FindAllOutletBySalesmanId(dataFilter entity.OrderQueryFilter) ([]model.OutletBySalesman, int64, int, error) {
	var ro []model.OutletBySalesman
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Table("sls.order").
		Select("distinct(sls.order.outlet_id)").
		Joins("left join mst.m_outlet mo on mo.outlet_id = sls.order.outlet_id AND mo.cust_id = ?", dataFilter.CustId).
		Where("sls.order.cust_id = ?", dataFilter.CustId)

	query := repository.Table("sls.order").
		Select("distinct(sls.order.outlet_id) as outlet_id, mo.outlet_code, mo.outlet_name").
		Joins("left join mst.m_outlet mo on mo.outlet_id = sls.order.outlet_id AND mo.cust_id = ?", dataFilter.CustId).
		Where("sls.order.cust_id = ?", dataFilter.CustId)

	if dataFilter.Query != "" {
		queryCount = queryCount.Where("(mo.outlet_code ILIKE ? OR mo.outlet_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query = query.Where("(mo.outlet_code ILIKE ? OR mo.outlet_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.order.data_status in ?", dataFilter.Status)
		query.Where("sls.order.data_status in ?", dataFilter.Status)
	}

	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("outlet_id DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryOrderImpl) FindProductByListID(ctx context.Context, productIDs []int64) (products []model.ProductOrder, err error) {
	err = repository.model(ctx).
		Select("*").
		Where("pro_id in ?", productIDs).
		Find(&products).Error

	return products, err
}

func (repository *RepositoryOrderImpl) GetSalesman(custID string, empID int64) (salesman model.Salesman, err error) {
	err = repository.
		Select("*").
		Where("cust_id = ? AND emp_id = ?", custID, empID).
		Take(&salesman).Error

	return salesman, err
}

func (repository *RepositoryOrderImpl) SummaryTotalBySalesmanAndDate(salesmanID int, startDate string, endDate string, dataStatus []int, parentCustId string) (summaryOrder model.SummaryOrder, err error) {

	err = repository.Table("sls.order").
		Select("SUM(total) as total_summary").
		Where("salesman_id = ? AND DATE(ro_date) BETWEEN ? AND ?", salesmanID, startDate, endDate).
		Where("data_status in ?", dataStatus).
		Where("is_del = false").
		Scan(&summaryOrder).Error
	return summaryOrder, err
}

func (repository *RepositoryOrderImpl) GetNextInvoiceNumber(c context.Context, date, custId string) (int, error) {
	// 1. Enforce Transaction
	tx := extractTx(c)
	if tx == nil {
		return 0, errors.New("GetNextInvoiceNumber requires an active transaction")
	}

	// 2. Acquire Lock (Use two-argument advisory lock to avoid collisions)
	// Key 1: hash of 'invoice_no'
	// Key 2: dateInt (YYMMDD)
	dateInt, err := strconv.Atoi(date)
	if err != nil {
		return 0, fmt.Errorf("invalid date format for lock: %v", err)
	}

	// hash('invoice_no') = 1918341611 (arbitrary stable ID)
	const invoiceLockKey = 1918341611
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?, ?)", invoiceLockKey, dateInt).Error; err != nil {
		return 0, err
	}

	// 3. Get Last Number
	var lastNo *string
	prefix := "INV" + date + "%"

	err = tx.Table("sls.order").
		Select("MAX(invoice_no)").
		Where("invoice_no LIKE ?", prefix).
		Where("cust_id = ?", custId).
		Row().Scan(&lastNo)

	if err != nil {
		return 0, err
	}

	if lastNo == nil {
		return 1, nil
	}

	// Extract last 4 digits
	if len(*lastNo) < 4 {
		return 1, nil
	}

	ln := *lastNo
	lastSeqStr := ln[len(ln)-4:]
	lastSeq, err := strconv.Atoi(lastSeqStr)
	if err != nil {
		return 1, nil
	}

	return lastSeq + 1, nil
}

func (repository *RepositoryOrderImpl) GetNextRoNumber(c context.Context, custId string, date string) (int, error) {
	// 1. Enforce Transaction
	tx := extractTx(c)
	if tx == nil {
		return 0, errors.New("GetNextRoNumber requires an active transaction")
	}

	// 2. Acquire Lock
	// We use a lock per customer and date to ensure serial generation
	// hash(custId + date) or similar.
	// For simplicity, we can use a hash of (custId_int + dateInt)
	// Since custId is string, let's use pg_advisory_xact_lock with 64-bit hash
	lockKey := str.HashStringToInt64(fmt.Sprintf("rono_%s_%s", custId, date))

	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
		return 0, err
	}

	// 3. Get Total RO for this customer and date
	// Original logic used Total RO + 1. We'll use the same but inside the lock.
	var total int64
	// date format in DB seems to be YYYY-MM-DD for ro_date
	err := tx.Table("sls.order").
		Where("cust_id = ?", custId).
		Where("ro_date = ?", date).
		Count(&total).Error

	if err != nil {
		return 0, err
	}

	return int(total), nil
}

func (repository *RepositoryOrderImpl) GetStockByWarehouseProductCustomer(ctx context.Context, warehouseID int64, productID int64, customerID string) (float64, error) {
	var stock float64
	err := repository.model(ctx).Table("inv.stock").
		Select("COALESCE(SUM(qty_in) - SUM(qty_out), 0)").
		Where("wh_id = ? AND pro_id = ? AND cust_id = ?", warehouseID, productID, customerID).
		Row().Scan(&stock)
	if err != nil {
		return 0, err
	}
	return stock, nil
}

func (repository *RepositoryOrderImpl) UpdateWarehouseStock(c context.Context, whId int64, proId int64, custId string, qty float64) error {
	err := repository.model(c).Exec(
		"UPDATE inv.warehouse_stock SET qty = qty - ?, updated_at = ? WHERE cust_id = ? AND wh_id = ? AND pro_id = ?",
		qty, time.Now().UTC().Unix(), custId, whId, proId,
	).Error
	return err
}

func (repository *RepositoryOrderImpl) GetInvoiceByNumbers(custID, parentCustID string, empID int64, isCollection bool) ([]model.OrderListInvoice, error) {
	var invoices []model.OrderListInvoice
	timeNow, err := times.GetCurrentTime()
	if err != nil {
		return nil, err
	}
	dateOnly := timeNow.Format("2006-01-02")

	q := `
	WITH invoice_balances AS (
		SELECT
			o.cust_id,
			o.invoice_no,
			o.invoice_date,
			o.due_date,
			o.outlet_id,
			o.salesman_id,
			o.ro_no,
			o.total AS invoice_amount,
			CASE
				WHEN o.opr_type = 'C' AND o.invoice_date::date = ? THEN COALESCE(pt.payment_amount, 0)
				ELSE COALESCE(paid_invoices.paid_amount, 0)
			END AS paid_amount,
			CASE
				WHEN o.opr_type = 'C' AND o.invoice_date::date = ? THEN COALESCE(pt.remaining_amount, o.total)
				ELSE (o.total - COALESCE(paid_invoices.paid_amount, 0))
			END AS remaining_amount
		FROM sls.order o
		LEFT JOIN acf.payment_trx pt
			ON pt.po_number = o.order_no
			AND pt.trx_source = 'C'
			AND pt.outlet_id = o.outlet_id
			AND pt.cust_id = o.cust_id
		LEFT JOIN (
			SELECT
				dd.invoice_no,
				dd.cust_id,
				SUM(dd.total_payment) AS paid_amount
			FROM acf.deposit_detail dd
			INNER JOIN acf.deposit d
				ON d.deposit_no = dd.deposit_no
				AND d.cust_id = dd.cust_id
				AND d.deposit_status IN (1, 2)
			GROUP BY dd.invoice_no, dd.cust_id
		) paid_invoices
			ON paid_invoices.invoice_no = o.invoice_no
			AND paid_invoices.cust_id = o.cust_id
		WHERE o.invoice_no IS NOT NULL
			AND o.cust_id = ?
	)
	SELECT
		ib.*,
		ot.outlet_code,
		ot.outlet_name,
		ot.ot_grp_id
	FROM invoice_balances ib
	LEFT JOIN mst.m_outlet ot
		ON ot.outlet_id = ib.outlet_id
		AND ot.cust_id = ?
	WHERE ib.remaining_amount > 0
		AND NOT EXISTS (
			SELECT 1 FROM acf.collection_det cd
			WHERE cd.invoice_no = ib.invoice_no
				AND cd.remaining_amount = ib.remaining_amount
				AND cd.cust_id = ib.cust_id
				AND NOT EXISTS (
					SELECT 1 FROM acf.deposit_detail dd
					JOIN acf.deposit d
						ON d.deposit_no = dd.deposit_no
						AND dd.invoice_no = ib.invoice_no
						AND dd.cust_id = ib.cust_id
					WHERE d.collection_no = cd.collection_no
				)
		)
	ORDER BY ib.remaining_amount DESC`

	err = repository.Raw(q, dateOnly, dateOnly, custID, custID).Scan(&invoices).Error
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (repository *RepositoryOrderImpl) GetSummarySales(custID string, salesmanID int, startDate string, endDate string) (float64, error) {
	var summarySales float64

	q := `
	SELECT (COALESCE(SUM(o.total), 0) - COALESCE(SUM(r.total), 0)) AS summary_sales
	FROM sls.order o
	LEFT JOIN sls.return r
		ON r.invoice_no = o.invoice_no
		AND r.salesman_id = o.salesman_id
		AND r.return_date BETWEEN ? AND ?
		AND r.data_status IN (1,2,3,4,5,6,7)
		AND r.is_del = false
	WHERE o.cust_id = ?
		AND o.salesman_id = ?
		AND o.ro_date BETWEEN ? AND ?
		AND o.data_status IN (1,2,3,4,5,6,7)
		AND o.is_del = false`

	err := repository.Raw(q, startDate, endDate, custID, salesmanID, startDate, endDate).Scan(&summarySales).Error
	if err != nil {
		return 0, err
	}

	return summarySales, nil
}
