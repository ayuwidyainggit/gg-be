package entity

type RejectRequest struct {
	ID          []int  `validate:"required,notEmptyIntSlice" json:"id"`
	ReasonID    int    `validate:"required"  json:"reason_id"`
	ReasonName  string `validate:"required"  json:"reason_name"`
	CurrentTime int64  `validate:"required"  json:"current_time"`
}

type RejectCancelRequest struct {
	ID []int `validate:"required,notEmptyIntSlice" json:"id"`
}

type Quantity struct {
	UnitID string `validate:"required" json:"unit_id" example:"CTN"`
	Stock  int64  `validate:"min=0" json:"stock" example:"2"`
}

type Product struct {
	ID  int        `validate:"required" json:"id" example:"1224"`
	Qty []Quantity `validate:"required,dive" json:"qty"`
}

type RejectPartialBody struct {
	ReasonID    int       `validate:"required" json:"reason_id"`
	ReasonName  string    `validate:"required" json:"reason_name"`
	OutletID    int       `validate:"required" json:"outlet_id"`
	ShipmentNo  string    `validate:"required" json:"shipment_no"`
	CurrentTime int64     `validate:"required"  json:"current_time"`
	Products    []Product `validate:"required,dive" json:"products"`
}

type RejectPartialRequest struct {
	Data []RejectPartialBody `validate:"required,dive" json:"data"`
}

type ReasonRejectResponse struct {
	RejectReasonID   int    `json:"reject_reason_id" example:"1"`
	RejectReasonCode string `json:"reject_reason_code" example:"0001"`
	RejectReasonName string `json:"reject_reason_name" example:"reject"`
}

type RejectResponse struct {
	ID            int    `json:"id" example:"1"`
	ProductID     int    `json:"product_id" example:"1"`
	ProductName   string `json:"product_name" example:"ASSORTED BISCUIT OPP RED 20 275G"`
	ProductStatus string `json:"product_status" example:"Reject"`
	Sku           string `json:"sku" example:"0001"`
	Qty1          int64  `json:"qty1" example:"0"`
	Qty2          int64  `json:"qty2" example:"0"`
	Qty3          int64  `json:"qty3" example:"0"`
	QtyReject1    int64  `json:"qty_reject_1" example:"10"`
	QtyReject2    int64  `json:"qty_reject_2" example:"10"`
	QtyReject3    int64  `json:"qty_reject_3" example:"10"`
	ConvUnit1     int64  `json:"conv_unit1" example:"0"`
	ConvUnit2     int64  `json:"conv_unit2" example:"20"`
	ConvUnit3     int64  `json:"conv_unit3" example:"10"`
	UnitId1       string `json:"unit_id1" example:"PCS"`
	UnitId2       string `json:"unit_id2" example:"CTN"`
	UnitId3       string `json:"unit_id3" example:"CTN"`
	CtgId1        string `json:"ctg_id1"`
	CtgId2        string `json:"ctg_id2"`
	CtgId3        string `json:"ctg_id3"`
	DriverID      int    `json:"driver_id" example:"1"`
	ReasonID      int    `json:"reason_id" example:"1"`
	ReasonName    string `json:"reason_name" example:"reject"`
}

type RejectPartialResponse struct {
	ID         int          `json:"id"`
	ReasonID   int          `json:"reason_id"`
	ReasonName string       `json:"reason_name"`
	Products   []ProductMap `json:"products"`
}

type RejectQueryFilter struct {
	ReasonID    int    `query:"reason_id"`
	DriverID    int    `query:"driver_id"`
	OutletID    int    `query:"outlet_id"`
	ProductName string `query:"product_name"`
	Sort        string `query:"sort"`
	IsActive    int    `query:"is_active"`
	ShipmentNo  string `json:"shipment_no"`
}

type ConversionRequest struct {
	OrderDetailID *int `json:"order_detail_id,omitempty"`
	ProId int `json:"pro_id"`
	Qty1  int `json:"qty1"`
	Qty2  int `json:"qty2"`
	Qty3  int `json:"qty3"`
}

type UpdateQtyFinalOrderRequest struct {
	DetailsFinal struct {
		Normal []ConversionRequest `json:"normal"`
	} `json:"details_final"`
}