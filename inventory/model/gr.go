package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	ITEM_TYPE_NORMAL = 1
	ITEM_TYPE_PROMO  = 2
)

type Gr struct {
	CustID          string         `gorm:"column:cust_id" json:"cust_id"`
	GrNo            string         `gorm:"column:gr_no" json:"gr_no"`
	GrDate          *time.Time     `gorm:"column:gr_date" json:"gr_date"`
	DeliveryDate    *time.Time     `gorm:"column:delivery_date" json:"delivery_date"`
	DeliveryNo      *string        `gorm:"column:delivery_no" json:"delivery_no"`
	InvoiceNo       *string        `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	VehicleNo       *string        `gorm:"column:vehicle_no" json:"vehicle_no"`
	PoNo            *string        `gorm:"column:po_no" json:"po_no"`
	PoDnNo          *string        `gorm:"column:po_dn_no" json:"po_dn_no"`
	SupID           *int64         `gorm:"column:sup_id" json:"sup_id"`
	WhID            *int64         `gorm:"column:wh_id" json:"wh_id"`
	Notes           *string        `gorm:"column:notes" json:"notes"`
	GoodReceiptType *string        `gorm:"column:good_receipt_type" json:"good_receipt_type"`
	WithReference   *bool          `gorm:"column:with_reference" json:"with_reference"`
	DeliveryFee     *float64       `gorm:"column:delivery_fee" json:"delivery_fee"`
	SoNo            *string        `gorm:"column:so_no" json:"so_no"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at,omitempty"`
	IsDel           bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy       *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed        bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy        *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt        time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (Gr) TableName() string {
	return "inv.gr"
}

func (m *Gr) BeforeCreate(trx *gorm.DB) (err error) {
	var grNo GrNo
	trCode := "GR"
	grDateStr := m.GrDate.Format("2006-01-02")
	grDateSubtr := grDateStr[2:4] + grDateStr[5:7] + grDateStr[8:10]
	// log.Println("grDateStr:", grDateStr)
	// log.Println("grDateSubtr:", grDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(gr_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn 
	FROM inv.gr
	WHERE substr(gr_no,3,6) = '%v' AND cust_id = '%v'`, grDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&grNo).Error
	if err != nil {
		return err
	}

	// log.Println("grNo:", grNo.GrNo)

	m.GrNo = trCode + grDateSubtr + grNo.GrNo
	// log.Println("m.GrNo:", m.GrNo)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (m *Gr) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now

	return nil
}

type GrNo struct {
	GrNo string `gorm:"column:get_no_fn"`
}

type GrList struct {
	CustID          string         `gorm:"column:cust_id" json:"cust_id"`
	GrNo            string         `gorm:"column:gr_no" json:"gr_no"`
	GrDate          *time.Time     `gorm:"column:gr_date" json:"gr_date"`
	DeliveryDate    *time.Time     `gorm:"column:delivery_date" json:"delivery_date"`
	DeliveryNo      *string        `gorm:"column:delivery_no" json:"delivery_no"`
	InvoiceNo       *string        `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate     *time.Time     `gorm:"column:invoice_date" json:"invoice_date"`
	VehicleNo       *string        `gorm:"column:vehicle_no" json:"vehicle_no"`
	PoNo            *string        `gorm:"column:po_no" json:"po_no"`
	PoDnNo          *string        `gorm:"column:po_dn_no" json:"po_dn_no"`
	SupId           *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode         *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName         *string        `gorm:"column:sup_name" json:"sup_name"`
	WhId            *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode          *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName          *string        `gorm:"column:wh_name" json:"wh_name"`
	Notes           *string        `gorm:"column:notes" json:"notes"`
	GoodReceiptType *string        `gorm:"column:good_receipt_type" json:"good_receipt_type"`
	WithReference   *bool          `gorm:"column:with_reference" json:"with_reference"`
	DeliveryFee     *float64       `gorm:"column:delivery_fee" json:"delivery_fee"`
	SoNo            *string        `gorm:"column:so_no" json:"so_no"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName   *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;default:null" json:"updated_at,omitempty"`
	IsDel           *bool          `gorm:"column:is_del" json:"is_del"`
	DeletedBy       *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed        bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy        *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt        time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (GrList) TableName() string {
	return "inv.gr"
}

type GrSupplier struct {
	SupId   *int64  `gorm:"column:sup_id" json:"sup_id"`
	SupCode *string `gorm:"column:sup_code" json:"sup_code"`
	SupName *string `gorm:"column:sup_name" json:"sup_name"`
}

func (GrSupplier) TableName() string {
	return "inv.gr"
}

type GrWarehouse struct {
	WhID   *int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode *string `gorm:"column:wh_code" json:"wh_code"`
	WhName *string `gorm:"column:wh_name" json:"wh_name"`
}

func (GrWarehouse) TableName() string {
	return "inv.gr"
}

type DistributorGr struct {
	CustID          *string `gorm:"column:cust_id" json:"cust_id"`
	DistributorId   *int64  `gorm:"column:distributor_id" json:"distributor_id"`
	DistributorCode *string `gorm:"column:distributor_code" json:"distributor_code"`
	DistributorName *string `gorm:"column:distributor_name" json:"distributor_name"`
}

func (DistributorGr) TableName() string {
	return "mst.m_distributor"
}

type GrLookup struct {
	CustID string     `gorm:"column:cust_id" json:"cust_id"`
	GrNo   string     `gorm:"column:gr_no" json:"gr_no"`
	GrDate *time.Time `gorm:"column:gr_date" json:"gr_date"`
}

func (GrLookup) TableName() string {
	return "inv.gr"
}

type GrBranchLookup struct {
	CustID       string     `gorm:"column:cust_id" json:"cust_id"`
	GrBranchNo   string     `gorm:"column:gr_branch_no" json:"gr_branch_no"`
	GrBranchDate *time.Time `gorm:"column:gr_branch_date" json:"gr_branch_date"`
}

func (GrBranchLookup) TableName() string {
	return "inv.gr_branch"
}
