package entity

// var dataStatusName = map[int64]string{
// 	1: "Need Review",
// 	2: "Processed",
// 	3: "Cancelled",
// 	4: "Invoicing",
// }

// var payTypeName = map[int64]string{
// 	1: "Cash",
// 	2: "Check",
// 	3: "Transfer",
// 	4: "Credit",
// }

type RoQueryFilter struct {
	SalesmanId   []int `query:"salesman_id"`
	OutletID     []int `query:"outlet_id"`
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

type CreateRoBody struct {
	CustId        string         `json:"cust_id"`
	RoDate        *string        `json:"ro_date"`
	ValDate       *string        `json:"val_date"`
	DueDate       *string        `json:"due_date"`
	SalesmanId    *int64         `json:"salesman_id"`
	WhId          *int64         `json:"wh_id"`
	OutletID      *int64         `json:"outlet_id"`
	DeliveryDate  *string        `json:"delivery_date"`
	OrderNo       *string        `json:"order_no"`
	PoNo          *string        `json:"po_no"`
	VehicleNo     *string        `json:"vehicle_no"`
	PayType       *int64         `json:"pay_type"`
	ReffNo        *string        `json:"reff_no"`
	MobileID      *int64         `json:"mobile_id"`
	SubTotal      *float64       `json:"sub_total"`
	Disc          *float64       `json:"disc"`
	DiscValue     *float64       `json:"disc_value"`
	PromoValue    *float64       `json:"promo_value"`
	CashDiscValue *float64       `json:"cash_disc_value"`
	TotDisc1      *float64       `json:"tot_disc1"`
	TotDisc2      *float64       `json:"tot_disc2"`
	Vat           *float64       `json:"vat"`
	VatValue      *float64       `json:"vat_value"`
	Total         *float64       `json:"total"`
	DataStatus    *int64         `json:"data_status"`
	CreatedBy     *int64         `json:"created_by"`
	DataSource    *int64         `json:"data_source"`
	Details       RoDetWithGroup `json:"details"`
}

type RoResponse struct {
	RoNo           string             `json:"ro_no"`
	RoDate         *string            `json:"ro_date"`
	ValDate        *string            `json:"val_date"`
	SalesmanId     *int64             `json:"salesman_id"`
	SalesName      string             `json:"sales_name"`
	WhId           *int64             `json:"wh_id"`
	WhCode         string             `json:"wh_code"`
	WhName         string             `json:"wh_name"`
	OutletID       *int64             `json:"outlet_id"`
	OutletCode     string             `json:"outlet_code"`
	OutletName     string             `json:"outlet_name"`
	DeliveryDate   *string            `json:"delivery_date"`
	OrderNo        *string            `json:"order_no"`
	PoNo           *string            `json:"po_no"`
	VehicleNo      *string            `json:"vehicle_no"`
	PayType        *int64             `json:"pay_type"`
	PayTypeName    string             `json:"pay_type_name"`
	ReffNo         *string            `json:"reff_no"`
	MobileID       *int64             `json:"mobile_id"`
	SubTotal       *float64           `json:"sub_total"`
	Disc           *float64           `json:"disc"`
	DiscValue      *float64           `json:"disc_value"`
	PromoValue     *float64           `json:"promo_value"`
	CashDiscValue  *float64           `json:"cash_disc_value"`
	TotDisc1       *float64           `json:"tot_disc1"`
	TotDisc2       *float64           `json:"tot_disc2"`
	Vat            *float64           `json:"vat"`
	VatValue       *float64           `json:"vat_value"`
	Total          *float64           `json:"total"`
	DataStatus     *int64             `json:"data_status"`
	DataStatusName string             `json:"data_status_name"`
	DataSource     *int64             `json:"data_source"`
	UpdatedAt      string             `json:"updated_at"`
	UpdatedByName  string             `json:"updated_by_name"`
	DueDate        *string            `json:"due_date"`
	Details        RoDetReadWithGroup `json:"details"`
}

func (ro RoResponse) GeneratePayTypeName() string {
	if ro.PayType != nil {
		return payTypeName[*ro.PayType]
	}
	return ""
}

func (ro RoResponse) GenerateDataStatusName() string {
	if ro.DataStatus != nil {
		return dataStatusName[*ro.DataStatus]
	}
	return ""
}

type DetailRoParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}
type DeleteRoParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type UpdateRoParams struct {
	RoNo string `params:"ro_no" validate:"required"`
}

type RoListResponse struct {
	RoNo           string   `json:"ro_no"`
	RoDate         *string  `json:"ro_date"`
	ValDate        *string  `json:"val_date"`
	SalesmanId     *int64   `json:"salesman_id"`
	SalesName      string   `json:"sales_name"`
	WhId           *int64   `json:"wh_id"`
	WhCode         string   `json:"wh_code"`
	WhName         string   `json:"wh_name"`
	OutletID       *int64   `json:"outlet_id"`
	OutletCode     string   `json:"outlet_code"`
	OutletName     string   `json:"outlet_name"`
	DeliveryDate   *string  `json:"delivery_date"`
	OrderNo        *string  `json:"order_no"`
	PoNo           *string  `json:"po_no"`
	VehicleNo      *string  `json:"vehicle_no"`
	PayType        *int64   `json:"pay_type"`
	PayTypeName    string   `json:"pay_type_name"`
	ReffNo         *string  `json:"reff_no"`
	MobileID       *int64   `json:"mobile_id"`
	SubTotal       *float64 `json:"sub_total"`
	Disc           *float64 `json:"disc"`
	DiscValue      *float64 `json:"disc_value"`
	PromoValue     *float64 `json:"promo_value"`
	CashDiscValue  *float64 `json:"cash_disc_value"`
	TotDisc1       *float64 `json:"tot_disc1"`
	TotDisc2       *float64 `json:"tot_disc2"`
	Vat            *float64 `json:"vat"`
	VatValue       *float64 `json:"vat_value"`
	Total          *float64 `json:"total"`
	DataStatus     *int64   `json:"data_status"`
	DataStatusName string   `json:"data_status_name"`
	UpdatedAt      string   `json:"updated_at"`
	UpdatedByName  string   `json:"updated_by_name"`
	DueDate        *string  `json:"due_date"`
}

func (ro RoListResponse) GenerateDataStatusName() string {
	if ro.DataStatus != nil {
		return dataStatusName[*ro.DataStatus]
	}
	return ""
}

func (ro RoListResponse) GeneratePayTypeName() string {
	if ro.PayType != nil {
		return payTypeName[*ro.PayType]
	}
	return ""
}

type UpdateRoBody struct {
	CustId        string               `json:"cust_id"`
	RoNo          string               `json:"ro_no"`
	RoDate        *string              `json:"ro_date"`
	ValDate       *string              `json:"val_date"`
	DueDate       *string              `json:"due_date"`
	SalesmanId    *int64               `json:"salesman_id"`
	WhId          *int64               `json:"wh_id"`
	OutletID      *int64               `json:"outlet_id"`
	DeliveryDate  *string              `json:"delivery_date"`
	OrderNo       *string              `json:"order_no"`
	PoNo          *string              `json:"po_no"`
	VehicleNo     *string              `json:"vehicle_no"`
	PayType       *int64               `json:"pay_type"`
	ReffNo        *string              `json:"reff_no"`
	MobileID      *int64               `json:"mobile_id"`
	SubTotal      *float64             `json:"sub_total"`
	Disc          *float64             `json:"disc"`
	DiscValue     *float64             `json:"disc_value"`
	PromoValue    *float64             `json:"promo_value"`
	CashDiscValue *float64             `json:"cash_disc_value"`
	TotDisc1      *float64             `json:"tot_disc1"`
	TotDisc2      *float64             `json:"tot_disc2"`
	Vat           *float64             `json:"vat"`
	VatValue      *float64             `json:"vat_value"`
	Total         *float64             `json:"total"`
	DataStatus    *int64               `json:"data_status"`
	CreatedBy     *int64               `json:"created_by"`
	CreatedAt     *string              `json:"created_at"`
	UpdatedBy     int64                `json:"updated_by"`
	Details       UpdateRoDetWithGroup `json:"details"`
}
