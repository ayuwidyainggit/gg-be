package model

import (
	"time"

	"gorm.io/gorm"
)

type StockDisposalDetail struct {
	CustID        string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	SdDetailID    int64      `gorm:"column:sd_detail_id;primaryKey;autoIncrement" json:"sd_detail_id"`
	SdID          int64      `gorm:"column:sd_id" json:"sd_id"`
	ProID         int64      `gorm:"column:pro_id" json:"pro_id"`
	FileName      *string    `gorm:"column:file_name" json:"file_name"`
	FileType      *string    `gorm:"column:file_type" json:"file_type"`
	MediaCategory *string    `gorm:"column:media_category" json:"media_category"`
	FileBase64    *string    `gorm:"column:file_base64" json:"file_base64"`
	FileUrl       *string    `gorm:"column:file_url" json:"file_url"`
	FileSize      *int64     `gorm:"column:file_size" json:"file_size"`
	UnitID1       string     `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2       string     `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3       string     `gorm:"column:unit_id3" json:"unit_id3"`
	Qty1          float64    `gorm:"column:qty1" json:"qty1"`
	Qty2          float64    `gorm:"column:qty2" json:"qty2"`
	Qty3          float64    `gorm:"column:qty3" json:"qty3"`
	PurchPrice1   float64    `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2   float64    `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3   float64    `gorm:"column:purch_price3" json:"purch_price3"`
	GrossPrice    float64    `gorm:"column:gross_price" json:"gross_price"`
	Vat           float64    `gorm:"column:vat" json:"vat"`
	VatValue      float64    `gorm:"column:vat_value" json:"vat_value"`
	SubTotal      float64    `gorm:"column:sub_total" json:"sub_total"`
	CreatedBy     int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy     *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	DeletedBy     *int64     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	IsDel         bool       `gorm:"column:is_del" json:"is_del"`
}

func (StockDisposalDetail) TableName() string {
	return "inv.stock_disposal_detail"
}

func (m *StockDisposalDetail) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	return nil
}

func (m *StockDisposalDetail) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = &now
	return nil
}

type StockDisposalDetailList struct {
	CustID        string     `gorm:"column:cust_id" json:"cust_id"`
	SdDetailID    int64      `gorm:"column:sd_detail_id" json:"sd_detail_id"`
	SdID          int64      `gorm:"column:sd_id" json:"sd_id"`
	ProID         int64      `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string     `gorm:"column:pro_code" json:"pro_code"`
	ProName       string     `gorm:"column:pro_name" json:"pro_name"`
	UnitID1       string     `gorm:"column:unit_id1" json:"unit_id1"`
	UnitID2       string     `gorm:"column:unit_id2" json:"unit_id2"`
	UnitID3       string     `gorm:"column:unit_id3" json:"unit_id3"`
	FileName      *string    `gorm:"column:file_name" json:"file_name"`
	FileType      *string    `gorm:"column:file_type" json:"file_type"`
	MediaCategory *string    `gorm:"column:media_category" json:"media_category"`
	FileBase64    *string    `gorm:"column:file_base64" json:"file_base64"`
	FileUrl       *string    `gorm:"column:file_url" json:"file_url"`
	FileSize      *int64     `gorm:"column:file_size" json:"file_size"`
	Qty1          float64    `gorm:"column:qty1" json:"qty1"`
	Qty2          float64    `gorm:"column:qty2" json:"qty2"`
	Qty3          float64    `gorm:"column:qty3" json:"qty3"`
	PurchPrice1   float64    `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2   float64    `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3   float64    `gorm:"column:purch_price3" json:"purch_price3"`
	GrossPrice    float64    `gorm:"column:gross_price" json:"gross_price"`
	Vat           float64    `gorm:"column:vat" json:"vat"`
	VatValue      float64    `gorm:"column:vat_value" json:"vat_value"`
	SubTotal      float64    `gorm:"column:sub_total" json:"sub_total"`
	CreatedBy     int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy     *int64     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	IsDel         bool       `gorm:"column:is_del" json:"is_del"`
}

func (StockDisposalDetailList) TableName() string {
	return "inv.stock_disposal_detail"
}
