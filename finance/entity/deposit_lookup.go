package entity

import "time"

type DepositLookupQueryFilter struct {
	CustId       string
	ParentCustId string
	SalesmanId   int64    `query:"salesman_id" validate:"required"`
	From         *int64   `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64   `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int      `query:"page"`
	Limit        int      `query:"limit" validate:""`
	Query        string   `query:"q"`
	Mode         []string `query:"mode"`
	Sort         string   `query:"sort"`
	IsActive     *int     `query:"is_active"`
	TrCode       string   `query:"tr_code"`
	CollectionNo *string  `query:"collection_no"`
	OutletId     []int    `query:"outlet_id"`
}

type CreateDepositLookupBody struct {
	CustID           string  `json:"cust_id"`
	DocNoBank        string  `json:"doc_no_bank"`
	OwnerID          int     `json:"owner_id"`
	SalesmanID       *int    `json:"salesman_id"`
	SupplierID       *int    `json:"sup_id"`
	OutletID         *int    `json:"outlet_id"`
	BankID           int     `json:"bank_id"`
	BankIDCollecting int     `json:"bank_id_collecting"`
	AccountNo        *string `json:"account_no"`
	TransferDate     *string `json:"transfer_date"`
	Amount           float64 `json:"amount"`
	StatusBank       int     `json:"status_bank_transfer"`
	CreatedBy        *int64  `json:"created_by"`
}

type DepositLookupResponse struct {
	CustID           string    `json:"cust_id"`
	DepositLookupNo  int       `json:"bank_transfer_no"`
	DocNoBank        string    `json:"doc_no_bank"`
	OwnerID          int       `json:"owner_id"`
	OwnerName        string    `json:"owner_name"`
	SupplierID       *int      `json:"sup_id"`
	SupplierName     *int      `json:"sup_name"`
	SalesmanID       *int      `json:"salesman_id"`
	SalesmanName     *string   `json:"sales_name"`
	OutletID         *int      `json:"outlet_id"`
	OutletName       *string   `json:"outlet_name"`
	BankID           int       `json:"bank_id"`
	BankName         string    `json:"bank_name"`
	BankIDCollecting int       `json:"bank_id_collecting"`
	AccountNo        *string   `json:"account_no"`
	TransferDate     *string   `json:"transfer_date"`
	Amount           float64   `json:"amount"`
	UsedAmount       float64   `json:"used_amount"`
	RemainingAmount  float64   `json:"remaining_amount"`
	StatusBank       int       `json:"status_bank_transfer"`
	StatusBankText   *string   `json:"status_bank_transfer_text"`
	CreatedBy        int64     `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedBy        int64     `json:"updated_by"`
	UpdatedAt        time.Time `json:"updated_at"`

	DepositStatus     int    `json:"deposit_status"`
	DepositStatusName string `json:"deposit_status_name"`
}

type UpdateDepositLookupBody struct {
	CustID           string  `json:"cust_id"`
	DocNoBank        string  `json:"doc_no_bank"`
	OwnerID          int     `json:"owner_id"`
	SalesmanID       *int    `json:"salesman_id"`
	SupplierID       *int    `json:"sup_id"`
	OutletID         *int    `json:"outlet_id"`
	BankID           int     `json:"bank_id"`
	BankIDCollecting int     `json:"bank_id_collecting"`
	AccountNo        *string `json:"account_no"`
	TransferDate     *string `json:"transfer_date"`
	Amount           float64 `json:"amount"`
	StatusBank       int     `json:"status_bank_transfer"`
	CreatedBy        *int64  `json:"created_by"`
	UpdatedBy        int64   `json:"updated_by"`
}

type DetailDepositLookupParams struct {
	DepositLookupNo int `params:"bank_transfer_no" validate:"required"`
}
type DeleteDepositLookupParams struct {
	DepositLookupNo int `params:"bank_transfer_no" validate:"required"`
}
type UpdateDepositLookupParams struct {
	DepositLookupNo int `params:"bank_transfer_no" validate:"required"`
}

type BankLookupDepositLookup struct {
	BankId   int    `json:"bank_id"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

type CollectionNoLookup struct {
	CollectionNo string `json:"collection_no"`
}
type DepositNoLookup struct {
	DepositNo string `json:"deposit_no"`
}
type DepositStatusLookup struct {
	DepositStatus     int    `json:"deposit_status"`
	DepositStatusName string `json:"deposit_status_name"`
}

var StatusDeposit = map[int]string{
	1: "In Review",
	2: "Approved",
	3: "Rejected",
}

func ConvDepositStatus(data map[int]string, param int) string {
	statusString, ok := data[int(param)]
	if !ok {
		statusString = "Unknown"
	}
	return statusString
}

type InvoiceByCollectionResponse struct {
	// WhId            *int64   `json:"wh_id"`
	// WhCode          string   `json:"wh_code"`
	// WhName          string   `json:"wh_name"`
	// WhLatitude      string   `json:"wh_latitude"`
	// WhLongitude     string   `json:"wh_longitude"`
	// OutletAddress   string   `json:"outlet_address"`
	// OutletLatitude  string   `json:"outlet_latitude"`
	// OutletLongitude string   `json:"outlet_longitude"`
	// DeliveryDate    *string  `json:"delivery_date"`
	// OrderNo         *string  `json:"order_no"`
	// PoNo            *string  `json:"po_no"`
	// VehicleNo       *string  `json:"vehicle_no"`
	// PayType         *int64   `json:"pay_type"`
	// PayTypeName     string   `json:"pay_type_name"`
	// ReffNo          *string  `json:"reff_no"`
	// MobileID        *int64   `json:"mobile_id"`
	// SubTotal        *float64 `json:"sub_total"`
	// Disc            *float64 `json:"disc"`
	// DiscValue       *float64 `json:"disc_value"`
	// PromoValue      *float64 `json:"promo_value"`
	// CashDiscValue   *float64 `json:"cash_disc_value"`
	// TotDisc1        *float64 `json:"tot_disc1"`
	// TotDisc2        *float64 `json:"tot_disc2"`
	// Vat             *float64 `json:"vat"`
	// VatValue        *float64 `json:"vat_value"`
	// Total           *float64 `json:"total"`
	// DataStatus      *int64   `json:"data_status"`
	// DataStatusName  string   `json:"data_status_name"`
	// DataSource      *int64   `json:"data_source"`
	// DueDate         *string  `json:"due_date"`

	CollectionNo        *string  `json:"collection_no"`
	InvoiceNo           *string  `json:"invoice_no"`
	InvoiceDate         *string  `json:"invoice_date"`
	RoNo                *string  `json:"ro_no"`
	OutletID            *int64   `json:"outlet_id"`
	OutletCode          *string  `json:"outlet_code"`
	OutletName          *string  `json:"outlet_name"`
	SalesmanId          *int64   `json:"salesman_id"`
	SalesmanCode        *string  `json:"salesman_code"`
	SalesmanName        *string  `json:"salesman_name"`
	InvoiceAmount       *float64 `json:"invoice_amount"`
	TotalInvoicePaymnet *float64 `json:"total_invoice_payment" default:"0.0"`
	RemainingAmount     *float64 `json:"remaining_amount"`
}

type DepositPaymentLookup struct {
	DocNo   string  `json:"doc_no"`
	Amount  float64 `json:"amount"`
	Balance float64 `json:"balance"`
}

// func (invoice InvoiceResponse) GeneratePayTypeName() string {
// 	if invoice.PayType != nil {
// 		return payTypeName[*invoice.PayType]
// 	}
// 	return ""
// }

// func (invoice InvoiceResponse) GenerateDataStatusName() string {
// 	if invoice.DataStatus != nil {
// 		return dataStatusName[*invoice.DataStatus]
// 	}
// 	return ""
// }
