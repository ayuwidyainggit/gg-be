package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type VanUl struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	VanUlNo    string         `gorm:"column:van_ul_no;primaryKey" json:"van_ul_no"`
	VanUlDate  *time.Time     `gorm:"column:van_ul_date" json:"van_ul_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID       *int64         `gorm:"column:wh_id" json:"wh_id"`
	SalesmanID *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	Notes      *string        `gorm:"column:notes" json:"notes"`
	TotEmbInc  *float64       `gorm:"column:tot_emb_inc" json:"tot_emb_inc"`
	TotEmbExc  *float64       `gorm:"column:tot_emb_exc" json:"tot_emb_exc"`
	SubTotal   *float64       `gorm:"column:sub_total" json:"sub_total"`
	Vat        *float64       `gorm:"column:vat" json:"vat"`
	VatValue   *float64       `gorm:"column:vat_value" json:"vat_value"`
	VatLg      *float64       `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue *float64       `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total      *float64       `gorm:"column:total" json:"total"`
	VatBg      *float64       `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue *float64       `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	WeekIDOb   int            `gorm:"column:week_id_ob" json:"week_id_ob"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel      *bool          `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed   bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy   *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt   time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (VanUl) TableName() string {
	return "inv.van_ul"
}

type VanUlNo struct {
	VanUlNo string `gorm:"column:get_no_fn"`
}

func (m *VanUl) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// if m.VanUlNo == "" {
	// 	m.VanUlNo = strconv.Itoa(int(intTmpsStr))
	// }
	// m.CreatedAt = time.Now()
	// m.UpdatedAt = time.Now()
	// m.UpdatedBy = m.CreatedBy
	// if m.IsDel == nil {
	// 	isDel := false
	// 	m.IsDel = &isDel
	// }
	// return nil
	var vanUlNo VanUlNo
	trCode := *m.TrCode
	vanUlDateStr := m.VanUlDate.Format("2006-01-02")
	vanUlDateSubtr := vanUlDateStr[2:4]

	queryStr := fmt.Sprintf(`SELECT 
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(van_ul_no,6,5),'99999')),0)+1, '00000')) AS get_no_fn 
	FROM inv.van_ul
	WHERE substr(van_ul_no,4,2) = '%v' AND cust_id = '%v' AND tr_code = '%v'`, vanUlDateSubtr, strings.ToUpper(m.CustID), strings.ToUpper(trCode))
	err = trx.Raw(queryStr).Scan(&vanUlNo).Error
	if err != nil {
		return err
	}

	m.VanUlNo = trCode + vanUlDateSubtr + vanUlNo.VanUlNo
	log.Println("m.VanUlNo:", m.VanUlNo)

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type VanUlRead struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	VanUlNo       string         `gorm:"column:van_ul_no;primaryKey" json:"van_ul_no"`
	VanUlDate     *time.Time     `gorm:"column:van_ul_date" json:"van_ul_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	WhID          *int64         `gorm:"column:wh_id" json:"wh_id"`
	WhCode        string         `gorm:"column:wh_code" json:"wh_code"`
	WhName        string         `gorm:"column:wh_name" json:"wh_name"`
	SalesmanID    *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName  string         `gorm:"column:salesman_name" json:"salesman_name"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
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
	WeekIDOb      int            `gorm:"column:week_id_ob" json:"week_id_ob"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         *bool          `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsClosed      bool           `gorm:"column:is_closed" json:"is_closed"`
	ClosedBy      *int64         `gorm:"column:closed_by" json:"closed_by"`
	ClosedAt      time.Time      `gorm:"column:closed_at" json:"closed_at"`
}

func (VanUlRead) TableName() string {
	return "inv.van_ul"
}
