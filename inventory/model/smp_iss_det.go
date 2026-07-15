package model

import "time"

type SampleIssueDet struct {
	SmpIssDetId *int       `gorm:"column:smp_iss_det_id;primaryKey" json:"smp_iss_det_id"`
	CustID      string     `gorm:"column:cust_id" json:"cust_id"`
	SmpIssNo    string     `gorm:"column:smp_iss_no" json:"smp_iss_no"`
	SeqNo       int        `gorm:"column:seq_no" json:"seq_no"`
	ProID       int        `gorm:"column:pro_id" json:"pro_id"`
	Qty         *float64   `gorm:"column:qty" json:"qty"`
	QtyStr      *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1  *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2  *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3  *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4  *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5  *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	UnitId1     *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4     *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5     *string    `gorm:"column:unit_id5" json:"unit_id5"`
	EmbInc      *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc      *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	BatchNo     *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"column:exp_date" json:"exp_date"`
	ConvUnit2   *float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   *float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4   *float64   `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5   *float64   `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (SampleIssueDet) TableName() string {
	return "inv.sample_issue_det"
}

type SampleIssueDetRead struct {
	SmpIssDetId *int       `gorm:"column:smp_iss_det_id;primaryKey" json:"smp_iss_det_id"`
	CustID      string     `gorm:"column:cust_id" json:"cust_id"`
	SmpIssNo    string     `gorm:"column:smp_iss_no" json:"smp_iss_no"`
	SeqNo       int        `gorm:"column:seq_no" json:"seq_no"`
	ProID       int        `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string     `gorm:"column:pro_code" json:"pro_code"`
	ProName     string     `gorm:"column:pro_name" json:"pro_name"`
	Qty         *float64   `gorm:"column:qty" json:"qty"`
	QtyStr      *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1  *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2  *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3  *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4  *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5  *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	UnitId1     *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4     *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5     *string    `gorm:"column:unit_id5" json:"unit_id5"`
	EmbInc      *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc      *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	BatchNo     *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"column:exp_date" json:"exp_date"`
	ConvUnit2   *float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   *float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4   *float64   `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5   *float64   `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (SampleIssueDetRead) TableName() string {
	return "inv.sample_issue_det"
}
