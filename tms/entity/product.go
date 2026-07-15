package entity

type ProductResponse struct {
	ID            int    `json:"id" example:"1"`
	ProductID     int    `json:"product_id" example:"1"`
	ProductName   string `json:"product_name" example:"ASSORTED BISCUIT OPP RED 20 275G"`
	ProductStatus string `json:"product_status" example:"Reject"`
	Sku           string `json:"sku" example:"0001"`
	Vat           int    `json:"vat" example:"1"`
	VatValue      int    `json:"vat_value" example:"1111"`
	Qty1          int64  `json:"qty1" example:"10"`
	Qty2          int64  `json:"qty2" example:"10"`
	Qty3          int64  `json:"qty3" example:"10"`
	ConvUnit1     int64  `json:"conv_unit1" example:"0"`
	ConvUnit2     int64  `json:"conv_unit2" example:"20"`
	ConvUnit3     int64  `json:"conv_unit3" example:"20"`
	UnitId1       string `json:"unit_id1" example:"PCS"`
	UnitId2       string `json:"unit_id2" example:"CTN"`
	UnitId3       string `json:"unit_id3" example:"CTN"`
	CtgId1        string `json:"ctg_id1"`
	CtgId2        string `json:"ctg_id2"`
	CtgId3        string `json:"ctg_id3"`
	DriverID      int    `json:"driver_id" example:"1"`
	OutletID      int    `json:"outlet_id" example:"34"`
	OutletCode    string `json:"outlet_code" example:"18"`
	Status        string `json:"status" example:"Delivery"`
	Salesman      string `json:"salesman" example:"John Doe"`
	ReturnNo      string `json:"return_no" example:"SR2409240002"`
	ReturnDate    string `json:"return_date" example:"2024-09-27"`
	ReasonName    string `json:"reason_name" example:"Reject"`
	ConditionName string `json:"condition_name" example:"good"`
}

type ProductMap struct {
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductStatus string `json:"product_status"`
	Sku           string `json:"sku"`
	UnitID1       string `json:"unit_id1"`
	UnitID2       string `json:"unit_id2"`
	UnitID3       string `json:"unit_id3"`
	ConvUnit1     int64  `json:"conv_unit1"`
	ConvUnit2     int64  `json:"conv_unit2"`
	ConvUnit3     int64  `json:"conv_unit3"`
	Stock1        int64  `json:"stock1"`
	Stock2        int64  `json:"stock2"`
	Stock3        int64  `json:"stock3"`
	CtgId1        string `json:"ctg_id1"`
	CtgId2        string `json:"ctg_id2"`
	CtgId3        string `json:"ctg_id3"`
}
