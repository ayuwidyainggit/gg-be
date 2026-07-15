package model

import (
	"errors"
	"fmt"
)

type Product struct {
	CustId      string  `json:"cust_id" gorm:"column:cust_id"`
	ProductId   int64   `json:"pro_id" gorm:"column:pro_id"`
	UnitId1     string  `json:"unit_id1" gorm:"column:unit_id1"`
	UnitId2     string  `json:"unit_id2" gorm:"column:unit_id2"`
	UnitId3     string  `json:"unit_id3" gorm:"column:unit_id3"`
	UnitId4     *string `json:"unit_id4" gorm:"column:unit_id4"`
	UnitId5     *string `json:"unit_id5" gorm:"column:unit_id5"`
	ConvUnit2   float32 `json:"conv_unit2" gorm:"column:conv_unit2"`
	ConvUnit3   float32 `json:"conv_unit3" gorm:"column:conv_unit3"`
	ConvUnit4   float32 `json:"conv_unit4" gorm:"column:conv_unit4"`
	ConvUnit5   float32 `json:"conv_unit5" gorm:"column:conv_unit5"`
	PurchPrice1 float64 `json:"purch_price1" gorm:"column:purch_price1"`
	PurchPrice2 float64 `json:"purch_price2" gorm:"column:purch_price2"`
	PurchPrice3 float64 `json:"purch_price3" gorm:"column:purch_price3"`
	PurchPrice4 float64 `json:"purch_price4" gorm:"column:purch_price4"`
	PurchPrice5 float64 `json:"purch_price5" gorm:"column:purch_price5"`
	SellPrice1  float64 `json:"sell_price1" gorm:"column:sell_price1"`
	SellPrice2  float64 `json:"sell_price2" gorm:"column:sell_price2"`
	SellPrice3  float64 `json:"sell_price3" gorm:"column:sell_price3"`
	SellPrice4  float64 `json:"sell_price4" gorm:"column:sell_price4"`
	SellPrice5  float64 `json:"sell_price5" gorm:"column:sell_price5"`
}

func (Product) TableName() string {
	return "mst.m_product"
}

type MapProduct map[int64]Product

func (m MapProduct) SetProduct(id int64, product Product) {
	m[id] = product
}

func (m MapProduct) GetByID(id int64) (product Product, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return product, errors.New(fmt.Sprintf("Product ID %v Not Found", id))
	}

	return val, nil
}

type ProductInfo struct {
	ConvUnit2 float64
	ConvUnit3 float64
	Subtotal  float64
}

type MapProductInfo map[int64]ProductInfo

func (m MapProductInfo) SetProduct(id int64, product ProductInfo) {
	m[id] = product
}

func (m MapProductInfo) GetByID(id int64) (product ProductInfo, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return product, errors.New(fmt.Sprintf("Product ID %v Not Found", id))
	}

	return val, nil
}
