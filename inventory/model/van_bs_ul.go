package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type VanBsUl struct {
	CustID      string         `gorm:"column:cust_id" json:"cust_id"`
	VanBsUlNo   string         `gorm:"column:van_bs_ul_no" json:"van_bs_ul_no"`
	VanBsUlDate *time.Time     `gorm:"column:van_bs_ul_date" json:"van_bs_ul_date"`
	TrCode      *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID        *int64         `gorm:"column:wh_id" json:"wh_id"`
	SalesmanID  *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	RefNo       *string        `gorm:"column:ref_no" json:"ref_no"`
	TotEmbInc   *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc   *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	SubTotal    *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat         *float64       `gorm:"column:vat" json:"vat"`
	VatValue    *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg       *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue  *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total       *float64       `gorm:"column:total" json:"total"`
	VatBg       *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue  *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus  *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel       bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy   *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed    bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy    *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt    time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (VanBsUl) TableName() string {
	return "inv.van_bs_ul"
}

type VanBsUlNo struct {
	VanBsUlNo string `gorm:"column:get_no_fn"`
}

func (m *VanBsUl) BeforeCreate(trx *gorm.DB) (err error) {
	var vanBsUlNo VanBsUlNo
	trCode := *m.TrCode
	vanBsUlDateStr := m.VanBsUlDate.Format("2006-01-02")
	vanBsUlDateSubtr := vanBsUlDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(van_bs_ul_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.van_bs_ul
	WHERE substr(van_bs_ul_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, vanBsUlDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&vanBsUlNo).Error
	if err != nil {
		return err
	}

	m.VanBsUlNo = trCode + vanBsUlDateSubtr + vanBsUlNo.VanBsUlNo
	log.Println("m.VanBsUlNo:", m.VanBsUlNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

type VanBsUlList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	VanBsUlNo     string         `gorm:"column:van_bs_ul_no" json:"van_bs_ul_no"`
	VanBsUlDate   *time.Time     `gorm:"column:van_bs_ul_date" json:"van_bs_ul_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        *string        `gorm:"column:wh_code" json:"wh_code"`
	WhName        *string        `gorm:"column:wh_name" json:"wh_name"`
	SalesmanID    *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  *string        `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  *string        `gorm:"column:salesman_name" json:"salesman_name"`
	RefNo         *string        `gorm:"column:ref_no" json:"ref_no"`
	TotEmbInc     *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc     *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	SubTotal      *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat           *float64       `gorm:"column:vat" json:"vat"`
	VatValue      *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg         *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue    *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total         *float64       `gorm:"column:total" json:"total"`
	VatBg         *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue    *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
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

func (VanBsUlList) TableName() string {
	return "inv.van_bs_ul"
}
