package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AccountPayablePayment struct {
	CustId                    string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayablePaymentNo   string         `gorm:"column:account_payable_payment_no" json:"account_payable_payment_no"`
	AccountPayablPaymenteDate *time.Time     `gorm:"column:account_payable_payment_date" json:"account_payable_payment_date"`
	SupId                     *int64         `gorm:"column:sup_id" json:"sup_id"`
	TotalDiscount             *float64       `gorm:"column:total_discount" json:"total_discount"`
	TotalPaymentBalance       *float64       `gorm:"column:total_payment_balance" json:"total_payment_balance"`
	TotalMaterai              *float64       `gorm:"column:total_materai" json:"total_materai"`
	TotalPayment              *float64       `gorm:"column:total_payment" json:"total_payment"`
	CreatedBy                 *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt                 time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy                 *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt                 time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy                 *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt                 gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AccountPayablePayment) TableName() string {
	return "acf.account_payable_payment"
}

type AccountPayablePaymentNo struct {
	AccountPayablePaymentNo string `gorm:"column:account_payable_payment_no_fn"`
}

func (m *AccountPayablePayment) BeforeCreate(trx *gorm.DB) (err error) {
	var AccountPayablePaymentNo AccountPayablePaymentNo
	trCode := "PY"

	ReturnDateStr := m.AccountPayablPaymenteDate.Format("2006-01-02")
	ReturnDateSubtr := ReturnDateStr[2:4] + ReturnDateStr[5:7] + ReturnDateStr[8:10]

	// log.Println("grDateStr:", grDateStr)
	// log.Println("grDateSubtr:", grDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(account_payable_payment_no,9,4),'9999')),0)+1, '0000')) AS account_payable_payment_no_fn 
	FROM acf.account_payable_payment
	WHERE substr(account_payable_payment_no,3,6) = '%v'  AND cust_id = '%v'`, ReturnDateSubtr, strings.ToUpper(m.CustId))
	// log.Println("QUERY ===>", queryStr)
	err = trx.Raw(queryStr).Scan(&AccountPayablePaymentNo).Error
	if err != nil {
		return err
	}

	// log.Println("grNo:", grNo.GrNo)

	m.AccountPayablePaymentNo = trCode + ReturnDateSubtr + AccountPayablePaymentNo.AccountPayablePaymentNo
	// log.Println("m.GrNo:", m.GrNo)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil

}

type AccountPayablePaymentList struct {
	CustId                    string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayablePaymentNo   *string        `gorm:"column:account_payable_payment_no" json:"account_payable_payment_no"`
	DocumentNo                *string        `gorm:"column:document_no" json:"document_no"`
	AccountPayablPaymenteDate *time.Time     `gorm:"column:account_payable_payment_date" json:"account_payable_payment_date"`
	SupId                     *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupName                   *string        `gorm:"column:sup_name" json:"sup_name"`
	SupCode                   *string        `gorm:"column:sup_code" json:"sup_code"`
	DistributorId             *int64         `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorName           *string        `gorm:"column:distributor" json:"distributor"`
	DistributorCode           *string        `gorm:"column:distributor_code" json:"distributor_code"`
	TotalDiscount             *float64       `gorm:"column:total_discount" json:"total_discount"`
	TotalPaymentBalance       *float64       `gorm:"column:total_payment_balance" json:"total_payment_balance"`
	TotalMaterai              *float64       `gorm:"column:total_materai" json:"total_materai"`
	TotalPayment              *float64       `gorm:"column:total_payment" json:"total_payment"`
	CreatedBy                 *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt                 time.Time      `gorm:"column:created_at" json:"created_at"`
	CreatedByName             *string        `gorm:"column:created_by_name" json:"created_by_name"`
	UpdatedBy                 *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt                 time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName             *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	DeletedBy                 *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt                 gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (AccountPayablePaymentList) TableName() string {
	return "acf.account_payable_payment"
}

type AccountPayablePaymentDetail struct {
	CustId                        string     `gorm:"column:cust_id" json:"cust_id"`
	AccountPayablePaymentDetailId int        `gorm:"column:account_payable_payment_detail_id;type:int8;default:nextval('acf.account_payable_payment_detail_id_seq'::regclass);not null" json:"account_payable_payment_detail_id"`
	AccountPayablePaymentNo       string     `gorm:"column:account_payable_payment_no" json:"account_payable_payment_no"`
	InvoiceNo                     string     `gorm:"column:invoice_no" json:"invoice_no"`
	InvDate                       *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvAmount                     float64    `gorm:"column:invoice_amount" json:"invoice_amount"`
	PaidAmount                    float64    `gorm:"column:paid_amount" json:"paid_amount"`
	RemainingAmount               float64    `gorm:"column:remaining_amount" json:"remaining_amount"`
	Discount                      float64    `gorm:"column:discount" json:"discount"`
	PaymentBalance                float64    `gorm:"column:payment_balance" json:"payment_balance"`
	Materai                       *float64   `gorm:"column:materai" json:"materai"`
	TotalPayment                  *float64   `gorm:"column:total_payment" json:"total_payment"`
}

func (AccountPayablePaymentDetail) TableName() string {
	return "acf.account_payable_payment_detail"
}

type AccountPayablePaymentDetailList struct {
	AccountPayablePaymentNo *string    `gorm:"column:account_payable_payment_no" json:"account_payable_payment_no"`
	InvoiceNo               string     `gorm:"column:invoice_no" json:"invoice_no"`
	InvDate                 *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvAmount               float64    `gorm:"column:invoice_amount" json:"invoice_amount"`
	PaidAmount              float64    `gorm:"column:paid_amount" json:"paid_amount"`
	RemainingAmount         float64    `gorm:"column:remaining_amount" json:"remaining_amount"`
	Discount                float64    `gorm:"column:discount" json:"discount"`
	PaymentBalance          float64    `gorm:"column:payment_balance" json:"payment_balance"`
	Materai                 *float64   `gorm:"column:materai" json:"materai"`
	TotalPayment            *float64   `gorm:"column:total_payment" json:"total_payment"`
}

func (AccountPayablePaymentDetailList) TableName() string {
	return "acf.account_payable_payment_detail"
}

type AccountPayablePaymentOptions struct {
	CustId                         string   `gorm:"column:cust_id" json:"cust_id"`
	AccountPayablePaymentOptionsId int      `gorm:"column:account_payable_payment_options_id;type:int8;default:nextval('acf.account_payable_payment_options_id_seq'::regclass);not null" json:"account_payable_payment_options_id"`
	AccountPayablePaymentNo        string   `gorm:"column:account_payable_payment_no" json:"account_payable_payment_no"`
	InvoiceNo                      string   `gorm:"column:invoice_no" json:"invoice_no"`
	PayType                        *int64   `gorm:"column:pay_type" json:"pay_type"`
	DocumentNo                     *string  `gorm:"column:document_no" json:"document_no"`
	Balance                        *float64 `gorm:"column:balance" json:"balance"`
	PaymentAmount                  *float64 `gorm:"column:payment_amount" json:"payment_amount"`
}

func (AccountPayablePaymentOptions) TableName() string {
	return "acf.account_payable_payment_options"
}

type AccountPayablePaymentOptionsList struct {
	AccountPayablePaymentNo *string    `gorm:"column:account_payable_payment_no" json:"account_payable_payment_no"`
	InvDate                 *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo               string     `gorm:"column:invoice_no" json:"invoice_no"`
	PayType                 *int64     `gorm:"column:pay_type" json:"pay_type"`
	DocumentNo              *string    `gorm:"column:document_no" json:"document_no"`
	Balance                 *float64   `gorm:"column:balance" json:"balance"`
	PaymentAmount           *float64   `gorm:"column:payment_amount" json:"payment_amount"`
	PaymentBalance          *float64   `gorm:"column:payment_balance" json:"payment_balance"`
	RemainingAmount         float64    `gorm:"column:remaining_amount" json:"remaining_amount"`
}

func (AccountPayablePaymentOptionsList) TableName() string {
	return "acf.account_payable_payment_options"
}

// type GrUpdate struct {
// 	InvoiceNo   *string    `gorm:"column:invoice_no" json:"invoice_no"`
// 	InvoiceDate *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
// 	IsAp        bool       `gorm:"column:is_ap" json:"is_ap"`
// }

// func (GrUpdate) TableName() string {
// 	return "inv.gr"
// }

type ApLookupSuppilerInvoiceReturnList struct {
	ID                 uint           `gorm:"column:account_payable_id;primaryKey" json:"account_payable_id"`
	CustId             string         `gorm:"column:cust_id" json:"cust_id"`
	AccountPayableDate *time.Time     `gorm:"column:account_payable_date" json:"account_payable_date"`
	ApType             string         `gorm:"column:ap_type" json:"ap_type"`
	SupId              *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupName            string         `gorm:"column:sup_name" json:"sup_name"`
	InvoiceNo          string         `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	InvoiceDate        *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	DocumentNo         string         `gorm:"column:document_no" json:"document_no"`
	TaxInvoiceDate     *time.Time     `gorm:"column:tax_invoice_date" json:"tax_invoice_date"`
	TaxInvoiceNo       *string        `gorm:"column:tax_invoice_no" json:"tax_invoice_no"`
	TaxReturnDate      *time.Time     `gorm:"column:tax_return_date" json:"tax_return_date"`
	TaxReturnNo        *string        `gorm:"column:tax_return_no" json:"tax_return_no"`
	DueDate            *time.Time     `gorm:"column:due_date" json:"due_date"`
	ReturnDate         *time.Time     `gorm:"column:return_date" json:"return_date"`
	Amount             float64        `gorm:"column:amount" json:"amount"`
	DiscountRp         *float64       `gorm:"column:discount_rp" json:"discount_rp"`
	DiscountPercent    *float64       `gorm:"column:discount_percent" json:"discount_percent"`
	SubTotal           *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat                *float64       `gorm:"column:vat" json:"vat"`
	VatValue           *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg              *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue         *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Materai            *float64       `gorm:"column:materai" json:"materai"`
	PaidAmount         float64        `gorm:"column:paid_amount" json:"paid_amount"`
	RemainingAmount    *float64       `gorm:"column:remaining_amount" json:"remaining_amount"`
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

func (ApLookupSuppilerInvoiceReturnList) TableName() string {
	return "acf.account_payable"
}
