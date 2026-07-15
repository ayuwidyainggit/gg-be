package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"inventory/pkg/str"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryStockOpnameImpl struct {
		*gorm.DB
	}

	ProductStatusSnapshot struct {
		ProStatus int
		IsActive  bool
	}
)

type StockOpnameRepository interface {
	Store(c context.Context, data *model.StockOpname) error
	StoreDetail(c context.Context, data *model.StockOpnameDetail) error
	StoreReport(c context.Context, data *model.StockOpnameReport) error
	FindAllByCustId(dataFilter entity.StockOpnameQueryFilter) ([]model.StockOpnameList, int64, int, error)
	FindWarehouseByIDAndCustID(whID int64, custID string) (model.StockOpnameWarehouse, error)
	FindEmployeeByIDAndCustID(empID int64, custID string) (model.StockOpnameEmployee, error)
	FindByNo(params entity.ReportStockOpanmeParams) (so model.StockOpnameList, err error)
	FindAllStockOpnameReportByNo(params entity.ReportStockOpanmeParams) (data []model.StockOpnameReport, err error)
	Update(c context.Context, params entity.UpdateStockOpnameParams, data model.StockOpname) error
	ProductList(dataFilter entity.StockOpnameProductListQueryFilter) ([]model.StockOpnameProductList, int64, int, error)
	GenerateDocNoV2(c context.Context, custID string, scheduleDate time.Time) (string, error)
	FindProductByIDs(c context.Context, proIDs []int64, custID string) ([]model.Product, error)
	UpdateProductStatus(c context.Context, proIDs []int64, custID string, status int) error
	UpdateProductStatusForCompletion(c context.Context, proIDs []int64, custID string) error
	GetProductStatusSnapshot(c context.Context, proIDs []int64, custID string) (map[int64]ProductStatusSnapshot, error)
	UpdateDetailsProductSnapshot(c context.Context, docNo, custID string, snapshot map[int64]ProductStatusSnapshot) error
	RestoreProductStatusFromSnapshot(c context.Context, docNo, custID string) error
	// V2 Methods
	StoreV2(c context.Context, data *model.StockOpnameV2) error
	StoreDetailV2(c context.Context, data *model.StockOpnameDetailV2) error
	// V2 List Method
	FindAllByCustIdV2(dataFilter entity.StockOpnameListV2QueryFilter) ([]model.StockOpnameListV2, int64, int, error)
	// V2 Detail Method
	FindDetailHeaderByDocNoV2(params entity.StockOpnameDetailV2Params) (model.StockOpnameDetailV2Header, error)
	FindDetailProductsByDocNoV2(params entity.StockOpnameDetailV2Params) ([]model.StockOpnameDetailV2Product, error)
	// V2 Update Status Method
	FindStockOpnameForUpdate(c context.Context, docNo, custID string) (model.StockOpnameForUpdate, error)
	UpdateStockOpnameStatusV2(c context.Context, docNo, custID string, status int, isProcess bool, updatedBy int64) error
	InsertStockOpnameLog(c context.Context, log *model.StockOpnameLog) error
	GetProductIDsFromDetail(c context.Context, docNo, custID string) ([]int64, error)
	CheckProductsHaveInvoice(c context.Context, proIDs []int64, custID string) (bool, error)
	CheckProductsHaveInvoiceNull(c context.Context, proIDs []int64, custID string) (bool, error)
	CheckProductsOrderStatus(c context.Context, proIDs []int64, custID string) (hasStatusInRange bool, hasStatusOutRange bool, err error)
	GetBlockingOrderRos(c context.Context, proIDs []int64, custID string) ([]string, error)
	GetStockOpnameDetailsForCompleted(c context.Context, docNo, custID, parentCustID string) ([]model.StockOpnameDetailForCompleted, error)
	GetProductPrices(c context.Context, proIDs []int64, custID, parentCustID string) (map[int64]model.ProductPrice, error)
	GetWarehouseStock(c context.Context, whID int64, proID int64, custID string) (model.WarehouseStock, error)
	// V2 Revised Methods
	UpdateStockOpnameIsRevised(c context.Context, docNo, custID string, isRevised bool, updatedBy int64) error
	UpdateStockOpnameDetailRevised(c context.Context, docNo, custID string, proID int64, qtyRevised1, qtyRevised2, qtyRevised3 float64, userRevised int64) error
	FindStockOpnameDetailByProID(c context.Context, docNo, custID string, proID int64) (model.StockOpnameDetailV2, error)
	// V2 Start Method
	FindStockOpnameForStart(c context.Context, docNo, custID string) (model.StockOpnameStartData, error)
	UpdateStockOpnameStart(c context.Context, docNo, custID string, startedBy int64) error
	// V2 Submit Method
	UpdateStockOpnameDetailsQtySO(c context.Context, docNo, custID string, details []model.StockOpnameDetailQtySO) error
	UpdateStockOpnameStatusToSubmit(c context.Context, docNo, custID string, updatedBy int64) error
	ValidateStockOpnameDetailIDs(c context.Context, docNo, custID string, detailIDs []int64) ([]int64, error)
	// Bulk upload
	InsertBulkUpload(c context.Context, data *model.StockOpnameBulkUpload) error
	InsertBulkUploadItems(c context.Context, items []model.StockOpnameBulkUploadItem) error
}

func NewStockOpnameRepo(db *gorm.DB) *RepositoryStockOpnameImpl {
	return &RepositoryStockOpnameImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryStockOpnameImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryStockOpnameImpl) Store(c context.Context, data *model.StockOpname) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) StoreDetail(c context.Context, data *model.StockOpnameDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) StoreReport(c context.Context, data *model.StockOpnameReport) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) FindAllByCustId(dataFilter entity.StockOpnameQueryFilter) ([]model.StockOpnameList, int64, int, error) {
	var (
		data  []model.StockOpnameList
		total int64
	)
	limit := 10
	if dataFilter.Limit != 0 {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("doc_no")
	query := repository.Select(`inv.stock_opname.*, wh.wh_code, wh.wh_name, wh.stock_type, emp.emp_code, emp.emp_name`).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = inv.stock_opname.wh_id AND wh.cust_id = ?", dataFilter.CustID).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = inv.stock_opname.assign_to_emp_id AND emp.cust_id = ?", dataFilter.CustID)

	queryCount.Where("inv.stock_opname.cust_id=?", dataFilter.CustID)
	query.Where("inv.stock_opname.cust_id=?", dataFilter.CustID)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where(`inv.stock_opname.scheduled_at BETWEEN ? AND ? 
					OR inv.stock_opname.scheduled_at BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To),
			str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To),
		)
		queryCount.Where(`inv.stock_opname.scheduled_at BETWEEN ? AND ? 
						OR inv.stock_opname.scheduled_at BETWEEN ? AND ?`,
			str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To),
			str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To),
		)
	}

	if dataFilter.DocNo != "" {
		queryCount.Where("inv.stock_opname.doc_no ILIKE ?", "%"+dataFilter.DocNo+"%")
		query.Where("inv.stock_opname.doc_no ILIKE ?", "%"+dataFilter.DocNo+"%")
	}

	if len(dataFilter.DataStatus) > 0 {
		queryCount.Where("inv.stock_opname.data_status IN ?", dataFilter.DataStatus)
		query.Where("inv.stock_opname.data_status IN ?", dataFilter.DataStatus)
	}

	query.Order("created_at DESC")
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		return data, total, 0, err
	}
	err = queryCount.Model(&data).Count(&total).Error
	if err != nil {
		return data, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return data, total, lastPage, nil
}

func (repository *RepositoryStockOpnameImpl) FindWarehouseByIDAndCustID(whID int64, custId string) (wh model.StockOpnameWarehouse, err error) {
	err = repository.
		Select("wh_id, wh_code, wh_name").
		Where("wh_id = ? AND cust_id = ?", whID, custId).
		Take(&wh).Error
	return wh, err
}

func (repository *RepositoryStockOpnameImpl) FindEmployeeByIDAndCustID(empID int64, custId string) (emp model.StockOpnameEmployee, err error) {
	err = repository.
		Select("emp_id, emp_code, emp_name").
		Where("emp_id = ? AND cust_id = ?", empID, custId).
		Take(&emp).Error
	return emp, err
}

func (repository *RepositoryStockOpnameImpl) FindByNo(params entity.ReportStockOpanmeParams) (data model.StockOpnameList, err error) {
	err = repository.
		Select("stock_opname.*, wh.wh_code, wh.wh_name, wh.stock_type, emp.emp_code, emp.emp_name").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = stock_opname.wh_id AND wh.cust_id = ?", params.CustID).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = stock_opname.assign_to_emp_id AND emp.cust_id = ?", params.CustID).
		Where("stock_opname.doc_no = ? AND stock_opname.cust_id = ?", params.DocNo, params.CustID).
		Take(&data).Error
	return data, err
}

func (repository *RepositoryStockOpnameImpl) FindAllStockOpnameReportByNo(params entity.ReportStockOpanmeParams) (data []model.StockOpnameReport, err error) {
	err = repository.
		Select("stock_opname_reports.*").
		Where("stock_opname_reports.doc_no = ? AND stock_opname_reports.cust_id = ?", params.DocNo, params.CustID).
		Take(&data).Error
	return data, err
}

func (repository *RepositoryStockOpnameImpl) Update(c context.Context, params entity.UpdateStockOpnameParams, data model.StockOpname) error {
	result := repository.model(c).Model(&data).Where("cust_id = ? AND doc_no=?", params.CustID, params.DocNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) ProductList(dataFilter entity.StockOpnameProductListQueryFilter) ([]model.StockOpnameProductList, int64, int, error) {
	var products []model.StockOpnameProductList
	var total int64

	limit := 5
	if dataFilter.Limit > 0 {
		limit = dataFilter.Limit
	}

	page := 1
	if dataFilter.Page > 0 {
		page = dataFilter.Page
	}

	// Base query for count - use subquery to count distinct products
	// This ensures count matches the GROUP BY in SELECT query
	queryCount := repository.Select("COUNT(DISTINCT p.pro_id)").
		Table("mst.m_product AS p").
		Where("p.is_del = false")

	// Add warehouse_stock JOIN with wh_id filter if provided
	if len(dataFilter.WhID) > 0 {
		queryCount = queryCount.Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = p.pro_id AND whs.wh_id IN ?", dataFilter.WhID)
	} else {
		queryCount = queryCount.Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = p.pro_id")
	}

	// Build warehouse JOIN with optional stock_type filter in JOIN condition
	if dataFilter.StockType != "" {
		queryCount = queryCount.Joins("LEFT JOIN mst.m_warehouse w ON w.wh_id = whs.wh_id AND w.stock_type = ?", dataFilter.StockType)
	} else {
		queryCount = queryCount.Joins("LEFT JOIN mst.m_warehouse w ON w.wh_id = whs.wh_id")
	}

	queryCount = queryCount.
		Joins("LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id").
		Joins("LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id").
		Joins("LEFT JOIN mst.m_brand b ON b.brand_id = sb1.brand_id").
		Joins("LEFT JOIN mst.m_product_line pl ON pl.pl_id = b.pl_id")

	// Base query for data - start from m_product as per requirement
	query := repository.Select(`p.pro_id,
		p.pro_code,
		p.pro_name,
		p.unit_id1 AS uom1,
		p.unit_id2 AS uom2,
		p.unit_id3 AS uom3,
		(SELECT MIN(u.unit_name) FROM mst.m_unit u WHERE u.unit_id = p.unit_id1) AS unit_name1,
		(SELECT MIN(u.unit_name) FROM mst.m_unit u WHERE u.unit_id = p.unit_id2) AS unit_name2,
		(SELECT MIN(u.unit_name) FROM mst.m_unit u WHERE u.unit_id = p.unit_id3) AS unit_name3,
		p.conv_unit2,
		p.conv_unit3,
		COALESCE(whs.qty, 0) AS qty,
		COALESCE(whs.wh_id, 0) AS wh_id,
		COALESCE(w.stock_type, '') AS stock_type,
		COALESCE(p.principal_id, 0) AS principal_id,
		COALESCE(pr.principal_name, '') AS principal_name,
		COALESCE(pl.pl_id, 0) AS pl_id,
		COALESCE(pl.pl_name, '') AS pl_name,
		COALESCE(b.brand_id, 0) AS brand_id,
		COALESCE(b.brand_name, '') AS brand_name,
		COALESCE(sb1.sbrand1_id, 0) AS sbrand1_id,
		COALESCE(sb1.sbrand1_name, '') AS sbrand1_name,
		p.is_active`).
		Table("mst.m_product AS p").
		Where("p.is_del = false")

	// Add warehouse_stock JOIN with wh_id filter if provided
	if len(dataFilter.WhID) > 0 {
		query = query.Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = p.pro_id AND whs.wh_id IN ?", dataFilter.WhID)
	} else {
		query = query.Joins("LEFT JOIN inv.warehouse_stock whs ON whs.pro_id = p.pro_id")
	}

	// Build warehouse JOIN with optional stock_type filter in JOIN condition
	if dataFilter.StockType != "" {
		query = query.Joins("LEFT JOIN mst.m_warehouse w ON w.wh_id = whs.wh_id AND w.stock_type = ?", dataFilter.StockType)
	} else {
		query = query.Joins("LEFT JOIN mst.m_warehouse w ON w.wh_id = whs.wh_id")
	}

	query = query.
		Joins("LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id").
		Joins("LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id").
		Joins("LEFT JOIN mst.m_brand b ON b.brand_id = sb1.brand_id").
		Joins("LEFT JOIN mst.m_product_line pl ON pl.pl_id = b.pl_id")

	// wh_id filter is already applied in JOIN condition above
	// stock_type filter is now in JOIN condition, not WHERE clause

	if len(dataFilter.PrincipalID) > 0 {
		queryCount.Where("p.principal_id IN ?", dataFilter.PrincipalID)
		query.Where("p.principal_id IN ?", dataFilter.PrincipalID)
	}

	if len(dataFilter.PLID) > 0 {
		queryCount.Where("pl.pl_id IN ?", dataFilter.PLID)
		query.Where("pl.pl_id IN ?", dataFilter.PLID)
	}

	if len(dataFilter.BrandID) > 0 {
		queryCount.Where("b.brand_id IN ?", dataFilter.BrandID)
		query.Where("b.brand_id IN ?", dataFilter.BrandID)
	}

	if len(dataFilter.SBrand1ID) > 0 {
		queryCount.Where("sb1.sbrand1_id IN ?", dataFilter.SBrand1ID)
		query.Where("sb1.sbrand1_id IN ?", dataFilter.SBrand1ID)
	}

	if !dataFilter.ZeroStock {
		queryCount.Where("COALESCE(whs.qty, 0) > 0")
		query.Where("COALESCE(whs.qty, 0) > 0")
	}

	if dataFilter.IsActive {
		queryCount.Where("p.is_active = true")
		query.Where("p.is_active = true")
	}

	if dataFilter.Query != "" {
		queryCount.Where("(p.pro_code ILIKE ? OR p.pro_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query.Where("(p.pro_code ILIKE ? OR p.pro_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	// Default sort by created_at DESC
	sortBy := "p.created_at DESC"
	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			col := sortParts[0]
			dir := strings.ToUpper(sortParts[1])
			if dir == "ASC" || dir == "DESC" {
				switch col {
				case "created_date":
					sortBy = "p.created_at " + dir
				case "created_at":
					sortBy = "p.created_at " + dir
				case "product_name":
					sortBy = "p.pro_name " + dir
				case "pro_name":
					sortBy = "p.pro_name " + dir
				case "pro_code":
					sortBy = "p.pro_code " + dir
				case "pro_id":
					sortBy = "p.pro_id " + dir
				default:
					sortBy = "p.created_at DESC"
				}
			}
		}
	}
	query.Order(sortBy)

	query = query.Group("p.pro_id, p.pro_code, p.pro_name, p.unit_id1, p.unit_id2, p.unit_id3, " +
		"p.conv_unit2, p.conv_unit3, " +
		"whs.qty, whs.wh_id, w.stock_type, p.principal_id, pr.principal_name, " +
		"pl.pl_id, pl.pl_name, b.brand_id, b.brand_name, sb1.sbrand1_id, sb1.sbrand1_name, p.is_active")

	// Count distinct products
	err := queryCount.Scan(&total).Error
	if err != nil {
		return products, total, 0, err
	}

	offset := (page - 1) * limit
	// Explicitly set table to override model's TableName()
	err = query.Table("mst.m_product AS p").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return products, total, lastPage, nil
}

func (repository *RepositoryStockOpnameImpl) GenerateDocNoV2(c context.Context, custID string, scheduleDate time.Time) (string, error) {
	var docNo string

	queryStr := `SELECT inv.generate_stock_opname_doc_no($1, $2)`

	err := repository.model(c).Raw(queryStr, strings.ToUpper(custID), scheduleDate).Scan(&docNo).Error
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "function") {
			return "", fmt.Errorf("PostgreSQL function inv.generate_stock_opname_doc_no not found. Please run migration: migration/inv.stock_opname/create_stock_opname_doc_no_sequence_function.sql. Error: %w", err)
		}
		return "", fmt.Errorf("failed to generate doc_no using sequence: %w", err)
	}

	return docNo, nil
}

func (repository *RepositoryStockOpnameImpl) FindProductByIDs(c context.Context, proIDs []int64, custID string) ([]model.Product, error) {
	var products []model.Product
	err := repository.model(c).
		Table("mst.m_product").
		Select("pro_id, purch_price1, purch_price2, purch_price3").
		Where("pro_id IN ?", proIDs).
		Find(&products).Error
	return products, err
}

func (repository *RepositoryStockOpnameImpl) UpdateProductStatus(c context.Context, proIDs []int64, custID string, status int) error {
	result := repository.model(c).
		Table("mst.m_product").
		Where("pro_id IN ? ", proIDs).
		Updates(map[string]interface{}{
			"pro_status": status,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) UpdateProductStatusForCompletion(c context.Context, proIDs []int64, custID string) error {
	result := repository.model(c).
		Table("mst.m_product").
		Where("pro_id IN ?", proIDs).
		Updates(map[string]interface{}{
			"pro_status": 1,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) GetProductStatusSnapshot(c context.Context, proIDs []int64, custID string) (map[int64]ProductStatusSnapshot, error) {
	if len(proIDs) == 0 {
		return map[int64]ProductStatusSnapshot{}, nil
	}

	type productStatusRow struct {
		ProID     int64 `gorm:"column:pro_id"`
		ProStatus int   `gorm:"column:pro_status"`
		IsActive  bool  `gorm:"column:is_active"`
	}

	var rows []productStatusRow
	err := repository.model(c).
		Table("mst.m_product").
		Select("pro_id, COALESCE(pro_status, 1) AS pro_status, COALESCE(is_active, true) AS is_active").
		Where("pro_id IN ? AND cust_id = ?", proIDs, custID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	snapshot := make(map[int64]ProductStatusSnapshot)
	for _, row := range rows {
		snapshot[row.ProID] = ProductStatusSnapshot{
			ProStatus: row.ProStatus,
			IsActive:  row.IsActive,
		}
	}
	return snapshot, nil
}

func (repository *RepositoryStockOpnameImpl) UpdateDetailsProductSnapshot(c context.Context, docNo, custID string, snapshot map[int64]ProductStatusSnapshot) error {
	for proID, s := range snapshot {
		result := repository.model(c).
			Table("inv.stock_opname_details").
			Where("doc_no = ? AND cust_id = ? AND pro_id = ?", docNo, custID, proID).
			Updates(map[string]interface{}{
				"pro_status_before": s.ProStatus,
				"is_active_before":  s.IsActive,
			})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) RestoreProductStatusFromSnapshot(c context.Context, docNo, custID string) error {
	type snapshotRow struct {
		ProID           int64 `gorm:"column:pro_id"`
		ProStatusBefore *int  `gorm:"column:pro_status_before"`
		IsActiveBefore  *bool `gorm:"column:is_active_before"`
	}

	var rows []snapshotRow
	err := repository.model(c).
		Table("inv.stock_opname_details").
		Select("pro_id, pro_status_before, is_active_before").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Find(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		proStatus := 1
		if row.ProStatusBefore != nil {
			proStatus = *row.ProStatusBefore
		}

		result := repository.model(c).
			Table("mst.m_product").
			Where("pro_id = ? AND cust_id = ?", row.ProID, custID).
			Updates(map[string]interface{}{
				"pro_status": proStatus,
			})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) StoreV2(c context.Context, data *model.StockOpnameV2) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) StoreDetailV2(c context.Context, data *model.StockOpnameDetailV2) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) FindAllByCustIdV2(dataFilter entity.StockOpnameListV2QueryFilter) ([]model.StockOpnameListV2, int64, int, error) {
	var (
		data  []model.StockOpnameListV2
		total int64
	)

	// Default values
	limit := 5
	if dataFilter.Limit > 0 {
		limit = dataFilter.Limit
	}

	page := 1
	if dataFilter.Page > 0 {
		page = dataFilter.Page
	}

	// Build base query
	queryCount := repository.Select("inv.stock_opname.doc_no").
		Table("inv.stock_opname").
		Where("inv.stock_opname.cust_id = ?", dataFilter.CustID)

	query := repository.Select(`
		inv.stock_opname.doc_no,
		TO_CHAR(inv.stock_opname.created_at, 'DD/MM/YYYY') AS created_date,
		inv.stock_opname.wh_id,
		wh.wh_code,
		wh.wh_name,
		COALESCE(u.user_name, '') AS created_by,
		COALESCE(u.user_fullname, u.user_name, '') AS user_name,
		COALESCE(TO_CHAR(inv.stock_opname.scheduled_at, 'DD/MM/YYYY'), '') AS scheduled_date,
		COALESCE(inv.stock_opname.emp_id, 0) AS emp_id,
		COALESCE(emp.emp_name, '') AS emp_name,
		inv.stock_opname.data_status AS status
	`).
		Table("inv.stock_opname").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = inv.stock_opname.wh_id AND wh.cust_id = ?", dataFilter.CustID).
		Joins("LEFT JOIN sys.m_user u ON u.user_id = inv.stock_opname.created_by").
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = inv.stock_opname.emp_id AND emp.cust_id = ?", dataFilter.CustID).
		Where("inv.stock_opname.cust_id = ?", dataFilter.CustID)

	// Role-based filtering: if not admin, filter by user (emp_id related to user)
	if !dataFilter.IsAdmin && dataFilter.UserID > 0 && len(dataFilter.EmpID) == 0 {
		queryCount.Where("inv.stock_opname.created_by = ?", dataFilter.UserID)
		query.Where("inv.stock_opname.created_by = ?", dataFilter.UserID)
	}

	if dataFilter.Query != "" {
		queryCount.Where("inv.stock_opname.doc_no ILIKE ?", "%"+dataFilter.Query+"%")
		query.Where("inv.stock_opname.doc_no ILIKE ?", "%"+dataFilter.Query+"%")
	}

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		startTime := str.UnixTimestampToUtcTime(*dataFilter.StartDate)
		endTime := str.UnixTimestampToUtcTime(*dataFilter.EndDate)
		queryCount.Where("inv.stock_opname.created_at BETWEEN ? AND ?", startTime, endTime)
		query.Where("inv.stock_opname.created_at BETWEEN ? AND ?", startTime, endTime)
	}

	if len(dataFilter.WhID) > 0 {
		queryCount.Where("inv.stock_opname.wh_id IN ?", dataFilter.WhID)
		query.Where("inv.stock_opname.wh_id IN ?", dataFilter.WhID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("inv.stock_opname.data_status IN ?", dataFilter.Status)
		query.Where("inv.stock_opname.data_status IN ?", dataFilter.Status)
	}

	if len(dataFilter.EmpID) > 0 {
		queryCount.Where("inv.stock_opname.emp_id IN ?", dataFilter.EmpID)
		query.Where("inv.stock_opname.emp_id IN ?", dataFilter.EmpID)
	}

	sortBy := "inv.stock_opname.created_at DESC"
	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			col := sortParts[0]
			dir := strings.ToUpper(sortParts[1])
			if dir == "ASC" || dir == "DESC" {
				switch col {
				case "created_date", "created_at":
					sortBy = "inv.stock_opname.created_at " + dir
				case "scheduled_date", "scheduled_at":
					sortBy = "inv.stock_opname.scheduled_at " + dir
				case "doc_no":
					sortBy = "inv.stock_opname.doc_no " + dir
				case "status":
					sortBy = "inv.stock_opname.data_status " + dir
				default:
					sortBy = "inv.stock_opname.created_at DESC"
				}
			}
		}
	}
	query.Order(sortBy)

	err := queryCount.Count(&total).Error
	if err != nil {
		return data, total, 0, err
	}

	offset := (page - 1) * limit
	err = query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		return data, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return data, total, lastPage, nil
}

func (repository *RepositoryStockOpnameImpl) FindDetailHeaderByDocNoV2(params entity.StockOpnameDetailV2Params) (model.StockOpnameDetailV2Header, error) {
	var data model.StockOpnameDetailV2Header

	err := repository.Select(`
		inv.stock_opname.doc_no,
		TO_CHAR(inv.stock_opname.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_date,
		inv.stock_opname.wh_id,
		COALESCE(wh.wh_code, '') AS wh_code,
		COALESCE(wh.wh_name, '') AS wh_name,
		COALESCE(wh.stock_type, 'G') AS stock_type,
		COALESCE(inv.stock_opname.created_by, 0) AS created_by,
		COALESCE(u.user_fullname, u.user_name, '') AS user_name,
		CASE 
			WHEN inv.stock_opname.scheduled_at IS NULL THEN NULL
			ELSE TO_CHAR(inv.stock_opname.scheduled_at, 'YYYY-MM-DD HH24:MI:SS+00')
		END AS scheduled_date,
		inv.stock_opname.data_status AS status,
		COALESCE(inv.stock_opname.is_revised, false) AS is_revised,
		inv.stock_opname.is_process,
		COALESCE(inv.stock_opname.input_by, 'Web') AS input_by,
		inv.stock_opname.division_id,
		COALESCE(div.division_name, '') AS division_name,
		COALESCE(inv.stock_opname.emp_id, 0) AS emp_id,
		COALESCE(emp.emp_name, '') AS emp_name,
		inv.stock_opname.notes
	`).
		Table("inv.stock_opname").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = inv.stock_opname.wh_id AND wh.cust_id = ?", params.CustID).
		Joins("LEFT JOIN sys.m_user u ON u.user_id = inv.stock_opname.created_by").
		Joins("LEFT JOIN mst.m_division div ON div.division_id = inv.stock_opname.division_id").
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = inv.stock_opname.emp_id").
		Where("inv.stock_opname.doc_no = ? AND inv.stock_opname.cust_id = ?", params.DocNo, params.CustID).
		Take(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, constant.ErrRecordNotFound
		}
		return data, err
	}

	return data, nil
}

func (repository *RepositoryStockOpnameImpl) FindDetailProductsByDocNoV2(params entity.StockOpnameDetailV2Params) ([]model.StockOpnameDetailV2Product, error) {
	var data []model.StockOpnameDetailV2Product

	err := repository.Select(`
		inv.stock_opname_details.stock_opname_detail_id,
		inv.stock_opname_details.pro_id,
		COALESCE(p.pro_code, '') AS pro_code,
		COALESCE(p.pro_name, '') AS pro_name,
		COALESCE(p.unit_id1, '') AS unit_id1,
		COALESCE(p.unit_id2, '') AS unit_id2,
		COALESCE(p.unit_id3, '') AS unit_id3,
		COALESCE(p.sell_price1, 0) AS sell_price1,
		COALESCE(p.sell_price2, 0) AS sell_price2,
		COALESCE(p.sell_price3, 0) AS sell_price3,
		COALESCE(u1.unit_name, '') AS unit_name1,
		COALESCE(u2.unit_name, '') AS unit_name2,
		COALESCE(u3.unit_name, '') AS unit_name3,
		COALESCE(p.conv_unit2, 0) AS conv_unit2,
		COALESCE(p.conv_unit3, 0) AS conv_unit3,
		COALESCE(inv.stock_opname_details.qty_stock1, 0) AS qty_stock1,
		COALESCE(inv.stock_opname_details.qty_stock2, 0) AS qty_stock2,
		COALESCE(inv.stock_opname_details.qty_stock3, 0) AS qty_stock3,
		CASE 
			WHEN inv.stock_opname_details.revised_date IS NOT NULL 
			THEN COALESCE(inv.stock_opname_details.qty_revised1, 0)
			ELSE COALESCE(inv.stock_opname_details.qty_so1, 0)
		END AS qty_opname1,
		CASE 
			WHEN inv.stock_opname_details.revised_date IS NOT NULL 
			THEN COALESCE(inv.stock_opname_details.qty_revised2, 0)
			ELSE COALESCE(inv.stock_opname_details.qty_so2, 0)
		END AS qty_opname2,
		CASE 
			WHEN inv.stock_opname_details.revised_date IS NOT NULL 
			THEN COALESCE(inv.stock_opname_details.qty_revised3, 0)
			ELSE COALESCE(inv.stock_opname_details.qty_so3, 0)
		END AS qty_opname3,
		inv.stock_opname_details.qty_revised1,
		inv.stock_opname_details.qty_revised2,
		inv.stock_opname_details.qty_revised3,
		COALESCE(inv.stock_opname_details.purch_price1, 0) AS purch_price1,
		COALESCE(inv.stock_opname_details.purch_price2, 0) AS purch_price2,
		COALESCE(inv.stock_opname_details.purch_price3, 0) AS purch_price3
	`).
		Table("inv.stock_opname_details").
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = inv.stock_opname_details.pro_id").
		Joins("LEFT JOIN mst.m_unit u1 ON u1.unit_id = p.unit_id1 AND u1.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_unit u2 ON u2.unit_id = p.unit_id2 AND u2.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_unit u3 ON u3.unit_id = p.unit_id3 AND u3.cust_id = ?", params.ParentCustID).
		Where("inv.stock_opname_details.doc_no = ? AND inv.stock_opname_details.cust_id = ?", params.DocNo, params.CustID).
		Order("p.pro_code ASC").
		Find(&data).Error

	return data, err
}

func (repository *RepositoryStockOpnameImpl) FindStockOpnameForUpdate(c context.Context, docNo, custID string) (model.StockOpnameForUpdate, error) {
	var data model.StockOpnameForUpdate

	err := repository.model(c).
		Select("cust_id, doc_no, data_status, created_by, is_process").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Take(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("record not found")
		}
		return data, err
	}

	return data, nil
}

func (repository *RepositoryStockOpnameImpl) UpdateStockOpnameStatusV2(c context.Context, docNo, custID string, status int, isProcess bool, updatedBy int64) error {
	now := time.Now()
	result := repository.model(c).
		Table("inv.stock_opname").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Updates(map[string]interface{}{
			"data_status": status,
			"is_process":  isProcess,
			"updated_by":  updatedBy,
			"updated_at":  now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) InsertStockOpnameLog(c context.Context, log *model.StockOpnameLog) error {
	return repository.model(c).Create(log).Error
}

func (repository *RepositoryStockOpnameImpl) GetProductIDsFromDetail(c context.Context, docNo, custID string) ([]int64, error) {
	var proIDs []int64

	err := repository.model(c).
		Table("inv.stock_opname_details").
		Select("pro_id").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Pluck("pro_id", &proIDs).Error

	return proIDs, err
}

func (repository *RepositoryStockOpnameImpl) CheckProductsHaveInvoice(c context.Context, proIDs []int64, custID string) (bool, error) {
	if len(proIDs) == 0 {
		return false, nil
	}

	var count int64

	err := repository.model(c).
		Table("sls.order_detail od").
		Joins("JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id").
		Where("od.pro_id IN ? AND od.cust_id = ? AND o.invoice_no IS NOT NULL AND o.invoice_no != ''", proIDs, custID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repository *RepositoryStockOpnameImpl) CheckProductsHaveInvoiceNull(c context.Context, proIDs []int64, custID string) (bool, error) {
	if len(proIDs) == 0 {
		return false, nil
	}

	var count int64

	err := repository.model(c).
		Table("sls.order_detail od").
		Joins("JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id").
		Where("od.pro_id IN ? AND od.cust_id = ? AND (o.invoice_no IS NULL OR o.invoice_no = '')", proIDs, custID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckProductsOrderStatus checks if products exist in sls.order_det with specific data_status ranges
// Returns:
func (repository *RepositoryStockOpnameImpl) CheckProductsOrderStatus(c context.Context, proIDs []int64, custID string) (hasStatusInRange bool, hasStatusOutRange bool, err error) {
	if len(proIDs) == 0 {
		return false, false, nil
	}

	var countInRange int64
	var countOutRange int64

	err = repository.model(c).
		Table("sls.order_detail od").
		Joins("JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id").
		Where("od.pro_id IN ? AND od.cust_id = ? AND o.data_status IN ?", proIDs, custID, []int{1, 2, 3, 4, 5}).
		Count(&countInRange).Error

	if err != nil {
		return false, false, err
	}

	err = repository.model(c).
		Table("sls.order_detail od").
		Joins("JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id").
		Where("od.pro_id IN ? AND od.cust_id = ? AND o.data_status NOT IN ?", proIDs, custID, []int{1, 2, 3, 4, 5}).
		Count(&countOutRange).Error

	if err != nil {
		return false, false, err
	}

	return countInRange > 0, countOutRange > 0, nil
}

func (repository *RepositoryStockOpnameImpl) GetBlockingOrderRos(c context.Context, proIDs []int64, custID string) ([]string, error) {
	if len(proIDs) == 0 {
		return []string{}, nil
	}

	var roNos []string

	err := repository.model(c).
		Table("sls.order_detail od").
		Select("DISTINCT o.ro_no").
		Joins("JOIN sls.order o ON o.ro_no = od.ro_no AND o.cust_id = od.cust_id").
		Where("od.pro_id IN ? AND od.cust_id = ? AND o.data_status IN ?", proIDs, custID, []int{1, 2, 3, 4, 5}).
		Pluck("o.ro_no", &roNos).Error

	if err != nil {
		return nil, err
	}

	return roNos, nil
}

func (repository *RepositoryStockOpnameImpl) GetStockOpnameDetailsForCompleted(c context.Context, docNo, custID, parentCustID string) ([]model.StockOpnameDetailForCompleted, error) {
	var data []model.StockOpnameDetailForCompleted

	q := repository.model(c).
		Select(`
			inv.stock_opname_details.stock_opname_detail_id,
			inv.stock_opname_details.pro_id,
			inv.stock_opname_details.qty_stock1,
			inv.stock_opname_details.qty_stock2,
			inv.stock_opname_details.qty_stock3,
			inv.stock_opname_details.qty_so1,
			inv.stock_opname_details.qty_so2,
			inv.stock_opname_details.qty_so3,
			inv.stock_opname_details.qty_revised1,
			inv.stock_opname_details.qty_revised2,
			inv.stock_opname_details.qty_revised3,
			inv.stock_opname_details.revised_date,
			COALESCE(p.conv_unit2, 0) AS conv_unit2,
			COALESCE(p.conv_unit3, 0) AS conv_unit3
		`).
		Table("inv.stock_opname_details")

	if parentCustID != "" && parentCustID != custID {
		q = q.Joins(`
			LEFT JOIN LATERAL (
				SELECT mp.conv_unit2, mp.conv_unit3
				FROM mst.m_product mp
				WHERE mp.pro_id = inv.stock_opname_details.pro_id
				  AND mp.cust_id IN (?, ?)
				ORDER BY CASE WHEN mp.cust_id = ? THEN 0 ELSE 1 END
				LIMIT 1
			) p ON true`, parentCustID, custID, custID)
	} else {
		q = q.Joins(`
			LEFT JOIN LATERAL (
				SELECT mp.conv_unit2, mp.conv_unit3
				FROM mst.m_product mp
				WHERE mp.pro_id = inv.stock_opname_details.pro_id
				  AND mp.cust_id = ?
				LIMIT 1
			) p ON true`, custID)
	}

	err := q.Where("inv.stock_opname_details.doc_no = ? AND inv.stock_opname_details.cust_id = ?", docNo, custID).
		Find(&data).Error

	return data, err
}

func (repository *RepositoryStockOpnameImpl) GetProductPrices(c context.Context, proIDs []int64, custID, parentCustID string) (map[int64]model.ProductPrice, error) {
	if len(proIDs) == 0 {
		return map[int64]model.ProductPrice{}, nil
	}

	var products []model.ProductPrice

	if parentCustID != "" && parentCustID != custID {
		err := repository.model(c).Raw(`
			SELECT DISTINCT ON (pro_id) pro_id, sell_price1, cogs
			FROM mst.m_product
			WHERE pro_id IN ? AND cust_id IN (?, ?)
			ORDER BY pro_id, CASE WHEN cust_id = ? THEN 0 ELSE 1 END
		`, proIDs, custID, parentCustID, custID).Scan(&products).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := repository.model(c).
			Table("mst.m_product").
			Select("pro_id, sell_price1, cogs").
			Where("pro_id IN ? AND cust_id = ?", proIDs, custID).
			Find(&products).Error
		if err != nil {
			return nil, err
		}
	}

	priceMap := make(map[int64]model.ProductPrice)
	for _, p := range products {
		priceMap[p.ProID] = p
	}

	return priceMap, nil
}

func (repository *RepositoryStockOpnameImpl) GetWarehouseStock(c context.Context, whID int64, proID int64, custID string) (model.WarehouseStock, error) {
	var data model.WarehouseStock

	err := repository.model(c).
		Table("inv.warehouse_stock").
		Where("wh_id = ? AND pro_id = ? AND cust_id = ?", whID, proID, custID).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return zero values if not found
			return model.WarehouseStock{
				CustID:        custID,
				WhID:          whID,
				ProID:         proID,
				Qty:           0,
				QtyOnOrder:    0,
				QtyOnShipping: 0,
				QtyBs:         0,
				QtyExp:        0,
			}, nil
		}
		return data, err
	}

	return data, nil
}

func (repository *RepositoryStockOpnameImpl) ValidateStockOpnameDetailIDs(c context.Context, docNo, custID string, detailIDs []int64) ([]int64, error) {
	var validIDs []int64

	err := repository.model(c).
		Table("inv.stock_opname_details").
		Select("stock_opname_detail_id").
		Where("doc_no = ? AND cust_id = ? AND stock_opname_detail_id IN ?", docNo, custID, detailIDs).
		Pluck("stock_opname_detail_id", &validIDs).Error

	if err != nil {
		return nil, err
	}

	return validIDs, nil
}

func (repository *RepositoryStockOpnameImpl) UpdateStockOpnameDetailsQtySO(c context.Context, docNo, custID string, details []model.StockOpnameDetailQtySO) error {
	for _, detail := range details {
		result := repository.model(c).
			Table("inv.stock_opname_details").
			Where("stock_opname_detail_id = ? AND doc_no = ? AND cust_id = ?", detail.StockOpnameDetID, docNo, custID).
			Updates(map[string]interface{}{
				"qty_so1": detail.QtySO1,
				"qty_so2": detail.QtySO2,
				"qty_so3": detail.QtySO3,
			})

		if result.Error != nil {
			return fmt.Errorf("failed to update stock_opname_detail_id %d: %w", detail.StockOpnameDetID, result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("stock_opname_detail_id %d not found in document %s", detail.StockOpnameDetID, docNo)
		}
	}

	return nil
}

func (repository *RepositoryStockOpnameImpl) UpdateStockOpnameStatusToSubmit(c context.Context, docNo, custID string, updatedBy int64) error {
	now := time.Now()
	result := repository.model(c).
		Table("inv.stock_opname").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Updates(map[string]interface{}{
			"data_status": entity.StockOpnameStatusSubmit, // 5
			"updated_by":  updatedBy,
			"updated_at":  now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return constant.ErrNoRowsAffected
	}

	return nil
}

func (repository *RepositoryStockOpnameImpl) UpdateStockOpnameIsRevised(c context.Context, docNo, custID string, isRevised bool, updatedBy int64) error {
	now := time.Now()
	result := repository.model(c).
		Table("inv.stock_opname").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Updates(map[string]interface{}{
			"is_revised": isRevised,
			"updated_by": updatedBy,
			"updated_at": now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) UpdateStockOpnameDetailRevised(c context.Context, docNo, custID string, proID int64, qtyRevised1, qtyRevised2, qtyRevised3 float64, userRevised int64) error {
	now := time.Now()
	result := repository.model(c).
		Table("inv.stock_opname_details").
		Where("doc_no = ? AND cust_id = ? AND pro_id = ?", docNo, custID, proID).
		Select("qty_revised1", "qty_revised2", "qty_revised3", "user_revised", "revised_date").
		Updates(map[string]interface{}{
			"qty_revised1": qtyRevised1,
			"qty_revised2": qtyRevised2,
			"qty_revised3": qtyRevised3,
			"user_revised": userRevised,
			"revised_date": now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) FindStockOpnameDetailByProID(c context.Context, docNo, custID string, proID int64) (model.StockOpnameDetailV2, error) {
	var data model.StockOpnameDetailV2

	err := repository.model(c).
		Where("doc_no = ? AND cust_id = ? AND pro_id = ?", docNo, custID, proID).
		Take(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, errors.New("stock opname detail not found")
		}
		return data, err
	}

	return data, nil
}

func (repository *RepositoryStockOpnameImpl) FindStockOpnameForStart(c context.Context, docNo, custID string) (model.StockOpnameStartData, error) {
	var data model.StockOpnameStartData

	err := repository.model(c).
		Select(`
			COALESCE(inv.stock_opname.cust_id, '') AS cust_id,
			COALESCE(inv.stock_opname.doc_no, '') AS doc_no,
			COALESCE(inv.stock_opname.wh_id, 0) AS wh_id,
			COALESCE(wh.wh_code, '') AS wh_code,
			COALESCE(wh.wh_name, '') AS wh_name,
			COALESCE(inv.stock_opname.data_status, 0) AS data_status
		`).
		Table("inv.stock_opname").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = inv.stock_opname.wh_id AND wh.cust_id = ?", custID).
		Where("inv.stock_opname.doc_no = ? AND inv.stock_opname.cust_id = ?", docNo, custID).
		Scan(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, constant.ErrRecordNotFound
		}
		return data, err
	}

	if data.DocNo == "" {
		return data, constant.ErrRecordNotFound
	}

	return data, nil
}

func (repository *RepositoryStockOpnameImpl) UpdateStockOpnameStart(c context.Context, docNo, custID string, startedBy int64) error {
	now := time.Now()
	result := repository.model(c).
		Table("inv.stock_opname").
		Where("doc_no = ? AND cust_id = ?", docNo, custID).
		Updates(map[string]interface{}{
			"data_status": entity.StockOpnameStatusOnGoing,
			"updated_by":  startedBy,
			"updated_at":  now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return constant.ErrNoRowsAffected
	}
	return nil
}

func (repository *RepositoryStockOpnameImpl) InsertBulkUpload(c context.Context, data *model.StockOpnameBulkUpload) error {
	return repository.model(c).Create(data).Error
}

func (repository *RepositoryStockOpnameImpl) InsertBulkUploadItems(c context.Context, items []model.StockOpnameBulkUploadItem) error {
	if len(items) == 0 {
		return nil
	}
	return repository.model(c).Create(&items).Error
}
