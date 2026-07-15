package model

import (
	"inventory/pkg"
	"strings"
	"time"

	"gorm.io/gorm"
)

// StockOpnameV2 model for v2 API - independent from v1
type StockOpnameV2 struct {
	CustID             string         `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	DocNo              string         `gorm:"column:doc_no;primaryKey" json:"doc_no"`
	WhID               int64          `gorm:"column:wh_id;not null" json:"wh_id"`
	Notes              *string        `gorm:"column:notes" json:"notes"`
	DataStatus         int            `gorm:"column:data_status;not null" json:"data_status"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          *time.Time     `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at" json:"updated_at,omitempty"`
	ScheduledAt        *time.Time     `gorm:"column:scheduled_at" json:"scheduled_at,omitempty"`
	AssignToEmpID      *int64         `gorm:"column:assign_to_emp_id" json:"assign_to_emp_id"`
	ProductHierarchy   int            `gorm:"column:product_hierarchy;not null" json:"product_hierarchy"`
	IncludeZeroStock   *bool          `gorm:"column:include_zero_stock" json:"include_zero_stock"`
	IsShowCurrentStock *bool          `gorm:"column:is_show_current_stock" json:"is_show_current_stock"`
	IsProcess          bool           `gorm:"column:is_process;not null;default:false" json:"is_process"`
	StockType          string         `gorm:"column:stock_type;not null;default:'G'" json:"stock_type"`
	DivisionID         *int64         `gorm:"column:division_id" json:"division_id"`
	PrincipalID        pkg.Int64Array `gorm:"column:principal_id;type:int8[]" json:"principal_id"`
	PLLane             pkg.Int64Array `gorm:"column:pl_lane;type:int8[]" json:"pl_lane"`
	BrandID            pkg.Int64Array `gorm:"column:brand_id;type:int8[]" json:"brand_id"`
	SBrand1ID          pkg.Int64Array `gorm:"column:sbrand1_id;type:int8[]" json:"sbrand1_id"`
	InputBy            string         `gorm:"column:input_by;not null;default:'Web'" json:"input_by"`
	EmpID              int64          `gorm:"column:emp_id;not null;default:0" json:"emp_id"`
	IsRevised          *bool          `gorm:"column:is_revised;default:false" json:"is_revised"`
}

func (StockOpnameV2) TableName() string {
	return "inv.stock_opname"
}

func (m *StockOpnameV2) BeforeCreate(trx *gorm.DB) (err error) {
	var numberDoc NumberDoc
	trCode := "SE"

	var dateForFormat time.Time
	if m.ScheduledAt != nil {
		dateForFormat = *m.ScheduledAt
	} else {
		dateForFormat = time.Now()
	}

	dateSubtr := dateForFormat.Format("060102") // YYMMDD format

	queryStr := `
		SELECT TRIM(
			to_char(
				COALESCE(MAX(TO_NUMBER(SUBSTR(doc_no, 9, 3), '999')), 0) + 1,
				'000'
			)
		) AS next_seq
		FROM inv.stock_opname
		WHERE substr(doc_no, 3, 6) = ? AND cust_id = ?
	`

	err = trx.Raw(queryStr, dateSubtr, strings.ToUpper(m.CustID)).Scan(&numberDoc).Error
	if err != nil {
		return err
	}

	m.DocNo = trCode + dateSubtr + numberDoc.NextSeq
	now := time.Now()
	m.CreatedAt = &now
	return nil
}

// StockOpnameDetailV2 model for v2 API - independent from v1
type StockOpnameDetailV2 struct {
	StockOpnameDetID int64      `gorm:"column:stock_opname_detail_id;primaryKey;autoIncrement" json:"stock_opname_detail_id"`
	CustID           string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	DocNo            string     `gorm:"column:doc_no;primaryKey" json:"doc_no"`
	ProID            int64      `gorm:"column:pro_id;primaryKey" json:"pro_id"`
	UnitID1          string     `gorm:"column:unit_id1;not null" json:"unit_id1"`
	UnitID2          string     `gorm:"column:unit_id2;not null" json:"unit_id2"`
	UnitID3          string     `gorm:"column:unit_id3;not null" json:"unit_id3"`
	ConvUnit1        float64    `gorm:"column:conv_unit1;not null;default:0" json:"conv_unit1"`
	ConvUnit2        float64    `gorm:"column:conv_unit2;not null;default:0" json:"conv_unit2"`
	ConvUnit3        float64    `gorm:"column:conv_unit3;not null;default:0" json:"conv_unit3"`
	Qty1             float64    `gorm:"column:qty1;not null;default:0" json:"qty1"`
	Qty2             float64    `gorm:"column:qty2;not null;default:0" json:"qty2"`
	Qty3             float64    `gorm:"column:qty3;not null;default:0" json:"qty3"`
	QtyStock1        float64    `gorm:"column:qty_stock1" json:"qty_stock1"`
	QtyStock2        float64    `gorm:"column:qty_stock2" json:"qty_stock2"`
	QtyStock3        float64    `gorm:"column:qty_stock3" json:"qty_stock3"`
	QtySO1           *float64   `gorm:"column:qty_so1;default:0" json:"qty_so1"`
	QtySO2           *float64   `gorm:"column:qty_so2;default:0" json:"qty_so2"`
	QtySO3           *float64   `gorm:"column:qty_so3;default:0" json:"qty_so3"`
	QtyOpname        *float64   `gorm:"column:qty_opname" json:"qty_opname"`
	QtyRevised1      *float64   `gorm:"column:qty_revised1" json:"qty_revised1"`
	QtyRevised2      *float64   `gorm:"column:qty_revised2" json:"qty_revised2"`
	QtyRevised3      *float64   `gorm:"column:qty_revised3" json:"qty_revised3"`
	StockOpnameDate  *time.Time `gorm:"column:stock_opname_date" json:"stock_opname_date,omitempty"`
	UserRevised      *int64     `gorm:"column:user_revised" json:"user_revised"`
	RevisedDate      *time.Time `gorm:"column:revised_date" json:"revised_date,omitempty"`
	PurchPrice1      float64    `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2      float64    `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3      float64    `gorm:"column:purch_price3" json:"purch_price3"`
	ProStatusBefore  *int       `gorm:"column:pro_status_before" json:"pro_status_before"`
	IsActiveBefore   *bool      `gorm:"column:is_active_before" json:"is_active_before"`
	CreatedBy        *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
}

func (StockOpnameDetailV2) TableName() string {
	return "inv.stock_opname_details"
}

// StockOpnameDetailV2Header for detail query result
type StockOpnameDetailV2Header struct {
	DocNo         string  `gorm:"column:doc_no"`
	CreatedDate   string  `gorm:"column:created_date"`
	WhID          int64   `gorm:"column:wh_id"`
	WhCode        string  `gorm:"column:wh_code"`
	WhName        string  `gorm:"column:wh_name"`
	StockType     string  `gorm:"column:stock_type"`
	CreatedBy     int64   `gorm:"column:created_by"`
	UserName      string  `gorm:"column:user_name"`
	ScheduledDate *string `gorm:"column:scheduled_date"`
	Status        int     `gorm:"column:status"`
	IsRevised     bool    `gorm:"column:is_revised"`
	IsProcess     bool    `gorm:"column:is_process"`
	InputBy       string  `gorm:"column:input_by"`
	DivisionID    *int64  `gorm:"column:division_id"`
	DivisionName  string  `gorm:"column:division_name"`
	EmpID         int64   `gorm:"column:emp_id"`
	EmpName       string  `gorm:"column:emp_name"`
	Notes         *string `gorm:"column:notes"`
}

func (StockOpnameDetailV2Header) TableName() string {
	return "inv.stock_opname"
}

// StockOpnameDetailV2Product for product list query result
type StockOpnameDetailV2Product struct {
	StockOpnameDetailID int64    `gorm:"column:stock_opname_detail_id"`
	ProID               int64    `gorm:"column:pro_id"`
	ProCode             string   `gorm:"column:pro_code"`
	ProName             string   `gorm:"column:pro_name"`
	UnitID1             string   `gorm:"column:unit_id1"`
	UnitID2             string   `gorm:"column:unit_id2"`
	UnitID3             string   `gorm:"column:unit_id3"`
	UnitName1           string   `gorm:"column:unit_name1"`
	UnitName2           string   `gorm:"column:unit_name2"`
	UnitName3           string   `gorm:"column:unit_name3"`
	ConvUnit2           float64  `gorm:"column:conv_unit2"`
	ConvUnit3           float64  `gorm:"column:conv_unit3"`
	QtyStock1           float64  `gorm:"column:qty_stock1"`
	QtyStock2           float64  `gorm:"column:qty_stock2"`
	QtyStock3           float64  `gorm:"column:qty_stock3"`
	QtyOpname1          float64  `gorm:"column:qty_opname1"`
	QtyOpname2          float64  `gorm:"column:qty_opname2"`
	QtyOpname3          float64  `gorm:"column:qty_opname3"`
	QtyRevised1         *float64 `gorm:"column:qty_revised1"`
	QtyRevised2         *float64 `gorm:"column:qty_revised2"`
	QtyRevised3         *float64 `gorm:"column:qty_revised3"`
	PurchPrice1         float64  `gorm:"column:purch_price1"`
	PurchPrice2         float64  `gorm:"column:purch_price2"`
	PurchPrice3         float64  `gorm:"column:purch_price3"`
	SellPrice1          float64  `gorm:"column:sell_price1"`
	SellPrice2          float64  `gorm:"column:sell_price2"`
	SellPrice3          float64  `gorm:"column:sell_price3"`
}

func (StockOpnameDetailV2Product) TableName() string {
	return "inv.stock_opname_details"
}

// StockOpnameLog model for logging status changes
type StockOpnameLog struct {
	IDStockOpnameLog int64      `gorm:"column:id_stock_opname_log;primaryKey;autoIncrement" json:"id_stock_opname_log"`
	Title            string     `gorm:"column:title;not null" json:"title"`
	ExecutionTime    time.Time  `gorm:"column:execution_time;not null" json:"execution_time"`
	OldStatus        int        `gorm:"column:old_status;not null" json:"old_status"`
	Status           int        `gorm:"column:status;not null" json:"status"`
	RefID            string     `gorm:"column:ref_id;not null" json:"ref_id"`
	TransactionCode  string     `gorm:"column:transaction_code;not null" json:"transaction_code"`
	RefTableName     string     `gorm:"column:ref_table_name;not null;default:'inv.stock_opname'" json:"ref_table_name"`
	TriggeredBy      string     `gorm:"column:triggered_by;not null;default:'MANUAL'" json:"triggered_by"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	CreatedBy        *int64     `gorm:"column:created_by" json:"created_by"`
	CustID           string     `gorm:"column:cust_id;not null" json:"cust_id"`
}

func (StockOpnameLog) TableName() string {
	return "inv.stock_opname_log"
}

// StockOpnameDetailForCompleted for completed stock opname calculation
type StockOpnameDetailForCompleted struct {
	StockOpnameDetailID int64      `gorm:"column:stock_opname_detail_id"`
	ProID               int64      `gorm:"column:pro_id"`
	QtyStock1           float64    `gorm:"column:qty_stock1"`
	QtyStock2           float64    `gorm:"column:qty_stock2"`
	QtyStock3           float64    `gorm:"column:qty_stock3"`
	QtySO1              *float64   `gorm:"column:qty_so1"`
	QtySO2              *float64   `gorm:"column:qty_so2"`
	QtySO3              *float64   `gorm:"column:qty_so3"`
	QtyRevised1         *float64   `gorm:"column:qty_revised1"`
	QtyRevised2         *float64   `gorm:"column:qty_revised2"`
	QtyRevised3         *float64   `gorm:"column:qty_revised3"`
	RevisedDate         *time.Time `gorm:"column:revised_date"`
	ConvUnit2           float64    `gorm:"column:conv_unit2"`
	ConvUnit3           float64    `gorm:"column:conv_unit3"`
}

// ProductPrice for product price information
type ProductPrice struct {
	ProID      int64   `gorm:"column:pro_id"`
	SellPrice1 float64 `gorm:"column:sell_price1"`
	Cogs       float64 `gorm:"column:cogs"`
}

// StockOpnameForUpdate model for fetching stock opname data for update
type StockOpnameForUpdate struct {
	CustID     string `gorm:"column:cust_id"`
	DocNo      string `gorm:"column:doc_no"`
	DataStatus int    `gorm:"column:data_status"`
	CreatedBy  *int64 `gorm:"column:created_by"`
	IsProcess  bool   `gorm:"column:is_process"`
}

func (StockOpnameForUpdate) TableName() string {
	return "inv.stock_opname"
}

// StockOpnameStartData model for fetching stock opname data for start
type StockOpnameStartData struct {
	CustID     string `gorm:"column:cust_id"`
	DocNo      string `gorm:"column:doc_no"`
	WhID       int64  `gorm:"column:wh_id"`
	WhCode     string `gorm:"column:wh_code"`
	WhName     string `gorm:"column:wh_name"`
	DataStatus int    `gorm:"column:data_status"`
}

func (StockOpnameStartData) TableName() string {
	return "inv.stock_opname"
}

// StockOpnameDetailQtySO for submit update
type StockOpnameDetailQtySO struct {
	StockOpnameDetID int64   `gorm:"column:stock_opname_detail_id"`
	QtySO1           float64 `gorm:"column:qty_so1"`
	QtySO2           float64 `gorm:"column:qty_so2"`
	QtySO3           float64 `gorm:"column:qty_so3"`
}

// StockOpnameBulkUpload model for inv.stock_opname_bulk_upload
type StockOpnameBulkUpload struct {
	UploadID   int    `gorm:"column:upload_id;primaryKey;autoIncrement"`
	DocNo      string `gorm:"column:doc_no;not null"`
	FilePath   string `gorm:"column:file_path;not null"`
	Status     *int   `gorm:"column:status"`
	TotalRow   *int   `gorm:"column:total_row"`
	ValidRow   *int   `gorm:"column:valid_row"`
	InvalidRow *int   `gorm:"column:invalid_row"`
	UploadedBy string `gorm:"column:uploaded_by;not null"`
	UploadedAt int64  `gorm:"column:uploaded_at;not null"`
}

func (StockOpnameBulkUpload) TableName() string {
	return "inv.stock_opname_bulk_upload"
}

// StockOpnameBulkUploadItem model for inv.stock_opname_bulk_upload_items
type StockOpnameBulkUploadItem struct {
	StockOpnameBulkUploadID int64   `gorm:"column:stock_opname_bulk_upload_id;primaryKey;autoIncrement"`
	UploadID                int     `gorm:"column:upload_id;not null"`
	ProductID               int64   `gorm:"column:product_id;not null"`
	QtySO1                  float64 `gorm:"column:qty_so1;not null;default:0"`
	QtySO2                  float64 `gorm:"column:qty_so2;not null;default:0"`
	QtySO3                  float64 `gorm:"column:qty_so3;not null;default:0"`
	QtyRevised1             float64 `gorm:"column:qty_revised1;not null;default:0"`
	QtyRevised2             float64 `gorm:"column:qty_revised2;not null;default:0"`
	QtyRevised3             float64 `gorm:"column:qty_revised3;not null;default:0"`
	UnitID1                 float64 `gorm:"column:unit_id1;not null;default:0"`
	UnitID2                 float64 `gorm:"column:unit_id2;not null;default:0"`
	UnitID3                 float64 `gorm:"column:unit_id3;not null;default:0"`
}

func (StockOpnameBulkUploadItem) TableName() string {
	return "inv.stock_opname_bulk_upload_items"
}
