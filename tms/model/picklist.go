package model

import "time"

// Picklist model
type Picklist struct {
	PicklistNo string `gorm:"primaryKey"`
	Driver     string
	Helper     string
	Vehicle    string
	UpdatedBy  string
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
	Orders     []OrderPicklist `gorm:"foreignKey:PicklistNo"`
	CustId     string          `gorm:"cust_id"`
}

// OrderPicklist model
type OrderPicklist struct {
	ID          uint `gorm:"type:int;primary_key"`
	PicklistNo  string
	OrderNo     string
	InvoiceNo   string
	OutletName  string
	Salesman    string
	InvoiceDate *time.Time
	DueDate     *time.Time
	TotalPrice  float64
	Ppn         float64
	Discount    float64
	TotalUnpaid float64
	TotalPromo  float64
	PaymentType string
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	Products    []OrderProduct `gorm:"foreignKey:OrderID;references:ID"`
	CustId      string         `gorm:"cust_id"`
}

// OrderProduct model
type OrderProduct struct {
	ID            uint `gorm:"type:int;primary_key"`
	OrderID       uint
	ProductName   string `json:"product_name" gorm:"column:product_name"`
	ProductCode   string `json:"product_code" gorm:"column:product_code"`
	ProductId     int    `json:"product_id" gorm:"column:product_id"`
	Quantity1     int    `json:"quantity_1" gorm:"column:quantity_1"`
	Quantity2     int    `json:"quantity_2" gorm:"column:quantity_2"`
	Quantity3     int    `json:"quantity_3" gorm:"column:quantity_3"`
	QuantityUnit1 string `json:"quantity_unit_1" gorm:"column:quantity_unit_1"`
	QuantityUnit2 string `json:"quantity_unit_2" gorm:"column:quantity_unit_2"`
	QuantityUnit3 string `json:"quantity_unit_3" gorm:"column:quantity_unit_3"`
	// Volume        float64
	// Weight        float64
	Unit1Price float64 `json:"unit_1_price" gorm:"column:unit_1_price"`
	Unit2Price float64 `json:"unit_2_price" gorm:"column:unit_3_price"`
	Unit3Price float64 `json:"unit_3_price" gorm:"column:unit_3_price"`
	Ppn        float64
	Volume1    float64   `json:"volume1"`
	Volume2    float64   `json:"volume2"`
	Volume3    float64   `json:"volume3"`
	Weight1    float64   `json:"weight1"`
	Weight2    float64   `json:"weight2"`
	Weight3    float64   `json:"weight3"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	CustId     string    `gorm:"cust_id"`
}

func (Picklist) TableName() string {
	return "picklist.picklist"
}

func (OrderPicklist) TableName() string {
	return "picklist.order_picklist"
}

func (OrderProduct) TableName() string {
	return "picklist.order_product"
}
