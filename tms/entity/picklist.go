package entity

import "time"

type Picklist struct {
	PicklistNo string          `json:"picklist_no"`
	Driver     string          `json:"driver"`
	Helper     string          `json:"helper"`
	Vehicle    string          `json:"vehicle"`
	Orders     []OrderResponse `json:"orders"`
	UpdatedBy  string          `json:"updated_by"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type PicklistFilter struct {
	Driver    string `json:"driver"`
	Vehicle   string `json:"vehicle"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type OrderPicklist struct {
	ID          uint           `json:"id"`
	PicklistNo  string         `json:"picklist_no"`
	OrderNo     string         `json:"order_no"`
	InvoiceNo   string         `json:"invoice_no,omitempty"`
	OutletName  string         `json:"outlet_name"`
	Salesman    string         `json:"salesman"`
	InvoiceDate string         `json:"invoice_date"`
	DueDate     string         `json:"due_date"`
	TotalPrice  float64        `json:"total_price"`
	Ppn         float64        `json:"ppn"`
	Discount    float64        `json:"discount"`
	TotalUnpaid float64        `json:"total_unpaid"`
	PaymentType string         `json:"payment_type"`
	Products    []OrderProduct `json:"products"`
}

type OrderProduct struct {
	ID            uint   `json:"id"`
	OrderID       uint   `json:"order_id"`
	ProductName   string `json:"product_name"`
	ProductCode   string `json:"product_code"`
	ProductId     int    `json:"product_id"`
	Quantity1     int    `json:"quantity_1"`
	Quantity2     int    `json:"quantity_2"`
	Quantity3     int    `json:"quantity_3"`
	QuantityUnit1 string `json:"quantity_unit_1"`
	QuantityUnit2 string `json:"quantity_unit_2"`
	QuantityUnit3 string `json:"quantity_unit_3"`
	// Volume        float64 `json:"volume"`
	// Weight        float64 `json:"weight"`
	Unit1Price float64 `json:"unit_1_price"`
	Unit2Price float64 `json:"unit_2_price"`
	Unit3Price float64 `json:"unit_3_price"`
	Ppn        float64 `json:"ppn"`
	Volume1    float64 `json:"volume1"`
	Volume2    float64 `json:"volume2"`
	Volume3    float64 `json:"volume3"`
	Weight1    float64 `json:"weight1"`
	Weight2    float64 `json:"weight2"`
	Weight3    float64 `json:"weight3"`
}

type CreatePicklistRequest struct {
	//PicklistNo string          `json:"picklist_no" validate:"required"`
	Driver    string                 `json:"driver"`
	Helper    string                 `json:"helper"`
	Vehicle   string                 `json:"vehicle"`
	Orders    []OrderPicklistRequest `json:"orders"`
	UpdatedBy string                 `json:"updated_by"`
}

type OrderPicklistRequest struct {
	PicklistNo  string                `json:"picklist_no"`
	OrderNo     string                `json:"order_no"`
	InvoiceNo   string                `json:"invoice_no,omitempty"`
	OutletName  string                `json:"outlet_name"`
	Salesman    string                `json:"salesman"`
	InvoiceDate string                `json:"invoice_date,omitempty"`
	DueDate     string                `json:"due_date,omitempty"`
	TotalPrice  float64               `json:"total_price"`
	Ppn         float64               `json:"ppn"`
	Discount    float64               `json:"discount"`
	TotalUnpaid float64               `json:"total_unpaid"`
	TotalPromo  float64               `json:"total_promo"`
	PaymentType string                `json:"payment_type"`
	Products    []OrderProductRequest `json:"products"`
}

type OrderProductRequest struct {
	// ID            uint    `json:"id"`
	OrderID       uint   `json:"order_id"`
	ProductName   string `json:"product_name"`
	ProductCode   string `json:"product_code"`
	ProductId     int    `json:"product_id"`
	Quantity1     int    `json:"quantity_1"`
	Quantity2     int    `json:"quantity_2"`
	Quantity3     int    `json:"quantity_3"`
	QuantityUnit1 string `json:"quantity_unit_1"`
	QuantityUnit2 string `json:"quantity_unit_2"`
	QuantityUnit3 string `json:"quantity_unit_3"`
	// Volume        float64 `json:"volume"`
	// Weight        float64 `json:"weight"`
	Unit1Price float64 `json:"unit_1_price"`
	Unit2Price float64 `json:"unit_2_price"`
	Unit3Price float64 `json:"unit_3_price"`
	Ppn        float64 `json:"ppn"`
	Volume1    float64 `json:"volume1"`
	Volume2    float64 `json:"volume2"`
	Volume3    float64 `json:"volume3"`
	Weight1    float64 `json:"weight1"`
	Weight2    float64 `json:"weight2"`
	Weight3    float64 `json:"weight3"`
}

type UpdatePicklistRequest struct {
	PicklistNo string          `json:"picklist_no" validate:"required"`
	Driver     string          `json:"driver"`
	Helper     string          `json:"helper"`
	Vehicle    string          `json:"vehicle"`
	Orders     []OrderPicklist `json:"orders"`
	UpdatedBy  string          `json:"updated_by"`
}

type DeletePicklistRequest struct {
	PicklistNo string `json:"picklist_no" validate:"required"`
}

type GetPicklistRequest struct {
	PicklistNo string `json:"picklist_no" validate:"required"`
}

type GetAllPicklistRequest struct {
	// Add any filters if needed
}

type PicklistResponse struct {
	PicklistNo string          `json:"picklist_no"`
	Driver     string          `json:"driver"`
	Helper     string          `json:"helper"`
	Vehicle    string          `json:"vehicle"`
	Orders     []OrderResponse `json:"orders"`
	UpdatedBy  string          `json:"updated_by"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type OrderResponse struct {
	OrderNo     string                 `json:"order_no"`
	InvoiceNo   string                 `json:"invoice_no,omitempty"`
	OutletName  string                 `json:"outlet_name"`
	Salesman    string                 `json:"salesman"`
	InvoiceDate *time.Time             `json:"invoice_date"`
	DueDate     *time.Time             `json:"due_date"`
	TotalPrice  float64                `json:"total_price"`
	Ppn         float64                `json:"ppn"`
	Discount    float64                `json:"discount"`
	TotalUnpaid float64                `json:"total_unpaid"`
	TotalPromo  float64                `json:"total_promo"`
	PaymentType string                 `json:"payment_type"`
	Products    []OrderProductResponse `json:"products"`
}

type OrderProductResponse struct {
	ProductName   string  `json:"product_name"`
	ProductCode   string  `json:"product_code"`
	ProductId     int     `json:"product_id"`
	Quantity1     int     `json:"quantity_1"`
	Quantity2     int     `json:"quantity_2"`
	Quantity3     int     `json:"quantity_3"`
	QuantityUnit1 string  `json:"quantity_unit_1"`
	QuantityUnit2 string  `json:"quantity_unit_2"`
	QuantityUnit3 string  `json:"quantity_unit_3"`
	Unit1Price    float64 `json:"unit_1_price"`
	Unit2Price    float64 `json:"unit_2_price"`
	Unit3Price    float64 `json:"unit_3_price"`
	Ppn           float64 `json:"ppn"`
	Volume1       float64 `json:"volume1"`
	Volume2       float64 `json:"volume2"`
	Volume3       float64 `json:"volume3"`
	Weight1       float64 `json:"weight1"`
	Weight2       float64 `json:"weight2"`
	Weight3       float64 `json:"weight3"`
}
