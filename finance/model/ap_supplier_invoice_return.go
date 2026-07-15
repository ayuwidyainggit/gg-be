package model

import (
	"time"

	"gorm.io/gorm"
)

type ApSupplierInvoiceReturn struct {
	ID                 *uint          `gorm:"column:account_payable_id;primaryKey"`
	CustId             string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayableDate *time.Time     `gorm:"column:account_payable_date" json:"account_payable_date"`
	ApType             *string        `gorm:"column:ap_type" json:"ap_type"`
	SupId              *int64         `gorm:"column:sup_id" json:"sup_id"`
	InvoiceNo          string         `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo         *string        `gorm:"column:document_no" json:"document_no"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	TaxInvoiceNo       *string        `gorm:"column:tax_invoice_no" json:"tax_invoice_no"`
	DueDate            *time.Time     `gorm:"column:due_date" json:"due_date"`
	Amount             float64        `gorm:"column:amount" json:"amount"`
	DiscountRp         *float64       `gorm:"column:discount_rp" json:"discount_rp"`
	DiscountPercent    *float64       `gorm:"column:discount_percent" json:"discount_percent"`
	SubTotal           float64        `gorm:"column:sub_total" json:"sub_total"`
	Vat                float64        `gorm:"column:vat" json:"vat"`
	VatValue           float64        `gorm:"column:vat_value" json:"vat_value"`
	VatLg              float64        `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         float64        `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Materai            float64        `gorm:"column:materai" json:"materai"`
	Total              float64        `gorm:"column:total" json:"total"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
}

func (ApSupplierInvoiceReturn) TableName() string {
	return "acf.account_payable"
}

func (m *ApSupplierInvoiceReturn) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.UpdatedBy = m.CreatedBy
	trx.Raw("SELECT nextval('acf.account_payable_id_seq'::regclass) AS account_payable_id").Scan(&m.ID)

	return nil
}

type ApSupplierInvoiceReturnupdate struct {
	CustId             string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayableDate *time.Time     `gorm:"column:account_payable_date" json:"account_payable_date"`
	ApType             *string        `gorm:"column:ap_type" json:"ap_type"`
	SupId              *int64         `gorm:"column:sup_id" json:"sup_id"`
	InvoiceNo          string         `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo         *string        `gorm:"column:document_no" json:"document_no"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	TaxInvoiceNo       *string        `gorm:"column:tax_invoice_no" json:"tax_invoice_no"`
	DueDate            *time.Time     `gorm:"column:due_date" json:"due_date"`
	Amount             float64        `gorm:"column:amount" json:"amount"`
	DiscountRp         *float64       `gorm:"column:discount_rp" json:"discount_rp"`
	DiscountPercent    *float64       `gorm:"column:discount_percent" json:"discount_percent"`
	SubTotal           float64        `gorm:"column:sub_total" json:"sub_total"`
	Vat                float64        `gorm:"column:vat" json:"vat"`
	VatValue           float64        `gorm:"column:vat_value" json:"vat_value"`
	VatLg              float64        `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         float64        `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Materai            float64        `gorm:"column:materai" json:"materai"`
	Total              float64        `gorm:"column:total" json:"total"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
}

func (ApSupplierInvoiceReturnupdate) TableName() string {
	return "acf.account_payable"
}

type ApSuppilerInvoiceReturnList struct {
	ID                 uint           `gorm:"column:account_payable_id;primaryKey" json:"account_payable_id"`
	CustId             string         `gorm:"column:cust_id" json:"cust_id"`
	PoNo               string         `gorm:"column:po_no_doc" json:"po_no"`
	AccountPayableDate *time.Time     `gorm:"column:account_payable_date" json:"account_payable_date"`
	ApType             string         `gorm:"column:ap_type" json:"ap_type"`
	SupId              *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupName            string         `gorm:"column:sup_name" json:"sup_name"`
	SupCode            string         `gorm:"column:sup_code" json:"sup_code"`
	DistributorId      *int64         `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorName    string         `gorm:"column:distributor" json:"distributor"`
	DistributorCode    string         `gorm:"column:distributor_code" json:"distributor_code"`
	InvoiceNo          string         `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo         string         `gorm:"column:document_no" json:"document_no"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	TaxInvoiceNo       *string        `gorm:"column:tax_invoice_no" json:"tax_invoice_no"`
	TaxReturnDate      *time.Time     `gorm:"column:tax_return_date" json:"tax_return_date"`
	TaxReturnNo        *string        `gorm:"column:tax_return_no" json:"tax_return_no"`
	DueDate            *time.Time     `gorm:"column:due_date" json:"due_date"`
	ReturnDate         *time.Time     `gorm:"column:return_date" json:"return_date"`
	Amount             *float64       `gorm:"column:amount" json:"amount"`
	DiscountRp         *float64       `gorm:"column:discount_rp" json:"discount_rp"`
	DiscountPercent    *float64       `gorm:"column:discount_percent" json:"discount_percent"`
	SubTotal           *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat                *float64       `gorm:"column:vat" json:"vat"`
	VatValue           *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg              *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Materai            *float64       `gorm:"column:materai" json:"materai"`
	Total              *float64       `gorm:"column:total" json:"total"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	CreatedByName      *string        `gorm:"column:created_by_name" json:"created_by_name"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName      *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
}

func (ApSuppilerInvoiceReturnList) TableName() string {
	return "acf.account_payable"
}

type AccountPayableProduct struct {
	AccountPayableID  *uint   `gorm:"column:account_payable_id" json:"account_payable_id"`
	CustId            string  `gorm:"column:cust_id" json:"cust_id"`
	ProId             int64   `gorm:"column:pro_id" json:"pro_id"`
	UnitPrice1        float64 `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2        float64 `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3        float64 `gorm:"column:unit_price3" json:"unit_price3"`
	ConvUnit2         float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3         float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	SubTotal          float64 `gorm:"column:sub_total" json:"sub_total"`
	Disc              float64 `gorm:"column:disc" json:"disc"`
	DiscValue         float64 `gorm:"column:disc_value" json:"disc_value"`
	SubTotalBeforePpn float64 `gorm:"column:sub_total_before_ppn" json:"sub_total_before_ppn"`
	Vat               float64 `gorm:"column:vat" json:"vat"`
	VatValue          float64 `gorm:"column:vat_value" json:"vat_value"`
	Total             float64 `gorm:"column:total" json:"total"`
	VatLg             float64 `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue        float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Qty               int     `gorm:"column:qty" json:"qty"`
	Type              int     `gorm:"column:item_type" json:"item_type"`
}

func (AccountPayableProduct) TableName() string {
	return "acf.account_payable_detail"
}

// func (m *AccountPayableProduct) BeforeCreate(trx *gorm.DB) (err error) {
// 	trx.Raw("SELECT nextval('acf.account_payable_detail_id_seq'::regclass) AS account_payable_detail_id").Scan(&m.AccountPayableID)

// 	return nil
// }

type AccountPayableProductList struct {
	AccountPayableDetailID int64    `gorm:"column:account_payable_detail_id" json:"account_payable_detail_id"`
	InvoiceNo              string   `gorm:"column:invoice_no" json:"invoice_no"`
	ProId                  int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode                string   `gorm:"column:pro_code" json:"pro_code"`
	ProName                string   `gorm:"column:pro_name" json:"pro_name"`
	ItemType               int      `gorm:"column:item_type" json:"item_type"`
	Qty1                   float64  `gorm:"column:qty1" json:"qty1"`
	Qty2                   float64  `gorm:"column:qty2" json:"qty2"`
	Qty3                   float64  `gorm:"column:qty3" json:"qty3"`
	UnitID1                string   `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2                string   `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3                string   `gorm:"column:unit_id3" json:"unit_id3"`
	UnitPrice1             float64  `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2             float64  `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3             float64  `gorm:"column:unit_price3" json:"unit_price3"`
	ConvUnit2              float64  `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3              float64  `gorm:"column:conv_unit3" json:"conv_unit3"`
	SubTotal               *float64 `gorm:"column:sub_total" json:"sub_total"`
	Disc                   float64  `gorm:"column:disc" json:"disc"`
	DiscValue              float64  `gorm:"column:disc_value" json:"disc_value"`
	SubTotalBeforePpn      *float64 `gorm:"column:sub_total_before_ppn" json:"sub_total_before_ppn"`
	Vat                    float64  `gorm:"column:vat" json:"vat"`
	VatValue               *float64 `gorm:"column:vat_value" json:"vat_value"`
	Total                  float64  `gorm:"column:total" json:"total"`
	VatLg                  float64  `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue             *float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Qty                    float64  `gorm:"column:qty" json:"qty"`
	QtyRemaining           float64  `gorm:"column:qty_remaining" json:"qty_remaining"`
}

func (AccountPayableProductList) TableName() string {
	return "acf.account_payable_detail"
}

type AccountPayableProductPromo struct {
	CustId     string  `gorm:"column:cust_id" json:"cust_id"`
	InvoiceNo  string  `gorm:"column:invoice_no" json:"invoice_no"`
	ProId      int64   `gorm:"column:pro_id" json:"pro_id"`
	UnitPrice1 float64 `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2 float64 `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3 float64 `gorm:"column:unit_price3" json:"unit_price3"`
	ConvUnit2  float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty        int     `gorm:"column:qty" json:"qty"`
}

func (AccountPayableProductPromo) TableName() string {
	return "acf.account_payable_product_promo"
}

type AccountPayableProductPromoList struct {
	InvoiceNo    *string  `gorm:"column:invoice_no" json:"invoice_no"`
	ProId        *int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode      *string  `gorm:"column:pro_code" json:"pro_code"`
	ProName      *string  `gorm:"column:pro_name" json:"pro_name"`
	Qty1         *float64 `gorm:"column:qty1" json:"qty1"`
	Qty2         *float64 `gorm:"column:qty2" json:"qty2"`
	Qty3         *float64 `gorm:"column:qty3" json:"qty3"`
	UnitPrice1   *float64 `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2   *float64 `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3   *float64 `gorm:"column:unit_price3" json:"unit_price3"`
	ConvUnit2    float64  `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3    float64  `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty          *float64 `gorm:"column:qty" json:"qty"`
	QtyRemaining float64  `gorm:"column:qty_remaining" json:"qty_remaining"`
}

func (AccountPayableProductPromoList) TableName() string {
	return "acf.account_payable_product_promo"
}

type GrUpdate struct {
	InvoiceNo   *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	IsAp        bool       `gorm:"column:is_ap" json:"is_ap"`
}

func (GrUpdate) TableName() string {
	return "inv.gr"
}

type GrList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	GrNo          string         `gorm:"column:gr_no" json:"gr_no"`
	GrDate        *time.Time     `gorm:"column:gr_date" json:"gr_date"`
	DeliveryDate  *time.Time     `gorm:"column:delivery_date" json:"delivery_date"`
	DeliveryNo    *string        `gorm:"column:delivery_no" json:"delivery_no"`
	InvoiceNo     *string        `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	VehicleNo     *string        `gorm:"column:vehicle_no" json:"vehicle_no"`
	PoNo          *string        `gorm:"column:po_no" json:"po_no"`
	PoDnNo        *string        `gorm:"column:po_dn_no" json:"po_dn_no"`
	SupId         *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string        `gorm:"column:sup_name" json:"sup_name"`
	WhId          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;default:null" json:"updated_at,omitempty"`
	IsDel         *bool          `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (GrList) TableName() string {
	return "inv.gr"
}

type GrbList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	GrNo          string         `gorm:"column:gr_branch_no" json:"gr_branch_no"`
	GrDate        *time.Time     `gorm:"column:gr_branch_date" json:"gr_branch_date"`
	DeliveryDate  *time.Time     `gorm:"column:delivery_date" json:"delivery_date"`
	DeliveryNo    *string        `gorm:"column:delivery_no" json:"delivery_no"`
	InvoiceNo     *string        `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	VehicleNo     *string        `gorm:"column:vehicle_no" json:"vehicle_no"`
	PoNo          *string        `gorm:"column:po_no" json:"po_no"`
	PoDnNo        *string        `gorm:"column:po_dn_no" json:"po_dn_no"`
	SupId         *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string        `gorm:"column:sup_name" json:"sup_name"`
	WhId          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;default:null" json:"updated_at,omitempty"`
	IsDel         *bool          `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (GrbList) TableName() string {
	return "inv.gr_branch"
}

type GrDetList struct {
	ID           int        `gorm:"column:gr_det_id;primaryKey" json:"gr_det_id"`
	CustID       string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrNo         string     `gorm:"column:gr_no;primaryKey" json:"gr_no"`
	SeqNo        int        `gorm:"column:seq_no" json:"seq_no"`
	ProID        int64      `gorm:"column:pro_id" json:"pro_id"`
	ProCode      string     `gorm:"column:pro_code" json:"pro_code"`
	ProName      string     `gorm:"column:pro_name" json:"pro_name"`
	ConvUnit2    float64    `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3    float64    `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4    float64    `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5    float64    `gorm:"column:conv_unit5" json:"conv_unit5"`
	ItemType     int        `gorm:"column:item_type" json:"item_type"`
	Qty          float64    `gorm:"column:qty" json:"qty"`
	QtyStr       *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitId1      *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2      *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3      *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4      *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5      *string    `gorm:"column:unit_id5" json:"unit_id5"`
	UnitPrice1   float64    `gorm:"column:unit_price1" json:"purch_price1"`
	UnitPrice2   float64    `gorm:"column:unit_price2" json:"purch_price2"`
	UnitPrice3   float64    `gorm:"column:unit_price3" json:"purch_price3"`
	EmbInc       *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc       *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo    *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo      *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate      *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat          float64    `gorm:"column:vat" json:"vat"`
	VatBg        float64    `gorm:"column:vat_bg" json:"vat_bg"`
	VatLgPurch   float64    `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate   *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax    *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	DiscP        *float64   `gorm:"column:disc_p" json:"disc_p"`
	QtyRemaining float64    `gorm:"column:qty_remaining" json:"qty_remaining"`
	QtyShip1     *float64   `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2     *float64   `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3     *float64   `gorm:"column:qty_ship3" json:"qty_ship3"`
	QtyShip      *float64   `gorm:"column:qty_ship" json:"qty_ship"`
	WhQty        float64    `gorm:"column:wh_qty" json:"wh_qty"`
	Discount     float64    `gorm:"column:discount" json:"discount"`
}

func (GrDetList) TableName() string {
	return "inv.gr_det"
}

type GrbDetList struct {
	ID           int        `gorm:"column:gr_det_id;primaryKey" json:"gr_det_id"`
	CustID       string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrNo         string     `gorm:"column:gr_branch_no;primaryKey" json:"gr_branch_no"`
	SeqNo        int        `gorm:"column:seq_no" json:"seq_no"`
	ProID        int64      `gorm:"column:pro_id" json:"pro_id"`
	ProCode      string     `gorm:"column:pro_code" json:"pro_code"`
	ProName      string     `gorm:"column:pro_name" json:"pro_name"`
	ConvUnit2    float64    `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3    float64    `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4    float64    `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5    float64    `gorm:"column:conv_unit5" json:"conv_unit5"`
	ItemType     int        `gorm:"column:item_type" json:"item_type"`
	Qty          float64    `gorm:"column:qty_received" json:"qty"`
	QtyStr       *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitId1      *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2      *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3      *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4      *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5      *string    `gorm:"column:unit_id5" json:"unit_id5"`
	UnitPrice1   float64    `gorm:"column:unit_price1" json:"purch_price1"`
	UnitPrice2   float64    `gorm:"column:unit_price2" json:"purch_price2"`
	UnitPrice3   float64    `gorm:"column:unit_price3" json:"purch_price3"`
	EmbInc       *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc       *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo    *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo      *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate      *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat          float64    `gorm:"column:vat" json:"vat"`
	VatBg        float64    `gorm:"column:vat_bg" json:"vat_bg"`
	VatLgPurch   float64    `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate   *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax    *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	DiscP        *float64   `gorm:"column:disc_p" json:"disc_p"`
	QtyRemaining float64    `gorm:"column:qty_remaining" json:"qty_remaining"`
	QtyShip1     *float64   `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2     *float64   `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3     *float64   `gorm:"column:qty_ship3" json:"qty_ship3"`
	QtyShip      *float64   `gorm:"column:qty_ship" json:"qty_ship"`
	WhQty        float64    `gorm:"column:wh_qty" json:"wh_qty"`
	Discount     float64    `gorm:"column:discount" json:"discount"`
}

func (GrbDetList) TableName() string {
	return "inv.gr_branch_det"
}

type WarehouseStock struct {
	ProID int64   `gorm:"column:pro_id" json:"pro_id"`
	Qty   float64 `gorm:"column:qty" json:"qty"`
}

func (WarehouseStock) TableName() string {
	return "inv.gr"
}

type WarehouseStockGrb struct {
	ProID int64   `gorm:"column:pro_id" json:"pro_id"`
	Qty   float64 `gorm:"column:qty" json:"qty"`
}

func (WarehouseStockGrb) TableName() string {
	return "inv.gr_branch"
}

type WarehouseStockFromReturn struct {
	ProID int64   `gorm:"column:pro_id" json:"pro_id"`
	Qty   float64 `gorm:"column:qty" json:"qty"`
}

func (WarehouseStockFromReturn) TableName() string {
	return "inv.supplier_returns"
}

type ApSuppilerInvoiceReturnVatInList struct {
	ID                 uint           `gorm:"column:account_payable_id;primaryKey" json:"account_payable_id"`
	CustId             string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayableDate *time.Time     `gorm:"column:account_payable_date" json:"account_payable_date"`
	ApType             string         `gorm:"column:ap_type" json:"ap_type"`
	SupId              *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupName            string         `gorm:"column:sup_name" json:"sup_name"`
	SupCode            string         `gorm:"column:sup_code" json:"sup_code"`
	Address            string         `gorm:"column:address" json:"address"`
	InvoiceNo          string         `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo         string         `gorm:"column:document_no" json:"document_no"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	TaxInvoiceNo       *string        `gorm:"column:tax_invoice_no" json:"tax_invoice_no"`
	TaxReturnDate      *time.Time     `gorm:"column:tax_return_date" json:"tax_return_date"`
	TaxReturnNo        *string        `gorm:"column:tax_return_no" json:"tax_return_no"`
	DueDate            *time.Time     `gorm:"column:due_date" json:"due_date"`
	ReturnDate         *time.Time     `gorm:"column:return_date" json:"return_date"`
	Amount             *float64       `gorm:"column:amount" json:"amount"`
	DiscountRp         *float64       `gorm:"column:discount_rp" json:"discount_rp"`
	DiscountPercent    *float64       `gorm:"column:discount_percent" json:"discount_percent"`
	SubTotal           *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat                *float64       `gorm:"column:vat" json:"vat"`
	VatValue           *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg              *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Materai            *float64       `gorm:"column:materai" json:"materai"`
	Total              *float64       `gorm:"column:total" json:"total"`
	CreatedBy          *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	CreatedByName      *string        `gorm:"column:created_by_name" json:"created_by_name"`
	UpdatedBy          *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName      *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy          *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsDel              bool           `gorm:"column:is_del" json:"is_del"`
	ExtractStatus      string         `gorm:"column:extract_status" json:"extract_status"`
	ExtractedAt        *time.Time     `gorm:"column:extracted_at" json:"extracted_at"`
	Npwp               string         `gorm:"column:npwp" json:"npwp"`
}

func (ApSuppilerInvoiceReturnVatInList) TableName() string {
	return "acf.account_payable"
}
