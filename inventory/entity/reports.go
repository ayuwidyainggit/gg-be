package entity

import "time"

// StockMovementReportQueryFilter for stock movement report query parameters
type StockMovementReportQueryFilter struct {
	CustID       string `query:"cust_id" json:"cust_id,omitempty"`
	ParentCustID string `query:"parent_cust_id" json:"parent_cust_id,omitempty"`
	Month        int    `query:"month" json:"month"`
	Year         int    `query:"year" json:"year"`
}

// StockMovementReportResponse response structure for stock movement report
type StockMovementReportResponse struct {
	WhTotalStock    []StockMovementWarehouseTotalStock `json:"wh_total_stock"`
	NetStockChanges StockMovementNetStockChanges       `json:"net_stock_changes"`
	StockMovement   []StockMovementTransactionType     `json:"stock_movement"`
	TopProductIn    []StockMovementTopProduct          `json:"top_product_in"`
	TopProductOut   []StockMovementTopProduct          `json:"top_product_out"`
}

// StockMovementWarehouseTotalStock warehouse total stock data
type StockMovementWarehouseTotalStock struct {
	WhID          int64  `json:"wh_id"`
	WhCode        string `json:"wh_code"`
	WhName        string `json:"wh_name"`
	OpeningStock  int64  `json:"opening_stock"`
	ChangingStock int64  `json:"changing_stock"`
	ClosingStock  int64  `json:"closing_stock"`
}

// StockMovementNetStockChanges net stock changes data
type StockMovementNetStockChanges struct {
	StockAwal  int64 `json:"stock_awal"`
	StockAkhir int64 `json:"stock_akhir"`
	GrowStock  int64 `json:"grow_stock"`
}

// StockMovementTransactionType transaction type data
type StockMovementTransactionType struct {
	TrCode  string `json:"tr_code"`
	TrName  string `json:"tr_name"`
	NoOfDoc int64  `json:"no_of_doc"`
}

// StockMovementTopProduct top product data
type StockMovementTopProduct struct {
	ProName     string `json:"pro_name"`
	QtyLargest  int64  `json:"qty_largest"`
	QtyMedium   int64  `json:"qty_medium"`
	QtySmallest int64  `json:"qty_smallest"`
}

type (
	PreviewDownloadStockMovementReportQueryFilter struct {
		CustID       string `query:"cust_id" json:"cust_id,omitempty"`
		ParentCustID string `query:"parent_cust_id" json:"parent_cust_id,omitempty"`

		StartDate        string   `query:"start_date" json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
		EndDate          string   `query:"end_date" json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
		WarehouseIDs     []int64  `query:"wh_id" json:"wh_id,omitempty" validate:"omitempty,dive,numeric"`
		TransactionTypes []string `query:"transaction_type" json:"transaction_type,omitempty"`
		SupIDs           []int64  `query:"sup_id" json:"sup_id,omitempty" validate:"omitempty,dive,numeric"`
		PrincipalIDs     []int64  `query:"principal_id" json:"principal_id,omitempty" validate:"omitempty,dive,numeric"`
		ProductIDs       []int64  `query:"pro_id" json:"pro_id,omitempty" validate:"omitempty,dive,numeric"`
		Page             int      `query:"page" json:"page,omitempty"`
		Limit            int      `query:"limit" json:"limit,omitempty"`

		// Filled from JWT middleware, not from query
		UserID       int64  `json:"-"`
		UserFullName string `json:"-"`
	}

	PreviewDownloadStockMovementReportResponse struct {
		DistributorID   int64  `json:"distributor_id"`
		DistributorCode string `json:"distributor_code"`
		DistributorName string `json:"distributor_name"`

		WhID   int64     `json:"wh_id"`
		WhCode string    `json:"wh_code"`
		WhName string    `json:"wh_name"`
		Date   time.Time `json:"date"`

		ProID   int64  `json:"pro_id"`
		ProCode string `json:"pro_code"`
		ProName string `json:"pro_name"`

		OpeningStock1 int64 `json:"opening_stock1"`
		OpeningStock2 int64 `json:"opening_stock2"`
		OpeningStock3 int64 `json:"opening_stock3"`

		ChangesStock1 int64 `json:"changes_stock1"`
		ChangesStock2 int64 `json:"changes_stock2"`
		ChangesStock3 int64 `json:"changes_stock3"`

		ClosingStock1 int64 `json:"closing_stock1"`
		ClosingStock2 int64 `json:"closing_stock2"`
		ClosingStock3 int64 `json:"closing_stock3"`

		TransactionType string `json:"transaction_type"`
		RefNo           string `json:"ref_no"`
		Remarks         string `json:"remarks"`
	}
)

type (
	DownloadStockMovementReportQueryFilter struct {
		CustID       string `query:"cust_id" json:"cust_id,omitempty"`
		ParentCustID string `query:"parent_cust_id" json:"parent_cust_id,omitempty"`

		StartDate        string   `query:"start_date" json:"start_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
		EndDate          string   `query:"end_date" json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
		WarehouseIDs     []int64  `query:"wh_id" json:"wh_id,omitempty" validate:"omitempty,dive,numeric"`
		TransactionTypes []string `query:"transaction_type" json:"transaction_type,omitempty"`
		SupIDs           []int64  `query:"sup_id" json:"sup_id,omitempty" validate:"omitempty,dive,numeric"`
		PrincipalIDs     []int64  `query:"principal_id" json:"principal_id,omitempty" validate:"omitempty,dive,numeric"`
		ProductIDs       []int64  `query:"pro_id" json:"pro_id,omitempty" validate:"omitempty,dive,numeric"`

		// Filled from JWT middleware, not from query
		UserID       int64  `json:"-"`
		UserFullName string `json:"-"`
	}

	DownloadStockMovementReportResponse struct {
		ReportID   string `json:"report_id"`
		ReportName string `json:"report_name"`
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		FileStatus int    `json:"file_status"`
		// FileStatusName string    `json:"file_status_name"`
		FileURL   string    `json:"file_url"`
		CreatedBy string    `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
	}
)
