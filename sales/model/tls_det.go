package model

type TlsDet struct {
	CustId   string  `gorm:"cust_id" json:"cust_id"`
	TlsId    int64   `gorm:"tls_id" json:"tls_id"`
	TlsDetId *int64  `gorm:"column:tls_det_id;primaryKey" json:"tls_det_id"`
	OutletId *int64  `gorm:"outlet_id" json:"outlet_id"`
	Notes    *string `gorm:"notes" json:"notes"`
}

func (TlsDet) TableName() string {
	return "sls.tls_det"
}

type TlsDetRead struct {
	CustId     string  `gorm:"cust_id" json:"cust_id"`
	TlsId      int64   `gorm:"tls_id" json:"tls_id"`
	TlsDetId   *int64  `gorm:"column:tls_det_id;primaryKey" json:"tls_det_id"`
	OutletId   *int64  `gorm:"outlet_id" json:"outlet_id"`
	OutletCode string  `gorm:"outlet_code" json:"outlet_code"`
	OutletName string  `gorm:"outlet_name" json:"outlet_name"`
	Notes      *string `gorm:"notes" json:"notes"`
}

func (TlsDetRead) TableName() string {
	return "sls.tls_det"
}
