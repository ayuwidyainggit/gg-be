package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WhTrf struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	WhTrfNo    string         `gorm:"column:wh_trf_no" json:"wh_trf_no"`
	WhTrfDate  *time.Time     `gorm:"column:wh_trf_date" json:"wh_trf_date"`
	WhIDFrom   *int64         `gorm:"column:wh_id_from" json:"wh_id_from"`
	WhIDTo     *int64         `gorm:"column:wh_id_to" json:"wh_id_to"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed   bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy   *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt   time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (WhTrf) TableName() string {
	return "inv.wh_trf"
}

type WhTrfNo struct {
	WhTrfNo string `gorm:"column:get_no_fn"`
}

func (m *WhTrf) BeforeCreate(trx *gorm.DB) (err error) {
	var whTrfNo WhTrfNo
	trCode := "ST"
	grDateStr := m.WhTrfDate.Format("2006-01-02")
	grDateSubtr := grDateStr[2:4] + grDateStr[5:7] + grDateStr[8:10]
	// log.Println("grDateStr:", grDateStr)
	// log.Println("grDateSubtr:", grDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(wh_trf_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn 
	FROM inv.wh_trf
	WHERE substr(wh_trf_no,3,6) = '%v' AND cust_id = '%v'`, grDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&whTrfNo).Error
	if err != nil {
		return err
	}

	m.WhTrfNo = trCode + grDateSubtr + whTrfNo.WhTrfNo
	log.Println("m.WhTrfNo:", m.WhTrfNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type WhTrfList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	WhTrfNo       string         `gorm:"column:wh_trf_no" json:"stock_trf_no"`
	WhTrfDate     *time.Time     `gorm:"column:wh_trf_date" json:"stock_trf_date"`
	WhIDFrom      int64          `gorm:"column:wh_id_from" json:"wh_id_from"`
	WhCodeFrom    string         `gorm:"column:wh_code_from" json:"wh_code_from"`
	WhNameFrom    string         `gorm:"column:wh_name_from" json:"wh_name_from"`
	WhIDTo        int64          `gorm:"column:wh_id_to" json:"wh_id_to"`
	WhCodeTo      string         `gorm:"column:wh_code_to" json:"wh_code_to"`
	WhNameTo      string         `gorm:"column:wh_name_to" json:"wh_name_to"`
	Notes         string         `gorm:"column:notes" json:"notes"`
	DataStatus    int64          `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (WhTrfList) TableName() string {
	return "inv.wh_trf"
}

type WarehouseStockTransfer struct {
	WhId   *int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode *string `gorm:"column:wh_code" json:"wh_code"`
	WhName *string `gorm:"column:wh_name" json:"wh_name"`
}

func (WarehouseStockTransfer) TableName() string {
	return "inv.wh_trf"
}
