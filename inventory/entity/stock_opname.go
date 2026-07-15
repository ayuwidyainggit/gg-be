package entity

import "time"

type CancelStockOpanmeParams struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	CancelBy     int64
}

type UpdateStockOpnameParams struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
}

type ReportStockOpanmeParams struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
}

type StockOpnameQueryFilter struct {
	CustID       string
	ParentCustID string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	WhID         int    `query:"wh_id"`
	DocNo        string `query:"doc_no"`
	DataStatus   []int  `query:"data_status"`
}

type StockOpnameList struct {
	CreatedAt            string `json:"created_at"`
	DocNo                string `json:"doc_no"`
	WhID                 int    `json:"wh_id"`
	WhCode               string `json:"wh_code"`
	WhName               string `json:"wh_name"`
	EmpCode              string `json:"emp_code"`
	EmpName              string `json:"emp_name"`
	StockType            string `json:"stock_type"`
	ScheduledAt          string `json:"scheduled_at"`
	DataStatus           int    `json:"data_status"`
	StatusDescription    string `json:"status_description"`
	AssignToEmpID        int64  `json:"assign_to_emp_id"`
	ProductHierarchy     int    `json:"product_hierarchy"`
	ProductHierarchyDesc string `json:"product_hierarchy_desc"`
	IncludeZeroStock     bool   `json:"include_zero_stock"`
	IsShowCurrentStock   bool   `json:"is_show_current_stock"`
}

type StockOpnameReports struct {
	CreatedAt            string              `json:"created_at"`
	DocNo                string              `json:"doc_no"`
	WhID                 int                 `json:"wh_id"`
	WhCode               string              `json:"wh_code"`
	WhName               string              `json:"wh_name"`
	EmpCode              string              `json:"emp_code"`
	EmpName              string              `json:"emp_name"`
	StockType            string              `json:"stock_type"`
	ScheduledAt          string              `json:"scheduled_at"`
	DataStatus           int                 `json:"data_status"`
	StatusDescription    string              `json:"status_description"`
	AssignToEmpID        int64               `json:"assign_to_emp_id"`
	ProductHierarchy     int                 `json:"product_hierarchy"`
	ProductHierarchyDesc string              `json:"product_hierarchy_desc"`
	IncludeZeroStock     bool                `json:"include_zero_stock"`
	IsShowCurrentStock   bool                `json:"is_show_current_stock"`
	StockOpnameReports   []StockOpnameReport `json:"stock_opname_reports"`
}

type StockOpnameReport struct {
	StockReportID     string    `json:"stock_report_id"`
	Status            int       `json:"status"`
	StatusDescription string    `json:"status_description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateStockOpname struct {
	CustID             string `json:"cust_id" validate:"required"`
	ParentCustID       string `json:"parent_cust_id" validate:"required"`
	DataStatus         int    `json:"data_status"`
	WhID               int64  `json:"wh_id" validate:"required"`
	Notes              string `json:"notes" validate:""`
	ScheduledAt        string `json:"scheduled_at" validate:"required"`
	AssignToEmpID      int64  `json:"assign_to_emp_id" validate:"required"`
	ProductHierarchy   int    `json:"product_hierarchy" validate:"required"`
	IncludeZeroStock   bool   `json:"include_zero_stock"`
	IsShowCurrentStock bool   `json:"is_show_current_stock"`
	CreatedBy          int64  `json:"created_by"`
	UpdatedBy          int64  `json:"updated_by"`
}

type CreateStockOpnameV2 struct {
	WhID             int64                        `json:"wh_id" validate:"required"`
	StockType        string                       `json:"stock_type" validate:"required,oneof=G E BS"`
	ScheduleDate     string                       `json:"schedule_date" validate:"required"`
	ProductHierarchy int                          `json:"product_hierarchy" validate:"required,min=1,max=6"`
	PrincipalID      []int64                      `json:"principal_id"`
	PLLane           []int64                      `json:"pl_lane"`
	BrandID          []int64                      `json:"brand_id"`
	SBrand1ID        []int64                      `json:"sbrand1_id"`
	IncludeZeroStock int                          `json:"include_zero_stock" validate:"oneof=0 1"`
	InputBy          string                       `json:"input_by" validate:"required,oneof=Mobile Manual Web"`
	DivisionID       int64                        `json:"division_id" validate:"required"`
	EmpID            int64                        `json:"emp_id" validate:"required"`
	ProductList      []CreateStockOpnameV2Product `json:"product_list" validate:"required,min=1,dive"`
	CustID           string                       `json:"-"`
	ParentCustID     string                       `json:"-"`
	CreatedBy        int64                        `json:"-"`
}

type CreateStockOpnameV2Product struct {
	ProID     int64    `json:"pro_id" validate:"required"`
	UnitID1   string   `json:"unit_id1" validate:"required,max=5"`
	UnitID2   string   `json:"unit_id2" validate:"required,max=5"`
	UnitID3   string   `json:"unit_id3" validate:"required,max=5"`
	ConvUnit1 *float64 `json:"conv_unit1"`
	ConvUnit2 *float64 `json:"conv_unit2"`
	ConvUnit3 *float64 `json:"conv_unit3"`
	Qty1      *float64 `json:"qty1" validate:"required,min=0"`
	Qty2      *float64 `json:"qty2" validate:"required,min=0"`
	Qty3      *float64 `json:"qty3" validate:"required,min=0"`
}

type StockOpnameStatusDescSlice []StockOpnameStatus

func (p StockOpnameStatusDescSlice) Len() int {
	return len(p)
}

func (p StockOpnameStatusDescSlice) Less(i, j int) bool {
	return p[i].StatusID < p[j].StatusID
}

func (p StockOpnameStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type StockOpnameStatus struct {
	StatusID   int    `json:"status_id"`
	StatusDesc string `json:"status_desc"`
}

var StockOpnameStatusDesc = map[int]string{
	1: "Schedule",
	2: "Assign",
	3: "Need Review",
	4: "On Going",
	5: "Submitted",
	6: "Completed",
	7: "Cancelled",
}

func (stockOpname StockOpnameList) GetStockOpnameStatusDesc() string {
	return StockOpnameStatusDesc[stockOpname.DataStatus]
}

var OpnameReportStatusDesc = map[int]string{
	1: "Assigned", 5: "On Going", 50: "Cancelled", 60: "Neeed Review", 70: "Rejected", 100: "Approved",
}

func (soRep StockOpnameReport) GetOpnameReportStatusDesc() string {
	return OpnameReportStatusDesc[soRep.Status]
}

// Product Hierarchy

type ProductHierarchyDescSlice []ProductHierarchy

func (p ProductHierarchyDescSlice) Len() int {
	return len(p)
}

func (p ProductHierarchyDescSlice) Less(i, j int) bool {
	return p[i].ProductHierarchyID < p[j].ProductHierarchyID
}

func (p ProductHierarchyDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type ProductHierarchy struct {
	ProductHierarchyID   int    `json:"id"`
	ProductHierarchyDesc string `json:"description"`
}

var ProductHierarchyDesc = map[int]string{
	1: "Product", 2: "Category", 3: "Brand", 4: "Product Line",
}

func (stockOpname StockOpnameList) GetProductHierarchyDesc() string {
	return ProductHierarchyDesc[stockOpname.ProductHierarchy]
}

// Stock Opname Product List
type StockOpnameProductListQueryFilter struct {
	CustID       string
	ParentCustID string
	Query        string `query:"q"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Sort         string `query:"sort"`
	WhID         []int  `query:"wh_id"`
	StockType    string `query:"stock_type"`
	PrincipalID  []int  `query:"principal_id"`
	PLID         []int  `query:"pl_id"`
	BrandID      []int  `query:"brand_id"`
	SBrand1ID    []int  `query:"sbrand1_id"`
	ZeroStock    bool   `query:"zero_stock"`
	IsActive     bool   `query:"is_active"`
}

type StockOpnameProductListResponse struct {
	ProID         int64   `json:"pro_id"`
	ProCode       string  `json:"pro_code"`
	ProName       string  `json:"pro_name"`
	Uom1          string  `json:"uom1"`
	Uom2          string  `json:"uom2"`
	Uom3          string  `json:"uom3"`
	UnitName1     string  `json:"unit_name1"`
	UnitName2     string  `json:"unit_name2"`
	UnitName3     string  `json:"unit_name3"`
	ConvUnit2     float64 `json:"conv_unit2"`
	ConvUnit3     float64 `json:"conv_unit3"`
	Qty1          float64 `json:"qty1"`
	Qty2          float64 `json:"qty2"`
	Qty3          float64 `json:"qty3"`
	WhID          int64   `json:"wh_id"`
	StockType     string  `json:"stock_type"`
	PrincipalID   int64   `json:"principal_id"`
	PrincipalName string  `json:"principal_name"`
	PLID          int64   `json:"pl_id"`
	PLName        string  `json:"pl_name"`
	BrandID       int64   `json:"brand_id"`
	BrandName     string  `json:"brand_name"`
	SBrand1ID     int64   `json:"sbrand1_id"`
	SBrand1Name   string  `json:"sbrand1_name"`
	IsActive      bool    `json:"is_active"`
}

// Stock Opname List V2 (underscore URL)
type StockOpnameListV2QueryFilter struct {
	CustID       string
	ParentCustID string
	UserID       int64
	IsAdmin      bool
	Query        string  `query:"q"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit"`
	Sort         string  `query:"sort"`
	StartDate    *int64  `query:"start_date"`
	EndDate      *int64  `query:"end_date"`
	WhID         []int64 `query:"wh_id"`
	Status       []int   `query:"status"`
	EmpID        []int64 `query:"emp_id"`
}

type StockOpnameListV2Response struct {
	DocNo         string `json:"doc_no"`
	CreatedDate   string `json:"created_date"`
	WhID          int64  `json:"wh_id"`
	WhCode        string `json:"wh_code"`
	WhName        string `json:"wh_name"`
	CreatedBy     string `json:"created_by"`
	User          string `json:"user"`
	ScheduledDate string `json:"scheduled_date"`
	EmpID         int64  `json:"emp_id"`
	EmpName       string `json:"emp_name"`
	Status        int    `json:"status"`
	StatusDesc    string `json:"status_desc"`
}

func (s StockOpnameListV2Response) GetStatusDesc() string {
	return StockOpnameStatusDesc[s.Status]
}

// Stock Opname Detail V2 Params
type StockOpnameDetailV2Params struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
}

// Stock Opname Detail V2 Response
type StockOpnameDetailV2Response struct {
	DocNo                 string                           `json:"doc_no"`
	CreatedDate           string                           `json:"created_date"`
	WhID                  int64                            `json:"wh_id"`
	WhCode                string                           `json:"wh_code"`
	WhName                string                           `json:"wh_name"`
	StockType             string                           `json:"stock_type"`
	CreatedBy             int64                            `json:"created_by"`
	UserName              string                           `json:"user_name"`
	ScheduledDate         *string                          `json:"scheduled_date"`
	Status                int                              `json:"status"`
	IsRevised             bool                             `json:"is_revised"`
	IsProcess             bool                             `json:"is_process"`
	StatusDesc            string                           `json:"status_desc"`
	SubTotalSystem        float64                          `json:"sub_total_system"`
	SubTotalPhysicalStock float64                          `json:"sub_total_physical_stock"`
	Difference            float64                          `json:"difference"`
	AssignTo              StockOpnameDetailV2AssignTo      `json:"assign_to"`
	ProductList           []StockOpnameDetailV2ProductItem `json:"product_list"`
	Notes                 string                           `json:"notes"`
}

func (s *StockOpnameDetailV2Response) SetStatusDesc() {
	s.StatusDesc = StockOpnameStatusDesc[s.Status]
}

type StockOpnameDetailV2AssignTo struct {
	InputBy      string `json:"input_by"`
	DivisionID   int64  `json:"division_id"`
	DivisionName string `json:"division_name"`
	EmpID        int64  `json:"emp_id"`
	EmpName      string `json:"emp_name"`
}

type StockOpnameDetailV2ProductItem struct {
	StockOpnameDetailID int64   `json:"stock_opname_detail_id"`
	ProID               int64   `json:"pro_id"`
	ProCode             string  `json:"pro_code"`
	ProName             string  `json:"pro_name"`
	UnitID1             string  `json:"unit_id1"`
	UnitID2             string  `json:"unit_id2"`
	UnitID3             string  `json:"unit_id3"`
	SellPrice1          float64 `json:"sell_price1"`
	SellPrice2          float64 `json:"sell_price2"`
	SellPrice3          float64 `json:"sell_price3"`
	UnitName1           string  `json:"unit_name1"`
	UnitName2           string  `json:"unit_name2"`
	UnitName3           string  `json:"unit_name3"`
	ConvUnit2           float64 `json:"conv_unit2"`
	ConvUnit3           float64 `json:"conv_unit3"`
	Qty1                float64 `json:"qty1"`
	Qty2                float64 `json:"qty2"`
	Qty3                float64 `json:"qty3"`
	QtyPhysical1        float64 `json:"qty_physical1"`
	QtyPhysical2        float64 `json:"qty_physical2"`
	QtyPhysical3        float64 `json:"qty_physical3"`
	DifferentStock1     float64 `json:"different_stock1"`
	DifferentStock2     float64 `json:"different_stock2"`
	DifferentStock3     float64 `json:"different_stock3"`
	DifferentPrice1     float64 `json:"different_price1"`
	DifferentPrice2     float64 `json:"different_price2"`
	DifferentPrice3     float64 `json:"different_price3"`
}

// Stock Opname Update Status V2 (Manual Reprocess)
type UpdateStockOpnameStatusV2Params struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	UserID       int64
}

type UpdateStockOpnameStatusV2Request struct {
	IsProcess   *bool `json:"is_process"`
	IsAssigne   *bool `json:"is_assigne"`
	IsCompleted *bool `json:"is_completed"`
	IsCancelled *bool `json:"is_cancelled"`
	OldStatus   int   `json:"old_status" validate:"required"`
}

const (
	StockOpnameStatusSchedule   = 1
	StockOpnameStatusAssign     = 2
	StockOpnameStatusNeedReview = 3
	StockOpnameStatusOnGoing    = 4
	StockOpnameStatusSubmit     = 5
	StockOpnameStatusCompleted  = 6
	StockOpnameStatusRejected   = 7
)

// Stock Opname Revised V2
type RevisedStockOpnameV2Params struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	UserID       int64
}

type RevisedStockOpnameV2Item struct {
	StockOpnameDetID int64   `json:"stockopname_det_id"`
	ProID            int64   `json:"pro_id" validate:"required"`
	QtyRevised1      float64 `json:"qty_revised1" validate:"min=0"`
	QtyRevised2      float64 `json:"qty_revised2" validate:"min=0"`
	QtyRevised3      float64 `json:"qty_revised3" validate:"min=0"`
}

type RevisedStockOpnameV2Request struct {
	Data []RevisedStockOpnameV2Item `json:"data" validate:"required,min=1,dive"`
}

// Stock Opname Template Download
type StockOpnameTemplateDownloadParams struct {
	DocNo        string `query:"doc_no"`
	CustID       string
	ParentCustID string
}

type StockOpnameTemplateDownloadResponse struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	ExpiresAt      string `json:"expires_at"`
	DownloadBase64 string `json:"download_base64"`
}

type StockOpnameDownloadParams struct {
	DocNo        string `query:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	UserID       int64
}

type StockOpnameDownloadResponse struct {
	ReportName string `json:"report_name"`
	FileBase64 string `json:"file_base64"`
}

// Stock Opname Start V2
type StartStockOpnameV2Params struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	UserID       int64
}

type StartStockOpnameV2Response struct {
	DocNo         string `json:"doc_no"`
	WhID          int64  `json:"wh_id"`
	WhCode        string `json:"wh_code"`
	WhName        string `json:"wh_name"`
	Status        int    `json:"status"`
	StatusDesc    string `json:"status_desc"`
	StartedAt     string `json:"started_at"`
	StartedBy     int64  `json:"started_by"`
	StartedByName string `json:"started_by_name"`
}

// Stock Opname Submit V2
type SubmitStockOpnameV2Params struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	UserID       int64
}

type SubmitStockOpnameV2Detail struct {
	StockOpnameDetID int64    `json:"stock_opname_detail_id" form:"stock_opname_detail_id" validate:"required"`
	ProID            int64    `json:"pro_id" form:"pro_id"`
	QtySO1           *float64 `json:"qty_so1" form:"qty_so1" validate:"required,min=0"`
	QtySO2           *float64 `json:"qty_so2" form:"qty_so2" validate:"required,min=0"`
	QtySO3           *float64 `json:"qty_so3" form:"qty_so3" validate:"required,min=0"`
}

type SubmitStockOpnameV2Request struct {
	Details []SubmitStockOpnameV2Detail `json:"details" form:"items" validate:"required,min=1,dive"`
}

type SubmitStockOpnameV2Response struct {
	DocNo          string                      `json:"doc_no"`
	UpdatedCount   int                         `json:"updated_count"`
	UpdatedDetails []SubmitStockOpnameV2Detail `json:"updated_details"`
}

// BulkUploadStockOpnameV2Params for POST /v1/stock_opname/bulk_upload/:doc_no
type BulkUploadStockOpnameV2Params struct {
	DocNo        string `params:"doc_no" validate:"required"`
	CustID       string
	ParentCustID string
	UserID       int64
}

// BulkUploadRowError error per baris Excel (untuk respons saat ada baris gagal)
type BulkUploadRowError struct {
	RowIndex   int    `json:"row_index"` // baris di Excel (1-based, baris 18 = 18)
	DetailID   int64  `json:"detail_id,omitempty"`
	ProID      int64  `json:"pro_id,omitempty"`
	ErrMessage string `json:"error_message"`
}

// BulkUploadStockOpnameV2Response response for bulk upload
type BulkUploadStockOpnameV2Response struct {
	DocNo             string                          `json:"doc_no"`
	Status            string                          `json:"status"` // FULL_SUCCESS, PARTIAL_SUCCESS, FAILED
	StatusDescription string                          `json:"status_desctiption"`
	TotalRow          int                             `json:"total_row"`
	SuccessRow        int                             `json:"success_row"`
	FailedRow         int                             `json:"failed_row"`
	ProcessedAt       string                          `json:"processed_at"`
	RowErrors         []BulkUploadRowError            `json:"row_errors,omitempty"` // baris Excel yang gagal + alasan
	SOData            []BulkUploadStockOpnameV2SOData `json:"so_data"`
}

// BulkUploadStockOpnameV2SOData stock opname header + product list in response
type BulkUploadStockOpnameV2SOData struct {
	DocNo         string                               `json:"doc_no"`
	CreatedDate   string                               `json:"created_date"`
	WhID          int64                                `json:"wh_id"`
	WhCode        string                               `json:"wh_code"`
	WhName        string                               `json:"wh_name"`
	StockType     string                               `json:"stock_type"`
	CreatedBy     int64                                `json:"created_by"`
	UserName      string                               `json:"user_name"`
	ScheduledDate string                               `json:"scheduled_date"`
	Status        int                                  `json:"status"`
	StatusDesc    string                               `json:"status_desc"`
	AssignTo      StockOpnameDetailV2AssignTo          `json:"assign_to"`
	ProductList   []BulkUploadStockOpnameV2ProductItem `json:"product_list"`
}

// BulkUploadStockOpnameV2ProductItem product item in bulk upload response
type BulkUploadStockOpnameV2ProductItem struct {
	StockOpnameDetailID int64   `json:"stock_opname_detail_id"`
	ProID               int64   `json:"pro_id"`
	ProCode             string  `json:"pro_code"`
	ProName             string  `json:"pro_name"`
	UnitID1             string  `json:"unit_id1"`
	UnitID2             string  `json:"unit_id2"`
	UnitID3             string  `json:"unit_id3"`
	UnitName1           string  `json:"unit_name1"`
	UnitName2           string  `json:"unit_name2"`
	UnitName3           string  `json:"unit_name3"`
	ConvUnit1           float64 `json:"conv_unit1"`
	ConvUnit2           float64 `json:"conv_unit2"`
	ConvUnit3           float64 `json:"conv_unit3"`
	Qty1                float64 `json:"qty1"`
	Qty2                float64 `json:"qty2"`
	Qty3                float64 `json:"qty3"`
	QtyPhysical1        float64 `json:"qty_physical1"`
	QtyPhysical2        float64 `json:"qty_physical2"`
	QtyPhysical3        float64 `json:"qty_physical3"`
	DifferentStock1     float64 `json:"different_stock1"`
	DifferentStock2     float64 `json:"different_stock2"`
	DifferentStock3     float64 `json:"different_stock3"`
	DifferentPrice1     float64 `json:"different_price1"`
	DifferentPrice2     float64 `json:"different_price2"`
	DifferentPrice3     float64 `json:"different_price3"`
	Status              string  `json:"status"` // "success" or "failed"
	ErrorMessage        *string `json:"error_message"`
}
