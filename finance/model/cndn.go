package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Cndn struct {
	CustID              string         `gorm:"column:cust_id" json:"cust_id"`
	CndnNo              string         `gorm:"column:cndn_no" json:"cndn_no"`
	CndnDate            *time.Time     `gorm:"column:cndn_date" json:"cndn_date"`
	OwnerId             int            `gorm:"column:owner_id" json:"owner_id"`
	OutletId            *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	CndnJenis           *string        `gorm:"column:cndn_jenis" json:"cndn_jenis"`
	CndnType            int64          `gorm:"column:cndn_type" json:"cndn_type"`
	Amount              *float64       `gorm:"column:amount" json:"amount"`
	UsedAmount          *float64       `gorm:"column:used_amount" json:"used_amount"`
	RemainingAmount     *float64       `gorm:"column:remaning_amount" json:"remaning_amount"`
	LastTransactionDate *time.Time     `gorm:"column:last_transaction_date" json:"last_transaction_date"`
	Notes               *string        `gorm:"column:notes" json:"notes"`
	CreatedBy           *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy           *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel               bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy           *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Cndn) TableName() string {
	return "acf.cndn"
}

type CndnNo struct {
	CndnNo string `gorm:"column:cndn_no_fn"`
}

func (m *Cndn) BeforeCreate(trx *gorm.DB) (err error) {
	var CndnNo CndnNo
	trCode := ""

	if *m.CndnJenis == "credit" {
		trCode = "CN"
	}
	if *m.CndnJenis == "debit" {
		trCode = "DN"
	}
	ReturnDateStr := m.CndnDate.Format("2006-01-02")
	ReturnDateSubtr := ReturnDateStr[2:4] + ReturnDateStr[5:7] + ReturnDateStr[8:10]

	// log.Println("grDateStr:", grDateStr)
	// log.Println("grDateSubtr:", grDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(cndn_no,9,4),'9999')),0)+1, '0000')) AS cndn_no_fn 
	FROM acf.cndn
	WHERE substr(cndn_no,3,6) = '%v' and cndn_jenis = '%v' AND cust_id = '%v'`, ReturnDateSubtr, *m.CndnJenis, strings.ToUpper(m.CustID))
	// log.Println("QUERY ===>", queryStr)
	err = trx.Raw(queryStr).Scan(&CndnNo).Error
	if err != nil {
		return err
	}

	// log.Println("grNo:", grNo.GrNo)

	m.CndnNo = trCode + ReturnDateSubtr + CndnNo.CndnNo
	// log.Println("m.GrNo:", m.GrNo)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type CndnGet struct {
	CustID              string         `gorm:"column:cust_id" json:"cust_id"`
	CndnNo              string         `gorm:"column:cndn_no" json:"cndn_no"`
	CndnDate            *time.Time     `gorm:"column:cndn_date" json:"cndn_date"`
	OwnerId             int            `gorm:"column:owner_id" json:"owner_id"`
	OutletId            *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode          *string        `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName          *string        `gorm:"column:outlet_name" json:"outlet_name"`
	CndnJenis           *string        `gorm:"column:cndn_jenis" json:"cndn_jenis"`
	CndnType            *string        `gorm:"column:cndn_type" json:"cndn_type"`
	Amount              *float64       `gorm:"column:amount" json:"amount"`
	UsedAmount          *float64       `gorm:"column:used_amount" json:"used_amount"`
	RemainingAmount     *float64       `gorm:"column:remaning_amount" json:"remaning_amount"`
	LastTransactionDate *time.Time     `gorm:"column:last_transaction_date" json:"last_transaction_date"`
	Notes               *string        `gorm:"column:notes" json:"notes"`
	CreatedBy           *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy           *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName       *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel               bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy           *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CndnGet) TableName() string {
	return "acf.cndn"
}

type CndnGetDetil struct {
	CustID              string         `gorm:"column:cust_id" json:"cust_id"`
	CndnNo              string         `gorm:"column:cndn_no" json:"cndn_no"`
	CndnDate            *time.Time     `gorm:"column:cndn_date" json:"cndn_date"`
	OwnerId             int            `gorm:"column:owner_id" json:"owner_id"`
	OutletId            *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode          *string        `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName          *string        `gorm:"column:outlet_name" json:"outlet_name"`
	CndnJenis           *string        `gorm:"column:cndn_jenis" json:"cndn_jenis"`
	CndnId              *int64         `gorm:"column:cndn_id" json:"cndn_id"`
	CndnCode            *string        `gorm:"column:cndn_code" json:"cndn_code"`
	CndnName            *string        `gorm:"column:cndn_name" json:"cndn_name"`
	Amount              *float64       `gorm:"column:amount" json:"amount"`
	UsedAmount          *float64       `gorm:"column:used_amount" json:"used_amount"`
	UsedAmountOutlet    float64        `gorm:"column:used_amount_outlet" json:"used_amount_outlet"`
	RemainingAmount     *float64       `gorm:"column:remaning_amount" json:"remaning_amount"`
	LastTransactionDate *time.Time     `gorm:"column:last_transaction_date" json:"last_transaction_date"`
	Notes               *string        `gorm:"column:notes" json:"notes"`
	CreatedBy           *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy           *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByName       *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel               bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy           *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CndnGetDetil) TableName() string {
	return "acf.cndn"
}
