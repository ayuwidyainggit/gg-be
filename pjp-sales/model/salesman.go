package model

type SalesmanDetail struct {
	CustId          string `gorm:"column:cust_id" json:"cust_id"`
	SalesmanId      int    `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode    string `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName    string `gorm:"column:salesman_name" json:"salesman_name"`
	AllowInputPrice bool   `gorm:"column:allow_input_price" json:"allow_input_price" `
	WhId            int    `gorm:"column:wh_id" json:"wh_id"`
}

func (SalesmanDetail) TableName() string {
	return "mst.m_salesman"
}
