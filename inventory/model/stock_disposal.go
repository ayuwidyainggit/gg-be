package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StockDisposal struct {
	CustID       string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	SdID         int64      `gorm:"column:sd_id;primaryKey;autoIncrement" json:"sd_id"`
	TrCode       string     `gorm:"column:tr_code" json:"tr_code"`
	DisposalDate *time.Time `gorm:"column:disposal_date" json:"disposal_date"`
	SdNumber     string     `gorm:"column:sd_number" json:"sd_number"`
	SupID        int64      `gorm:"column:sup_id" json:"sup_id"`
	WhID         int64      `gorm:"column:wh_id" json:"wh_id"`
	StockType    string     `gorm:"column:stock_type" json:"stock_type"`
	GrNo         *string    `gorm:"column:gr_no" json:"gr_no"`
	Note         string     `gorm:"column:note;not null" json:"note"`
	SubTotal     float64    `gorm:"column:sub_total" json:"sub_total"`
	VatValue     float64    `gorm:"column:vat_value" json:"vat_value"`
	Total        float64    `gorm:"column:total" json:"total"`
	CreatedBy    int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy    *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedBy    *int64     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	IsDel        bool       `gorm:"column:is_del" json:"is_del"`
}

func (StockDisposal) TableName() string {
	return "inv.stock_disposal"
}

func (m *StockDisposal) BeforeCreate(trx *gorm.DB) (err error) {
	var sdNumber SdNumber
	trCode := "SD"
	disposalDateStr := m.DisposalDate.Format("2006-01-02")
	disposalDateSubtr := disposalDateStr[2:4] + disposalDateStr[5:7] + disposalDateStr[8:10]

	queryStr := fmt.Sprintf(`SELECT
		TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(sd_number,9,3),'999')),0)+1, '000')) AS get_no_fn
		FROM inv.stock_disposal
		WHERE substr(sd_number,3,6) = '%v' AND cust_id = '%v'`, disposalDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&sdNumber).Error
	if err != nil {
		return err
	}

	log.Println("sdNumber:", sdNumber.SdNumber)
	m.SdNumber = trCode + disposalDateSubtr + sdNumber.SdNumber
	log.Println("m.SdNumber:", m.SdNumber)
	m.TrCode = trCode
	m.CreatedAt = time.Now()
	now := time.Now()
	m.UpdatedAt = &now
	if m.CreatedBy != 0 {
		m.UpdatedBy = &m.CreatedBy
	}
	return nil
}

func (m *StockDisposal) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = &now
	return nil
}

type SdNumber struct {
	SdNumber string `gorm:"column:get_no_fn"`
}

type StockDisposalList struct {
	CustID             string     `gorm:"column:cust_id" json:"cust_id"`
	SdID               int64      `gorm:"column:sd_id" json:"sd_id"`
	TrCode             string     `gorm:"column:tr_code" json:"tr_code"`
	DisposalDate       *time.Time `gorm:"column:disposal_date" json:"disposal_date"`
	SdNumber           string     `gorm:"column:sd_number" json:"sd_number"`
	SupID              int64      `gorm:"column:sup_id" json:"sup_id"`
	SupCode            *string    `gorm:"column:sup_code" json:"sup_code"`
	SupName            *string    `gorm:"column:sup_name" json:"sup_name"`
	WhID               int64      `gorm:"column:wh_id" json:"wh_id"`
	WhCode             *string    `gorm:"column:wh_code" json:"wh_code"`
	WhName             *string    `gorm:"column:wh_name" json:"wh_name"`
	StockType          string     `gorm:"column:stock_type" json:"stock_type"`
	GrNo               *string    `gorm:"column:gr_no" json:"gr_no"`
	Note               string     `gorm:"column:note" json:"note"`
	SubTotal           float64    `gorm:"column:sub_total" json:"sub_total"`
	VatValue           float64    `gorm:"column:vat_value" json:"vat_value"`
	Total              float64    `gorm:"column:total" json:"total"`
	CalculatedSubtotal float64    `gorm:"column:calculated_subtotal" json:"calculated_subtotal"`
	CreatedBy          int64      `gorm:"column:created_by" json:"created_by"`
	CreatedByName      *string    `gorm:"column:created_by_name" json:"created_by_name"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy          *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName      *string    `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt          *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	IsDel              bool       `gorm:"column:is_del" json:"is_del"`
}

func (StockDisposalList) TableName() string {
	return "inv.stock_disposal"
}

type StockDisposalProductLookup struct {
	ProID           int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode         string  `gorm:"column:pro_code" json:"pro_code"`
	ProName         string  `gorm:"column:pro_name" json:"pro_name"`
	Vat             float64 `gorm:"column:vat" json:"vat"`
	ConvUnit2       int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       int     `gorm:"column:conv_unit3" json:"conv_unit3"`
	UnitID1         string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2         string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3         string  `gorm:"column:unit_id3" json:"unit_id3"`
	PurchPrice1     float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2     float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3     float64 `gorm:"column:purch_price3" json:"purch_price3"`
	MinStockQty     float64 `gorm:"column:min_stock_qty" json:"min_stock_qty"`
	SafStockQty     float64 `gorm:"column:saf_stock_qty" json:"saf_stock_qty"`
	Qty1            float64 `gorm:"column:qty1" json:"qty1"`
	Qty2            float64 `gorm:"column:qty2" json:"qty2"`
	Qty3            float64 `gorm:"column:qty3" json:"qty3"`
	TotalQty        float64 `gorm:"column:total_qty" json:"total_qty"`
	InTransitStock1 float64 `gorm:"column:in_transit_stock1" json:"in_transit_stock1"`
	InTransitStock2 float64 `gorm:"column:in_transit_stock2" json:"in_transit_stock2"`
	InTransitStock3 float64 `gorm:"column:in_transit_stock3" json:"in_transit_stock3"`
}
