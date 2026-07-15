package entity

const (
	SUBMIT         = 1
	PROCESSEDORDER = 2
	COMPLETEDORDER = 3
	REJECTED       = 0
)

var dataStatusName = map[int64]string{
	SUBMIT:         "Submit ",
	PROCESSEDORDER: "Processed",
	COMPLETEDORDER: "Completed",
	REJECTED:       "Rejected ",
}

var DataStatus = dataStatusName

type StatusList struct {
	StatusOrderBooking     int    ` json:"status_order_booking_id"`
	StatusOrderBookingName string ` json:"status_order_booking_name"`
}

type OrderQueryFilter struct {
	Status       []int `query:"status"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type CreateOrderBody struct {
	CustId       string `json:"cust_id"`
	ParentCustId string `json:"parent_cust_id"`

	SupId         int64                ` json:"sup_id"`
	SubTotal      float64              ` json:"sub_total"`
	SubTotalAlloc *float64             ` json:"sub_total_alloc"`
	Vat           float64              ` json:"vat"`
	VatValue      float64              ` json:"vat_value"`
	VatValueAlloc *float64             ` json:"vat_value_alloc"`
	Total         float64              ` json:"total"`
	TotalAlloc    *float64             ` json:"total_alloc"`
	Details       []CreateOrderDetBody `json:"details"`

	CreatedBy *int64 ` json:"created_by"`
}

type CreateOrderDetBody struct {
	OrderBookingId int      `json:"order_booking_id"`
	ProId          int      `json:"pro_id"`
	ProCode        *string  `json:"pro_code"`
	ProName        *string  `json:"pro_name"`
	SupId          *int     `json:"sup_id"`
	SupName        *string  `json:"sup_name"`
	ItemType       int      `json:"item_type"`
	Qty            *float64 `json:"qty_bo"`
	QtyAlloc       *float64 `json:"qty_alloc"`
	Qty1           *float64 `json:"qty1"`
	Qty2           *float64 `json:"qty2"`
	Qty3           *float64 `json:"qty3"`
	Qty4           *float64 `json:"qty4"`
	Qty5           *float64 `json:"qty5"`
	Qty1Alloc      *float64 `json:"qty1_alloc"`
	Qty2Alloc      *float64 `json:"qty2_alloc"`
	Qty3Alloc      *float64 `json:"qty3_alloc"`
	Qty4Alloc      *float64 `json:"qty4_alloc"`
	Qty5Alloc      *float64 `json:"qty5_alloc"`
	Qty1Total      float64  `json:"qty1_total"`
	Qty2Total      float64  `json:"qty2_total"`
	Qty3Total      float64  `json:"qty3_total"`
	Qty4Total      float64  `json:"qty4_total"`
	Qty5Total      float64  `json:"qty5_total"`
	PurchPrice1    *float64 `json:"purch_price1"`
	PurchPrice2    *float64 `json:"purch_price2"`
	PurchPrice3    *float64 `json:"purch_price3"`
	PurchPrice4    *float64 `json:"purch_price4"`
	PurchPrice5    *float64 `json:"purch_price5"`
	SellPrice1     *float64 `json:"sell_price1"`
	SellPrice2     *float64 `json:"sell_price2"`
	SellPrice3     *float64 `json:"sell_price3"`
	SellPrice4     *float64 `json:"sell_price4"`
	SellPrice5     *float64 `json:"sell_price5"`
	Amount         *float64 `json:"amount"`
	AmountAlloc    *float64 `json:"amount_alloc"`
	Vat            *float64 `json:"vat"`
	VatValue       *float64 `json:"vat_value"`
	VatValueAlloc  *float64 `json:"vat_value_alloc"`
	UnitId1        *string  `json:"unit_id1"`
	UnitId2        *string  `json:"unit_id2"`
	UnitId3        *string  `json:"unit_id3"`
	UnitId4        *string  `json:"unit_id4"`
	UnitId5        *string  `json:"unit_id5"`
	ConvUnit2      *int     `json:"conv_unit2"`
	ConvUnit3      *int     `json:"conv_unit3"`
	ConvUnit4      *int     `json:"conv_unit4"`
	ConvUnit5      *int     `json:"conv_unit5"`
}

type OrderBookingDetailResponse struct {
	OrderBookingDetailId *int `json:"order_booking_detail_id"`
	CreateOrderDetBody
}

type OrderBookingDetailFinalResponse struct {
	OrderBookingDetailId *int `json:"order_booking_detail_id"`
	CreateOrderDetBody
	UnitPrice1   *float64 `json:"unit_price1"`
	UnitPrice2   *float64 `json:"unit_price2"`
	UnitPrice3   *float64 `json:"unit_price3"`
	UnitPrice4   *float64 `json:"unit_price4"`
	UnitPrice5   *float64 `json:"unit_price5"`
	QtyReceived  float64  `json:"qty_received"`
	QtyReceived1 *float64 `json:"qty_received1"`
	QtyReceived2 *float64 `json:"qty_received2"`
	QtyReceived3 *float64 `json:"qty_received3"`
	QtyReceived4 *float64 `json:"qty_received4"`
	QtyReceived5 *float64 `json:"qty_received5"`
}

type CreateOrderBookingResponse struct {
	OrderBookingId int `json:"order_booking_id"`
}

type OrderBookingResponse struct {
	OrderBookingId       int                               `json:"order_booking_id"`
	CustId               string                            `json:"cust_id"`
	ParentCustId         string                            `json:"parent_cust_id"`
	PoNo                 string                            `json:"po_no"`
	SoPo                 string                            `json:"so_po"`
	GrBranchNo           string                            `json:"gr_branch_no"`
	DistributorId        int                               `json:"distributor_id"`
	DistributorName      string                            `json:"distributor_name"`
	DistributorCode      string                            `json:"distributor_code"`
	DistributorAddress   string                            `json:"distributor_address"`
	CreatedBy            int64                             ` json:"created_by"`
	CreatedByName        string                            ` json:"created_by_name"`
	UpdatedBy            int64                             ` json:"updated_by"`
	UpdatedByName        string                            ` json:"updated_by_name"`
	SupId                int64                             ` json:"sup_id"`
	SupName              string                            ` json:"sup_name"`
	SupCode              string                            ` json:"sup_code"`
	SubTotal             float64                           ` json:"sub_total"`
	SubTotalAlloc        *float64                          ` json:"sub_total_alloc"`
	SubTotalTotal        float64                           ` json:"sub_total_total"`
	SubTotalFinal        *float64                          ` json:"sub_total_final"`
	Vat                  float64                           ` json:"vat"`
	VatValue             float64                           ` json:"vat_value"`
	VatValueAlloc        *float64                          ` json:"vat_value_alloc"`
	VatValueTotal        float64                           ` json:"vat_value_total"`
	VatValueFinal        *float64                          ` json:"vat_value_final"`
	DeliveryFee          *float64                          ` json:"delivery_fee"`
	DeliveryFeeFinal     *float64                          ` json:"delivery_fee_final"`
	Total                float64                           ` json:"total"`
	TotalAlloc           *float64                          ` json:"total_alloc"`
	TotalTotal           float64                           ` json:"total_total"`
	TotalTotalFinal      *float64                          ` json:"total_total_final"`
	Status               int64                             ` json:"status_order_booking"`
	StatusName           string                            ` json:"status_name"`
	TypeApproval         *int                              ` json:"type_approval"`
	CreatedAt            string                            ` json:"created_at"`
	Details              []OrderBookingDetailResponse      `json:"details"`
	OrderBooking         []OrderBookingDetailResponse      `json:"order_booking"`
	OrderBookingFinal    []OrderBookingDetailFinalResponse `json:"order_booking_final"`
	OrderBookingApproval []OrderBookingDetailResponse      `json:"order_booking_approval"`
}

func (ro OrderBookingResponse) GenerateDataStatusName() string {
	return dataStatusName[ro.Status]
}

type DetailOrderBookingParams struct {
	OrderBookingId int `params:"order_booking_id" validate:"required"`
}
type DeleteOrderBookingParams struct {
	OrderBookingId int `params:"order_booking_id" validate:"required"`
}

type UpdateOrderBookingParams struct {
	OrderBookingId int `params:"order_booking_id" validate:"required"`
}

type OrderBookingListResponse struct {
	OrderBookingId     int      `json:"order_booking_id"`
	CustId             string   `json:"cust_id"`
	PoNo               string   `json:"po_no"`
	SoPo               string   `json:"so_po"`
	SupId              int64    ` json:"sup_id"`
	SupName            string   ` json:"sup_name"`
	CreditLimit        float64  ` json:"credit_limit"`
	DistributorName    string   ` json:"distributor_name"`
	DistributorAddress string   ` json:"distributor_address"`
	CreatedBy          int64    ` json:"created_by"`
	CreatedByName      string   ` json:"created_by_name"`
	UpdatedBy          int64    ` json:"updated_by"`
	UpdatedByName      string   ` json:"updated_by_name"`
	SubTotal           float64  ` json:"sub_total"`
	SubTotalAlloc      *float64 ` json:"sub_total_alloc"`
	Vat                float64  ` json:"vat"`
	VatValue           float64  ` json:"vat_value"`
	VatValueAlloc      *float64 ` json:"vat_value_alloc"`
	Total              float64  ` json:"total"`
	TotalAlloc         *float64 ` json:"total_alloc"`
	TotalTotal         float64  ` json:"total_total"`
	Status             int64    ` json:"status_order_booking"`
	StatusName         string   ` json:"status_name"`
	CreatedAt          string   ` json:"created_at"`
}

func (ro OrderBookingListResponse) GenerateDataStatusName() string {
	return dataStatusName[ro.Status]
}

type CreateConversionBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	ProductId int64  `json:"pro_id"`
	Qty1      int64  `json:"qty1" validate:"numeric"`
	Qty2      int64  `json:"qty2" validate:"numeric"`
	Qty3      int64  `json:"qty3" validate:"numeric"`
}

type OrderConversionResponse struct {
	Qty1     int64 `json:"qty1"`
	Qty2     int64 `json:"qty2"`
	Qty3     int64 `json:"qty3"`
	TotalQty int64 `json:"total_qty"`
}

type UpdateOrderBookingStatus struct {
	CustId         string `json:"cust_id"`
	OrderBookingId int    `json:"order_booking_id"`
}
