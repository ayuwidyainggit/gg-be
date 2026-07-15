package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WhSo struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	WhSoNo     string         `gorm:"column:wh_so_no;primaryKey" json:"wh_so_no"`
	WhSoDate   *time.Time     `gorm:"column:wh_so_date" json:"wh_so_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID       *int64         `gorm:"column:wh_id" json:"wh_id"`
	SoType     *int64         `gorm:"column:so_type" json:"so_type"`
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

func (WhSo) TableName() string {
	return "inv.wh_so"
}

type WhSoNo struct {
	WhSoNo string `gorm:"column:get_no_fn"`
}

func (m *WhSo) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.WhSoNo = strconv.Itoa(int(intTmpsStr))
	// m.CreatedAt = time.Now()
	// m.UpdatedBy = m.CreatedBy

	var whSoNo WhSoNo
	trCode := *m.TrCode
	whSoDateStr := m.WhSoDate.Format("2006-01-02")
	whSoDateSubtr := whSoDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(wh_so_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.wh_so
	WHERE substr(wh_so_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, whSoDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&whSoNo).Error
	if err != nil {
		return err
	}

	m.WhSoNo = trCode + whSoDateSubtr + whSoNo.WhSoNo
	log.Println("m.WhSoNo:", m.WhSoNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

type WhSoList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	WhSoNo        string         `gorm:"column:wh_so_no;primaryKey" json:"wh_so_no"`
	WhSoDate      *time.Time     `gorm:"column:wh_so_date" json:"wh_so_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	SoType        *int64         `gorm:"column:so_type" json:"so_type"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (WhSoList) TableName() string {
	return "inv.wh_so"
}
