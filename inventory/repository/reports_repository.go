package repository

import (
	"context"
	"fmt"
	"inventory/model"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type (
	RepositoryReportsImpl struct {
		*gorm.DB
	}
)

type ReportsRepository interface {
	GetStockMovementWarehouseTotalStock(ctx context.Context, custID string, month, year int) ([]model.StockMovementWarehouseTotalStock, error)
	GetStockMovementTransactionTypes(ctx context.Context, custID string, month, year int) ([]model.StockMovementTransactionType, error)
	GetTopProductsByStockIn(ctx context.Context, custID string, month, year int) ([]model.StockMovementTopProduct, error)
	GetTopProductsByStockOut(ctx context.Context, custID string, month, year int) ([]model.StockMovementTopProduct, error)
	GetNetStockChangesCurrentMonth(ctx context.Context, custID string, month, year int) (int64, error)
	GetNetStockChangesPreviousMonth(ctx context.Context, custID string, month, year int) (int64, error)
	GetStockLedger(ctx context.Context, req *model.StockLedgerRequest) (data []*model.StockLedgerRow, err error)
	CreateReportList(ctx context.Context, report *model.ReportList) error
	UpdateReportListFile(ctx context.Context, reportID string, fileStatus int, fileBase64, fileURL string) error
}

func NewReportsRepo(db *gorm.DB) *RepositoryReportsImpl {
	return &RepositoryReportsImpl{db}
}

func (repo *RepositoryReportsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

// GetStockMovementWarehouseTotalStock gets warehouse total stock data
func (repo *RepositoryReportsImpl) GetStockMovementWarehouseTotalStock(ctx context.Context, custID string, month, year int) ([]model.StockMovementWarehouseTotalStock, error) {
	var results []model.StockMovementWarehouseTotalStock

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	var endDate time.Time
	if month == 12 {
		endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		endDate = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	}

	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	query := repo.model(ctx).
		Select(`
			mw.wh_id,
			mw.wh_code,
			mw.wh_name,
			(COALESCE(sa.total_qty_in, 0) - COALESCE(sa.total_qty_out, 0))::BIGINT AS opening_stock,
			(COALESCE(sc.total_qty_in, 0) - COALESCE(sc.total_qty_out, 0))::BIGINT AS changing_stock,
			((COALESCE(sa.total_qty_in, 0) - COALESCE(sa.total_qty_out, 0)) + (COALESCE(sc.total_qty_in, 0) - COALESCE(sc.total_qty_out, 0)))::BIGINT AS closing_stock
		`).
		Table("mst.m_warehouse mw").
		Joins(fmt.Sprintf(`LEFT JOIN (
			SELECT
				s.wh_id,
				SUM(s.qty_in) AS total_qty_in,
				SUM(s.qty_out) AS total_qty_out
			FROM inv.stock s
			JOIN mst.m_product mp ON mp.pro_id = s.pro_id
			WHERE s.cust_id = '%s' AND s.stock_date < DATE '%s'
			GROUP BY s.wh_id
		) sa ON sa.wh_id = mw.wh_id`, custID, startDateStr)).
		Joins(fmt.Sprintf(`LEFT JOIN (
			SELECT
				s.wh_id,
				SUM(s.qty_in) AS total_qty_in,
				SUM(s.qty_out) AS total_qty_out
			FROM inv.stock s
			JOIN mst.m_product mp ON mp.pro_id = s.pro_id
			WHERE s.cust_id = '%s' AND s.stock_date >= DATE '%s' AND s.stock_date < DATE '%s'
			GROUP BY s.wh_id
		) sc ON sc.wh_id = mw.wh_id`, custID, startDateStr, endDateStr)).
		Where("mw.cust_id = ?", custID).
		Order("mw.wh_id")

	err := query.Scan(&results).Error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return results, nil
}

// GetStockMovementTransactionTypes gets transaction types with document counts
func (repo *RepositoryReportsImpl) GetStockMovementTransactionTypes(ctx context.Context, custID string, month, year int) ([]model.StockMovementTransactionType, error) {
	var results []model.StockMovementTransactionType

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	var endDate time.Time
	if month == 12 {
		endDate = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		endDate = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	}

	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	query := repo.model(ctx).
		Select(`
			t.tr_code,
			t.tr_name,
			COUNT(DISTINCT s.tr_no) AS no_of_doc
		`).
		Table("sys.m_trans t").
		Joins(fmt.Sprintf("LEFT JOIN inv.stock s ON s.tr_code = t.tr_code AND s.cust_id = '%s' AND s.stock_date >= DATE '%s' AND s.stock_date < DATE '%s'", custID, startDateStr, endDateStr)).
		Group("t.tr_code, t.tr_name").
		Order("t.tr_code")

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetTopProductsByStockIn gets top 5 products by stock in
func (repo *RepositoryReportsImpl) GetTopProductsByStockIn(ctx context.Context, custID string, month, year int) ([]model.StockMovementTopProduct, error) {
	results := make([]model.StockMovementTopProduct, 0)

	query := repo.model(ctx).
		Select(`
			p.pro_id,
			p.pro_code,
			p.pro_name,
			COALESCE(SUM(s.qty_in), 0) AS total_qty,
			p.conv_unit2,
			p.conv_unit3
		`).
		Table("inv.stock s").
		Joins("JOIN mst.m_product p ON p.pro_id = s.pro_id AND p.cust_id = ?", custID).
		Where("s.cust_id = ?", custID).
		Where("EXTRACT(YEAR FROM s.stock_date) = ?", year).
		Where("EXTRACT(MONTH FROM s.stock_date) = ?", month).
		Group("p.pro_id, p.pro_code, p.pro_name, p.conv_unit2, p.conv_unit3").
		Having("SUM(s.qty_in) > 0").
		Order("total_qty DESC").
		Limit(5)

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if results == nil {
		return make([]model.StockMovementTopProduct, 0), nil
	}

	return results, nil
}

// GetTopProductsByStockOut gets top 5 products by stock out
func (repo *RepositoryReportsImpl) GetTopProductsByStockOut(ctx context.Context, custID string, month, year int) ([]model.StockMovementTopProduct, error) {
	results := make([]model.StockMovementTopProduct, 0)

	query := repo.model(ctx).
		Select(`
			p.pro_id,
			p.pro_code,
			p.pro_name,
			COALESCE(SUM(s.qty_out), 0) AS total_qty,
			p.conv_unit2,
			p.conv_unit3
		`).
		Table("inv.stock s").
		Joins("JOIN mst.m_product p ON p.pro_id = s.pro_id AND p.cust_id = ?", custID).
		Where("s.cust_id = ?", custID).
		Where("EXTRACT(YEAR FROM s.stock_date) = ?", year).
		Where("EXTRACT(MONTH FROM s.stock_date) = ?", month).
		Group("p.pro_id, p.pro_code, p.pro_name, p.conv_unit2, p.conv_unit3").
		Having("SUM(s.qty_out) > 0").
		Order("total_qty DESC").
		Limit(5)

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if results == nil {
		return make([]model.StockMovementTopProduct, 0), nil
	}

	return results, nil
}

// GetNetStockChangesCurrentMonth gets net stock changes for current month
func (repo *RepositoryReportsImpl) GetNetStockChangesCurrentMonth(ctx context.Context, custID string, month, year int) (int64, error) {
	var result struct {
		Stock int64 `gorm:"column:stock"`
	}

	query := repo.model(ctx).
		Select("COALESCE(SUM(s.qty_in) - SUM(s.qty_out), 0)::BIGINT AS stock").
		Table("inv.stock s").
		Where("s.cust_id = ?", custID).
		Where("EXTRACT(YEAR FROM s.stock_date) = ?", year).
		Where("EXTRACT(MONTH FROM s.stock_date) = ?", month)

	err := query.Scan(&result).Error
	if err != nil {
		return 0, err
	}

	return result.Stock, nil
}

// GetNetStockChangesPreviousMonth gets net stock changes for previous month
func (repo *RepositoryReportsImpl) GetNetStockChangesPreviousMonth(ctx context.Context, custID string, month, year int) (int64, error) {
	var result struct {
		Stock int64 `gorm:"column:stock"`
	}

	prevMonth := month - 1
	prevYear := year
	if prevMonth < 1 {
		prevMonth = 12
		prevYear = year - 1
	}

	query := repo.model(ctx).
		Select("COALESCE(SUM(s.qty_in) - SUM(s.qty_out), 0) AS stock").
		Table("inv.stock s").
		Where("s.cust_id = ?", custID).
		Where("EXTRACT(YEAR FROM s.stock_date) = ?", prevYear).
		Where("EXTRACT(MONTH FROM s.stock_date) = ?", prevMonth)

	err := query.Scan(&result).Error
	if err != nil {
		return 0, err
	}

	return result.Stock, nil
}

func (repo *RepositoryReportsImpl) GetStockLedger(ctx context.Context, req *model.StockLedgerRequest) (data []*model.StockLedgerRow, err error) {
	query := `WITH
params AS (
   SELECT
       COALESCE(CAST(@start_date AS date), DATE '1900-01-01') AS start_dt,
       COALESCE(CAST(@end_date AS date), CURRENT_DATE) AS end_dt,
       COALESCE(CAST(@warehouse_ids AS INT[]), CAST(ARRAY[] AS INT[])) AS wh_ids,
       COALESCE(CAST(@product_ids AS BIGINT[]), CAST(ARRAY[] AS BIGINT[])) AS product_ids,
       COALESCE(CAST(@sup_ids AS BIGINT[]), CAST(ARRAY[] AS BIGINT[])) AS sup_ids,
       COALESCE(CAST(@transaction_types AS TEXT[]), CAST(ARRAY[] AS TEXT[])) AS transaction_types
),
filtered_products AS (
   SELECT
       p.pro_id,
       p.pro_code,
       p.pro_name,
       p.cust_id,
       p.unit_id1, p.unit_id2, p.unit_id3,
       p.conv_unit2, p.conv_unit3,
       p.sbrand1_id,
	   p.sup_id,
       sb1.brand_id,
       b.pl_id,
       b.brand_name,
       sb1.sbrand1_name,
       pr.principal_name,
       pl.pl_name
   FROM mst.m_product p
   JOIN mst.m_sub_brand1 sb1 ON p.sbrand1_id = sb1.sbrand1_id
   JOIN mst.m_brand b ON sb1.brand_id = b.brand_id
   JOIN mst.m_product_line pl ON b.pl_id = pl.pl_id
   LEFT JOIN mst.m_principal pr ON p.principal_id = pr.principal_id
   JOIN params par ON TRUE
   WHERE
	(
		array_length(par.product_ids,1) IS NULL
		OR p.pro_id = ANY(par.product_ids)
	)
	AND (
		array_length(par.sup_ids,1) IS NULL
		OR p.sup_id = ANY(par.sup_ids)
	)
),
opening_balance AS (
   SELECT
       s.wh_id,
       s.pro_id,
       SUM(s.qty_in - s.qty_out) AS opening_qty_total
   FROM inv.stock s
   JOIN params p ON TRUE
   WHERE
       s.stock_date < p.start_dt
       AND s.wh_id = ANY(p.wh_ids)
       AND s.pro_id IN (SELECT pro_id FROM filtered_products)
	   AND s.tr_no NOT LIKE '%-CO'
   GROUP BY s.wh_id, s.pro_id
),
transactions AS (
   SELECT
       s.stock_id,
       s.stock_date,
       s.wh_id,
       s.pro_id,
       s.tr_code,
       t.tr_name,
       s.tr_no,
       s.qty_in,
       s.qty_out,
       (s.qty_in - s.qty_out) AS change_qty,
       s.created_at,
	   w.wh_name,
	   w.wh_code,
	   d.distributor_id,
	   d.distributor_code,
	   d.distributor_name,
       fp.unit_id1, fp.unit_id2, fp.unit_id3,
       fp.conv_unit2, fp.conv_unit3,
       fp.pro_code, fp.pro_name,
       fp.brand_name, fp.sbrand1_name, fp.principal_name, fp.sup_id
   FROM inv.stock s
   JOIN sys.m_trans t ON s.tr_code = t.tr_code
   JOIN mst.m_warehouse w ON s.wh_id = w.wh_id AND s.cust_id = w.cust_id
   JOIN mst.m_distributor d ON w.cust_id = d.cust_id
   JOIN filtered_products fp ON s.pro_id = fp.pro_id
   JOIN params p ON TRUE
   WHERE
       s.stock_date BETWEEN p.start_dt AND p.end_dt
	   AND s.tr_no NOT LIKE '%-CO'
       AND s.wh_id = ANY(p.wh_ids)
	   AND (
			array_length(p.transaction_types,1) IS NULL
			OR s.tr_code = ANY(p.transaction_types)
		)
		-- FILTER TAMBAHAN: Skip jika transaksinya tidak mengubah qty sama sekali (hanya numpang lewat)
		AND (s.qty_in > 0 OR s.qty_out > 0)

),
ledger_calc AS (
   SELECT
       t.*,
       COALESCE(ob.opening_qty_total, 0) AS initial_opening,
       SUM(t.change_qty) OVER (
           PARTITION BY t.wh_id, t.pro_id
           ORDER BY t.stock_date, t.stock_id
       ) AS running_change
   FROM transactions t
   LEFT JOIN opening_balance ob
       ON t.wh_id = ob.wh_id AND t.pro_id = ob.pro_id
)
SELECT
	COUNT(*) OVER() AS total_record,
	l.distributor_id,
	l.distributor_code,
	l.distributor_name,
	l.wh_id,
	l.wh_code,
	l.wh_name AS warehouse,
	l.stock_date AS date,
	l.pro_id,
	l.pro_code AS product_code,
	l.pro_name AS product_name,

   (l.initial_opening + l.running_change - l.change_qty) as opening_stock,
   l.change_qty as updates,
   (l.initial_opening + l.running_change) as closing_stock,

  	-- Opening Stock Converted
	FLOOR(CAST((l.initial_opening + l.running_change - l.change_qty) AS numeric) / NULLIF(CAST((l.conv_unit2 * l.conv_unit3) AS numeric), 0)) as opening_stock_large,
	FLOOR(MOD(CAST((l.initial_opening + l.running_change - l.change_qty) AS numeric), NULLIF(CAST((l.conv_unit2 * l.conv_unit3) AS numeric), 0)) / NULLIF(CAST(l.conv_unit2 AS numeric), 0)) as opening_stock_medium,
	MOD(CAST((l.initial_opening + l.running_change - l.change_qty) AS numeric), NULLIF(CAST(l.conv_unit2 AS numeric), 0)) as opening_stock_small,
  	-- Updates Converted
    TRUNC(CAST(l.change_qty AS numeric) / NULLIF(CAST((l.conv_unit2 * l.conv_unit3) AS numeric), 0)) as updates_large,
	TRUNC(MOD(CAST(l.change_qty AS numeric), NULLIF(CAST((l.conv_unit2 * l.conv_unit3) AS numeric), 0)) / NULLIF(CAST(l.conv_unit2 AS numeric), 0)) as updates_medium,
	MOD(CAST(l.change_qty AS numeric), NULLIF(CAST(l.conv_unit2 AS numeric), 0)) as updates_small,
	-- Closing Stock Converted
	FLOOR(CAST((l.initial_opening + l.running_change) AS numeric) / NULLIF(CAST((l.conv_unit2 * l.conv_unit3) AS numeric), 0)) as closing_stock_large,
	FLOOR(MOD(CAST((l.initial_opening + l.running_change) AS numeric), NULLIF(CAST((l.conv_unit2 * l.conv_unit3) AS numeric), 0)) / NULLIF(CAST(l.conv_unit2 AS numeric), 0)) as closing_stock_medium,
	MOD(CAST((l.initial_opening + l.running_change) AS numeric), NULLIF(CAST(l.conv_unit2 AS numeric), 0)) as closing_stock_small,

	l.tr_name AS transaction_type,
	l.tr_no AS reference_no,
	l.sup_id as supplier_id

FROM ledger_calc l

ORDER BY l.wh_name, l.pro_name, l.stock_date, l.stock_id`

	queryParams := map[string]interface{}{
		"start_date":        req.StartDate,
		"end_date":          req.EndDate,
		"warehouse_ids":     pq.Array(req.WarehouseIDs),
		"transaction_types": pq.Array(req.TransactionTypes),
		"product_ids":       pq.Array(req.ProductIDs),
		"sup_ids":           pq.Array(req.SupIDs),
		// "principal_ids":     pq.Array(req.PrincipalIDs),
	}

	if req.Limit > 0 {
		page := req.Page
		if page <= 0 {
			page = 1
		}
		queryParams["limit"] = req.Limit
		queryParams["offset"] = (page - 1) * req.Limit
		query += `
LIMIT @limit OFFSET @offset`
	}

	db := repo.model(ctx)
	err = db.Raw(query, queryParams).Scan(&data).Error

	if err != nil {
		log.Error("GetStockLedger query error: ", err)
		return nil, err
	}

	return data, nil
}

// CreateReportList inserts a new report.list entry (initial state, without file content).
func (repo *RepositoryReportsImpl) CreateReportList(ctx context.Context, report *model.ReportList) error {
	if report == nil {
		return fmt.Errorf("report payload is nil")
	}

	if err := repo.model(ctx).Create(report).Error; err != nil {
		log.Error("CreateReportList error: ", err)
		return err
	}

	return nil
}

// UpdateReportListFile updates file-related fields (status, url, base64) for a given report_id.
func (repo *RepositoryReportsImpl) UpdateReportListFile(ctx context.Context, reportID string, fileStatus int, fileBase64, fileURL string) error {
	if reportID == "" {
		return fmt.Errorf("reportID is required")
	}

	updates := map[string]interface{}{
		"file_status": fileStatus,
		"file_url":    fileURL,
		"file_base64": fileBase64,
	}

	if err := repo.model(ctx).
		Table("report.list").
		Where("report_id = ?", reportID).
		Updates(updates).Error; err != nil {
		log.Error("UpdateReportListFile error: ", err)
		return err
	}

	return nil
}
