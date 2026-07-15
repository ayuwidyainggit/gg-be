package entity

// Request Entities

type ReplenishmentQueryFilter struct {
	CustId                 string
	ParentCustId           string
	IsPrincipal            bool
	DistributorIDFromToken *int64
	UserID                 int64
	EmpID                  int64
	IsPICUser              bool
	Page                   int    `query:"page"`
	Limit                  int    `query:"limit" validate:"required"`
	Sort                   string `query:"sort"`
	StartDate              *int64 `query:"start_date" validate:"omitempty,gte=1000000000"`
	EndDate                *int64 `query:"end_date" validate:"omitempty,lte=9999999999,gtefield=StartDate"`
	StartDeliveryDate      *int64 `query:"start_delivery_date" validate:"omitempty,gte=1000000000"`
	EndDeliveryDate        *int64 `query:"end_delivery_date" validate:"omitempty,lte=9999999999,gtefield=StartDeliveryDate"`
	Distributor            *int64 `query:"distributor" validate:"omitempty"`
	SupID                  string `query:"sup_id"`
	Status                 string `query:"status"`
	PoNo                   string `query:"po_no"`
	SoNo                   string `query:"so_no"`
	SupIDParsed            []int64
	StatusParsed           []int
}

// Response Entities

type ReplenishmentListResponse struct {
	ReplenishmentID string                 `json:"replenishment_id"`
	ReplenishmentNo string                 `json:"replenishment_no"`
	Date            string                 `json:"date"`
	DeliveryDate    string                 `json:"delivery_date"`
	SupID           string                 `json:"sup_id"`
	SupCode         string                 `json:"sup_code"`
	SupName         string                 `json:"sup_name"`
	WhID            string                 `json:"wh_id"`
	WhCode          string                 `json:"wh_code"`
	WhName          string                 `json:"wh_name"`
	Status          string                 `json:"status"`
	DistributorID   *int64                 `json:"distributor_id,omitempty"`
	DistributorCode string                 `json:"distributor_code"`
	DistributorName string                 `json:"distributor_name"`
	DistributorAddr string                 `json:"distributor_address"`
	CreatedBy       ReplenishmentUserInfo  `json:"created_by"`
	UpdatedBy       *ReplenishmentUserInfo `json:"updated_by"`
}

type ReplenishmentUserInfo struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Create Request Entities

type CreateReplenishmentOrderBody struct {
	CustID            string                           `json:"cust_id" validate:"required"`
	ParentCustID      string                           `json:"parent_cust_id"`
	CreatedBy         int64                            `json:"created_by"`
	CreatedEmpID      int64                            `json:"-"`
	IsPic             *bool                            `json:"is_pic" validate:"required"`
	DistributorID     *int64                           `json:"distributor_id" validate:"omitempty,gte=0"`
	SupID             int64                            `json:"sup_id" validate:"required"`
	WhID              int64                            `json:"wh_id" validate:"required"`
	DeliveryType      string                           `json:"delivery_type" validate:"omitempty,oneof=Full Partial"`
	ReplenishmentType string                           `json:"replenishment_type" validate:"required,oneof='Replenishment' 'Replenishment Event'"`
	SoStartDate       *string                          `json:"so_start_date" validate:"omitempty,ddMmYyyyDate"`
	SoEndDate         *string                          `json:"so_end_date" validate:"omitempty,ddMmYyyyDate"`
	DeliveryDate      *string                          `json:"delivery_date" validate:"omitempty,ddMmYyyyDate"`
	Note              string                           `json:"note" validate:"omitempty,max=255"`
	Data              []CreateReplenishmentProductBody `json:"data" validate:"required,min=1,dive"`
}

type CreateReplenishmentProductBody struct {
	ProID            int64    `json:"pro_id" validate:"required"`
	OrderBookingQty1 *float64 `json:"order_booking_qty1" validate:"required,gte=0"`
	OrderBookingQty2 *float64 `json:"order_booking_qty2" validate:"required,gte=0"`
	OrderBookingQty3 *float64 `json:"order_booking_qty3" validate:"required,gte=0"`
	PurchPrice1      *float64 `json:"purch_price1" validate:"required"`
	PurchPrice2      *float64 `json:"purch_price2" validate:"required"`
	PurchPrice3      *float64 `json:"purch_price3" validate:"required"`
	EstimatedPrice   *float64 `json:"estimated_price" validate:"required"`
}

type ReplenishmentDetailsProductId struct {
	ProductIds []ProductId `json:"data[].pro_id" validate:"unique,dive"`
}

// Detail Request Entities

type DetailReplenishmentParams struct {
	ReplenishmentID string `params:"replenishment_id" validate:"required"`
	Type            string `query:"type" validate:"omitempty,oneof=replenishment final"`
	Status          string `query:"status"`
}

// Detail Response Entities

type ReplenishmentDetailResponse struct {
	ReplenishmentID    string                               `json:"replenishment_id"`
	ReplenishmentNo    string                               `json:"replenishment_no"`
	Date               string                               `json:"date"`
	SoStartDate        *string                              `json:"so_start_date"`
	SoEndDate          *string                              `json:"so_end_date"`
	DeliveryDate       *string                              `json:"delivery_date"`
	ReplenishmentType  string                               `json:"replenishment_type"`
	DeliveryType       string                               `json:"delivery_type"`
	SupID              string                               `json:"sup_id"`
	SupCode            string                               `json:"sup_code"`
	SupName            string                               `json:"sup_name"`
	WhID               string                               `json:"wh_id"`
	WhCode             string                               `json:"wh_code"`
	WhName             string                               `json:"wh_name"`
	Status             string                               `json:"status"`
	DeliveryFee        int64                                `json:"delivery_fee"`
	Notes              string                               `json:"notes"`
	DistributorID      *int64                               `json:"distributor_id,omitempty"`
	DistributorCode    *string                              `json:"distributor_code,omitempty"`
	DistributorName    *string                              `json:"distributor_name,omitempty"`
	Address            *string                              `json:"address,omitempty"`
	ProductData        []ReplenishmentDetailProductResponse `json:"product_data"`
	FinalReplanishment []ReplenishmentFinalResponse         `json:"final_replanishment"`
	GoodReceipt        []ReplenishmentGoodReceiptResponse   `json:"good_receipt"`
}

type ReplenishmentDetailProductResponse struct {
	ReplenishmentDetID  string   `json:"replenishment_det_id"`
	ProductID           int64    `json:"product_id"`
	ProductCode         string   `json:"product_code"`
	ProductName         string   `json:"product_name"`
	UnitID1             string   `json:"unit_id1"`
	UnitID2             string   `json:"unit_id2"`
	UnitID3             string   `json:"unit_id3"`
	PurchPrice1         float64  `json:"purch_price1"`
	PurchPrice2         float64  `json:"purch_price2"`
	PurchPrice3         float64  `json:"purch_price3"`
	QtyRo1              float64  `json:"qty_ro1"`
	QtyRo2              float64  `json:"qty_ro2"`
	QtyRo3              float64  `json:"qty_ro3"`
	QtyFinal1           *float64 `json:"qty_final1,omitempty"`
	QtyFinal2           *float64 `json:"qty_final2,omitempty"`
	QtyFinal3           *float64 `json:"qty_final3,omitempty"`
	StockReceived1      *float64 `json:"stock_received1,omitempty"`
	StockReceived2      *float64 `json:"stock_received2,omitempty"`
	StockReceived3      *float64 `json:"stock_received3,omitempty"`
	QtyOrderAllocation1 float64  `json:"qty_order_allocation1"`
	QtyOrderAllocation2 float64  `json:"qty_order_allocation2"`
	QtyOrderAllocation3 float64  `json:"qty_order_allocation3"`
	QtyOrderApproval1   float64  `json:"qty_order_approval1"`
	QtyOrderApproval2   float64  `json:"qty_order_approval2"`
	QtyOrderApproval3   float64  `json:"qty_order_approval3"`
	QtyTotal1           float64  `json:"qty_total1"`
	QtyTotal2           float64  `json:"qty_total2"`
	QtyTotal3           float64  `json:"qty_total3"`
	SubTotal            float64  `json:"sub_total"`
	EstimatedPrice      float64  `json:"estimated_price"`
	Vat                 float64  `json:"vat"`
	InTransitStock1     float64  `json:"in_transit_stock1"`
	InTransitStock2     float64  `json:"in_transit_stock2"`
	InTransitStock3     float64  `json:"in_transit_stock3"`
	SafStockQty         float64  `json:"saf_stock_qty"`
	MinStockQty         float64  `json:"min_stock_qty"`
	Qty1                float64  `json:"qty1"`
	Qty2                float64  `json:"qty2"`
	Qty3                float64  `json:"qty3"`
}

type ReplenishmentFinalResponse struct {
	ReplenishmentDetID  string   `json:"replenishment_det_id"`
	ProductID           int64    `json:"product_id"`
	ProductCode         string   `json:"product_code"`
	ProductName         string   `json:"product_name"`
	UnitID1             *string  `json:"unit_id1,omitempty"`
	UnitID2             *string  `json:"unit_id2,omitempty"`
	UnitID3             *string  `json:"unit_id3,omitempty"`
	PurchPriceDelivery1 *float64 `json:"purch_price_delivery1,omitempty"`
	PurchPriceDelivery2 *float64 `json:"purch_price_delivery2,omitempty"`
	PurchPriceDelivery3 *float64 `json:"purch_price_delivery3,omitempty"`
	FinalOrder1         *float64 `json:"final_order1,omitempty"`
	FinalOrder2         *float64 `json:"final_order2,omitempty"`
	FinalOrder3         *float64 `json:"final_order3,omitempty"`
	GrPrice1            *float64 `json:"gr_price1,omitempty"`
	GrPrice2            *float64 `json:"gr_price2,omitempty"`
	GrPrice3            *float64 `json:"gr_price3,omitempty"`
	StockReceived1      *float64 `json:"stock_received1,omitempty"`
	StockReceived2      *float64 `json:"stock_received2,omitempty"`
	StockReceived3      *float64 `json:"stock_received3,omitempty"`
	SubTotal            float64  `json:"sub_total"`
	Vat                 *float64 `json:"vat,omitempty"`
}

type ReplenishmentGoodReceiptResponse struct {
	PoNo           string  `json:"po_no"`
	ProID          int64   `json:"pro_id"`
	ProCode        string  `json:"pro_code"`
	ProName        string  `json:"pro_name"`
	UnitPrice1     float64 `json:"unit_price1"`
	UnitPrice2     float64 `json:"unit_price2"`
	UnitPrice3     float64 `json:"unit_price3"`
	QtyReceived1   float64 `json:"qty_received1"`
	QtyReceived2   float64 `json:"qty_received2"`
	QtyReceived3   float64 `json:"qty_received3"`
	EstimatedPrice float64 `json:"estimated_price"`
	Vat            float64 `json:"vat"`
}

// Product List Query Filter
type ReplenishmentProductQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required,max=500"`
	Sort         string `query:"sort"`
	StartDate    *int64 `query:"start_date" validate:"omitempty,gte=1000000000"`
	EndDate      *int64 `query:"end_date" validate:"omitempty,lte=9999999999,gtefield=StartDate"`
	WhID         int64  `query:"wh_id" validate:"required"`
}

// Product List Response
type ReplenishmentProductListResponse struct {
	ProID           int64   `json:"pro_id"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	PurchPrice1     float64 `json:"purch_price1"`
	PurchPrice2     float64 `json:"purch_price2"`
	PurchPrice3     float64 `json:"purch_price3"`
	UnitID1         string  `json:"unit_id1"`
	UnitID2         string  `json:"unit_id2"`
	UnitID3         string  `json:"unit_id3"`
	Vat             float64 `json:"vat"`
	Qty1            float64 `json:"qty1"`
	Qty2            float64 `json:"qty2"`
	Qty3            float64 `json:"qty3"`
	EstimatedPrice  float64 `json:"estimated_price"`
	InTransitStock1 float64 `json:"in_transit_stock1"`
	InTransitStock2 float64 `json:"in_transit_stock2"`
	InTransitStock3 float64 `json:"in_transit_stock3"`
}

// Product GR List Query Filter
type ProductGrListQueryFilter struct {
	CustId          string
	ParentCustId    string
	Q               string `query:"q"`
	Page            int    `query:"page" validate:"required"`
	Limit           int    `query:"limit" validate:"required"`
	ReplenishmentNo string `query:"replenishment_no" validate:"required"`
}

// Product GR List Response
type ProductGrListResponse struct {
	ReplenishmentNo   string                        `json:"replanishment_no"`
	ReplenishmentType string                        `json:"replenishment_type"`
	Details           []ProductGrListDetailResponse `json:"details"`
}

type ProductGrListDetailResponse struct {
	ReplenishmentDetID int64    `json:"replanishment_det_id"`
	ProID              int64    `json:"pro_id"`
	ProCode            string   `json:"pro_code"`
	ProName            string   `json:"pro_name"`
	Vat                *float64 `json:"vat,omitempty"`
	StockShipment1     float64  `json:"stock_shipment1"`
	StockShipment2     float64  `json:"stock_shipment2"`
	StockShipment3     float64  `json:"stock_shipment3"`
	StockReceived1     float64  `json:"stock_reveived1"`
	StockReceived2     float64  `json:"stock_reveived2"`
	StockReceived3     float64  `json:"stock_reveived3"`
	UnitID1            string   `json:"unit_id1"`
	UnitID2            string   `json:"unit_id2"`
	UnitID3            string   `json:"unit_id3"`
	PurchPrice1        float64  `json:"purch_price1"`
	PurchPrice2        float64  `json:"purch_price2"`
	PurchPrice3        float64  `json:"purch_price3"`
	SubTotal           float64  `json:"sub_total"`
}

// PO List Query Filter
type PoListQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int `query:"page" validate:"required"`
	Limit        int `query:"limit" validate:"required"`
}

// PO List Response
type PoListResponse struct {
	ReplenishmentNo   string `json:"replanishment_no"`
	ReplenishmentType string `json:"replenishment_type"`
	WhID              int64  `json:"wh_id"`
	WhCode            string `json:"wh_code"`
	WhName            string `json:"wh_name"`
	SupID             int64  `json:"sup_id"`
	SupCode           string `json:"sup_code"`
	SupName           string `json:"sup_name"`
}

type ReplenishmentApprovalListQueryFilter struct {
	CustId                 string
	ParentCustId           string
	IsPrincipal            bool
	DistributorIDFromToken *int64
	UserID                 int64
	EmpID                  int64
	Page                   int    `query:"page"`
	Limit                  int    `query:"limit" validate:"required,min=1"`
	Q                      string `query:"q"`
	StartDate              *int64 `query:"start_date" validate:"omitempty,gte=1000000000"`
	EndDate                *int64 `query:"end_date" validate:"omitempty,lte=9999999999,gtefield=StartDate"`
	DeliveryDateStart      *int64 `query:"delivery_date_start" validate:"omitempty,gte=1000000000"`
	DeliveryDateEnd        *int64 `query:"delivery_date_end" validate:"omitempty,lte=9999999999,gtefield=DeliveryDateStart"`
	Status                 string `query:"status"`
	StatusParsed           []int
}

type ReplenishmentApprovalListResponse struct {
	ReplenishmentID int64   `json:"replenishment_id"`
	Date            string  `json:"date"`
	DeliveryDate    *string `json:"delivery_date,omitempty"`
	ReplenishmentNo string  `json:"replenishment_no"`
	SupID           int64   `json:"sup_id"`
	SupCode         string  `json:"sup_code"`
	SupName         string  `json:"sup_name"`
	DistributorID   *int64  `json:"distributor_id,omitempty"`
	DistributorCode string  `json:"distributor_code"`
	DistributorName string  `json:"distributor_name"`
	Address         string  `json:"address"`
	WhID            int64   `json:"wh_id"`
	WhCode          string  `json:"wh_code"`
	WhName          string  `json:"wh_name"`
	CreatedBy       int64   `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	CreatedAt       string  `json:"created_at"`
	UpdatedBy       *int64  `json:"updated_by,omitempty"`
	UpdatedByName   string  `json:"updated_by_name,omitempty"`
	UpdatedAt       *string `json:"updated_at,omitempty"`
}

// Replenishment Approval Product Query Filter
type ReplenishmentApprovalProductQueryFilter struct {
	CustId        string
	ParentCustId  string
	UserID        int64
	Q             string `query:"q"`
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required,max=500"`
	SupID         *int64 `query:"sup_id" validate:"omitempty,gte=0"`
	WhID          int64  `query:"wh_id" validate:"omitempty,gte=0"`
	DistributorID *int64 `query:"distributor_id" validate:"omitempty,gte=0"`
	ZeroStock     *bool  `query:"zero_stock" validate:"omitempty"`
}

// Replenishment Approval Product Response
type ReplenishmentApprovalProductResponse struct {
	ProID           int64    `json:"pro_id"`
	ProCode         string   `json:"pro_code"`
	ProName         string   `json:"pro_name"`
	Ripening        *float64 `json:"ripening,omitempty"`
	InTransitStock1 float64  `json:"in_transit_stock1"`
	InTransitStock2 float64  `json:"in_transit_stock2"`
	InTransitStock3 float64  `json:"in_transit_stock3"`
	SafStockQty     *float64 `json:"saf_stock_qty,omitempty"`
	MinStockQty     *float64 `json:"min_stock_qty,omitempty"`
	Vat             *float64 `json:"vat,omitempty"`
	ConvUnit2       *float64 `json:"conv_unit2,omitempty"`
	ConvUnit3       *float64 `json:"conv_unit3,omitempty"`
	UnitID1         *string  `json:"unit_id1,omitempty"`
	UnitID2         *string  `json:"unit_id2,omitempty"`
	UnitID3         *string  `json:"unit_id3,omitempty"`
	PurchPrice1     *float64 `json:"purch_price1,omitempty"`
	PurchPrice2     *float64 `json:"purch_price2,omitempty"`
	PurchPrice3     *float64 `json:"purch_price3,omitempty"`
	Qty1            *float64 `json:"qty1,omitempty"`
	Qty2            *float64 `json:"qty2,omitempty"`
	Qty3            *float64 `json:"qty3,omitempty"`
	TotalQty        *float64 `json:"total_qty,omitempty"`
}

// Update Approval Request Entities
type UpdateReplenishmentApprovalRequest struct {
	Approval    *bool                               `json:"approval" validate:"required"`
	Remarks     *string                             `json:"remarks"`
	Data        []UpdateReplenishmentApprovalDetail `json:"data"`
	DeletedData []int64                             `json:"deleted_data"`
}

type UpdateReplenishmentApprovalDetail struct {
	ReplenishmentDetID *int64   `json:"replenishment_det_id"`
	ProID              *int64   `json:"pro_id"`
	ReturnReasonID     []int64  `json:"return_reason_id"`
	OrderBookingQty1   *float64 `json:"order_booking_qty1"`
	OrderBookingQty2   *float64 `json:"order_booking_qty2"`
	OrderBookingQty3   *float64 `json:"order_booking_qty3"`
	PurchPrice1        *float64 `json:"purch_price1"`
	PurchPrice2        *float64 `json:"purch_price2"`
	PurchPrice3        *float64 `json:"purch_price3"`
	EstimatedPrice     *float64 `json:"estimated_price"`
}

// BatchReplenishmentApprovalRequest is used by PATCH /v1/replenishment-approval/batch
type BatchReplenishmentApprovalRequest struct {
	Remarks *string                               `json:"remarks"`
	Data    []BatchReplenishmentApprovalOrderItem `json:"data" validate:"required,min=1,dive"`
}

// BatchReplenishmentApprovalOrderItem is one replenishment order in the batch (JSON uses replanishment_* spelling per API spec).
type BatchReplenishmentApprovalOrderItem struct {
	ReplenishmentID int64                                 `json:"replenishment_id" validate:"required,gt=0"`
	Approval        *bool                                 `json:"approval" validate:"required"`
	Remarks         *string                               `json:"remarks"`
	Details         []BatchReplenishmentApprovalDetailRow `json:"details"`
	DeletedData     []int64                               `json:"deleted_data"`
}

// BatchReplenishmentApprovalDetailRow matches bulk API field names (replanishment_det_id).
type BatchReplenishmentApprovalDetailRow struct {
	ReplenishmentDetID *int64   `json:"replanishment_det_id"`
	ProID              *int64   `json:"pro_id"`
	ReturnReasonID     []int64  `json:"return_reason_id"`
	OrderBookingQty1   *float64 `json:"order_booking_qty1"`
	OrderBookingQty2   *float64 `json:"order_booking_qty2"`
	OrderBookingQty3   *float64 `json:"order_booking_qty3"`
	PurchPrice1        *float64 `json:"purch_price1"`
	PurchPrice2        *float64 `json:"purch_price2"`
	PurchPrice3        *float64 `json:"purch_price3"`
	EstimatedPrice     *float64 `json:"estimated_price"`
}

// ToUpdateReplenishmentApprovalDetail maps batch row to the single-approval payload shape.
func (r BatchReplenishmentApprovalDetailRow) ToUpdateReplenishmentApprovalDetail() UpdateReplenishmentApprovalDetail {
	return UpdateReplenishmentApprovalDetail{
		ReplenishmentDetID: r.ReplenishmentDetID,
		ProID:              r.ProID,
		ReturnReasonID:     r.ReturnReasonID,
		OrderBookingQty1:   r.OrderBookingQty1,
		OrderBookingQty2:   r.OrderBookingQty2,
		OrderBookingQty3:   r.OrderBookingQty3,
		PurchPrice1:        r.PurchPrice1,
		PurchPrice2:        r.PurchPrice2,
		PurchPrice3:        r.PurchPrice3,
		EstimatedPrice:     r.EstimatedPrice,
	}
}

// BatchReplenishmentApprovalResultItem is one line in the batch response `data` array.
type BatchReplenishmentApprovalResultItem struct {
	ReplenishmentID int64  `json:"replenishment_id"`
	ReplenishmentNo string `json:"replenishment_no,omitempty"`
	Status          string `json:"status"`
	Message         string `json:"message,omitempty"`
}

type SummarizeReplanishmentQuery struct {
	ReplanishmentID []int64 `query:"replanishment_id" validate:"required,min=1,dive,gt=0"`
}

type SummarizeReplanishmentResponse struct {
	ReplanishmentID int64                              `json:"replanishment_id"`
	ReplanishmentNo string                             `json:"replanishment_no"`
	DisributorID    int64                              `json:"disributor_id"`
	DistributorCode string                             `json:"distributor_code"`
	DistributorName string                             `json:"distributor_name"`
	SupID           int64                              `json:"sup_id"`
	SupCode         string                             `json:"sup_code"`
	SupName         string                             `json:"sup_name"`
	WhID            int64                              `json:"wh_id"`
	WhCode          string                             `json:"wh_code"`
	WhName          string                             `json:"wh_name"`
	Details         []SummarizeReplanishmentDetailItem `json:"details"`
}

type SummarizeReplanishmentDetailItem struct {
	ReplanishmentDetailID int64   `json:"replanishment_detail_id"`
	ProID                 int64   `json:"pro_id"`
	ProCode               string  `json:"pro_code"`
	ProName               string  `json:"pro_name"`
	WhStockLarge          int64   `json:"wh_stock_large"`
	WhStockMedium         int64   `json:"wh_stock_medium"`
	WhStockSmall          int64   `json:"wh_stock_small"`
	OptimumQty            int64   `json:"optimum_qty"`
	Ripening              int64   `json:"ripening"`
	ReturnReasonID        *int64  `json:"return_reason_id,omitempty"`
	UnitID1               string  `json:"unit_id1"`
	UnitID2               string  `json:"unit_id2"`
	UnitID3               string  `json:"unit_id3"`
	PurchPrice1           float64 `json:"purch_price1"`
	PurchPrice2           float64 `json:"purch_price2"`
	PurchPrice3           float64 `json:"purch_price3"`
	QtyRo1                float64 `json:"qty_ro1"`
	QtyRo2                float64 `json:"qty_ro2"`
	QtyRo3                float64 `json:"qty_ro3"`
	EstimatedPrice        float64 `json:"estimated_price"`
	RequestID             string  `json:"request_id"`
}
