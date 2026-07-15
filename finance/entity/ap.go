package entity

type CreateApBody struct {
	CustID            string                   `json:"cust_id"`
	ParentCustID      string                   `json:"parent_cust_id"`
	ApDate            *string                  `json:"ap_date" validate:"required"`
	TrCode            *string                  `json:"tr_code" validate:"required,len=3"`
	SupID             *int64                   `json:"sup_id" validate:"required"`
	InvNo             *string                  `json:"inv_no" validate:"required,max=30"`
	InvDate           *string                  `json:"inv_date" validate:"required"`
	InvDueDate        *string                  `json:"inv_due_date"`
	TaxInvNo          *string                  `json:"tax_inv_no" validate:"max=30"`
	TaxInvDate        *string                  `json:"tax_inv_date"`
	TaxReturnNo       *string                  `json:"tax_return_no" validate:"max=30"`
	TaxReturnDate     *string                  `json:"tax_return_date"`
	SubTotal          *float64                 `json:"sub_total" validate:"required"`
	MoneyPromo        *float64                 `json:"money_promo"`
	InvDisc           *float64                 `json:"inv_disc"`
	InvDiscValue      *float64                 `json:"inv_disc_value"`
	SubTotalBtax      *float64                 `json:"sub_total_btax"`
	Vat               *float64                 `json:"vat"`
	VatValue          *float64                 `json:"vat_value"`
	VatLg             *float64                 `json:"vat_lg"`
	VatLgValue        *float64                 `json:"vat_lg_value"`
	VatBg             *float64                 `json:"vat_bg"`
	VatBgValue        *float64                 `json:"vat_bg_value"`
	DutyStamp         *float64                 `json:"duty_stamp"`
	Total             *float64                 `json:"total" validate:"required"`
	TotalDiff         *float64                 `json:"total_diff"`
	TotalRound        *float64                 `json:"total_round"`
	ApPaID            *float64                 `json:"ap_paid"`
	DataStatus        *int64                   `json:"data_status"`
	CreatedBy         *int64                   `json:"created_by"`
	IsPosted          *bool                    `json:"is_posted"`
	Details           []CreateApDetBody        `json:"details"`
	QtyPromoDetails   []CreateApQtyPromoBody   `json:"qty_promo_details"`
	MoneyPromoDetails []CreateApMoneyPromoBody `json:"money_promo_details"`
}
type ApResponse struct {
	ApNo              string                 `json:"ap_no"`
	ApDate            *string                `json:"ap_date"`
	TrCode            *string                `json:"tr_code"`
	SupID             *int64                 `json:"sup_id"`
	SupCode           string                 `json:"sup_code"`
	SupName           string                 `json:"sup_name"`
	InvNo             *string                `json:"inv_no"`
	InvDate           *string                `json:"inv_date"`
	InvDueDate        *string                `json:"inv_due_date"`
	TaxInvNo          *string                `json:"tax_inv_no"`
	TaxInvDate        *string                `json:"tax_inv_date"`
	TaxReturnNo       *string                `json:"tax_return_no"`
	TaxReturnDate     *string                `json:"tax_return_date"`
	SubTotal          *float64               `json:"sub_total"`
	MoneyPromo        *float64               `json:"money_promo"`
	InvDisc           *float64               `json:"inv_disc"`
	InvDiscValue      *float64               `json:"inv_disc_value"`
	SubTotalBtax      *float64               `json:"sub_total_btax"`
	Vat               *float64               `json:"vat"`
	VatValue          *float64               `json:"vat_value"`
	VatLg             *float64               `json:"vat_lg"`
	VatLgValue        *float64               `json:"vat_lg_value"`
	VatBg             *float64               `json:"vat_bg"`
	VatBgValue        *float64               `json:"vat_bg_value"`
	DutyStamp         *float64               `json:"duty_stamp"`
	Total             *float64               `json:"total"`
	TotalDiff         *float64               `json:"total_diff"`
	TotalRound        *float64               `json:"total_round"`
	ApPaID            *float64               `json:"ap_paid"`
	DataStatus        *int64                 `json:"data_status"`
	UpdatedByName     string                 `json:"updated_by_name"`
	UpdatedAt         string                 `json:"updated_at"`
	IsPosted          *bool                  `json:"is_posted"`
	Details           []ApDetResponse        `json:"details"`
	QtyPromoDetails   []ApQtyPromoResponse   `json:"qty_promo_details"`
	MoneyPromoDetails []ApMoneyPromoResponse `json:"money_promo_details"`
}
type ApListResponse struct {
	ApNo          string   `json:"ap_no"`
	ApDate        *string  `json:"ap_date"`
	TrCode        *string  `json:"tr_code"`
	SupID         *int64   `json:"sup_id"`
	SupCode       string   `json:"sup_code"`
	SupName       string   `json:"sup_name"`
	InvNo         *string  `json:"inv_no"`
	InvDate       *string  `json:"inv_date"`
	InvDueDate    *string  `json:"inv_due_date"`
	TaxInvNo      *string  `json:"tax_inv_no"`
	TaxInvDate    *string  `json:"tax_inv_date"`
	TaxReturnNo   *string  `json:"tax_return_no"`
	TaxReturnDate *string  `json:"tax_return_date"`
	SubTotal      *float64 `json:"sub_total"`
	MoneyPromo    *float64 `json:"money_promo"`
	InvDisc       *float64 `json:"inv_disc"`
	InvDiscValue  *float64 `json:"inv_disc_value"`
	SubTotalBtax  *float64 `json:"sub_total_btax"`
	Vat           *float64 `json:"vat"`
	VatValue      *float64 `json:"vat_value"`
	VatLg         *float64 `json:"vat_lg"`
	VatLgValue    *float64 `json:"vat_lg_value"`
	VatBg         *float64 `json:"vat_bg"`
	VatBgValue    *float64 `json:"vat_bg_value"`
	DutyStamp     *float64 `json:"duty_stamp"`
	Total         *float64 `json:"total"`
	TotalDiff     *float64 `json:"total_diff"`
	TotalRound    *float64 `json:"total_round"`
	ApPaID        *float64 `json:"ap_paid"`
	DataStatus    *int64   `json:"data_status"`
	CreatedBy     *int64   `json:"created_by"`
	UpdatedByName string   `json:"updated_by_name"`
	UpdatedAt     string   `json:"updated_at"`
	IsPosted      *bool    `json:"is_posted"`
}

type UpdateApBody struct {
	CustID            string                   `json:"cust_id"`
	ParentCustID      string                   `json:"parent_cust_id"`
	ApDate            *string                  `json:"ap_date"`
	TrCode            *string                  `json:"tr_code" validate:"len=3"`
	SupID             *int64                   `json:"sup_id"`
	InvNo             *string                  `json:"inv_no"`
	InvDate           *string                  `json:"inv_date"`
	InvDueDate        *string                  `json:"inv_due_date"`
	TaxInvNo          *string                  `json:"tax_inv_no"`
	TaxInvDate        *string                  `json:"tax_inv_date"`
	TaxReturnNo       *string                  `json:"tax_return_no"`
	TaxReturnDate     *string                  `json:"tax_return_date"`
	SubTotal          *float64                 `json:"sub_total"`
	MoneyPromo        *float64                 `json:"money_promo"`
	InvDisc           *float64                 `json:"inv_disc"`
	InvDiscValue      *float64                 `json:"inv_disc_value"`
	SubTotalBtax      *float64                 `json:"sub_total_btax"`
	Vat               *float64                 `json:"vat"`
	VatValue          *float64                 `json:"vat_value"`
	VatLg             *float64                 `json:"vat_lg"`
	VatLgValue        *float64                 `json:"vat_lg_value"`
	VatBg             *float64                 `json:"vat_bg"`
	VatBgValue        *float64                 `json:"vat_bg_value"`
	DutyStamp         *float64                 `json:"duty_stamp"`
	Total             *float64                 `json:"total"`
	TotalDiff         *float64                 `json:"total_diff"`
	TotalRound        *float64                 `json:"total_round"`
	ApPaID            *float64                 `json:"ap_paid"`
	DataStatus        *int64                   `json:"data_status"`
	CreatedBy         *int64                   `json:"created_by"`
	UpdatedBy         int64                    `json:"updated_by"`
	IsPosted          *bool                    `json:"is_posted"`
	Details           []UpdateApDetBody        `json:"details"`
	QtyPromoDetails   []UpdateApQtyPromoBody   `json:"qty_promo_details"`
	MoneyPromoDetails []UpdateApMoneyPromoBody `json:"money_promo_details"`
}
type DetailApParams struct {
	ApNo string `params:"ap_no" validate:"required" json:"ap_no"`
}

type DeleteApParams struct {
	ApNo string `params:"ap_no" validate:"required" json:"ap_no"`
}
type UpdateApParams struct {
	ApNo string `params:"ap_no" validate:"required" json:"ap_no"`
}
