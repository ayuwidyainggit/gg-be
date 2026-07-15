package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type VanSo struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	VanSoNo    string         `gorm:"column:van_so_no;primaryKey" json:"van_so_no"`
	VanSoDate  *time.Time     `gorm:"column:van_so_date" json:"van_so_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	EmpID      *int64         `gorm:"column:emp_id" json:"emp_id"`
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

func (VanSo) TableName() string {
	return "inv.van_so"
}

type VanSoNo struct {
	VanSoNo string `gorm:"column:get_no_fn"`
}

func (m *VanSo) BeforeCreate(trx *gorm.DB) (err error) {
	var vanSoNo VanSoNo
	trCode := *m.TrCode
	vanSoDateStr := m.VanSoDate.Format("2006-01-02")
	vanSoDateSubtr := vanSoDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(van_so_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.van_so
	WHERE substr(van_so_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, vanSoDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&vanSoNo).Error
	if err != nil {
		return err
	}

	m.VanSoNo = trCode + vanSoDateSubtr + vanSoNo.VanSoNo
	log.Println("m.VanSoNo:", m.VanSoNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type VanSoList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	VanSoNo       string         `gorm:"column:van_so_no;primaryKey" json:"van_so_no"`
	VanSoDate     *time.Time     `gorm:"column:van_so_date" json:"van_so_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	EmpID         *int64         `gorm:"column:emp_id" json:"emp_id"`
	EmpCode       *string        `gorm:"column:emp_code" json:"emp_code"`
	EmpName       *string        `gorm:"column:emp_name" json:"emp_name"`
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

func (VanSoList) TableName() string {
	return "inv.van_so"
}
