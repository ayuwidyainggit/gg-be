package model

type SupplierReturnGet struct {
	CustID           string  `gorm:"column:cust_id" json:"cust_id"`
	SupplierReturnNo string  `gorm:"column:supplier_return_no" json:"supplier_return_no"`
	GrNO             string  `gorm:"column:gr_no" json:"gr_no"`
	InvoiceNo        *string `gorm:"column:invoice_no" json:"invoice_no"`
}

func (SupplierReturnGet) TableName() string {
	return "inv.supplier_returns"
}
