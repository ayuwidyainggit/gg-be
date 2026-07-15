package model

import (
	"database/sql"
	"time"
)

type SecondarySalesReport struct {
	CustID          string    `gorm:"column:cust_id" json:"cust_id"`
	DistributorCode string    `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName string    `gorm:"column:distributor_name" json:"distributor_name"`
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	OutletID        int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanID      int64     `gorm:"column:salesman_id" json:"salesman_id"`
	EmpCode         string    `gorm:"column:emp_code" json:"emp_code"`
	EmpName         string    `gorm:"column:emp_name" json:"emp_name"`
	ProID           int64     `gorm:"column:pro_id" json:"pro_id"`
	SupCode         string    `gorm:"column:sup_code" json:"sup_code"`
	SupName         string    `gorm:"column:sup_name" json:"sup_name"`
	ProCode         string    `gorm:"column:pro_code" json:"pro_code"`
	ProName         string    `gorm:"column:pro_name" json:"pro_name"`
	UnitID1         string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         string    `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2       float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty1Final       float64   `gorm:"column:qty1_final" json:"qty1_final"`
	Qty2Final       float64   `gorm:"column:qty2_final" json:"qty2_final"`
	Qty3Final       float64   `gorm:"column:qty3_final" json:"qty3_final"`
	GrossSales      float64   `gorm:"column:gross_sales" json:"gross_sales"`
	SpecialDiscount float64   `gorm:"column:special_discount" json:"special_discount"`
	Discount        float64   `gorm:"column:discount" json:"discount"`
	NetSalesExcPPN  float64   `gorm:"column:net_sales_exc_ppn" json:"net_sales_exc_ppn"`
	PPN             float64   `gorm:"column:ppn" json:"ppn"`
	NetSalesIncPPN  float64   `gorm:"column:net_sales_inc_ppn" json:"net_sales_inc_ppn"`
	Qty1Return      float64   `gorm:"column:qty1_return" json:"qty1_return"`
	Qty2Return      float64   `gorm:"column:qty2_return" json:"qty2_return"`
	Qty3Return      float64   `gorm:"column:qty3_return" json:"qty3_return"`
	SellPrice1      float64   `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64   `gorm:"column:sell_price3" json:"sell_price3"`
}

func (SecondarySalesReport) TableName() string {
	return "sls.order_detail"
}

type SecondarySalesReportUnion struct {
	DistributorCode string    `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName string    `gorm:"column:distributor_name" json:"distributor_name"`
	TrxType         string    `gorm:"column:trx_type" json:"trx_type"` // <-- tambahan
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo      string    `gorm:"column:document_no" json:"document_no"`
	DocumentDate    time.Time `gorm:"column:document_date" json:"document_date"`
	OutletID        int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanID      int64     `gorm:"column:salesman_id" json:"salesman_id"`
	EmpCode         string    `gorm:"column:emp_code" json:"emp_code"`
	EmpName         string    `gorm:"column:emp_name" json:"emp_name"`
	ProductID       int64     `gorm:"column:product_id" json:"product_id"` // <-- samakan dengan query
	SupCode         string    `gorm:"column:sup_code" json:"sup_code"`
	SupName         string    `gorm:"column:sup_name" json:"sup_name"`
	ProCode         string    `gorm:"column:pro_code" json:"pro_code"`
	ProName         string    `gorm:"column:pro_name" json:"pro_name"`
	UnitID1         string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         string    `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2       float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty1            float64   `gorm:"column:qty1" json:"qty1"`
	Qty2            float64   `gorm:"column:qty2" json:"qty2"`
	Qty3            float64   `gorm:"column:qty3" json:"qty3"`
	GrossSales      float64   `gorm:"column:gross_sales" json:"gross_sales"`
	SpecialDiscount float64   `gorm:"column:special_discount" json:"special_discount"`
	Discount        float64   `gorm:"column:discount" json:"discount"`
	NetSalesExcPPN  float64   `gorm:"column:net_sales_exc_ppn" json:"net_sales_exc_ppn"`
	PPN             float64   `gorm:"column:ppn" json:"ppn"`
	NetSalesIncPPN  float64   `gorm:"column:net_sales_inc_ppn" json:"net_sales_inc_ppn"`
	SellPrice1      float64   `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64   `gorm:"column:sell_price3" json:"sell_price3"`
}

type ReportList struct {
	CustID     string    `gorm:"column:cust_id" json:"cust_id"`
	ReportID   string    `gorm:"column:report_id" json:"report_id"`
	ReportName string    `gorm:"column:report_name" json:"report_name"`
	StartDate  time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate    time.Time `gorm:"column:end_date" json:"end_date"`
	FileStatus int       `gorm:"column:file_status" json:"file_status"`
	FileURL    string    `gorm:"column:file_url" json:"file_url"`
	FileBase64 string    `gorm:"column:file_base64" json:"file_base64"`
	CreatedBy  string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ReportList) TableName() string {
	return "report.list"
}

type SalesActivityReport struct {
	CustID          string        `gorm:"column:cust_id" json:"cust_id"`
	DistributorCode string        `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName string        `gorm:"column:distributor_name" json:"distributor_name"`
	PJPCode         string        `gorm:"column:pjp_code" json:"pjp_code"`
	EmpCode         string        `gorm:"column:emp_code" json:"emp_code"`
	SalesmanName    string        `gorm:"column:salesman_name" json:"salesman_name"`
	OutletID        int64         `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      string        `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      string        `gorm:"column:outlet_name" json:"outlet_name"`
	RODate          time.Time     `gorm:"column:ro_date" json:"ro_date"`
	ArriveAt        int64         `gorm:"column:arrive_at" json:"arrive_at"`
	LeaveAt         int64         `gorm:"column:leave_at" json:"leave_at"`
	TotalOrder      float64       `gorm:"column:total_order" json:"total_order"`
	TotalReturn     float64       `gorm:"column:total_return" json:"total_return"`
	TotalPayment    float64       `gorm:"column:total_payment" json:"total_payment"`
	Longitude       float64       `gorm:"column:longitude" json:"longitude"`
	Latitude        float64       `gorm:"column:latitude" json:"latitude"`
	ActualLongitude float64       `gorm:"column:actual_longitude" json:"actual_longitude"`
	ActualLatitude  float64       `gorm:"column:actual_latitude" json:"actual_latitude"`
	LocationStatus  sql.NullInt32 `gorm:"column:location_status" json:"location_status"`
}

func (SalesActivityReport) TableName() string {
	return "sls.order"
}

type SalesActivityReportRow struct {
	BusinessUnitCode    string        `gorm:"column:business_unit_code"`
	BusinessUnitName    string        `gorm:"column:business_unit_name"`
	DistributorCode     string        `gorm:"column:distributor_code"`
	DistributorName     string        `gorm:"column:distributor_name"`
	PJPCode             string        `gorm:"column:pjp_code"`
	EmpCode             string        `gorm:"column:emp_code"`
	SalesmanName        string        `gorm:"column:salesman_name"`
	OutletCode          string        `gorm:"column:outlet_code"`
	OutletPrincipalCode string        `gorm:"column:outlet_principal_code"`
	OutletName          string        `gorm:"column:outlet_name"`
	VisitDate           time.Time     `gorm:"column:visit_date"`
	ClockInTime         string        `gorm:"column:clock_in_time"`
	ClockOutTime        string        `gorm:"column:clock_out_time"`
	CheckinTime         string        `gorm:"column:checkin_time"`
	CheckoutTime        string        `gorm:"column:checkout_time"`
	DurationMinutes     int64         `gorm:"column:duration_in_minutes"`
	IsPlanned           bool          `gorm:"column:is_planned"`
	SkipAt              sql.NullInt64 `gorm:"column:skip_at"`
	LocationStatus      sql.NullInt32 `gorm:"column:location_status"`
	PjpStatus           string        `gorm:"column:pjp_status"`
	VisitStatus         string        `gorm:"column:visit_status"`
	Compliance          string        `gorm:"column:compliance"`
	SalesValue          float64       `gorm:"column:sales_value"`
	ReturnValue         float64       `gorm:"column:return_value"`
	PaymentValue        float64       `gorm:"column:payment_value"`
	LocationMaster      string        `gorm:"column:location_master"`
	LocationActual      string        `gorm:"column:location_actual"`
	GeotagStatusLabel   string        `gorm:"column:geotag_status"`
	Remarks             string        `gorm:"column:remarks"`
}

type DimProductCategory struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	Code string `gorm:"column:code" json:"code"`
	Name string `gorm:"column:name" json:"name"`
}

func (DimProductCategory) TableName() string {
	return "report.dim_product_categories"
}

type DimProduct struct {
	ID         int64   `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	CategoryID int64   `gorm:"column:category_id" json:"category_id"`
	Code       string  `gorm:"column:code" json:"code"`
	Name       string  `gorm:"column:name" json:"name"`
	UnitID1    string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2    string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3    string  `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2  float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
}

func (DimProduct) TableName() string {
	return "report.dim_products"
}

type DimOutlet struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	Code string `gorm:"column:code" json:"code"`
	Name string `gorm:"column:name" json:"name"`
}

func (DimOutlet) TableName() string {
	return "report.dim_outlets"
}

type DimSalesman struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	Code string `gorm:"column:code" json:"code"`
	Name string `gorm:"column:name" json:"name"`
}

func (DimSalesman) TableName() string {
	return "report.dim_salesmans"
}

type DimDate struct {
	ID    *int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Day   int    `gorm:"column:day" json:"day"`
	Month int    `gorm:"column:month" json:"month"`
	Year  int    `gorm:"column:year" json:"year"`
}

func (DimDate) TableName() string {
	return "report.dim_dates"
}

type FactOrder struct {
	CustID          string    `gorm:"column:cust_id" json:"cust_id"`
	RoNo            string    `gorm:"column:ro_no" json:"ro_no"`
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	DateID          int64     `gorm:"column:date_id" json:"date_id"`
	SalesmanID      int64     `gorm:"column:salesman_id" json:"salesman_id"`
	OutletID        int64     `gorm:"column:outlet_id" json:"outlet_id"`
	ProID           int64     `gorm:"column:pro_id" json:"pro_id"`
	Qty             float64   `gorm:"column:qty" json:"qty"`
	Qty1            float64   `gorm:"column:qty1" json:"qty1"`
	Qty2            float64   `gorm:"column:qty2" json:"qty2"`
	Qty3            float64   `gorm:"column:qty3" json:"qty3"`
	ItemType        int       `gorm:"column:item_type" json:"item_type"`
	GrossSales      float64   `gorm:"column:gross_sale" json:"gross_sale"`
	SpecialDiscount float64   `gorm:"column:special_discount" json:"special_discount"`
	Discount        float64   `gorm:"column:discount" json:"discount"`
	NetSalesExcPPN  float64   `gorm:"column:net_sales_exclude_ppn" json:"net_sales_exclude_ppn"`
	PPN             float64   `gorm:"column:ppn" json:"ppn"`
	NetSalesIncPPN  float64   `gorm:"column:net_sales_include_ppn" json:"net_sales_include_ppn"`
	SellPrice1      float64   `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64   `gorm:"column:sell_price3" json:"sell_price3"`
	ExtractedAt     time.Time `gorm:"column:extracted_at" json:"extracted_at"`
}

func (FactOrder) TableName() string {
	return "report.fact_orders"
}

type FactReturn struct {
	CustID          string    `gorm:"column:cust_id" json:"cust_id"`
	ReturnNo        string    `gorm:"column:return_no" json:"return_no"`
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	DateID          int64     `gorm:"column:date_id" json:"date_id"`
	SalesmanID      int64     `gorm:"column:salesman_id" json:"salesman_id"`
	OutletID        int64     `gorm:"column:outlet_id" json:"outlet_id"`
	ProID           int64     `gorm:"column:pro_id" json:"pro_id"`
	Qty             float64   `gorm:"column:qty" json:"qty"`
	Qty1            float64   `gorm:"column:qty1" json:"qty1"`
	Qty2            float64   `gorm:"column:qty2" json:"qty2"`
	Qty3            float64   `gorm:"column:qty3" json:"qty3"`
	ItemType        int       `gorm:"column:item_type" json:"item_type"`
	GrossSales      float64   `gorm:"column:gross_sale" json:"gross_sale"`
	SpecialDiscount float64   `gorm:"column:special_discount" json:"special_discount"`
	Discount        float64   `gorm:"column:discount" json:"discount"`
	NetSalesExcPPN  float64   `gorm:"column:net_sales_exclude_ppn" json:"net_sales_exclude_ppn"`
	PPN             float64   `gorm:"column:ppn" json:"ppn"`
	NetSalesIncPPN  float64   `gorm:"column:net_sales_include_ppn" json:"net_sales_include_ppn"`
	SellPrice1      float64   `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64   `gorm:"column:sell_price3" json:"sell_price3"`
	ExtractedAt     time.Time `gorm:"column:extracted_at" json:"extracted_at"`
}

func (FactReturn) TableName() string {
	return "report.fact_returns"
}

type SecondarySalesReportUnionReport struct {
	CustID          string    `gorm:"column:cust_id" json:"cust_id"`
	RoNo            string    `gorm:"column:ro_no" json:"ro_no"`
	RoDate          string    `gorm:"column:ro_date" json:"ro_date"`
	DistributorCode string    `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName string    `gorm:"column:distributor_name" json:"distributor_name"`
	TrxType         string    `gorm:"column:trx_type" json:"trx_type"` // <-- tambahan
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo      string    `gorm:"column:document_no" json:"document_no"`
	DocumentDate    time.Time `gorm:"column:document_date" json:"document_date"`
	OutletID        int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanID      int64     `gorm:"column:salesman_id" json:"salesman_id"`
	EmpCode         string    `gorm:"column:emp_code" json:"emp_code"`
	EmpName         string    `gorm:"column:emp_name" json:"emp_name"`
	ProductID       int64     `gorm:"column:product_id" json:"product_id"` // <-- samakan dengan query
	SupCode         string    `gorm:"column:sup_code" json:"sup_code"`
	SupName         string    `gorm:"column:sup_name" json:"sup_name"`
	ProCode         string    `gorm:"column:pro_code" json:"pro_code"`
	ProName         string    `gorm:"column:pro_name" json:"pro_name"`
	UnitID1         string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         string    `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2       float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty1            float64   `gorm:"column:qty1" json:"qty1"`
	Qty2            float64   `gorm:"column:qty2" json:"qty2"`
	Qty3            float64   `gorm:"column:qty3" json:"qty3"`
	Qty             float64   `gorm:"column:qty" json:"qty"`
	GrossSales      float64   `gorm:"column:gross_sales" json:"gross_sales"`
	SpecialDiscount float64   `gorm:"column:special_discount" json:"special_discount"`
	Discount        float64   `gorm:"column:discount" json:"discount"`
	NetSalesExcPPN  float64   `gorm:"column:net_sales_exc_ppn" json:"net_sales_exc_ppn"`
	PPN             float64   `gorm:"column:ppn" json:"ppn"`
	NetSalesIncPPN  float64   `gorm:"column:net_sales_inc_ppn" json:"net_sales_inc_ppn"`
	SellPrice1      float64   `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64   `gorm:"column:sell_price3" json:"sell_price3"`
	PcatID          int64     `gorm:"column:pcat_id" json:"pcat_id"`
	PcatCode        string    `gorm:"column:pcat_code" json:"pcat_code"`
	PcatName        string    `gorm:"column:pcat_name" json:"pcat_name"`
	ItemType        int       `gorm:"column:item_type" json:"item_type"`
}

func (SecondarySalesReportUnionReport) TableName() string {
	return "sls.order_detail"
}

type SecondarySalesReportUnionReturn struct {
	CustID          string    `gorm:"column:cust_id" json:"cust_id"`
	ReturnNo        string    `gorm:"column:return_no" json:"return_no"`
	ReturnDate      string    `gorm:"column:return_date" json:"return_date"`
	DistributorCode string    `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName string    `gorm:"column:distributor_name" json:"distributor_name"`
	TrxType         string    `gorm:"column:trx_type" json:"trx_type"` // <-- tambahan
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo      string    `gorm:"column:document_no" json:"document_no"`
	DocumentDate    time.Time `gorm:"column:document_date" json:"document_date"`
	OutletID        int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      string    `gorm:"column:outlet_name" json:"outlet_name"`
	SalesmanID      int64     `gorm:"column:salesman_id" json:"salesman_id"`
	EmpCode         string    `gorm:"column:emp_code" json:"emp_code"`
	EmpName         string    `gorm:"column:emp_name" json:"emp_name"`
	ProductID       int64     `gorm:"column:product_id" json:"product_id"` // <-- samakan dengan query
	SupCode         string    `gorm:"column:sup_code" json:"sup_code"`
	SupName         string    `gorm:"column:sup_name" json:"sup_name"`
	ProCode         string    `gorm:"column:pro_code" json:"pro_code"`
	ProName         string    `gorm:"column:pro_name" json:"pro_name"`
	UnitID1         string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         string    `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2       float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty1            float64   `gorm:"column:qty1" json:"qty1"`
	Qty2            float64   `gorm:"column:qty2" json:"qty2"`
	Qty3            float64   `gorm:"column:qty3" json:"qty3"`
	Qty             float64   `gorm:"column:qty" json:"qty"`
	GrossSales      float64   `gorm:"column:gross_sales" json:"gross_sales"`
	SpecialDiscount float64   `gorm:"column:special_discount" json:"special_discount"`
	Discount        float64   `gorm:"column:discount" json:"discount"`
	NetSalesExcPPN  float64   `gorm:"column:net_sales_exc_ppn" json:"net_sales_exc_ppn"`
	PPN             float64   `gorm:"column:ppn" json:"ppn"`
	NetSalesIncPPN  float64   `gorm:"column:net_sales_inc_ppn" json:"net_sales_inc_ppn"`
	SellPrice1      float64   `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64   `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64   `gorm:"column:sell_price3" json:"sell_price3"`
	PcatID          int64     `gorm:"column:pcat_id" json:"pcat_id"`
	PcatCode        string    `gorm:"column:pcat_code" json:"pcat_code"`
	PcatName        string    `gorm:"column:pcat_name" json:"pcat_name"`
	ItemType        int       `gorm:"column:item_type" json:"item_type"`
}

func (SecondarySalesReportUnionReturn) TableName() string {
	return "sls.return_detail"
}

type SecondarySalesReportOrderCustID struct {
	CustID string `gorm:"column:cust_id" json:"cust_id"`
}

func (SecondarySalesReportOrderCustID) TableName() string {
	return "sls.order"
}

type SecondarySalesReportReturnCustID struct {
	CustID string `gorm:"column:cust_id" json:"cust_id"`
}

func (SecondarySalesReportReturnCustID) TableName() string {
	return "sls.return"
}

type SumReportByMonthModel struct {
	TotalGrossSales    float64    `gorm:"column:total_gross_sale" json:"total_gross_sale"`
	TotalDiscountPromo float64    `gorm:"column:total_discount_promo" json:"total_discount_promo"`
	TotalPPN           float64    `gorm:"column:total_ppn" json:"total_ppn"`
	NetSalesExcPPN     float64    `gorm:"column:net_sales_exc_ppn" json:"net_sales_exc_ppn"`
	NetSales           float64    `gorm:"column:net_sales" json:"net_sales"`
	TotalSalesman      int32      `gorm:"column:total_salesman" json:"total_salesman"`
	TotalOutlet        int32      `gorm:"column:total_outlet" json:"total_outlet"`
	TotalProduct       int32      `gorm:"column:total_product" json:"total_product"`
	Qty                int64      `gorm:"column:qty" json:"qty"`
	QtyReturn          int64      `gorm:"column:qty_return" json:"qty_return"`
	ReturnRate         float64    `gorm:"column:return_rate" json:"return_rate"`
	NetSalesReturn     float64    `gorm:"column:net_sales_return" json:"net_sales_return"`
	LastUpdate         *time.Time `gorm:"column:last_update" json:"last_update"`
}

func (SumReportByMonthModel) TableName() string {
	return "report.fact_orders"
}

type SumReportReturnByMonthModel struct {
	Qty        int64      `gorm:"column:qty" json:"qty"`
	NetSales   float64    `gorm:"column:net_sales" json:"net_sales"`
	LastUpdate *time.Time `gorm:"column:last_update" json:"last_update"`
}

func (SumReportReturnByMonthModel) TableName() string {
	return "report.fact_returns"
}

type SalesmanActivitySumByMonthModel struct {
	TotalSales    float64    `gorm:"column:total_sales" json:"total_sales"`
	TotalReturn   float64    `gorm:"column:total_return" json:"total_return"`
	TotalSalesman int32      `gorm:"column:total_salesman" json:"total_salesman"`
	LastUpdate    *time.Time `gorm:"column:last_update" json:"last_update"`
}

type TrendSalesSecondarySalesModel struct {
	Month              int     `gorm:"column:month" json:"month"`
	TotalGrossSales    float64 `gorm:"column:total_gross_sale" json:"total_gross_sale"`
	TotalDiscountPromo float64 `gorm:"column:total_discount_promo" json:"total_discount_promo"`
	NetSales           float64 `gorm:"column:net_sales" json:"net_sales"`
}

func (TrendSalesSecondarySalesModel) TableName() string {
	return "report.fact_orders"
}

type ActivityReportTrendSalesModel struct {
	Month        int     `gorm:"column:month" json:"month"`
	TotalInvoice float64 `gorm:"column:total_invoice" json:"total_invoice"`
	TotalReturn  float64 `gorm:"column:total_return" json:"total_return"`
	NetSales     float64 `gorm:"column:net_sales" json:"net_sales"`
}

func (ActivityReportTrendSalesModel) TableName() string {
	return "report.fact_orders"
}

type SecondarySalesReportGroup struct {
	ID       int     `gorm:"column:id"`
	Code     string  `gorm:"column:code"`
	Name     string  `gorm:"column:name"`
	NetSales float64 `gorm:"column:net_sales"`
}

func (SecondarySalesReportGroup) TableName() string {
	return "report.fact_orders"
}

type ReturnReportGroup struct {
	ID       int     `gorm:"column:id"`
	Code     string  `gorm:"column:code"`
	Name     string  `gorm:"column:name"`
	NetSales float64 `gorm:"column:net_sales"`
}

func (ReturnReportGroup) TableName() string {
	return "report.fact_returns"
}

type ActivityReportGeotagRow struct {
	SalesmanCode       int64   `gorm:"column:salesman_code"`
	SalesmanName       string  `gorm:"column:salesman_name"`
	TotalVisit         int64   `gorm:"column:total_visit"`
	GeotagMatchCount   int64   `gorm:"column:geotag_match_count"`
	GeotagUnmatchCount int64   `gorm:"column:geotag_unmatch_count"`
	GeotagMatchPct     float64 `gorm:"column:geotag_match_pct"`
	GeotagUnmatchPct   float64 `gorm:"column:geotag_unmatch_pct"`
}

func (ActivityReportGeotagRow) TableName() string {
	return "pjp.outlet_visit_list"
}

type SalesActivityReportSalesmanList struct {
	SalesmanID   int64  `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode string `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName string `gorm:"column:salesman_name" json:"salesman_name"`
}

func (SalesActivityReportSalesmanList) TableName() string {
	return "sls.order"
}
