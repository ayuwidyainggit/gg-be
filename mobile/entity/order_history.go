package entity

type OrderHistoryListResponse struct {
	RoNo   string  `json:"ro_no"`
	RoDate *string `json:"ro_date"`
	// ValDate    *string  `json:"val_date"`
	SalesmanId *int64           `json:"salesman_id"`
	SalesName  string           `json:"sales_name"`
	WhId       *int64           `json:"wh_id"`
	WhCode     string           `json:"wh_code"`
	WhName     string           `json:"wh_name"`
	OutletID   *int64           `json:"outlet_id"`
	OutletCode string           `json:"outlet_code"`
	OutletName string           `json:"outlet_name"`
	SubTotal   *float64         `json:"sub_total"`
	IsPayment  *bool            `json:"is_payment"`
	ProductImg []ProductImgRead `json:"product_img" validate:"dive"`
}

type ProductImgRead struct {
	ImgUrl *string `json:"img_url"`
}

type OrderHistoryDetailResponse struct {
	RoNo           string                       `json:"ro_no"`
	OrderNo        *string                      `json:"order_no"`
	RoDate         *string                      `json:"ro_date"`
	ValDate        *string                      `json:"val_date"`
	SalesmanId     *int64                       `json:"salesman_id"`
	SalesName      string                       `json:"sales_name"`
	WhId           *int64                       `json:"wh_id"`
	WhCode         string                       `json:"wh_code"`
	WhName         string                       `json:"wh_name"`
	OutletID       *int64                       `json:"outlet_id"`
	OutletCode     string                       `json:"outlet_code"`
	OutletName     string                       `json:"outlet_name"`
	OutletAddress1 string                       `json:"outlet_address1"`
	OutletAddress2 string                       `json:"outlet_address2"`
	DeliveryDate   *string                      `json:"delivery_date"`
	PoNo           *string                      `json:"po_no"`
	VehicleNo      *string                      `json:"vehicle_no"`
	PayType        *int64                       `json:"pay_type"`
	PayTypeName    string                       `json:"pay_type_name"`
	ReffNo         *string                      `json:"reff_no"`
	MobileID       *int64                       `json:"mobile_id"`
	SubTotal       *float64                     `json:"sub_total"`
	Disc           *float64                     `json:"disc"`
	DiscValue      *float64                     `json:"disc_value"`
	PromoValue     *float64                     `json:"promo_value"`
	CashDiscValue  *float64                     `json:"cash_disc_value"`
	TotDisc1       *float64                     `json:"tot_disc1"`
	TotDisc2       *float64                     `json:"tot_disc2"`
	Vat            *float64                     `json:"vat"`
	VatValue       *float64                     `json:"vat_value"`
	Total          *float64                     `json:"total"`
	DataStatus     *int64                       `json:"data_status"`
	DataStatusName string                       `json:"data_status_name"`
	DataSource     *int64                       `json:"data_source"`
	UpdatedAt      string                       `json:"updated_at"`
	UpdatedByName  string                       `json:"updated_by_name"`
	DueDate        *string                      `json:"due_date"`
	ProducDetails  []OrderHistoryProductDetails `json:"product_details"`
	TrCode         *string                      `json:"tr_code"`
	IsClosed       bool                         `json:"is_closed"`
	Notes          *string                      `json:"notes"`
	InvoiceNo      *string                      `json:"invoice_no"`
	InvoiceDate    *string                      `json:"invoice_date"`
	IsPrinted      *bool                        `json:"is_printed"`
	PrintedBy      *int64                       `json:"printed_by"`
	PrintedByName  *string                      `json:"printed_by_name"`
	PrintedAt      *string                      `json:"printed_at"`
}

type OrderHistoryProductDetails struct {
	ProId      int64   `json:"pro_id"`
	ProCode    string  `json:"pro_code"`
	ProName    string  `json:"pro_name"`
	BarCode    string  `json:"bar_code"`
	UnitId1    string  `json:"unit_id1"`
	UnitId2    string  `json:"unit_id2"`
	UnitId3    string  `json:"unit_id3"`
	UnitId4    string  `json:"unit_id4"`
	UnitId5    string  `json:"unit_id5"`
	UnitName1  string  `json:"unit_name1"`
	UnitName2  string  `json:"unit_name2"`
	UnitName3  string  `json:"unit_name3"`
	UnitName4  string  `json:"unit_name4"`
	UnitName5  string  `json:"unit_name5"`
	ConvUnit1  float32 `json:"conv_unit1"`
	ConvUnit2  float32 `json:"conv_unit2"`
	ConvUnit3  float32 `json:"conv_unit3"`
	ConvUnit4  float32 `json:"conv_unit4"`
	ConvUnit5  float32 `json:"conv_unit5"`
	Stock1     float64 `json:"stock1"`
	Stock2     float64 `json:"stock2"`
	Stock3     float64 `json:"stock3"`
	Stock4     float64 `json:"stock4"`
	Stock5     float64 `json:"stock5"`
	PoFormula  int     `json:"po_formula"`
	Vat        float64 `json:"vat"`
	VatBg      float64 `json:"vat_bg"`
	VatLgPurch float64 `json:"vat_lg_purch"`
	VatLgSell  float64 `json:"vat_lg_sell"`
	ExciseRate float64 `json:"excise_rate"`
	ExciseTax  float64 `json:"excise_tax"`
	ImageUrl   string  `json:"image_url"`
	Cogs       float64 `json:"cogs"`
	Price1     float64 `json:"price1"`
	Price2     float64 `json:"price2"`
	Price3     float64 `json:"price3"`
	Price4     float64 `json:"price4"`
	Price5     float64 `json:"price5"`
	PromoCode  string  `json:"promo_code"`
	CtgId1     string  `json:"ctg_id1"`
	CtgId2     string  `json:"ctg_id2"`
	CtgId3     string  `json:"ctg_id3"`
}

func (ro OrderHistoryDetailResponse) GeneratePayTypeNameHistory() string {
	if ro.PayType != nil {
		return payTypeName[*ro.PayType]
	}
	return ""
}

func (ro OrderHistoryDetailResponse) GenerateDataStatusNameHistory() string {
	if ro.DataStatus != nil {
		return dataStatusName[*ro.DataStatus]
	}
	return ""
}
