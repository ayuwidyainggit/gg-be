package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StockOpname struct {
	CustID             string    `gorm:"column:cust_id" json:"cust_id"`
	DocNo              string    `gorm:"column:doc_no" json:"doc_no"`
	WhID               int64     `gorm:"column:wh_id" json:"wh_id"`
	Notes              string    `gorm:"column:notes" json:"notes"`
	DataStatus         int       `gorm:"column:data_status" json:"data_status"`
	CreatedBy          *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy          *int64    `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	ScheduledAt        time.Time `gorm:"column:scheduled_at" json:"scheduled_at,omitempty"`
	AssignToEmpID      int64     `gorm:"column:assign_to_emp_id" json:"assign_to_emp_id"`
	ProductHierarchy   int       `gorm:"column:product_hierarchy" json:"product_hierarchy"`
	Include_ZeroStock  bool      `gorm:"column:include_zero_stock" json:"include_zero_stock"`
	IsShowCurrentStock bool      `gorm:"column:is_show_current_stock" json:"is_show_current_stock"`
}

func (StockOpname) TableName() string {
	return "inv.stock_opname"
}

func (m *StockOpname) BeforeCreate(trx *gorm.DB) (err error) {
	var docNo DocNo

	dateNow := time.Now()
	trCode := "SE"
	stockOpnameDateStr := dateNow.Format("2006-01-02")
	stockOpnameDateSubtr := stockOpnameDateStr[2:4] + stockOpnameDateStr[5:7] + stockOpnameDateStr[8:10]
	// log.Println("stockOpnameDateStr:", stockOpnameDateStr)
	// log.Println("stockOpnameDateSubtr:", stockOpnameDateSubtr)

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(doc_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn 
	FROM inv.stock_opname
	WHERE substr(doc_no,3,6) = '%v' AND cust_id = '%v'`, stockOpnameDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&docNo).Error
	if err != nil {
		return err
	}

	log.Println("docNo:", docNo.DocNo)
	m.DocNo = trCode + stockOpnameDateSubtr + docNo.DocNo
	log.Println("m.DocNo:", m.DocNo)
	m.CreatedAt = time.Now()
	return nil
}

type DocNo struct {
	DocNo string `gorm:"column:get_no_fn"`
}

type StockOpnameList struct {
	CustID             string    `gorm:"column:cust_id" json:"cust_id"`
	DocNo              string    `gorm:"column:doc_no" json:"doc_no"`
	WhID               int64     `gorm:"column:wh_id" json:"wh_id"`
	WhCode             string    `gorm:"column:wh_code" json:"wh_code"`
	WhName             string    `gorm:"column:wh_name" json:"wh_name"`
	EmpCode            string    `gorm:"column:emp_code" json:"emp_code"`
	EmpName            string    `gorm:"column:emp_name" json:"emp_name"`
	StockType          string    `gorm:"column:stock_type" json:"stock_type"`
	Notes              string    `gorm:"column:notes" json:"notes"`
	DataStatus         int       `gorm:"column:data_status" json:"data_status"`
	CreatedBy          *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedBy          *int64    `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
	ScheduledAt        time.Time `gorm:"column:scheduled_at" json:"scheduled_at,omitempty"`
	AssignToEmpID      int64     `gorm:"column:assign_to_emp_id" json:"assign_to_emp_id"`
	ProductHierarchy   int       `gorm:"column:product_hierarchy" json:"product_hierarchy"`
	Include_ZeroStock  bool      `gorm:"column:include_zero_stock" json:"include_zero_stock"`
	IsShowCurrentStock bool      `gorm:"column:is_show_current_stock" json:"is_show_current_stock"`
}

func (StockOpnameList) TableName() string {
	return "inv.stock_opname"
}

type StockOpnameWarehouse struct {
	WhID   int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode string `gorm:"column:wh_code" json:"wh_code"`
	WhName string `gorm:"column:wh_name" json:"wh_name"`
}

func (StockOpnameWarehouse) TableName() string {
	return "mst.m_warehouse"
}

type StockOpnameEmployee struct {
	EmpID   int64  `gorm:"column:emp_id" json:"emp_id"`
	EmpCode string `gorm:"column:emp_code" json:"emp_code"`
	EmpName string `gorm:"column:emp_name" json:"emp_name"`
}

func (StockOpnameEmployee) TableName() string {
	return "mst.m_employee"
}

type StockOpnameProductList struct {
	ProID         int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string  `gorm:"column:pro_code" json:"pro_code"`
	ProName       string  `gorm:"column:pro_name" json:"pro_name"`
	Uom1          string  `gorm:"column:uom1" json:"uom1"`
	Uom2          string  `gorm:"column:uom2" json:"uom2"`
	Uom3          string  `gorm:"column:uom3" json:"uom3"`
	UnitName1     string  `gorm:"column:unit_name1" json:"unit_name1"`
	UnitName2     string  `gorm:"column:unit_name2" json:"unit_name2"`
	UnitName3     string  `gorm:"column:unit_name3" json:"unit_name3"`
	ConvUnit2     float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3     float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	Qty           float64 `gorm:"column:qty" json:"qty"`
	WhID          int64   `gorm:"column:wh_id" json:"wh_id"`
	StockType     string  `gorm:"column:stock_type" json:"stock_type"`
	PrincipalID   int64   `gorm:"column:principal_id" json:"principal_id"`
	PrincipalName string  `gorm:"column:principal_name" json:"principal_name"`
	PLID          int64   `gorm:"column:pl_id" json:"pl_id"`
	PLName        string  `gorm:"column:pl_name" json:"pl_name"`
	BrandID       int64   `gorm:"column:brand_id" json:"brand_id"`
	BrandName     string  `gorm:"column:brand_name" json:"brand_name"`
	SBrand1ID     int64   `gorm:"column:sbrand1_id" json:"sbrand1_id"`
	SBrand1Name   string  `gorm:"column:sbrand1_name" json:"sbrand1_name"`
	IsActive      bool    `gorm:"column:is_active" json:"is_active"`
}

func (StockOpnameProductList) TableName() string {
	return "inv.warehouse_stock"
}

// StockOpnameListV2 model for v2 list endpoint
type StockOpnameListV2 struct {
	DocNo         string `gorm:"column:doc_no" json:"doc_no"`
	CreatedDate   string `gorm:"column:created_date" json:"created_date"`
	WhID          int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode        string `gorm:"column:wh_code" json:"wh_code"`
	WhName        string `gorm:"column:wh_name" json:"wh_name"`
	CreatedBy     string `gorm:"column:created_by" json:"created_by"`
	UserName      string `gorm:"column:user_name" json:"user_name"`
	ScheduledDate string `gorm:"column:scheduled_date" json:"scheduled_date"`
	EmpID         int64  `gorm:"column:emp_id" json:"emp_id"`
	EmpName       string `gorm:"column:emp_name" json:"emp_name"`
	Status        int    `gorm:"column:status" json:"status"`
}

func (StockOpnameListV2) TableName() string {
	return "inv.stock_opname"
}
