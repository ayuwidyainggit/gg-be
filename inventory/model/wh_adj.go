package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WhAdj struct {
	CustID    string         `gorm:"column:cust_id" json:"cust_id"`
	AdjNo     string         `gorm:"column:adj_no;primaryKey" json:"stock_adjustment_no"`
	AdjDate   time.Time      `gorm:"column:adj_date" json:"stock_adjustment_date"`
	WhID      int64          `gorm:"column:wh_id" json:"wh_id"`
	Notes     string         `gorm:"column:notes" json:"notes"`
	CreatedBy int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedBy int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel     bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed  bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy  int64          `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt  *time.Time     `gorm:"column:closed_at" json:"closed_at"`
}

func (WhAdj) TableName() string {
	return "inv.wh_adj"
}

type AdjNo struct {
	AdjNo string `gorm:"column:get_no_fn"`
}

func (m *WhAdj) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	var grNo GrNo
	trCode := "SA"
	grDateStr := m.AdjDate.Format("2006-01-02")
	grDateSubtr := grDateStr[2:4] + grDateStr[5:7] + grDateStr[8:10]
	// log.Println("grDateStr:", grDateStr)
	// log.Println("grDateSubtr:", grDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(adj_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn 
	FROM inv.wh_adj
	WHERE substr(adj_no,3,6) = '%v' AND cust_id = '%v'`, grDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&grNo).Error
	if err != nil {
		return err
	}

	m.AdjNo = trCode + grDateSubtr + grNo.GrNo
	log.Println("m.AdjNo:", m.AdjNo)

	m.CreatedAt = &now
	m.UpdatedAt = &now
	m.UpdatedBy = m.CreatedBy

	return nil
}

type WhAdjList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	AdjNo         string         `gorm:"column:adj_no;primaryKey" json:"stock_adjustment_no"`
	AdjDate       *time.Time     `gorm:"column:adj_date" json:"stock_adjustment_date"`
	WhID          int64          `gorm:"column:wh_id" json:"wh_id"`
	WhCode        string         `gorm:"column:wh_code" json:"wh_code"`
	WhName        string         `gorm:"column:wh_name" json:"wh_name"`
	StockType     string         `gorm:"column:stock_type" json:"stock_type"`
	Notes         string         `gorm:"column:notes" json:"notes"`
	DataStatus    int64          `gorm:"column:data_status" json:"status"`
	CreatedBy     int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      int64          `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (WhAdjList) TableName() string {
	return "inv.wh_adj"
}

type WarehouseAdjustment struct {
	WhId   *int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode *string `gorm:"column:wh_code" json:"wh_code"`
	WhName *string `gorm:"column:wh_name" json:"wh_name"`
}

func (WarehouseAdjustment) TableName() string {
	return "inv.wh_adj"
}
