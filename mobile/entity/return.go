package entity

type ReturnReasonsRequest struct {
}

type ReturnQueryFilter struct {
	CustId       string
	ParentCustId string
	SalesmanId   []int  `query:"salesman_id"`
	OutletID     []int  `query:"outlet_id"`
	Status       []int  `query:"status"`
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type WarehouseQueryFilter struct {
	CustId       string
	ParentCustId string
	SalesmanId   int    `query:"salesman_id" validate:"required"`
	ItemCnd      int    `query:"item_cnd" validate:"required"`
	StockType    string `query:"stock_type"`
	WhId         int    `query:"wh_id"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:""`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

/*
	type CreateReturnBodyOld struct {
		CustID        string                   `json:"cust_id"`
		ReturnNo      string                   `json:"return_no"`
		ReturnDate    *string                  `json:"return_date"`
		SysDate       *string                  `json:"sys_date"`
		ReturnType    *int64                   `json:"return_type"`
		OutletID      *int64                   `json:"outlet_id"`
		OutletTaxNo   *string                  `json:"outlet_tax_no"`
		PoNo          *string                  `json:"po_no"`
		VehicleNo     *string                  `json:"vehicle_no"`
		InvoiceNo     *string                  `json:"invoice_no"`
		InvoiceDate   *string                  `json:"invoice_date"`
		DeliveryDate  *string                  `json:"delivery_date"`
		PayType       *int64                   `json:"pay_type"`
		SumNo         *string                  `json:"sum_no"`
		DataSource    *int64                   `json:"data_source"`
		MobileID      *int64                   `json:"mobile_id"`
		SubTotal      *float64                 `json:"sub_total"`
		Disc          *float64                 `json:"disc"`
		DiscValue     *float64                 `json:"disc_value"`
		PromoValue    *float64                 `json:"promo_value"`
		CashDiscValue *float64                 `json:"cash_disc_value"`
		TotDisc1      *float64                 `json:"tot_disc_1"`
		TotDisc2      *float64                 `json:"tot_disc_2"`
		Vat           *float64                 `json:"vat"`
		VatValue      *float64                 `json:"vat_value"`
		Total         *float64                 `json:"total"`
		DataStatus    *int64                   `json:"data_status"`
		CreatedBy     *int64                   `json:"created_by"`
		Details       CreateReturnDetBodyGroup `json:"details"`
	}
*/
type ReturnResponse struct {
	RefferenceNo   *string                `json:"refference_no"`
	ReturnNo       string                 `json:"return_no"`
	ReturnDate     *string                `json:"return_date"`
	InvoiceNo      *string                `json:"invoice_no"`
	InvoiceDate    *string                `json:"invoice_date"`
	SalesmanID     *int64                 `json:"salesman_id"`
	SalesmanCode   *string                `json:"salesman_code"`
	SalesmanName   *string                `json:"salesman_name"`
	OutletID       *int64                 `json:"outlet_id"`
	OutletCode     *string                `json:"outlet_code"`
	OutletName     *string                `json:"outlet_name"`
	TprCashValue   *float64               `json:"tpr_cash_value"`
	TprItemValue   *float64               `json:"tpr_item_value"`
	Discount       *float64               `json:"discount"`
	DiscountValue  *float64               `json:"discount_value"`
	Vat            *float64               `json:"vat"`
	VatValue       *float64               `json:"vat_value"`
	SubTotal       *float64               `json:"sub_total"`
	Total          *float64               `json:"total"`
	DataStatus     *int64                 `json:"data_status"`
	DataStatusName *string                `json:"data_status_name"`
	Details        []ReturnDetailResponse `json:"details"`
}

type ReturnListResponse struct {
	RefferenceNo   *string `json:"refference_no"`
	ReturnNo       string  `json:"return_no"`
	ReturnDate     string  `json:"return_date"`
	InvoiceNo      *string `json:"invoice_no"`
	InvoiceDate    *string `json:"invoice_date"`
	SalesmanID     *int64  `json:"salesman_id"`
	SalesmanCode   *string `json:"salesman_code"`
	SalesmanName   *string `json:"salesman_name"`
	OutletID       *int64  `json:"outlet_id"`
	OutletCode     *string `json:"outlet_code"`
	OutletName     *string `json:"outlet_name"`
	DataStatus     *int64  `json:"data_status"`
	DataStatusName *string `json:"data_status_name"`
	CreatedBy      *int64  `json:"created_by"`
	CreatedByName  *string `json:"created_by_name"`
	CreatedAt      *string `json:"created_at"`
	ReviewedBy     *int64  `json:"reviewed_by"`
	ReviewedByName *string `json:"reviewed_by_name"`
	ReviewedAt     *string `json:"reviewed_at"`

	// ProCode    string  `json:"pro_code"`
	// ProName    string  `json:"pro_name"`
	// UnitId1    string  `json:"unit_id1"`
	// UnitId2    string  `json:"unit_id2"`
	// UnitId3    string  `json:"unit_id3"`
	// Qty        int64   `json:"qty"`
	// Price      float64 `json:"price"`
}

type ReturnShipmentListResponse struct {
	// RefferenceNo *string `json:"refference_no"`
	ReturnNo        string  `json:"return_no"`
	ReturnDate      string  `json:"return_date"`
	InvoiceNo       *string `json:"invoice_no"`
	InvoiceDate     *string `json:"invoice_date"`
	SalesmanID      *int64  `json:"salesman_id"`
	SalesmanCode    *string `json:"salesman_code"`
	SalesmanName    *string `json:"salesman_name"`
	OutletID        *int64  `json:"outlet_id"`
	OutletCode      *string `json:"outlet_code"`
	OutletName      *string `json:"outlet_name"`
	OutletAddress   *string `json:"outlet_address"`
	OutletLatitude  *string `json:"outlet_latitude"`
	OutletLongitude *string `json:"outlet_longitude"`
	DataStatus      *int64  `json:"data_status"`
	DataStatusName  *string `json:"data_status_name"`
	// CreatedBy      *int64                 `json:"created_by"`
	// CreatedByName  *string                `json:"created_by_name"`
	// CreatedAt      *string                `json:"created_at"`
	// ReviewedBy     *int64                 `json:"reviewed_by"`
	// ReviewedByName *string                `json:"reviewed_by_name"`
	// ReviewedAt     *string                `json:"reviewed_at"`
	TotalVolume float64                `json:"total_volume"`
	TotalWeight float64                `json:"total_weight"`
	Details     []ReturnDetailResponse `json:"details"`
}

/*
	type ReturnListResponseOld struct {
		ReturnNo      string   `json:"return_no"`
		ReturnDate    *string  `json:"return_date"`
		SysDate       *string  `json:"sys_date"`
		ReturnType    *int64   `json:"return_type"`
		OutletID      *int64   `json:"outlet_id"`
		OutletCode    string   `json:"outlet_code"`
		OutletName    string   `json:"outlet_name"`
		OutletTaxNo   *string  `json:"outlet_tax_no"`
		PoNo          *string  `json:"po_no"`
		VehicleNo     *string  `json:"vehicle_no"`
		InvoiceNo     *string  `json:"invoice_no"`
		InvoiceDate   *string  `json:"invoice_date"`
		DeliveryDate  *string  `json:"delivery_date"`
		PayType       *int64   `json:"pay_type"`
		SumNo         *string  `json:"sum_no"`
		DataSource    *int64   `json:"data_source"`
		MobileID      *int64   `json:"mobile_id"`
		SubTotal      *float64 `json:"sub_total"`
		Disc          *float64 `json:"disc"`
		DiscValue     *float64 `json:"disc_value"`
		PromoValue    *float64 `json:"promo_value"`
		CashDiscValue *float64 `json:"cash_disc_value"`
		TotDisc1      *float64 `json:"tot_disc_1"`
		TotDisc2      *float64 `json:"tot_disc_2"`
		Vat           *float64 `json:"vat"`
		VatValue      *float64 `json:"vat_value"`
		Total         *float64 `json:"total"`
		DataStatus    *int64   `json:"data_status"`
		UpdatedAt     string   `json:"updated_at"`
		UpdatedByName string   `json:"updated_by_name"`
	}

	type UpdateReturnBody struct {
		CustID        string                   `json:"cust_id"`
		ReturnNo      string                   `json:"return_no"`
		ReturnDate    *string                  `json:"return_date"`
		SysDate       *string                  `json:"sys_date"`
		ReturnType    *int64                   `json:"return_type"`
		OutletID      *int64                   `json:"outlet_id"`
		OutletTaxNo   *string                  `json:"outlet_tax_no"`
		PoNo          *string                  `json:"po_no"`
		VehicleNo     *string                  `json:"vehicle_no"`
		InvoiceNo     *string                  `json:"invoice_no"`
		InvoiceDate   *string                  `json:"invoice_date"`
		DeliveryDate  *string                  `json:"delivery_date"`
		PayType       *int64                   `json:"pay_type"`
		SumNo         *string                  `json:"sum_no"`
		DataSource    *int64                   `json:"data_source"`
		MobileID      *int64                   `json:"mobile_id"`
		SubTotal      *float64                 `json:"sub_total"`
		Disc          *float64                 `json:"disc"`
		DiscValue     *float64                 `json:"disc_value"`
		PromoValue    *float64                 `json:"promo_value"`
		CashDiscValue *float64                 `json:"cash_disc_value"`
		TotDisc1      *float64                 `json:"tot_disc_1"`
		TotDisc2      *float64                 `json:"tot_disc_2"`
		Vat           *float64                 `json:"vat"`
		VatValue      *float64                 `json:"vat_value"`
		Total         *float64                 `json:"total"`
		DataStatus    *int64                   `json:"data_status"`
		UpdatedBy     int64                    `json:"updated_by"`
		Details       UpdateReturnDetBodyGroup `json:"details"`
	}
*/
type DetailReturnParams struct {
	ReturnNo string `params:"return_no" validate:"required"`
}

type UpdateReturnParams struct {
	ReturnNo string `params:"return_no" validate:"required"`
}
type ApproveReturnParams struct {
	ReturnNo string `params:"return_no" validate:"required"`
}

type CancelReturnParams struct {
	ReturnNo string `params:"return_no" validate:"required"`
}

var dataReturnStatusName = map[int64]string{
	1: "In Review",
	2: "Need Review",
	3: "Processed",
	4: "In Pickup",
	5: "Picked Up",
	6: "Completed",
	9: "Cancelled",
}

func (rtn ReturnListResponse) GenerateReturnStatusName() string {
	if rtn.DataStatus != nil {
		return dataReturnStatusName[*rtn.DataStatus]
	}
	return ""
}

func (rtn ReturnShipmentListResponse) GenerateReturnStatusName() string {
	if rtn.DataStatus != nil {
		return dataReturnStatusName[*rtn.DataStatus]
	}
	return ""
}

type ReturnStatusesLookupResponse struct {
	ReturnStatus     *int64  `json:"return_status"`
	ReturnStatusName *string `json:"return_status_name"`
}

func (returnStatus ReturnStatusesLookupResponse) GenerateDataReturnStatusName() string {
	if returnStatus.ReturnStatus != nil {
		return dataReturnStatusName[*returnStatus.ReturnStatus]
	}
	return ""
}

type SalesmansLookupResponse struct {
	SalesmanId   int     `json:"salesman_id"`
	SalesmanCode *string `json:"salesman_code"`
	SalesmanName *string `json:"salesman_name"`
}

type OutletsLookupResponse struct {
	OutletId   int     `json:"outlet_id"`
	OutletCode *string `json:"outlet_code"`
	OutletName *string `json:"outlet_name"`
}

func (rtn ReturnResponse) GenerateReturnStatusName() string {
	if rtn.DataStatus != nil {
		return dataReturnStatusName[*rtn.DataStatus]
	}
	return ""
}

type ProductListResponse struct {
	OrderDetailID *int64   `json:"order_detail_id"`
	InvoiceNo     *string  `json:"invoice_no"`
	InvoiceDate   *string  `json:"invoice_date"`
	ProductID     int64    `json:"product_id"`
	ProductCode   *string  `json:"product_code"`
	ProductName   *string  `json:"product_name"`
	InvoiceQty1   float64  `json:"invoice_qty1"`
	InvoiceQty2   float64  `json:"invoice_qty2"`
	InvoiceQty3   float64  `json:"invoice_qty3"`
	RemainingQty1 float64  `json:"remaining_qty1"`
	RemainingQty2 float64  `json:"remaining_qty2"`
	RemainingQty3 float64  `json:"remaining_qty3"`
	SellPrice1    float64  `json:"sell_price1"`
	SellPrice2    float64  `json:"sell_price2"`
	SellPrice3    float64  `json:"sell_price3"`
	SubTotal1     float64  `json:"sub_total1"`
	SubTotal2     float64  `json:"sub_total2"`
	SubTotal3     float64  `json:"sub_total3"`
	UnitId1       string   `json:"unit_id1"`
	UnitId2       string   `json:"unit_id2"`
	UnitId3       string   `json:"unit_id3"`
	UnitName1     *string  `json:"unit_name1"`
	UnitName2     *string  `json:"unit_name2"`
	UnitName3     *string  `json:"unit_name3"`
	ConvUnit2     *float64 `json:"conv_unit2"`
	ConvUnit3     *float64 `json:"conv_unit3"`
	Vat           *float64 `json:"vat"`
	Total         float64  `json:"total"`
}

type ProductListQueryFilter struct {
	SalesmanId   []int  `query:"salesman_id" validate:"required"`
	OutletID     []int  `query:"outlet_id" validate:"required"`
	SearchBy     string `query:"search_by" validate:"required"`
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type ReturnReasonsLookupResponse struct {
	ReturnReasonId   int     `json:"return_reason_id"`
	ReturnReasonName *string `json:"return_reason_name"`
}

type ProductConditionsLookupResponse struct {
	ProductConditionId   int    `json:"product_condition_id"`
	ProductConditionName string `json:"product_condition_name"`
}

type WarehousesLookupResponse struct {
	WhId   int     `json:"wh_id"`
	WhCode *string `json:"wh_code"`
	WhName *string `json:"wh_name"`
}

var dataProductConditionName = map[int64]string{
	1: "Good",
	2: "Bad",
	3: "Expired",
}

/*
func (productCon ProductConditionsLookupResponse) GetProductConditionList() (proCons map[int64]string) {

		// var proCon ProductConditionsLookupResponse
		// for id, dataProductCondition := range dataProductConditionName {
		// 	proCon.ProductConditionId = int(id)
		// 	proCon.ProductConditionName = dataProductCondition

		// 	productConditions = append(productConditions, proCon)
		// }
		proCons = dataProductConditionName

		return proCons
	}
*/
type ProductsLookupCreateResponse struct {
	ProductID   *int64   `json:"product_id"`
	ProductCode *string  `json:"product_code"`
	ProductName *string  `json:"product_name"`
	SellPrice1  *float64 `json:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3"`
	UnitId1     *string  `json:"unit_id1"`
	UnitId2     *string  `json:"unit_id2"`
	UnitId3     *string  `json:"unit_id3"`
	UnitName1   *string  `json:"unit_name1"`
	UnitName2   *string  `json:"unit_name2"`
	UnitName3   *string  `json:"unit_name3"`
	ConvUnit2   *float64 `json:"conv_unit2"`
	ConvUnit3   *float64 `json:"conv_unit3"`
	Vat         *float64 `json:"vat"`
}

type CreateReturnBody struct {
	CustID       string  `json:"cust_id"`
	RefferenceNo *string `json:"refference_no"`
	ReturnNo     *string `json:"return_no"`
	// ReturnDate    *string                  `json:"return_date"`
	// InvoiceNo     *string                  `json:"invoice_no"`
	// InvoiceDate   *string                  `json:"invoice_date"`
	// SalesmanID    *int64                   `json:"salesman_id"`
	// OutletID      *int64                   `json:"outlet_id"`
	// TprCash       *float64                 `json:"tpr_cash"`
	// TprCashValue  *float64                 `json:"tpr_cash_value"`
	// TprItem       *float64                 `json:"tpr_item"`
	// TprItemValue  *float64                 `json:"tpr_item_value"`
	// Discount      *float64                 `json:"discount"`
	// DiscountValue *float64                 `json:"discount_value"`
	// Vat           *float64                 `json:"vat"`
	// VatValue      *float64                 `json:"vat_value"`
	// Total         *float64                 `json:"total"`
	// DataStatus    *int64                   `json:"data_status"`
	CreatedBy *int64 `json:"created_by"`
	// CreatedAt *string                  `json:"created_at"`
	Details []CreateReturnDetailBody `json:"details"`
}

type UpdateReturnBody struct {
	CustID    string                   `json:"cust_id"`
	UpdatedBy int64                    `json:"updated_by"`
	Details   []UpdateReturnDetailBody `json:"details"`
}

type ApproveReturnBody struct {
	CustID    string                    `json:"cust_id"`
	UpdatedBy int64                     `json:"updated_by"`
	Details   []ApproveReturnDetailBody `json:"details"`
}

type CancelReturnBody struct {
	CustID    string `json:"cust_id"`
	UpdatedBy int64  `json:"updated_by"`
}

type UpdateStatusReturnBody struct {
	CustID    string                         `json:"cust_id"`
	UpdatedBy int64                          `json:"updated_by"`
	Returns   []UpdateStatusReturnDetailBody `json:"returns"`
}

type UpdateQuantityReturnBody struct {
	CustID    string                     `json:"cust_id"`
	UpdatedBy int64                      `json:"updated_by"`
	Details   []UpdateReturnQuantityBody `json:"details"`
}
